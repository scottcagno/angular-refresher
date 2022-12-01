package users

import (
	"net/http"
	"strconv"

	"github.com/scottcagno/angular-refresher/pkg/web/api"
)

type Controller struct {
	*UserRepository
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	id, found := api.GetParam(r, "id")
	if !found {
		// handle, get all
		users, err := c.Find(func(u *User) bool { return u != nil })
		if err != nil {
			api.WriteJSON(w, http.StatusExpectationFailed, err)
			return
		}
		api.WriteJSON(w, http.StatusOK, users)
		return
	}
	// handle get one
	uid, err := strconv.Atoi(id)
	if err != nil {
		api.WriteJSON(w, http.StatusExpectationFailed, err)
		return
	}
	user, err := c.Find(func(u *User) bool { return u.ID == uid })
	if err != nil || len(user) != 1 {
		api.WriteJSON(w, http.StatusExpectationFailed, err)
		return
	}
	api.WriteJSON(w, http.StatusOK, user)
	return
}

func (c *Controller) Add(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := api.ReadJSON(r, &newUser)
	if err != nil {
		api.WriteJSON(w, http.StatusExpectationFailed, err)
		return
	}
	if newUser.ID == 0 {
		newUser.ID = c.nextID
		c.nextID++
	}
	err = c.Insert(newUser.ID, &newUser)
	if err != nil {
		api.WriteJSON(w, http.StatusExpectationFailed, err)
		return
	}
	api.WriteJSON(w, http.StatusCreated, newUser)
}

func (c *Controller) Set(w http.ResponseWriter, r *http.Request) {
	// get id we need to update
	id, found := api.GetParam(r, "id")
	if !found {
		api.WriteJSON(w, http.StatusNotFound, "error: id required, but not found")
		return
	}
	// get the updated room
	var updateUser User
	err := api.ReadJSON(r, &updateUser)
	if err != nil {
		api.WriteJSON(w, http.StatusExpectationFailed, err)
		return
	}
	uid, err := strconv.Atoi(id)
	if err != nil {
		api.WriteJSON(w, http.StatusExpectationFailed, err)
		return
	}
	err = c.Update(uid, &updateUser)
	if err != nil {
		api.WriteJSON(w, http.StatusExpectationFailed, err)
		return
	}
	api.WriteJSON(w, http.StatusOK, updateUser)
}

func (c *Controller) Del(w http.ResponseWriter, r *http.Request) {
	// get id we need to update
	id, found := api.GetParam(r, "id")
	if !found {
		api.WriteJSON(w, http.StatusNotFound, "error: id required, but not found")
		return
	}
	// locate room using id
	uid, err := strconv.Atoi(id)
	if err != nil {
		api.WriteJSON(w, http.StatusExpectationFailed, err)
		return
	}
	err = c.Delete(uid)
	if err != nil {
		api.WriteJSON(w, http.StatusExpectationFailed, err)
		return
	}
	api.WriteJSON(w, http.StatusOK, nil)
}

func (c *Controller) Custom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, found := api.GetParam(r, "id")
		if !found {
			api.WriteJSON(w, http.StatusExpectationFailed, "required an ID")
			return
		}
		// handle get one
		uid, err := strconv.Atoi(id)
		if err != nil {
			api.WriteJSON(w, http.StatusExpectationFailed, err)
			return
		}
		user, err := c.Find(func(u *User) bool { return u.ID == uid })
		if err != nil || len(user) != 1 {
			api.WriteJSON(w, http.StatusExpectationFailed, err)
			return
		}
		user[0].Password = "reset"
		err = c.Update(user[0].ID, user[0])
		if err != nil {
			api.WriteJSON(w, http.StatusExpectationFailed, err)
			return
		}
		api.WriteJSON(w, http.StatusOK, user)
		return
	}
}
