package models

import "encoding/json"

type Link struct {
	Frontend string `json:"frontend"`
	Backend  string `json:"backend"`
}

type Project struct {
	Id              string          `json:"id"`
	Title           string          `json:"title"`
	TypeProyect     string          `json:"typeProyect"`
	Description     string          `json:"description"`
	Tecnologies     json.RawMessage `json:"tecnologies"`
	Characteristics json.RawMessage `json:"characteristics"`
	Learning        json.RawMessage `json:"learning"`
	Image           string          `json:"image"`
	ImagenesProyect json.RawMessage `json:"imagenesProyect"`
	Link            string          `json:"link"`
	CreatedAt       string          `json:"createdAt"`
	Links           json.RawMessage `json:"links"`
	Status          string          `json:"status"`
	CounterLikes    int             `json:"counter_likes"`
}
