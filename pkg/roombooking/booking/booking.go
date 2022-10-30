package booking

import (
	"time"
)

type Booking struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	User         string    `json:"user"` // user id
	Room         string    `json:"room"` // room id
	Date         string    `json:"date"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Participants int       `json:"participants"`
}
