package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"NewsRSS/Comments/pkg/postgres"

	"github.com/gorilla/mux"
)

type API struct {
	r  *mux.Router
	db *postgres.Storage
}

// Конструктор API.
func New(db *postgres.Storage) *API {
	a := API{}
	a.db = db
	a.r = mux.NewRouter()
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
	// получить комментарии к новости /comments/news/n
	api.r.HandleFunc("/comments/news/{n}", api.comments).Methods(http.MethodGet, http.MethodOptions)
	// добавить комментарий к новости или к комментарию
	api.r.HandleFunc("/comments", api.addComment).Methods(http.MethodPost, http.MethodOptions)
}

// comments получение комментариев к новости
func (api *API) comments(w http.ResponseWriter, r *http.Request) {
	s := mux.Vars(r)["n"]
	n, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Получение данных из БД.
	c, err := api.db.Comments(n)
	// НУЖно ПРЕОБразоВАть в дерево комментариев
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Отправка данных клиенту в формате JSON.
	json.NewEncoder(w).Encode(c)
}

// addComment добавление комментария
func (api *API) addComment(w http.ResponseWriter, r *http.Request) {
	var c postgres.Comment
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// проверка на недопустимые фразы
	//c.Content
	err = api.db.AddComment(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
