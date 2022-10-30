package rooms

import (
	"time"
)

type Room struct {
	ID    string    `json:"id"`
	Title string    `json:"title"`
	Time  time.Time `json:"time"`
}
