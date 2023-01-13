package mongo

import (
	"GoNews/pkg/storage"
	"context"
	"log"
	"os"
	"testing"
)

var s *Storage

func TestMain(m *testing.M) {
	var err error
	s, err = New("mongodb://192.168.1.62:27017/")
	if err != nil {
		log.Fatal(err)
	}
	// не забываем закрывать ресурсы
	defer s.db.Disconnect(context.Background())
	// проверка связи с БД
	err = s.db.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}

func TestStorage_Posts(t *testing.T) {
	var err error
	p := storage.Post{
		ID:         1,
		Title:      "Mongo Title 1",
		Content:    "Mongo content 1",
		CreatedAt:  168340123,
		AuthorName: "Alexander",
	}
	err = s.AddPost(p)
	if err != nil {
		log.Fatal(err)
	}
	p = storage.Post{
		ID:         2,
		Title:      "Mongo Title 2",
		Content:    "Mongo content 2",
		CreatedAt:  168340123,
		AuthorName: "Semen",
	}
	err = s.AddPost(p)
	if err != nil {
		log.Fatal(err)
	}
	p.Title = "Up"
	p.ID = 2
	err = s.UpdatePost(p)
	if err != nil {
		log.Fatal(err)
	}
	p.ID = 1
	err = s.DeletePost(p)
	if err != nil {
		log.Fatal(err)
	}
	posts, err := s.Posts()
	if err != nil {
		log.Fatal(err)
	}
	t.Log(posts)
}
