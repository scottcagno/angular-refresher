package users

import (
	"errors"
	"sort"
)

type UserRepository struct {
	Users  []User `json:"data"`
	nextID int    `json:"omitempty"`
}

func NewUserRepoistory() *UserRepository {
	return &UserRepository{
		Users: make([]User, 0),
	}
}

func (u *UserRepository) FindOne(id int) (*User, error) {
	users := UserSet(u.Users)
	at, found := sort.Find(
		len(users), func(i int) int {
			if id < users[i].ID {
				return -1
			}
			if id > users[i].ID {
				return 1
			}
			return 0
		},
	)
	if !found {
		return nil, errors.New("error: could not find user")
	}
	return &users[at], nil
}

func (u *UserRepository) FindAll() ([]User, error) {
	if len(u.Users) == 0 {
		return []User{}, errors.New("error: the user repository is empty")
	}
	return u.Users, nil
}

func (u *UserRepository) Save(user *User) error {
	if user == nil {
		return errors.New("cannot save a nil user")
	}
	if user.ID == 0 {
		return u.insertOnly(user)
	}
	return u.updateOnly(user)
}

func (u *UserRepository) updateOnly(user *User) error {
	if user == nil {
		return errors.New("error: cannot update a nil user")
	}
	if user.ID == 0 {
		return errors.New("error: cannot update a new user without an ID")
	}
	users := UserSet(u.Users)
	at, found := sort.Find(
		len(users), func(i int) int {
			if user.ID < users[i].ID {
				return -1
			}
			if user.ID > users[i].ID {
				return 1
			}
			return 0
		},
	)
	if !found {
		return errors.New("error: cannot update, could not find existing user")
	}
	users[at] = *user
	return nil
}

func (u *UserRepository) insertOnly(user *User) error {
	if user == nil {
		return errors.New("error: cannot insert a nil user")
	}
	if user.ID != 0 {
		return errors.New("error: cannot insert a user that already has an ID")
	}
	user.ID = u.nextID
	u.nextID++
	u.Users = append(u.Users, *user)
	sort.Sort(UserSet(u.Users))
	return nil
}

func (u *UserRepository) DeleteOne(id int) error {
	users := UserSet(u.Users)
	at, found := sort.Find(
		len(users), func(i int) int {
			if id < users[i].ID {
				return -1
			}
			if id > users[i].ID {
				return 1
			}
			return 0
		},
	)
	if !found {
		return errors.New("error: could not find user")
	}
	success := users.deleteAt(at)
	if !success {
		return errors.New("error: unable to delete user")
	}
	return nil
}
