package booking

import (
	"time"

	"github.com/scottcagno/angular-refresher/cmd/roombooking/internal/booking/rooms"
	"github.com/scottcagno/angular-refresher/cmd/roombooking/internal/booking/users"
)

type Booking struct {
	ID           string     `json:"id"`
	Title        string     `json:"title"`
	User         users.User `json:"user"` // user id
	Room         rooms.Room `json:"room"` // room id
	Date         string     `json:"date"`
	StartTime    time.Time  `json:"start_time"`
	EndTime      time.Time  `json:"end_time"`
	Participants int        `json:"participants"`
}
