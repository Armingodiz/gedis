package models

type KeyValue struct {
	UserName string `json:"user_name" db:"user_name"`
	Key      string `json:"key" db:"key"`
	Value    string `db:"value" json:"value"`
}
