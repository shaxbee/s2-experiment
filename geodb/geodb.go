package geodb

import (
	sq "github.com/elgris/sqrl"
	"github.com/golang/geo/s2"
	"github.com/jmoiron/sqlx"
)

type GeoDB struct {
	coverer *s2.RegionCoverer
}

type Config struct {
	MinLevel int
	MaxLevel int
	MaxCells int
}

func New(c Config) *GeoDB {
	return &GeoDB{
		coverer: &s2.RegionCoverer{
			MinLevel: c.MinLevel,
			MaxLevel: c.MaxLevel,
			MaxCells: c.MaxCells,
		},
	}
}

func (g *GeoDB) Within(q *sq.SelectBuilder, r s2.Region) *sq.SelectBuilder {
	q = q.Columns("lat", "lng", "cell_id")
	c := g.coverer.Covering(r)
	for _, x := range c {
		q = q.Where(match, int64(x.RangeMin()), int64(x.RangeMax()))
	}

	return q
}

func (g *GeoDB) QueryWithin(db *sqlx.DB, q *sq.SelectBuilder, r s2.Region) (*sqlx.Rows, error) {
	s, v, err := g.Within(q, r).ToSql()
	if err != nil {
		return nil, err
	}

	return db.Queryx(s, v...)
}

var match = "(cell_id BETWEEN ? AND ?)"
