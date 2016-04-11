package main

import (
	"fmt"

	sq "github.com/elgris/sqrl"
	"github.com/golang/geo/s2"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shaxbee/s2-experiment/geodb"
)

type job struct {
	geodb.Location
	Title    string
	Distance float64
}

func newJob(t string, lat, lng float64) job {
	return job{
		Location: geodb.NewLocation(lat, lng),
		Title:    t,
	}
}

func prepare() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", "file:dummy.db?mode=memory&cache=shared")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE jobs (title TEXT, lat DECIMAL(10, 6), lng DECIMAL(10, 6), cell_id BIGINT)")
	if err != nil {
		return nil, err
	}

	t, err := db.Beginx()
	if err != nil {
		return nil, err
	}

	insert, err := t.PrepareNamed("INSERT INTO jobs (title, lat, lng, cell_id) VALUES (:title, :lat, :lng, :cell_id)")
	if err != nil {
		return nil, err
	}

	jobs := []job{
		newJob("Walk dog", 14.550983, 121.043799),
		newJob("Museum guide", 14.552091, 121.045548),
		newJob("Buy groceries", 14.554980, 121.049341),
	}

	for _, j := range jobs {
		if _, err := insert.Exec(&j); err != nil {
			return nil, err
		}
	}

	if err := t.Commit(); err != nil {
		return nil, err
	}

	return db, nil
}

var geo = geodb.New(geodb.Config{})

func findNearbyJobs(db *sqlx.DB, ll s2.LatLng, d float64, l int) ([]job, error) {
	search := geodb.Near(ll, d)
	r, err := geo.QueryWithin(db, sq.Select("title").From("jobs"), search)
	if err != nil {
		return nil, err
	}

	s := []job{}
	f := geodb.FromCap(search, 10)
	for i := 0; r.Next(); i++ {
		j := job{}
		if err := r.StructScan(&j); err != nil {
			return nil, err
		}
		f.Add(j.LatLng(), i)
		s = append(s, j)
	}

	es := f.Elements()
	jobs := make([]job, len(es))
	for i, e := range f.Elements() {
		j := s[e.Index]
		j.Distance = e.Distance
		jobs[i] = j
	}

	return jobs, nil
}

func main() {
	db, err := prepare()
	if err != nil {
		panic(err)
	}

	jobs, err := findNearbyJobs(db, s2.LatLngFromDegrees(14.550983, 121.043799), 0.743, 10)
	if err != nil {
		panic(err)
	}

	for _, j := range jobs {
		fmt.Printf("%s, %dm\n", j.Title, int(j.Distance*1000.0))
	}
}
