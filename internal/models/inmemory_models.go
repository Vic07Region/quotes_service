package models

import "time"

type Quote struct {
	ID        int       `json:"id"`
	Quote     string    `json:"quote"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}
