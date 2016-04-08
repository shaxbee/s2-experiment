package geodb

import (
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

const EarthCircumreference = 24901.0
const EarthRadius = 6371.0

type Distance func(s2.LatLng, s2.LatLng) float64

func WGS84Haversine(a, b s2.LatLng) float64 {
	return EarthRadius * a.Distance(b).Radians()
}

func NearSphere(ll s2.LatLng, d float64) s2.Cap {
	return s2.CapFromCenterAngle(s2.PointFromLatLng(ll), s1.Angle(d/EarthRadius)*s1.Radian)
}

var DefaultDistance = WGS84Haversine
