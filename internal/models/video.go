package models

type Video struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Status      bool   `json:"status"`
	CreatedAt   string `json:"created_at"`
}
