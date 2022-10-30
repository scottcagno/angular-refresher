package users

import (
	"github.com/scottcagno/angular-refresher/pkg/web/api"
)

type UserRepository struct {
	items []User
}

func (r *UserRepository) FindAll() (any, error) {
	if len(r.items) == 0 {
		return nil, api.ErrNone
	}
	return r.items, nil
}

func (r *UserRepository) FindOne(key string) (any, error) {
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

func (r *UserRepository) Insert(v any) error {
	newItem, ok := v.(User)
	if !ok {
		return api.ErrBadType
	}
	r.items = append(r.items, newItem)
	return nil
}

func (r *UserRepository) Update(key string, v any) error {
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
	updatedItem, ok := v.(User)
	if !ok {
		return api.ErrBadType
	}
	r.items[at] = updatedItem
	return nil
}

func (r *UserRepository) Delete(key string) error {
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
	r.items[len(r.items)-1] = User{} // or the zero value of T
	r.items = r.items[:len(r.items)-1]
	return nil
}

func (r *UserRepository) Size() int {
	return len(r.items)
}

func (r *UserRepository) Init(data map[string]any) {
	user1 := User{
		ID:   "1",
		Name: "Dick Chesterwood",
	}
	user2 := User{
		ID:   "2",
		Name: "Matt Greencroft",
	}
	r.items = append(r.items, user1)
	r.items = append(r.items, user2)
}

func (r *UserRepository) GetRepositorySet() any {
	return r.items
}
