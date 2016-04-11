package geodb

import "github.com/golang/geo/s2"

type Location struct {
	Lat    float64 `db:"lat"`
	Lng    float64 `db:"lng"`
	CellID int64   `db:"cell_id"`
}

func (l *Location) LatLng() s2.LatLng {
	return s2.LatLngFromDegrees(l.Lat, l.Lng)
}

func NewLocation(lat, lng float64) Location {
	id := s2.CellIDFromLatLng(s2.LatLngFromDegrees(lat, lng))
	return Location{
		lat,
		lng,
		int64(id),
	}
}
