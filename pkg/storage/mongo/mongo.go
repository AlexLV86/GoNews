package mongo

import (
	"context"
	"fmt"

	"github.com/fatih/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"GoNews/pkg/storage"
)

type Storage struct {
	db             *mongo.Client
	databaseName   string //"data"  // имя учебной БД
	collectionName string //"posts" // имя коллекции в учебной БД
}

func New(constr string) (*Storage, error) {
	mongoOpts := options.Client().ApplyURI("mongodb://192.168.1.62:27017/")
	client, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db:             client,
		databaseName:   "data",
		collectionName: "posts",
	}
	return &s, nil
}

// Posts получение всех публикаций
func (s *Storage) Posts() ([]storage.Post, error) {
	collection := s.db.Database(s.databaseName).Collection(s.collectionName)
	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	var posts []storage.Post
	for cur.Next(context.Background()) {
		var p storage.Post
		err := cur.Decode(&p)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, cur.Err()
}

// AddPost создание новой публикации
func (s *Storage) AddPost(p storage.Post) error {
	collection := s.db.Database(s.databaseName).Collection(s.collectionName)
	// создаем мапу по структуре, ключи имена полей структуры
	m := structs.Map(p)
	var addFields bson.D
	for k, v := range m {
		temp := bson.D{{k, v}}
		addFields = append(addFields, temp...)

	}
	_, err := collection.InsertOne(context.Background(), addFields)
	if err != nil {
		return err
	}
	return nil
}

// UpdatePost обновление публикации
func (s *Storage) UpdatePost(p storage.Post) error {
	m := structs.Map(p)
	var updateFields bson.D
	for k, v := range m {
		if v != "" && v != 0 && v != int64(0) {
			temp := bson.D{{k, v}}
			updateFields = append(updateFields, temp...)
		}
	}
	fmt.Println(updateFields)
	if len(updateFields) == 0 {
		return fmt.Errorf("empty data")
	}
	filter := bson.D{{"ID", p.ID}}
	update := bson.D{{"$set", updateFields}}
	collection := s.db.Database(s.databaseName).Collection(s.collectionName)
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

// DeletePost удаление публикации по ID
func (s *Storage) DeletePost(p storage.Post) error {
	collection := s.db.Database(s.databaseName).Collection(s.collectionName)
	filter := bson.D{{"ID", p.ID}}
	//filter := bson.M{"id": p.ID}
	res, err := collection.DeleteOne(context.Background(), filter)
	fmt.Println(filter, res.DeletedCount)
	if err != nil {
		return err
	}
	return nil
}
