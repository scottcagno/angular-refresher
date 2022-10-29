package booking

import (
	"time"

	"github.com/scottcagno/angular-refresher/pkg/web/api"
)

type BookingRepository struct {
	items []Booking
}

func (r *BookingRepository) FindAll() (any, error) {
	if len(r.items) == 0 {
		return nil, api.ErrNone
	}
	return r.items, nil
}

func (r *BookingRepository) FindOne(key string) (any, error) {
	if len(r.items) == 0 {
		return nil, api.ErrNone
	}
	for _, item := range r.items {
		if item.ID == key {
			return item, nil
		}
	}
	return nil, api.ErrNoMatchFound
}

func (r *BookingRepository) Insert(v any) error {
	newItem, ok := v.(Booking)
	if !ok {
		return api.ErrBadType
	}
	r.items = append(r.items, newItem)
	return nil
}

func (r *BookingRepository) Update(key string, v any) error {
	if len(r.items) == 0 {
		return api.ErrNone
	}
	at := -1
	for i, item := range r.items {
		if item.ID == key {
			at = i
			break
		}
	}
	if at == -1 {
		return api.ErrNoMatchFound
	}
	updatedItem, ok := v.(Booking)
	if !ok {
		return api.ErrBadType
	}
	r.items[at] = updatedItem
	return nil
}

func (r *BookingRepository) Delete(key string) error {
	if len(r.items) == 0 {
		return api.ErrNone
	}
	at := -1
	for i, item := range r.items {
		if item.ID == key {
			at = i
			break
		}
	}
	if at == -1 {
		return api.ErrNoMatchFound
	}
	if at < len(r.items)-1 {
		copy(r.items[at:], r.items[at+1:])
	}
	r.items[len(r.items)-1] = Booking{} // or the zero value of T
	r.items = r.items[:len(r.items)-1]
	return nil
}

func (r *BookingRepository) Size() int {
	return len(r.items)
}

func (r *BookingRepository) Init() {
	booking1 := Booking{
		ID:           "1",
		Title:        "Conference call with CEO",
		User:         "2",
		Room:         "1",
		Date:         time.Now().String(),
		StartTime:    time.Now().Add(30 * time.Minute),
		EndTime:      time.Now().Add(3 * time.Hour),
		Participants: 4,
	}
	booking2 := Booking{
		ID:           "2",
		Title:        "Some important meeting",
		User:         "1",
		Room:         "2",
		Date:         time.Now().Add(29 * time.Hour).String(),
		StartTime:    time.Now().Add(30 * time.Hour),
		EndTime:      time.Now().Add(31 * time.Hour),
		Participants: 7,
	}
	r.items = append(r.items, booking1)
	r.items = append(r.items, booking2)
}
