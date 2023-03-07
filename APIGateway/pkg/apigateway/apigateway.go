package apigateway

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	censor "NewsRSS/APIGateway/pkg/client/censorservice"
	comments "NewsRSS/APIGateway/pkg/client/commentservice"
	news "NewsRSS/APIGateway/pkg/client/newsservice"

	"github.com/gorilla/mux"
)

type API struct {
	r              *mux.Router
	commentService comments.CommentService
	newsService    news.NewsService
	censorService  censor.CensorService
}

// Конструктор API.
func New(comment comments.CommentService, news news.NewsService, censor censor.CensorService) *API {
	a := API{r: mux.NewRouter()}
	a.commentService = comment
	a.newsService = news
	a.censorService = censor
	a.endpoints()
	return &a
}

// Router возвращает маршрутизатор для использования
// в качестве аргумента HTTP-сервера.
func (api *API) Router() *mux.Router {
	return api.r
}

// Регистрация методов API в маршрутизаторе запросов.
func (api *API) endpoints() {
	api.r.Use(api.headersMiddleware)
	api.r.Use(api.logMiddleware)
	// получить новости постранично page страница s фильтр для поиска по заголовку
	api.r.HandleFunc("/news", api.newsPage).Methods(http.MethodGet, http.MethodOptions)
	// получить новости по фильтру filter?sort=date&direction=desc&count=10&offset=0
	api.r.HandleFunc("/news/filter", api.newsFilter).Methods(http.MethodGet, http.MethodOptions)
	// получить детальную новость по id
	api.r.HandleFunc("/news/{id}", api.newsDetail).Methods(http.MethodGet, http.MethodOptions)
	// добавить комментарий к новости
	api.r.HandleFunc("/news/{id}", api.addComment).Methods(http.MethodPost, http.MethodOptions)

}

// newsPage вывод новостей списком
func (api *API) newsPage(w http.ResponseWriter, r *http.Request) {
	// если параметр был передан, вернется строка со значением.
	// Если не был - в переменной будет пустая строка
	pageParam := r.URL.Query().Get("page")
	if pageParam == "" {
		pageParam = "1"
	}
	// параметр page - это число, поэтому нужно сконвертировать
	// строку в число при помощи пакета strconv
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		// обработка ошибки
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filter := r.URL.Query().Get("s")
	ctx := r.Context()
	// получение комментариев к новости
	data, err := api.newsService.Posts(ctx, page, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	// далее запрос к бд с получение послендних новостей на странице page
	//json.NewEncoder(w).Encode(page)
}

// newsFilter получение новостей по фильтру с параметрами get.
func (api *API) newsFilter(w http.ResponseWriter, r *http.Request) {
	sortParam := r.URL.Query().Get("sort")
	directionParam := r.URL.Query().Get("direction")
	countParam := r.URL.Query().Get("count")
	count, err := strconv.Atoi(countParam)
	if err != nil {
		// обработка ошибки
		count = 10
	}
	offsetParam := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		// обработка ошибки
		offset = 0
	}

	// далее запрос к бд с получение новостей по фильтру
	json.NewEncoder(w).Encode(fmt.Sprintf("%s %s %d %d", sortParam, directionParam, count, offset))
}

// newsDetail детализированная информация по новости.
func (api *API) newsDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	s := mux.Vars(r)["id"]
	newsID, err := strconv.Atoi(s)
	if err != nil {
		//http.Error(w, "Ошибка страница не существует", http.StatusBadRequest)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	ctx := r.Context()
	chData := make(chan interface{}, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	// получение комментариев к новости
	go func() {
		defer wg.Done()
		data, err := api.commentService.Comments(ctx, newsID)
		if err != nil {
			chData <- err
			return
		}
		var c []comments.Comment
		err = json.Unmarshal(data, &c)
		if err != nil {
			chData <- err
			return
		}
		chData <- c
	}()
	// запрос к сервису новостей для получения новости по newsID
	go func() {
		defer wg.Done()
		data, err := api.newsService.PostsDetail(ctx, newsID)
		if err != nil {
			chData <- err
			return
		}
		var p news.Post
		err = json.Unmarshal(data, &p)
		if err != nil {
			chData <- err
			return
		}
		chData <- p
	}()
	wg.Wait()
	close(chData)
	var p news.Post
	var treeComments *[]treeComm
	for i := range chData {
		switch data := i.(type) {
		case news.Post:
			p = data
		case []comments.Comment:
			c := data
			treeComments = TreeComments(c)
		case error:
			http.Error(w, data.Error(), http.StatusInternalServerError)
			return
		}
	}
	type detailNews struct {
		Post     news.Post
		Comments *[]treeComm
	}
	dNews := detailNews{Post: p, Comments: treeComments}
	// совместить ответ от комментариев (создать дерево комментариев) и новостей в одну структуру
	json.NewEncoder(w).Encode(dNews)
}

// addComment добавление комментария.
func (api *API) addComment(w http.ResponseWriter, r *http.Request) {
	s := mux.Vars(r)["id"]
	newsID, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var c comments.Comment
	err = json.NewDecoder(r.Body).Decode(&c)
	//fmt.Println(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if c.NewsID != newsID {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// отправляем запрос на цензурирование комментария
	err = api.censorService.Censor(r.Context(), c.Content)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// отправляем запрос на добавление комментария
	err = api.commentService.AddComment(r.Context(), c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type treeComm struct {
	Comm comments.Comment `json:"Parent_comment"`
	C    *[]treeComm      `json:"Children_comment"`
}

type indexComment struct {
	index    int // индекс родителя в подмассиве дерева
	parentid int // id родителя по БД, если 0 то корень
}

func TreeComments(c []comments.Comment) *[]treeComm {
	// итоговое дерево комментариев
	listComm := &[]treeComm{}
	// координаты каждого комментария для ускорения поиска
	coord := make(map[int]indexComment)
	// стек id до родительского комментария
	stack := make([]indexComment, 0)
	for _, v := range c {
		iCom := indexComment{}
		// указываем индекс родителя, если есть или 0 если это комментарий верхнего уровня
		iCom.parentid = v.ParentID
		// это дочерний комментарий
		if v.ParentID != 0 {
			index := indexComment{parentid: v.ParentID}
			ok := false
			for {
				index, ok = coord[index.parentid]
				if !ok {
					fmt.Println("Странная ошибка ", v.ParentID)
					return nil
				}
				stack = append(stack, index)
				if index.parentid == 0 {
					break
				}
			}
			var lC *[]treeComm
			index = indexComment{}
			lC = listComm
			for len(stack) > 0 {
				index = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				lC = (*lC)[index.index].C
			}
			*lC = append(*lC, treeComm{Comm: v, C: &[]treeComm{}})
			iCom.index = len(*lC) - 1
			coord[v.ID] = iCom
		} else {
			*listComm = append(*listComm, treeComm{Comm: v, C: &[]treeComm{}})
			iCom.index = len(*listComm) - 1
			coord[v.ID] = iCom
		}
	}
	return listComm
}
