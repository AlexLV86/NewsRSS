package api

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	ctxRequestID = "request_id"
)

// HeadersMiddleware устанавливает заголовки ответа сервера.
// установка request_id и передача его с контекстом
func (api *API) headersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		request_id := r.URL.Query().Get(ctxRequestID)
		if request_id == "" {
			request_id = RandString(rand.Intn(14) + 6)
		}
		ctx := r.Context()
		// устанавливаем request_id в контекст
		ctx = context.WithValue(ctx, ctxRequestID, request_id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// logMiddleware - логируем запрос в файл
func (api *API) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &StatusRecorder{
			ResponseWriter: w,
			Status:         200,
		}
		ctx := r.Context()
		next.ServeHTTP(rec, r.WithContext(ctx))
		str := time.Now().Format("2006-01-02 15:04:05") + "; "
		str += r.RemoteAddr + "; "
		str += strconv.FormatInt(int64(rec.Status), 10) + "; "
		str += ctx.Value(ctxRequestID).(string)
		f, err := os.OpenFile("log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("error opening file: %v\n", err)
			return
		}
		defer f.Close()
		log.SetOutput(f)
		log.Println(str)
	})
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

const charset = "0123456789"

func RandString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
