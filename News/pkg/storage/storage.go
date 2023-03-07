package storage

// Публикация, получаемая из RSS.
type Post struct {
	ID      int    // номер записи
	Title   string // заголовок публикации
	Content string // содержание публикации
	PubTime int64  // время публикации
	Link    string // ссылка на источник
}

// Pagination, пагинация страниц для новостей.
type Pagination struct {
	PageNum    int // номер страницы
	PageItems  int // кол-во элементов на странице
	TotalPages int // кол-во страниц
}

type PostsPagination struct {
	Posts []Post
	Pages Pagination
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	PostsDetail(id int) (Post, error)                        // получение детальной информации о новости
	Posts(page int, filter string) (*PostsPagination, error) // получение заданного кол-ва публикаций и фильтр новостей по заголовку
	AddPosts([]Post) (int, error)                            // добавление публикаций в БД
}
