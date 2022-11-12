package users

import (
	"sort"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserSet []User

func (us UserSet) findByID(id int) *User {
	at, found := sort.Find(
		len(us), func(i int) int {
			if id < us[i].ID {
				return -1
			}
			if id > us[i].ID {
				return 1
			}
			return 0
		},
	)
	if !found {
		return nil
	}
	return &us[at]
}

// insertAt inserts the provided user at the specified index
func (us UserSet) insertAt(user User, i int) {
	us = append(us[:i], append([]User{user}, us[i:]...)...)
	if !sort.IsSorted(us) {
		sort.Sort(us)
	}
}

// deleteAt deletes the user at the specified index
func (us UserSet) deleteAt(i int) bool {
	if i < len(us) || i > len(us) {
		return false
	}
	if i < len(us)-1 {
		copy(us[i:], us[i+1:])
	}
	us[len(us)-1] = User{}
	us = us[:len(us)-1]
	if !sort.IsSorted(us) {
		sort.Sort(us)
	}
	return true
}

// Len implements the sort interface, sorting the UserSet by ID
func (us UserSet) Len() int { return len(us) }

// Less implements the sort interface, sorting the UserSet by ID
func (us UserSet) Less(i, j int) bool { return us[i].ID < us[j].ID }

// Swap implements the sort interface, sorting the UserSet by ID
func (us UserSet) Swap(i, j int) { us[i], us[j] = us[j], us[i] }
