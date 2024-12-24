package uniex

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/url"
	"sort"
	"sync"

	"github.com/gocarina/gocsv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config
type Config struct {
	MongoDB string // mongodb uri, default: mongodb://127.0.0.1:27117
	Format  string // export format, default: csv [csv|json]
	Scope   string // export scope, default: client [client|infra]
}

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
}

// Stat
type Stat struct {
	MAC      string `bson:"mac"`
	NAME     string `bson:"name"`
	HOSTNAME string `bson:"hostname"`
	IP       string `bson:"ip"`
	TIME     int64  `bson:"time"`
}

// Setup defaults and Sanitize Config
func (c *Config) Setup() (*Config, error) {

	// validate db connection
	switch c.MongoDB {
	case "":
		c.MongoDB = "mongodb://127.0.0.1:27117"
	}
	uri, err := url.Parse(c.MongoDB)
	if err != nil {
		return c, errors.New("invalid mongodb uri: " + c.MongoDB + " error: " + err.Error())
	}
	if uri.Scheme != "mongodb" {
		return c, errors.New("invalid mongodb uri scheme, need mongodb, got: " + uri.Scheme + " error: " + err.Error())
	}
	if _, err := net.LookupIP(uri.Hostname()); err != nil {
		return c, errors.New("unable to dns lookup mongodb hostname: " + uri.Hostname() + " error: " + err.Error())
	}

	// validate output format
	switch c.Format {
	case "csv":
	case "json":
	case "":
		c.Format = "csv"
	default:
		return c, errors.New("invalid export format, need: [csv|json], got: " + c.Format)
	}

	// validate search scope
	switch c.Scope {
	case "client":
	case "infra":
	case "":
		c.Scope = "client"
	default:
		return c, errors.New("invalid export scope, need: [client|infra], got: " + c.Scope)
	}

	// success
	return c, nil
}

// Export Data
func (c *Config) Export() (string, error) {

	// init
	var err error

	// setup default and sanitize input
	if c, err = c.Setup(); err != nil {
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
	var org int64
	for _, d := range devices {
		org = d.LASTSEEN
		d.LASTSEEN = 0
		for _, s := range stats {
			if d.MAC == s.MAC {
				if d.LASTSEEN < s.TIME {
					d.LASTSEEN = s.TIME
					d.NAME = s.NAME
					d.HOSTNAME = s.HOSTNAME
				}
			}
		}
		if org > d.LASTSEEN {
			d.LASTSEEN = org
		}

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
