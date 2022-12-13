package middleware

// type AuthToken = string
// type AuthUser struct {
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }
//
// type BasicAuthConfig struct {
// 	LoginPath    string
// 	LogoutPath   string
// 	ValidatePath string
// 	AuthTokens   []AuthToken
// 	AuthUsers    []AuthUser
// }
//
// var defaultBasicAuthConfig = &BasicAuthConfig{
// 	LoginPath:    "/login",
// 	LogoutPath:   "/logout",
// 	ValidatePath: "/validate",
// 	AuthTokens:   make([]AuthToken, 0),
// 	AuthUsers: []AuthUser{
// 		{
// 			Username: "admin",
// 			Password: "secret",
// 		},
// 	},
// }
//
// type BasicAuthController struct {
// 	conf    *BasicAuthConfig
// 	protect string
// }
//
// func NewBasicAuthController(c *BasicAuthConfig) *BasicAuthController {
// 	if c == nil {
// 		c = defaultBasicAuthConfig
// 	}
// 	return &BasicAuthController{
// 		conf: c,
// 	}
// }
//
// func (bac *BasicAuthController) validate(r *http.Request) bool {
// 	u, p, ok := r.BasicAuth()
// 	return bac.hasAuthUser(u, p) && ok
// }
//
// func (bac *BasicAuthController) hasAuthUser(username, password string) bool {
// 	for _, au := range bac.conf.AuthUsers {
// 		if username == au.Username && password == au.Password {
// 			return true
// 		}
// 	}
// 	return false
// }
//
// func (bac *BasicAuthController) Secure(next http.Handler) http.Handler {
// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		switch r.URL.Path {
// 		case bac.conf.LoginPath:
// 			if bac.validate(r) {
// 				next.ServeHTTP(w, r)
// 				return
// 			}
// 			u, p := r.FormValue("username"), r.FormValue("password")
// 			if bac.hasAuthUser(u, p) {
// 				next.ServeHTTP(w, r)
// 				return
// 			}
// 			http.Redirect(w, r, bac.conf.LoginPath, http.StatusOK)
// 			return
// 		case bac.conf.LogoutPath:
// 			r.SetBasicAuth("", "")
// 			http.Redirect(w, r, bac.conf.LoginPath, http.StatusOK)
// 			return
// 		case bac.conf.ValidatePath:
// 			if bac.validate(r) {
// 				next.ServeHTTP(w, r)
// 				return
// 			}
// 		}
// 		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
// 		return
// 	}
// 	return http.HandlerFunc(fn)
// }
