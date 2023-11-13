package entity

type Photo struct {
	ID   int    `json:"id"`
	Data []byte `json:"data"`
}
