package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var cens = [...]string{"qwerty", "йцукен", "zxvbnm"}

type API struct {
	r *mux.Router
}

// Конструктор API.
func New() *API {
	a := API{}
	a.r = mux.NewRouter()
	a.endpoints()
	return &a
}

func (api *API) Router() *mux.Router {
	return api.r
}

// Регистрация методов API в маршрутизаторе запросов.
func (api *API) endpoints() {
	api.r.Use(api.headersMiddleware)
	api.r.Use(api.logMiddleware)
	// получить комментарии к новости /comments/news/n
	api.r.HandleFunc("/censor", api.censor).Methods(http.MethodPost, http.MethodOptions)
}

// qwerty, йцукен, zxvbnm.
// censor проверка комментария на недопустимые значения
func (api *API) censor(w http.ResponseWriter, r *http.Request) {
	var comment string
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	comment = strings.ToLower(comment)
	for _, c := range cens {
		if strings.Contains(comment, c) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	}
	w.WriteHeader(http.StatusOK)
}
