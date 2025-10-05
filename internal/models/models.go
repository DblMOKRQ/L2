package models

import "time"

type Event struct {
	UserID int64
	Date   time.Time
	Event  string
}

type EventRequest struct {
	UserID int64  `json:"user_id"`
	Date   string `json:"date"`
	Event  string `json:"event,omitempty"`
}
