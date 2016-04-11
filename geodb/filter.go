package geodb

import (
	"container/heap"
	"sort"

	"github.com/golang/geo/s2"
)

type Filter struct {
	contains func(s2.LatLng) bool
	limit    int
	origin   s2.LatLng
	data     elements
}

// Element represents pair of distance and index
type Element struct {
	Distance float64
	Index    int
}

func FromCap(c s2.Cap, limit int) *Filter {
	return &Filter{
		contains: func(ll s2.LatLng) bool {
			return c.ContainsPoint(s2.PointFromLatLng(ll))
		},
		limit:  limit,
		origin: s2.LatLngFromPoint(c.Center()),
		data:   make([]Element, 0, limit),
	}
}

func FromRect(r s2.Rect, limit int) *Filter {
	return &Filter{
		contains: r.ContainsLatLng,
		limit:    limit,
		origin:   r.Center(),
		data:     make([]Element, 0, limit),
	}
}

func FromLoop(l s2.Loop, ll s2.LatLng, limit int) *Filter {
	return &Filter{
		contains: func(ll s2.LatLng) bool {
			return l.ContainsPoint(s2.PointFromLatLng(ll))
		},
		limit:  limit,
		origin: ll,
		data:   make([]Element, 0, limit),
	}
}

// Add registers given position and index in filter
// Element is added if point is withing specified region and
// is not further away than furthest existing element if filter is full
func (f *Filter) Add(ll s2.LatLng, i int) {
	// ll is out of bounds
	if !f.contains(ll) {
		return
	}

	d := Distance(f.origin, ll)
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

type GetLatLng func(interface{}) s2.LatLng
type SetDistance func(interface{}, float64)

func (f *Filter) Filter(s []interface{}, ll GetLatLng, d SetDistance) []interface{} {
	for i, v := range s {
		f.Add(ll(v), i)
	}

	es := f.Elements()
	r := make([]interface{}, len(es))
	for i, e := range es {
		v := s[e.Index]
		d(v, e.Distance)
		r[i] = v
	}

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
