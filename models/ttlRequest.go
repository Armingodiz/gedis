package models

type TtlRequest struct {
	Key string `json:"key"`
	Ttl int    `json:"ttl"`
}
