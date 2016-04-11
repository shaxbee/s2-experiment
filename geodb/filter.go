package geodb

import (
	"container/heap"
	"sort"

	"github.com/golang/geo/s2"
)

type Filter struct {
	distance Distance
	contains func(s2.LatLng) bool
	limit    int
	center   s2.LatLng
	data     elements
}

// Element represents pair of distance and index
type Element struct {
	Distance float64
	Value    interface{}
}

func FromCap(c s2.Cap, d Distance, limit int) *Filter {
	return &Filter{
		distance: d,
		contains: func(ll s2.LatLng) bool {
			return c.ContainsPoint(s2.PointFromLatLng(ll))
		},
		limit:  limit,
		center: s2.LatLngFromPoint(c.Center()),
		data:   make([]Element, 0, limit),
	}
}

func FromRect(r s2.Rect, d Distance, limit int) *Filter {
	return &Filter{
		distance: d,
		contains: r.ContainsLatLng,
		limit:    limit,
		center:   r.Center(),
		data:     make([]Element, 0, limit),
	}
}

// Add registers given position and index in filter
// Element is added if point is withing specified region and
// is not further away than furthest existing element if filter is full
func (f *Filter) Add(ll s2.LatLng, i interface{}) {
	// ll is out of bounds
	if !f.contains(ll) {
		return
	}

	d := f.distance(f.center, ll)
	if f.limit == len(f.data) {
		// discard if distance > maxDistance
		if d > f.data[len(f.data)-1].Distance {
			return
		}
		heap.Pop(&f.data)
	}

	heap.Push(&f.data, Element{d, i})
}

// Elements returns indices and distance of values that fitted criteria
// Result is ordered by distance
func (f *Filter) Elements() []Element {
	r := f.data
	sort.Sort(r)
	return r
}

type elements []Element

func (e elements) Len() int {
	return len(e)
}

func (e elements) Less(i, j int) bool {
	return e[i].Distance < e[j].Distance
}

func (e elements) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e *elements) Push(x interface{}) {
	*e = append(*e, x.(Element))
}

func (e *elements) Pop() interface{} {
	t := *e
	n := len(*e) - 1
	*e = t[0:n]
	return t[n]
}
