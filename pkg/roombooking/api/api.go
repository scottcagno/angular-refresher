package api

import (
	"fmt"
	"net/http"
	"path"
	"sort"
	"strings"
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
	s.data[len(s.data)-1] = nil // or the zero value of T
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

func (s *Service[T]) Insert(v *T) *T {
	s.data = append(s.data, *v)
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

type RESTController[T Resource] struct {
	service *Service[T]
}

func NewController[T Resource]() *Controller[T] {
	return &Controller[T]{
		service: NewService[T](),
	}
}
