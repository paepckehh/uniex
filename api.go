package uniex

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// global
const SemVer = "v0.1.3"

// Config
type Config struct {
	MongoDB string // [UNIEX_MONGODB] mongodb uri, default: mongodb://127.0.0.1:27117
	Format  string // [UNIEX_FORMAT]  export format, default: csv [csv|json]
	Scope   string // [UNIEX_SCOPE]   export scope, default: client [client] TODO: infra
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
	SWITCHMAC     string `bson:"sw_mac"`
	SWITCHPORT    string `bson:"sw_port"`
	LASTSEEN_UNIX int64  `bson:"time"`
}

// New Config
func NewConfig() *Config {
	c := &Config{}
	c.setup()
	return c
}

// Export Data
func (c *Config) Export() ([]byte, error) {

	// init
	var err error

	// setup default and sanitize input
	if c, err = c.setup(); err != nil {
		return nil, err
	}

	// setup unifi mongodb connection
	db, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(c.MongoDB))
	if err != nil {
		return nil, errors.New("unable to connect to mongodb: " + c.MongoDB + ", " + err.Error())
	}

	// prep global clean exit
	defer func() {
		if err := db.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// setup output
	var out []byte

	// fetch
	switch c.Scope {
	case "client":
		out, err = c.clientInventory(db)
	default:
		panic("internal error, unsupported scope") // unreachable
	}
	if err != nil {
		return nil, err
	}
	return out, nil
}
