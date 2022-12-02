package booking

import (
	"time"

	"github.com/scottcagno/angular-refresher/cmd/roombooking/internal/booking/rooms"
	"github.com/scottcagno/angular-refresher/cmd/roombooking/internal/booking/users"
)

const dateFormat = `2006-01-02`

type Booking struct {
	ID           int        `json:"id"`
	Title        string     `json:"title"`
	User         users.User `json:"user"` // user id
	Room         rooms.Room `json:"room"` // room id
	Date         string     `json:"date"`
	StartTime    time.Time  `json:"start_time"`
	EndTime      time.Time  `json:"end_time"`
	Participants int        `json:"participants"`
}

const Day = time.Duration(24 * time.Hour)

const (
	Sub = iota - 1
	Today
	Add
)

func GetDate(action int, days int) string {
	date := time.Now()
	if days != 0 && (action == Sub || action < Today) {
		date = date.Add(-(time.Duration(days) * 24 * time.Hour))
	}
	if days != 0 && (action == Add || action > Today) {
		date = date.Add(time.Duration(days) * 24 * time.Hour)
	}
	return date.Format(dateFormat)
}
