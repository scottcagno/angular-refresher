package rooms

import (
	"time"
)

type Room struct {
	ID    string
	Title string
	Time  time.Time
}
