-- перед созданием удалить все таблицы и создать их заново
DROP TABLE IF EXISTS comments;
-- удаляю и последовательности, чтобы облегчить тестирование
DROP SEQUENCE IF EXISTS comments_id_seq;

CREATE TABLE IF NOT EXISTS comments (
	id SERIAL PRIMARY KEY,
	parent_id INT DEFAULT 0, -- 0 комментарий к новости. иначе id родительского комментария
	news_id INT NOT NULL, -- id новости из сервиса 
	content TEXT NOT NULL,
	pubdate BIGINT NOT NULL
	DEFAULT extract(epoch from now())-- дата создания комментария 
);
