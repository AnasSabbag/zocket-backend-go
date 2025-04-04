package models

// Task struct (Define it once here)
type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type Message struct {
	Type string `json:"type"`
	Task Task   `json:"task"`
}