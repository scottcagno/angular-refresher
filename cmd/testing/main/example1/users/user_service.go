package users

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/scottcagno/angular-refresher/cmd/testing/main/example1/custom"
)

type UserService struct {
	userRepo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	// initialize and return user service
	return &UserService{
		userRepo: repo,
	}
}

func (u *UserService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.URL.Query().Has("id") {
			u.getUserByID()
			return
		}
		u.getAllUsers()
		return
	case http.MethodPut:
		u.updateUser()
		return
	case http.MethodPost:
		u.addNewUser()
		return
	case http.MethodDelete:
		u.deleteUserByID()
		return
	default:
		custom.NotFound(w, r)
		return
	}
}

func (u *UserService) getUserByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get id from request
		if !r.URL.Query().Has("id") {
			custom.WriteErrorJSON(w, r, http.StatusBadRequest, errors.New("`id` is required"))
			return
		}
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			custom.WriteErrorJSON(w, r, http.StatusExpectationFailed, err)
			return
		}
		// find user by the id provided
		user, err := u.userRepo.FindOne(id)
		if err != nil {
			custom.WriteErrorJSON(w, r, http.StatusExpectationFailed, err)
			return
		}
		// send user back to client
		custom.WriteJSON(w, r, http.StatusOK, user)
	}
}

func (u *UserService) getAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// return all the users
		users, err := u.userRepo.FindAll()
		if err != nil {
			custom.WriteErrorJSON(w, r, http.StatusExpectationFailed, err)
			return
		}
		// send users back to client
		custom.WriteJSON(w, r, http.StatusOK, users)
	}
}

func (u *UserService) addNewUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read json from request
		var newUser User
		custom.ReadJSON(w, r, &newUser)
		// add to our data set
		err := u.userRepo.Save(&newUser)
		if err != nil {
			custom.WriteErrorJSON(w, r, http.StatusExpectationFailed, err)
			return
		}
		// respond with 200 OK
		custom.WriteJSON(
			w, r, http.StatusOK, custom.JSON{
				"code":   http.StatusOK,
				"status": http.StatusText(http.StatusOK),
			},
		)
	}
}

func (u *UserService) updateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read json from request
		var existingUser User
		custom.ReadJSON(w, r, &existingUser)
		// add to our data set
		err := u.userRepo.Save(&existingUser)
		if err != nil {
			custom.WriteErrorJSON(w, r, http.StatusExpectationFailed, err)
			return
		}
		// respond with 200 OK
		custom.WriteJSON(
			w, r, http.StatusOK, custom.JSON{
				"code":   http.StatusOK,
				"status": http.StatusText(http.StatusOK),
			},
		)
	}
}

func (u *UserService) deleteUserByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get id from request
		if !r.URL.Query().Has("id") {
			custom.WriteErrorJSON(w, r, http.StatusBadRequest, errors.New("`id` is required"))
			return
		}
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			custom.WriteErrorJSON(w, r, http.StatusExpectationFailed, err)
			return
		}
		// remove user by the id provided
		err = u.userRepo.DeleteOne(id)
		if err != nil {
			custom.WriteErrorJSON(w, r, http.StatusExpectationFailed, err)
			return
		}
		// respond with 200 OK
		custom.WriteJSON(
			w, r, http.StatusOK, custom.JSON{
				"code":   http.StatusOK,
				"status": http.StatusText(http.StatusOK),
			},
		)
	}
}
