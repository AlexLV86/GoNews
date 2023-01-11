-- перед созданием удалить все таблицы и создать их заново
DROP TABLE IF EXISTS posts, authors;
-- удаляю и последовательности, чтобы облегчить тестирование
DROP SEQUENCE IF EXISTS posts_id_seq, authors_id_seq;

CREATE TABLE IF NOT EXISTS authors (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS posts (
	id SERIAL PRIMARY KEY,
	author_id INT REFERENCES authors(id) ON DELETE CASCADE DEFAULT 0, -- автор из таблицы authors и каскадное удаление статей связанных с автором
	title TEXT DEFAULT 'Без названия',
	content TEXT NOT NULL,
	created_at BIGINT NOT NULL 
	DEFAULT extract(epoch from now()) -- дата создания статьи по умолчанию заполняется текущим временем
);

-- добавим пользователя по умолчанию
--INSERT INTO users (id, name) VALUES (0, 'default');
