-- перед созданием удалить все таблицы и создать их заново
DROP TABLE IF EXISTS posts;
-- удаляю и последовательности, чтобы облегчить тестирование
DROP SEQUENCE IF EXISTS posts_id_seq;

CREATE TABLE IF NOT EXISTS posts (
	id SERIAL PRIMARY KEY,
	title TEXT DEFAULT 'Без названия',
	content TEXT NOT NULL,
	pubdate BIGINT NOT NULL, -- дата создания статьи 
	link TEXT UNIQUE NOT NULL 
);
