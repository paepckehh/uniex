package uniex

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/gocarina/gocsv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// global
const SemVer = "v0.0.6"

// Config
type Config struct {
	MongoDB string // [UNIEX_MONGODB] mongodb uri, default: mongodb://127.0.0.1:27117
	Format  string // [UNIEX_FORMAT]  export format, default: csv [csv|json]
	Scope   string // [UNIEX_SCOPE]   export scope, default: client [client|infra]
}

// Device
type Device struct {
	NAME           string `bson:"name"`
	NOTE           string `bson:"note"`
	HOSTNAME       string // hostname (from stats)
	IP             string // ip address (from stats)
	MAC            string `bson:"mac"`
	OUI            string `bson:"oui"`
	SWITCHNAME     string `bson:"last_uplink_name"`
	SWITCHMAC      string `bson:"last_uplink_mac"`
	SWITCHPORT     string // switch port (from stats)
	VLANNETWORK    string `bson:"last_connection_network_name"`
	FIRSTSEEN      string // rfc3339 timestamp (calculated)
	LASTSEEN       string // rfc3339 timestamp (calculated)
	FIRSTSEEN_UNIX int64  `bson:"first_seen"`
	LASTSEEN_UNIX  int64  `bson:"last_seen"`
}

// Stat
type Stat struct {
	MAC           string `bson:"mac"`
	HOSTNAME      string `bson:"hostname"`
	IP            string `bson:"ip"`
	LASTSEEN_UNIX int64  `bson:"time"`
}

// Export Data
func (c *Config) Export() (string, error) {

	// init
	var err error

	// setup default and sanitize input
	if c, err = c.setup(); err != nil {
		return "", err
	}

	// setup unifi mongodb connection
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(c.MongoDB))
	if err != nil {
		return "", errors.New("unable to connect to mongodb: " + c.MongoDB + ", " + err.Error())
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

	// get all device records
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

	// fetch all stats snipets
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
	var out string
	switch c.Format {
	case "csv":
		out, err = gocsv.MarshalString(&devices)
		if err != nil {
			return "", err
		}
	case "json":
		j, err := json.Marshal(&devices)
		if err != nil {
			return "", err
		}
		out = string(j)
	default:
		panic("internal error, unsupported output format") // unreachable
	}
	return out, nil
}
