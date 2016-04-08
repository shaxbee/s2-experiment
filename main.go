package main

import (
	"fmt"

	"github.com/golang/geo/s2"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shaxbee/s2-experiment/geodb"
	"github.com/shaxbee/s2-experiment/geodb/filter"
)

type job struct {
	Title  string
	Lat    float64
	Lng    float64
	CellID int64
}

func (j *job) LatLng() s2.LatLng {
	return s2.LatLngFromDegrees(j.Lat, j.Lng)
}

func newJob(t string, ll s2.LatLng) job {
	return job{
		t,
		ll.Lat.Degrees(),
		ll.Lng.Degrees(),
		int64(s2.CellIDFromLatLng(ll)),
	}
}

func main() {
	db, err := sqlx.Open("sqlite3", "file:dummy.db?mode=memory&cache=shared")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE jobs (title TEXT, lat DECIMAL(10, 6), lng DECIMAL(10, 6), cell_id BIGINT)")
	if err != nil {
		panic(err)
	}

	insert, err := db.PrepareNamed("INSERT INTO jobs (title, lat, lng, cell_id) VALUES (:title, :lat, :lng, :cellid)")
	if err != nil {
		panic(err)
	}

	jobs := []job{
		newJob("Walk dog", s2.LatLngFromDegrees(14.550983, 121.043799)),
		newJob("Museum guide", s2.LatLngFromDegrees(14.552091, 121.045548)),
		newJob("Buy groceries", s2.LatLngFromDegrees(14.554980, 121.049341)),
	}

	for _, j := range jobs {
		_, err := insert.Exec(&j)
		if err != nil {
			panic(err)
		}
	}

	search := geodb.NearSphere(jobs[0].LatLng(), 0.25)

	g := geodb.New(geodb.Config{})
	where, vals := g.Select(search)

	sel, err := db.Preparex(fmt.Sprint("SELECT title, lat, lng FROM jobs WHERE ", where))
	if err != nil {
		panic(err)
	}

	r, err := sel.Queryx(vals...)
	if err != nil {
		panic(err)
	}

	f := filter.FromCap(search, geodb.DefaultDistance, 10)

	for r.Next() {
		j := &job{}
		r.StructScan(j)
		f.Add(j.LatLng(), j)
	}

	for _, e := range f.Elements() {
		j := e.Value.(*job)
		fmt.Println(j.Title, e.Distance)
	}
}
