package main

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/gocarina/gocsv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Device
type Device struct {
	NAME        string // host
	HOSTNAME    string // hostname
	IP          string // ip address
	MAC         string `bson:"mac"`
	OUI         string `bson:"oui"`
	SWITCH      string `bson:"last_uplink_name"`
	SWITCHPORT  string // switch port
	VLANNETWORK string `bson:"last_connection_network_name"`
	FIRSTSEEN   int64  `bson:"assoc_time"`
	LASTSEEN    int64  // last stat timestamp
	SITEID      string `bson:"site_id"`
}

// Stat
type Stat struct {
	MAC      string `bson:"mac"`
	NAME     string `bson:"name"`
	HOSTNAME string `bson:"hostname"`
	IP       string `bson:"ip"`
	TIME     int64  `bson:"time"`
}

// main
func main() {

	// setup connection
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://127.0.0.1:27117"))
	if err != nil {
		panic(err)
	}

	// prep global clean exit
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// setup
	var wg sync.WaitGroup
	var devices []Device
	var stats []Stat

	// get all devices
	wg.Add(1)
	go func() {

		// clean exit
		defer wg.Done()

		// setup user db query context
		c := client.Database("ace").Collection("user")

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

	// get all stats
	wg.Add(1)
	go func() {

		// clean exit
		defer wg.Done()

		// setup user db query context
		c := client.Database("ace_stat").Collection("stat_archive")

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

	// parste last stats data into device record
	var org int64
	for _, d := range devices {
		org = d.LASTSEEN
		d.LASTSEEN = 0
		for _, s := range stats {
			if d.LASTSEEN < s.TIME {
				d.LASTSEEN = s.TIME
				d.NAME = s.NAME
				d.HOSTNAME = s.HOSTNAME
			}
		}
		if org > d.LASTSEEN {
			d.LASTSEEN = org
		}

	}

	// sort devices by mac
	sort.Slice(devices, func(i, j int) bool {
		return devices[i].NAME < devices[j].NAME
	})

	// write as csv
	csv, err := gocsv.MarshalString(&devices)
	if err != nil {
		panic(err)
	}

	// out
	fmt.Println(csv)
}
