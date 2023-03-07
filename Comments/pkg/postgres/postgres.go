package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Comment - комментарий. добавить для json поля
type Comment struct {
	ID       int
	ParentID int
	NewsID   int
	Content  string
	PubDate  int64
}

// Comment - комментарий.
// type treeComm struct {
// 	Comment
// 	ch *treeComm
// }

// var listComm *treeComm

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

// Comments получение всех коментариев для новости
func (s *Storage) Comments(newsid int) ([]Comment, error) {
	query := `SELECT comments.id, comments.parent_id, comments.news_id, 
	comments.content, comments.pubdate
	FROM comments WHERE comments.news_id=$1;`
	rows, err := s.db.Query(context.Background(), query, newsid)
	if err != nil {
		return nil, err
	}
	var comments []Comment
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var c Comment
		err = rows.Scan(
			&c.ID,
			&c.ParentID,
			&c.NewsID,
			&c.Content,
			&c.PubDate,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		comments = append(comments, c)
	}
	// ВАЖНО не забыть проверить rows.Err()
	return comments, rows.Err()
}

// AddComment добавление нового комментария
func (s *Storage) AddComment(c Comment) error {
	cmd, err := s.db.Exec(context.Background(), `
	INSERT INTO comments (parent_id, content, news_id)
	VALUES ($1, $2, $3);
	`,
		c.ParentID, c.Content, c.NewsID)

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("entry not added")
	}
	return err
}
