package postgres

import (
	"GoNews/pkg/storage"
	"log"
	"os"
	"testing"
)

var s *Storage

func TestMain(m *testing.M) {
	pwd := os.Getenv("dbpass")
	if pwd == "" {
		m.Run()
	}
	var err error
	s, err = New("postgres://postgres:" + pwd + "@192.168.1.62:/gonews")
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}

func TestStorage_Posts(t *testing.T) {
	post := storage.Post{Title: "Тестовая статья", Content: "Прекрасная тестовая статья для постгрес!",
		AuthorID: 2, CreatedAt: 167893}
	err := s.AddPost(post)
	if err != nil {
		t.Fatal(err)
	}

	data, err := s.Posts()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
	// несуществующий автор
	post.AuthorID = 5
	err = s.AddPost(post)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStorage_UpdatePost(t *testing.T) {
	post := storage.Post{Title: "Update123",
		ID: 2}
	err := s.UpdatePost(post)
	if err != nil {
		t.Fatal(err)
	}
	post.Title = "Update 3"
	post.Content = "Update 3 content"
	post.ID = 3
	err = s.UpdatePost(post)
	if err != nil {
		t.Fatal(err)
	}
	post.ID = 1
	err = s.DeletePost(post)
	if err != nil {
		t.Fatal(err)
	}
	data, err := s.Posts()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}
