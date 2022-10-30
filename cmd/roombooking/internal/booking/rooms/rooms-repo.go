package rooms

import (
	"time"

	"github.com/scottcagno/angular-refresher/pkg/web/api"
)

type RoomRepository struct {
	items []Room
}

func (r *RoomRepository) FindAll() (any, error) {
	if len(r.items) == 0 {
		return nil, api.ErrNone
	}
	return r.items, nil
}

func (r *RoomRepository) FindOne(key string) (any, error) {
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

func (r *RoomRepository) Insert(v any) error {
	newItem, ok := v.(Room)
	if !ok {
		return api.ErrBadType
	}
	r.items = append(r.items, newItem)
	return nil
}

func (r *RoomRepository) Update(key string, v any) error {
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
	updatedItem, ok := v.(Room)
	if !ok {
		return api.ErrBadType
	}
	r.items[at] = updatedItem
	return nil
}

func (r *RoomRepository) Delete(key string) error {
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
	r.items[len(r.items)-1] = Room{} // or the zero value of T
	r.items = r.items[:len(r.items)-1]
	return nil
}

func (r *RoomRepository) Size() int {
	return len(r.items)
}

func (r *RoomRepository) Init(data map[string]any) {
	// add initial data to room repo
	room1 := Room{
		ID:    "1",
		Title: "Room number one",
		Time:  time.Now(),
	}
	room2 := Room{
		ID:    "2",
		Title: "Room number two",
		Time:  time.Now().Add(5 * time.Hour),
	}
	room3 := Room{
		ID:    "3",
		Title: "Room number three",
		Time:  time.Now().Add(2 * time.Hour),
	}
	r.items = append(r.items, room1)
	r.items = append(r.items, room2)
	r.items = append(r.items, room3)
}

func (r *RoomRepository) GetRepositorySet() any {
	return r.items
}
