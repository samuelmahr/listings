package main

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/gocarina/gocsv"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
)

//
//id,street,status,price,bedrooms,bathrooms,sq_ft,lat,lng
type listing struct {
	ID            int64   `csv:"id"`
	Street        string  `csv:"street"`
	Status        string  `csv:"status"`
	Price         int     `csv:"price"`
	Bedrooms      int     `csv:"bedrooms"`
	Bathrooms     int     `csv:"bathrooms"`
	SquareFootage int     `csv:"sq_ft"`
	Latitude      float64 `csv:"lat"`
	Longitude     float64 `csv:"lng"`
}

func main() {
	csvFile, err := os.Open("./listing-details.csv")
	if err != nil {
		panic(err)
	}

	fmt.Println("file read")
	defer csvFile.Close()

	listings := make([]listing, 0)

	if err = gocsv.UnmarshalFile(csvFile, &listings); err != nil {
		panic(err)
	}

	fmt.Println("unmarshaled")
	query := sq.Insert("features.listings").Columns("id", "street", "status", "price", "bedrooms", "bathrooms", "sq_ft", "lat", "lng").PlaceholderFormat(sq.Dollar)
	fmt.Println("building insert query")
	records := 0
	for _, l := range listings {
		query = query.Values(l.ID, l.Street, l.Status, l.Price, l.Bedrooms, l.Bathrooms, l.SquareFootage, l.Latitude, l.Longitude)
		records += 1

		// insert 100 at a time
		if records == 100 {
			sql, args, err := query.ToSql()
			if err != nil {
				panic(err)
			}

			if len(args) == 0 {
				return
			}

			// fmt.Printf("Query: %s\n", sql)
			// fmt.Printf("args: %#v\n", args)

			db, err := sqlx.Connect("postgres", "postgres://master:Passw0rd@localhost:5432/market?sslmode=disable&TimeZone=utc")
			if err != nil {
				panic(err)
			}

			fmt.Println("inserting...")
			_, err = db.Exec(sql, args...)
			if err != nil {
				panic(err)
			}

			// start query over
			query = sq.Insert("features.listings").Columns("id", "street", "status", "price", "bedrooms", "bathrooms", "sq_ft", "lat", "lng").PlaceholderFormat(sq.Dollar)
			records = 0
			fmt.Println("building insert query")
		}
	}

	sql, args, err := query.ToSql()
	if err != nil {
		panic(err)
	}

	if len(args) == 0 {
		return
	}

	// fmt.Printf("Query: %s\n", sql)
	// fmt.Printf("args: %#v\n", args)

	// insert remaining
	db, err := sqlx.Connect("postgres", "postgres://master:Passw0rd@localhost:5432/market?sslmode=disable&TimeZone=utc")
	if err != nil {
		panic(err)
	}

	fmt.Println("inserting...")
	_, err = db.Exec(sql, args...)
	if err != nil {
		panic(err)
	}

}
