package geodb

import (
	"bytes"
	"fmt"

	"github.com/golang/geo/s2"
)

type GeoDB struct {
	cellID  string
	coverer *s2.RegionCoverer
}

type Config struct {
	CellID   string
	MinLevel int
	MaxLevel int
	MaxCells int
}

const (
	DefaultCellID = "cell_id"
)

func New(c Config) *GeoDB {
	if c.CellID == "" {
		c.CellID = DefaultCellID
	}

	return &GeoDB{
		cellID: c.CellID,
		coverer: &s2.RegionCoverer{
			MinLevel: c.MinLevel,
			MaxLevel: c.MaxLevel,
			MaxCells: c.MaxCells,
		},
	}
}

func (g *GeoDB) Select(r s2.Region) (string, []interface{}) {
	c := g.coverer.Covering(r)
	p := make([]interface{}, len(c)*2)
	for i, c := range c {
		p[i*2], p[i*2+1] = int64(c.RangeMin()), int64(c.RangeMax())
	}

	return g.query(len(c)), p
}

func (g *GeoDB) query(n int) string {
	if n == 0 {
		return ""
	}

	match := []byte(fmt.Sprint("(", g.cellID, " BETWEEN ? AND ?)"))
	or := []byte(" OR ")

	b := bytes.NewBuffer(make([]byte, 0, len(match)*n+len(or)*(n-1)))

	b.Write(match)
	for i := 0; i < n-1; i++ {
		b.Write(or)
		b.Write(match)
	}

	return b.String()
}
