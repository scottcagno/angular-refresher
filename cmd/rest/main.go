package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/scottcagno/angular-refresher/pkg/http/rest"
)

func main() {

	config := &rest.Config{
		StaticHandler: rest.HandleStatic("static", "."),
		ErrHandler:    nil,
		MetricsOn:     true,
		LoggingLevel:  rest.LevelInfo,
	}

	app := rest.NewRouter(config)

	us := new(UserService)

	app.Get("/api/users", us.allUsers())
	app.Get("/api/users/", us.getOneUser("id"))

	log.Panicln(http.ListenAndServe(":9191", app))

}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserService struct {
	Users []User `json:"users"`
}

func (u *UserService) allUsers() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusOK
		payload := rest.NewRaw(u.Users)
		rest.WriteResponse(w, code, payload)
		return
	}
	return http.HandlerFunc(fn)
}

func (u *UserService) getOneUser(param string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !r.URL.Query().Has(param) {
			code := http.StatusExpectationFailed
			errPayload := rest.NewError(code, "could not marshal json data")
			rest.WriteResponse(w, code, errPayload)
			return
		}
		var user User
		id := r.URL.Query().Get(param)
		for i := range u.Users {
			if strconv.Itoa(u.Users[i].ID) == id {
				user = u.Users[i]
				break
			}
		}
		code := http.StatusOK
		payload := rest.NewRaw(user)
		rest.WriteResponse(w, code, payload)
		return
	}
	return http.HandlerFunc(fn)
}
