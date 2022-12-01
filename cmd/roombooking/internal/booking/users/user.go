package users

import (
	"github.com/scottcagno/angular-refresher/pkg/web/api"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password,omitempty"`
}

func (u *User) RequestMappingFunc(mapping api.RequestMapping) int {
	return 1
}

type mappingFunc = func(meth, path string, params ...string) int

// func DoSomething(meth, path string, params ...string) int {
// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		keys := make([]string, 0, len(r.URL.Query()))
// 		for k := range r.URL.Query() {
// 			keys = append(keys, k)
// 		}
// 	}
// 	return 1
// }

//
// func (u *User) RequestMappingFunc(method, path string, requiredParams ...string) int {
// 	// TODO implement me
// 	panic("implement me")
// }
//
// func (u *User) GetAllUsers(mapping api.RequestMappingFunc, next http.Handler) http.Handler {
// 	// TODO implement me
// 	panic("implement me")
// }
//
// func (u *User) GetUser(mapping api.RequestMappingFunc, next http.Handler) http.Handler {
//
// 	return nil
// }

func (u *User) SetPassword(password string) {
	u.Password = password
}
