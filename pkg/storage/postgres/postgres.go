package postgres

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"

	"GoNews/pkg/storage"
)

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

// Posts получение всех публикаций
func (s *Storage) Posts() ([]storage.Post, error) {
	query := `SELECT posts.id, posts.title, 
	posts.content, posts.author_id, authors.name, posts.created_at 
	FROM authors, posts WHERE authors.id=posts.author_id;`
	rows, err := s.db.Query(context.Background(), query)
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
			&p.Content,
			&p.AuthorID,
			&p.AuthorName,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		posts = append(posts, p)
	}
	// ВАЖНО не забыть проверить rows.Err()
	return posts, rows.Err()
}

// AddPost создание новой публикации
func (s *Storage) AddPost(p storage.Post) error {
	/*_, err := s.db.Exec(context.Background(), `
	INSERT INTO posts (title, content, author_id, created_at)
	VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING;
	`,
	p.Title, p.Content, p.AuthorID, p.CreatedAt)*/
	// cmd, err := s.db.Exec(context.Background(), `
	// 	INSERT INTO posts (title, content, created_at, author_id)
	// 	VALUES ($1, $2, $3,
	// 		(SELECT authors.id from authors where authors.id=$4));
	// 	`,
	// 	p.Title, p.Content, p.CreatedAt, p.AuthorID)
	// insert into posts (title, content, author_id)
	// select 'value for column123', 'value for column2',
	// authors.id from authors where authors.id=15;
	// проверяю через select есть ли такой автор в таблице авторов
	cmd, err := s.db.Exec(context.Background(), `
		INSERT INTO posts (title, content, created_at, author_id) 
		SELECT $1, $2, $3, authors.id FROM authors where authors.id=$4;
		`,
		p.Title, p.Content, p.CreatedAt, p.AuthorID)
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("entry not added")
	}
	return err
}

// UpdatePost обновление публикации
func (s *Storage) UpdatePost(p storage.Post) error {
	var values []interface{}
	values = append(values, p.ID)
	set := ""
	i := 1
	if p.Title != "" {
		i++
		set += "title=$" + strconv.Itoa(i) + ","
		values = append(values, p.Title)
	}
	if p.Content != "" {
		i++
		set += "content=$" + strconv.Itoa(i) + ","
		values = append(values, p.Content)
	}
	if p.CreatedAt != 0 {
		i++
		set += "created_at=$" + strconv.Itoa(i) + ","
		values = append(values, p.CreatedAt)
	}
	// не переданы значения для обновления
	if set == "" {
		return fmt.Errorf("empty data")
	}
	// убирает последнюю запятую
	set = set[:len(set)-1]
	query := "UPDATE posts SET " + set + "WHERE id=$1;"
	_, err := s.db.Exec(context.Background(), query, values...)
	return err
}

// DeletePost удаление публикации по ID
func (s *Storage) DeletePost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
		DELETE FROM posts WHERE id=$1;
		`,
		p.ID,
	)
	return err
}
