package filter

import (
	"fmt"
	"testing"
)

func TestKeys(t *testing.T) {
	fmt.Printf("mapIntInt keys: %#v\n", Keys(mapIntInt))
	fmt.Printf("mapIntStr keys: %#v\n", Keys(mapIntStr))
	fmt.Printf("mapStrInt keys: %#v\n", Keys(mapStrInt))
	fmt.Printf("mapIntAny keys: %#v\n", Keys(mapIntAny))
}

func TestVals(t *testing.T) {
	fmt.Printf("mapIntInt vals: %#v\n", Vals(mapIntInt))
	fmt.Printf("mapIntStr vals: %#v\n", Vals(mapIntStr))
	fmt.Printf("mapStrInt vals: %#v\n", Vals(mapStrInt))
	fmt.Printf("mapIntAny vals: %#v\n", Vals(mapIntAny))
}

func TestFilter(t *testing.T) {
	arrInt := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Printf("arrInt before: %#v\n", arrInt)
	evens := Filter(arrInt, func(i int, e int) bool { return e%2 == 0 })
	fmt.Printf("filtered evens: %#v\n", evens)
	odds := Filter(arrInt, func(i int, e int) bool { return e%2 == 1 })
	fmt.Printf("filtered odds: %#v\n", odds)
	fmt.Printf("origional arrInt: %#v\n", arrInt)
}

var mapIntInt = map[int]int{
	1: 11,
	9: 99,
	5: 55,
	3: 33,
	2: 22,
	7: 77,
	0: 00,
	4: 44,
	8: 88,
	6: 66,
}

var mapIntStr = map[int]string{
	1: "11",
	9: "99",
	5: "55",
	3: "33",
	2: "22",
	7: "77",
	0: "00",
	4: "44",
	8: "88",
	6: "66",
}

var mapStrInt = map[string]int{
	"1": 11,
	"9": 99,
	"5": 55,
	"3": 33,
	"2": 22,
	"7": 77,
	"0": 00,
	"4": 44,
	"8": 88,
	"6": 66,
}

var mapIntAny = map[int]any{
	1: "the value one",
	9: []any{99, "problemns"},
	5: "when I retire",
	3: 33.3333,
	2: struct{ number int }{22},
	7: "7 is the perfect number",
	0: nil,
	4: []int{4},
	8: []byte{0x88},
	6: "the number is six",
}
