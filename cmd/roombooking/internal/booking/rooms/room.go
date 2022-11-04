package rooms

type Room struct {
	ID         int               `json:"id"`
	Name       string            `json:"name"`
	Location   string            `json:"location"`
	Capacities []*LayoutCapacity `json:"capacities"`
}

func NewRoom(id int, name string, location string) *Room {
	return &Room{
		ID:         id,
		Name:       name,
		Location:   location,
		Capacities: make([]*LayoutCapacity, 0),
	}
}

func (r *Room) AddLayoutCapacity(capacity *LayoutCapacity) {
	r.Capacities = append(r.Capacities, capacity)
}

type LayoutCapacity struct {
	Layout   LayoutType `json:"layout"`
	Capacity int        `json:"capacity"`
}

func NewLayoutCapacity(layout LayoutType, capacity int) *LayoutCapacity {
	return &LayoutCapacity{
		Layout:   layout,
		Capacity: capacity,
	}
}

type LayoutType string

const (
	Layout_THEATER LayoutType = "Theater"
	Layout_USHAPE  LayoutType = "U-Shape"
	Layout_BOARD   LayoutType = "Board Meeting"
)
