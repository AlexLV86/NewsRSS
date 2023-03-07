package api

import (
	"News/pkg/storage"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// API приложения.
type API struct {
	r  *mux.Router       // маршрутизатор запросов
	db storage.Interface // база данных
}

// Конструктор API.
func New(db storage.Interface) *API {

	api := API{}
	api.db = db
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

// Router возвращает маршрутизатор запросов.
func (api *API) Router() *mux.Router {
	return api.r
}

// Регистрация методов API в маршрутизаторе запросов.
func (api *API) endpoints() {
	api.r.Use(api.headersMiddleware)
	api.r.Use(api.logMiddleware)
	// получить  детализированную информацию по новости
	api.r.HandleFunc("/news/{id}", api.newsDetail).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc("/news", api.newsPage).Methods(http.MethodGet, http.MethodOptions)
	// веб-приложение
	api.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("cmd/gonews/webapp"))))
}

// newsDetail получение детализированной информации о новости
func (api *API) newsDetail(w http.ResponseWriter, r *http.Request) {
	// Считывание параметра {id} из пути запроса.
	// Например, /news/10.
	s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Получение данных из БД.
	p, err := api.db.PostsDetail(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Отправка данных клиенту в формате JSON.
	json.NewEncoder(w).Encode(p)
}

// newsPage возвращает список новостей.
// page - номер страницы
// s - текстовое поля для поиска по новостям
func (api *API) newsPage(w http.ResponseWriter, r *http.Request) {
	pageParam := r.URL.Query().Get("page")
	if pageParam == "" {
		pageParam = "1"
	}
	// параметр page - это число, поэтому нужно сконвертировать
	// строку в число при помощи пакета strconv
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filter := r.URL.Query().Get("s")
	// Получение данных из БД.
	p, err := api.db.Posts(page, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Отправка данных клиенту в формате JSON.
	json.NewEncoder(w).Encode(p)
}
