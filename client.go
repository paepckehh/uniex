package uniex

import (
	"context"
	"encoding/json"
	"sort"
	"sync"
	"time"

	"github.com/gocarina/gocsv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *Config) clientInventory(db *mongo.Client) ([]byte, error) {

	// setup
	var (
		err     error
		wg      sync.WaitGroup
		devices []Device
		stats   []Stat
	)

	// get all device records
	wg.Add(1)
	go func() {

		// clean exit
		defer wg.Done()

		// setup user db query context
		c := db.Database("ace").Collection("user")

		// setup query
		q, err := c.Find(context.TODO(), bson.M{})
		if err != nil {
			panic(err)
		}

		// perform query
		if err := q.All(context.TODO(), &devices); err != nil {
			panic(err)
		}
	}()

	// fetch all stats snipets
	wg.Add(1)
	go func() {

		// clean exit
		defer wg.Done()

		// setup user db query context
		c := db.Database("ace_stat").Collection("stat_archive")

		// setup query
		q, err := c.Find(context.TODO(), bson.M{})
		if err != nil {
			panic(err)
		}

		// perform query
		if err := q.All(context.TODO(), &stats); err != nil {
			panic(err)
		}
	}()

	// wait till all queries done
	wg.Wait()

	// parste all stats, add missing data into device records
	var ts int64
	for i, d := range devices {
		d.LASTSEEN_UNIX = d.LASTSEEN_UNIX * 1000 // stats stamps have higher time resolution
		ts = d.LASTSEEN_UNIX
		d.LASTSEEN_UNIX = 0
		for _, s := range stats {
			// find latest matching stats record, get data
			if d.MAC == s.MAC {
				if d.LASTSEEN_UNIX < s.LASTSEEN_UNIX {
					d.LASTSEEN_UNIX = s.LASTSEEN_UNIX
					d.HOSTNAME = s.HOSTNAME
					d.IP = s.IP
				}
			}
		}
		if ts > d.LASTSEEN_UNIX {
			d.LASTSEEN_UNIX = ts
		}
		d.LASTSEEN_UNIX = d.LASTSEEN_UNIX / 1000
		d.LASTSEEN = time.Unix(d.LASTSEEN_UNIX, 0).Format(time.RFC3339)
		d.FIRSTSEEN = time.Unix(d.FIRSTSEEN_UNIX, 0).Format(time.RFC3339)
		devices[i] = d // write back in the array
	}

	// sort devices by name
	sort.Slice(devices, func(i, j int) bool {
		return devices[i].NAME < devices[j].NAME
	})

	// output
	var out []byte
	switch c.Format {
	case "csv":
		out, err = gocsv.MarshalBytes(&devices)
	case "json":
		out, err = json.Marshal(&devices)
	default:
		panic("internal error, unsupported output format") // unreachable
	}
	if err != nil {
		return nil, err
	}
	return out, nil
}
