package postgres

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"News/pkg/storage"
)

// кол-во элементов на одной странице
const numItemsPage = 15

// Хранилище данных.
type Storage struct {
	db *pgxpool.Pool
}

// Конструктор, принимает строку подключения к БД.
func New(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

// PostsDetail получение детальной информации по id новости
func (s *Storage) PostsDetail(id int) (storage.Post, error) {
	query := `SELECT posts.id, posts.title, 
	posts.content, posts.pubdate, posts.link 
	FROM posts WHERE posts.id=$1;`
	row := s.db.QueryRow(context.Background(), query, id)
	var p storage.Post
	err := row.Scan(
		&p.ID,
		&p.Title,
		&p.Content,
		&p.PubTime,
		&p.Link,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return storage.Post{}, nil
		}
		return p, err
	}
	// ВАЖНО не забыть проверить rows.Err()
	return p, nil
}

// Posts получение n публикаций
func (s *Storage) Posts(page int, filter string) (*storage.PostsPagination, error) {
	// запрос на общее кол-во страниц
	query := `SELECT COUNT(*) FROM posts WHERE posts.title ILIKE '%' || $1 || '%'`
	row := s.db.QueryRow(context.Background(), query, filter)
	pagin := storage.Pagination{}
	var cntPosts int
	err := row.Scan(&cntPosts)
	if err != nil {
		return nil, err
	}
	// расчет кол-ва страниц
	pagin.TotalPages = cntPosts / numItemsPage
	if cntPosts%numItemsPage != 0 {
		pagin.TotalPages++
	}
	pagin.PageItems = numItemsPage
	pagin.PageNum = page
	query = `SELECT posts.id, posts.title,
	posts.pubdate, posts.link
	FROM posts WHERE posts.title
	ILIKE '%' || $3 || '%'
	ORDER BY posts.pubdate DESC
	LIMIT $1 OFFSET $2;`
	rows, err := s.db.Query(context.Background(), query, pagin.PageItems, (pagin.PageNum-1)*pagin.PageItems, filter)
	if err != nil {
		return nil, err
	}
	var posts []storage.Post
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		posts = append(posts, p)
	}
	// ВАЖНО не забыть проверить rows.Err()
	return &storage.PostsPagination{Posts: posts, Pages: pagin}, rows.Err()
}

// AddPosts добавляет новые публикации
func (s *Storage) AddPosts(posts []storage.Post) (int, error) {
	var err error
	var cmd pgconn.CommandTag
	rows := 0
	for _, p := range posts {
		cmd, err = s.db.Exec(context.Background(), `
		INSERT INTO posts (title, content, pubdate, link) 
		VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING;
		`,
			p.Title, p.Content, p.PubTime, p.Link)
		if err != nil {
			return 0, err
		}
		rows += int(cmd.RowsAffected())
	}
	return rows, err
}
