package booking

import (
	"time"
)

type Booking struct {
	ID           string
	Title        string
	User         string // user id
	Room         string // room id
	Date         string
	StartTime    time.Time
	EndTime      time.Time
	Participants int
}
