package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"
	"sort"
	"strings"
	"unicode"
)

const version = "v1"

type RestAPI struct {
	addr string
	base string
}

func NewRestAPI(base, addr string) *RestAPI {
	return &RestAPI{
		addr: addr,
		base: path.Join(base, version),
	}
}

func (api *RestAPI) Addr() string {
	return api.addr
}

func (api *RestAPI) BasePath() string {
	return api.base
}

func (api *RestAPI) Version() string {
	return version
}

func (api *RestAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.URL.Path == api.base {
		w.WriteHeader(200)
		fmt.Fprintf(w, "%s", api.base)
		return
	}
}

// type Service[T any] struct {
// 	Query(q string) []T
// 	FindAll() []T
// 	FindOne(id string) T
// 	Insert(v T) T
// 	Update(v T) T
// 	Delete(id string) T
// }

type Resource interface {
	GetID() string
	SetID(id string)
}

type Repository interface {
	Insert(t any) error
	Update(t any) error
	Return(q string) ([]any, error)
	Delete(q string) error
}

type Service[T Resource] struct {
	data []T
}

func NewService[T Resource]() *Service[T] {
	return &Service[T]{
		data: make([]T, 0),
	}
}

func (s *Service[T]) sortData() {
	sort.Slice(s.data, func(i int, j int) bool { return s.data[i].GetID() < s.data[j].GetID() })
}

func (s *Service[T]) findData(id string) (*T, int) {
	findFn := func(i int) int {
		return strings.Compare(id, s.data[i].GetID())
	}
	i, found := sort.Find(len(s.data), findFn)
	if !found {
		return nil, -1
	}
	return &s.data[i], i
}

func (s *Service[T]) deleteAt(i int) *T {
	if i < len(s.data)-1 {
		copy(s.data[i:], s.data[i+1:])
	}
	var old T
	old = (s.data)[len(s.data)-1]
	s.data[len(s.data)-1] = *new(T) // or the zero value of T
	s.data = s.data[:len(s.data)-1]
	return &old
}

func (s *Service[T]) Query(q string) []T {
	return nil
}

func (s *Service[T]) FindAll() []T {
	return s.data
}

func (s *Service[T]) FindOne(id string) *T {
	t, _ := s.findData(id)
	if t == nil {
		return nil
	}
	return t
}

func (s *Service[T]) Insert(v T) T {
	s.data = append(s.data, v)
	s.sortData()
	return v
}

func (s *Service[T]) Update(v *T) *T {
	t, i := s.findData((*v).GetID())
	if t == nil {
		return nil
	}
	s.data[i] = *v
	s.sortData()
	return v
}

func (s *Service[T]) Delete(id string) *T {
	t, i := s.findData(id)
	if t == nil {
		return nil
	}
	old := s.deleteAt(i)
	return old
}

func (s *Service[T]) String() string {
	return fmt.Sprintf("%T, %#v\n", s.data, s.data)
}

type Controller[T Resource] struct {
	Service *Service[T]
	mux     *http.ServeMux
	routes  [][]string
}

func NewController[T Resource](mux *http.ServeMux) *Controller[T] {
	if mux == nil {
		mux = http.NewServeMux()
	}
	return &Controller[T]{
		Service: NewService[T](),
		mux:     mux,
		routes:  make([][]string, 0),
	}
}

func (c *Controller[T]) addRoute(method string, path string) {
	c.routes = append(c.routes, []string{method, path})
}

func (c *Controller[T]) Any(path string, handler http.Handler) {
	c.mux.Handle(path, handler)
	c.addRoute("*", path)
}

func withMethod(method string, next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			code := http.StatusMethodNotAllowed
			http.Error(w, http.StatusText(code), code)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (c *Controller[T]) Get(path string, handler http.Handler) {
	c.mux.Handle(path, withMethod(http.MethodGet, handler))
	c.addRoute(http.MethodGet, path)
}

func (c *Controller[T]) Post(path string, handler http.Handler) {
	c.mux.Handle(path, withMethod(http.MethodPost, handler))
	c.addRoute(http.MethodPost, path)
}

func (c *Controller[T]) Put(path string, handler http.Handler) {
	c.mux.Handle(path, withMethod(http.MethodPut, handler))
	c.addRoute(http.MethodPut, path)
}

func (c *Controller[T]) Delete(path string, handler http.Handler) {
	c.mux.Handle(path, withMethod(http.MethodDelete, handler))
	c.addRoute(http.MethodDelete, path)
}

func (c *Controller[T]) Patch(path string, handler http.Handler) {
	c.mux.Handle(path, withMethod(http.MethodPatch, handler))
	c.addRoute(http.MethodPatch, path)
}

func (c *Controller[T]) Options(path string, handler http.Handler) {
	c.mux.Handle(path, withMethod(http.MethodOptions, handler))
	c.addRoute(http.MethodOptions, path)
}

func trimTrailingSlash(s string, c byte) string {
	for len(s) > 0 && s[len(s)-1] == c {
		s = s[:len(s)-1]
	}
	return s
}

func (c *Controller[T]) InitCRUD(base string) {
	// trim trailing slash
	base = trimTrailingSlash(base, '/')
	// get all T
	c.Get(base, c.DefaultGetAllHandler())
	// get specific T
	c.Get(path.Join(base, "/"), c.DefaultGetOneHandler())
	// add a new T
	c.Post(base, c.DefaultAddNewHandler())
	// update a specific T
	c.Put(base, c.DefaultUpdateHandler())
	// delete a T
	c.Delete(path.Join(base, "/"), c.DefaultDeleteOneHandler())
}

func (c *Controller[T]) StatsHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(c.routes)
		if err != nil {
			log.Printf("status handler error: %s\n", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, "%s", b)
	}
	return http.HandlerFunc(fn)
}

func (c *Controller[T]) DefaultGetAllHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		all := c.Service.FindAll()
		b, err := json.Marshal(all)
		if err != nil {
			log.Printf("default get all handler: %s\n", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, "%s", b)
	}
	return http.HandlerFunc(fn)
}

func GetURLParams(r *http.Request, expect int) (string, error) {
	str := strings.FieldsFunc(r.URL.Path, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})
	if len(str) != expect {
		return "", errors.New("expectation did not match result")
	}
	return str[expect-1], nil
}

func (c *Controller[T]) DefaultGetOneHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		id, err := GetURLParams(r, 3)
		if err != nil {
			log.Printf("default get one handler: %s\n", err)
		}
		one := c.Service.FindOne(id)
		b, err := json.Marshal(one)
		if err != nil {
			log.Printf("default get one handler: %s\n", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, "%s", b)
	}
	return http.HandlerFunc(fn)
}

func (c *Controller[T]) DefaultAddNewHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println("running default add new handler...")
	}
	return http.HandlerFunc(fn)
}

func (c *Controller[T]) DefaultUpdateHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println("running default update handler...")
	}
	return http.HandlerFunc(fn)
}

func (c *Controller[T]) DefaultDeleteOneHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		id, err := GetURLParams(r, 3)
		if err != nil {
			log.Printf("default delete one handler: %s\n", err)
		}
		one := c.Service.Delete(id)
		b, err := json.Marshal(one)
		if err != nil {
			log.Printf("default delete one handler: %s\n", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, "%s", b)
	}
	return http.HandlerFunc(fn)
}
