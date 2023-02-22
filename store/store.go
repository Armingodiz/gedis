package store

import (
	"encoding/json"
	"gedis/db"
	"gedis/models"
)

func NewStore(db *db.DB) Store {
	return &PostgresStore{db}
}

type Store interface {
	CreateUser(user models.User) error
	GetUser(userName string) (models.User, error)
	CreateKv(kv models.KeyValue) error
	GetKv(userName, key string) (models.KeyValue, error)
	GetKvs(userName string) ([]models.KeyValue, error)
	DeleteKv(key string) error
}

type PostgresStore struct {
	db *db.DB
}

func (s *PostgresStore) CreateUser(user models.User) error {
	_, err := s.db.Connection.Exec("INSERT INTO users (user_name, password, urls) VALUES ($1, $2, $3)", user.Username, user.Password, nil)
	return err
}

func (s *PostgresStore) GetUser(userName string) (models.User, error) {
	var user models.User
	var body interface{}
	err := s.db.Connection.QueryRow("SELECT user_name, password, urls FROM users WHERE user_name = $1", userName).Scan(&user.Username, &user.Password, &body)
	if err != nil {
		return user, err
	}
	if body != nil {
		var urls map[string]string
		json.Unmarshal(body.([]byte), &urls)
		user.Urls = urls
	}
	return user, err
}

func (s *PostgresStore) CreateKv(kv models.KeyValue) error {
	_, err := s.db.Connection.Exec("Insert into key_values (user_name, key, value) values ($1, $2, $3)", kv.UserName, kv.Key, kv.Value)
	return err
}

func (s *PostgresStore) GetKvs(userName string) ([]models.KeyValue, error) {
	var kvs []models.KeyValue
	rows, err := s.db.Connection.Query("SELECT user_name, key, value FROM key_values WHERE user_name = $1", userName)
	if err != nil {
		return kvs, err
	}
	for rows.Next() {
		var kv models.KeyValue
		err = rows.Scan(&kv.UserName, &kv.Key, &kv.Value)
		if err != nil {
			return kvs, err
		}
		kvs = append(kvs, kv)
	}
	return kvs, err
}

func (s *PostgresStore) GetKv(userName, key string) (models.KeyValue, error) {
	var kv models.KeyValue
	err := s.db.Connection.QueryRow("SELECT user_name, key, value FROM key_values WHERE user_name = $1 and key = $2", userName, key).Scan(&kv.UserName, &kv.Key, &kv.Value)
	return kv, err
}

func (s *PostgresStore) DeleteKv(key string) error {
	_, err := s.db.Connection.Exec("delete from key_values where key = $1", key)
	return err
}
