package api

import (
	"log"

	"github.com/scottcagno/angular-refresher/pkg/web"
)

type UserStore struct {
	store *MemoryStore[string, web.SystemUser]
}

func NewUserStore() *UserStore {
	return &UserStore{
		store: NewMemoryStore[string, web.SystemUser](),
	}
}

func (us *UserStore) AddUser(username, password, role string) {
	err := us.store.Add(username, web.SystemUser{
		Username: username,
		Password: password,
		Role:     role,
	})
	if err != nil {
		log.Printf("[UserStore] user could not be added: %s\n", err)
	}
}

func (us *UserStore) UpdateUser(username string, user web.SystemUser) {
	us.store.Set(username, user)
}

func (us *UserStore) GetUser(username string) *web.SystemUser {
	user, err := us.store.Get(username)
	if err != nil {
		log.Printf("[UserStore] user could not be found: %s\n", err)
		return nil
	}
	return &user
}

func (us *UserStore) DeleteUser(username string) {
	_, err := us.store.Del(username)
	if err != nil {
		log.Printf("[UserStore] user could not be found: %s\n", err)
	}
}
