package memdb

import (
	"News/pkg/storage"
	"strings"
	"sync"
)

// База данных заказов.
type DB struct {
	m    sync.Mutex           //мьютекс для синхронизации доступа
	id   int                  // текущее значение ID для нового заказа
	post map[int]storage.Post // БД заказов
}

// кол-во элементов на одной странице
const numItemsPage = 15

// Конструктор БД.
func New() *DB {
	db := DB{
		id:   1, // первый номер заказа
		post: map[int]storage.Post{},
	}
	return &db
}

func (db *DB) PostsDetail(id int) (storage.Post, error) {
	p, _ := db.post[id]
	return p, nil
}

// Posts возвращает заданное кол-во статей
func (db *DB) Posts(page int, filter string) (*storage.PostsPagination, error) {
	pagin := storage.Pagination{}
	db.m.Lock()
	defer db.m.Unlock()
	cntPosts := len(db.post)
	// расчет кол-ва страниц
	pagin.TotalPages = cntPosts / numItemsPage
	if cntPosts%numItemsPage != 0 {
		pagin.TotalPages++
	}
	pagin.PageItems = numItemsPage
	pagin.PageNum = page
	data := make([]storage.Post, 0, numItemsPage)
	i := 0
	for _, p := range db.post {
		// если нет текущей подстроки в заголовке новости
		if !strings.Contains(strings.ToLower(p.Title), strings.ToLower(filter)) {
			continue
		}
		if i >= numItemsPage {
			break
		}
		i++
		// отправляю  n новостей без сортировки по дате
		data = append(data, p)
	}
	return &storage.PostsPagination{data, pagin}, nil
}

// AddPost добавляет статью без проверки на уникальность
func (db *DB) AddPosts(posts []storage.Post) (int, error) {
	db.m.Lock()
	defer db.m.Unlock()
	for _, p := range posts {
		p.ID = db.id
		db.post[p.ID] = p
		db.id++
	}
	return len(posts), nil
}

// размер БД
func (db *DB) Len() int {
	return len(db.post)
}
