package geodb

import (
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

const EarthRadius = 6371.0

func Distance(a, b s2.LatLng) float64 {
	return EarthRadius * a.Distance(b).Radians()
}

func Near(ll s2.LatLng, d float64) s2.Cap {
	return s2.CapFromCenterAngle(s2.PointFromLatLng(ll), s1.Angle(d/EarthRadius)*s1.Radian)
}
