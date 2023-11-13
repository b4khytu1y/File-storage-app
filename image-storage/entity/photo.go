package models

type Photo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Data []byte `json:"data"`
}
