package uniex

import (
	"context"
	"errors"
	"net"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Setup defaults and sanitize config
func (c *Config) setup() error {

	// parse input
	// validate input

	// db target
	switch c.MongoDB {
	case "":
		c.MongoDB = "mongodb://127.0.0.1:27117"
	}
	uri, err := url.Parse(c.MongoDB)
	if err != nil {
		return errors.New("invalid mongodb uri: " + c.MongoDB + " error: " + err.Error())
	}
	if uri.Scheme != "mongodb" {
		return errors.New("invalid mongodb uri scheme, need mongodb, got: " + uri.Scheme + " error: " + err.Error())
	}

	// validate output format
	switch c.Format {
	case "csv":
	case "json":
	case "":
		c.Format = "csv"
	default:
		return errors.New("invalid export format, need: [csv|json], got: " + c.Format)
	}

	// validate search scope
	switch c.Scope {
	case "client":
	case "":
		c.Scope = "client"
	default:
		return errors.New("invalid export scope, need: [client|infra], got: " + c.Scope)
	}

	// input parsing done
	// validate live db connection

	// lookup target
	if _, err := net.LookupIP(uri.Hostname()); err != nil {
		return errors.New("unable to dns lookup mongodb hostname: " + uri.Hostname() + " error: " + err.Error())
	}

	// setup test db client connection
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	if c.MongoDB == "mongodb://127.0.0.1:27117" {
		// reduce timeout for localhost db
		ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	}
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri.String()))
	if err != nil {
		return errors.New("mongodb connection client setup error:" + err.Error())
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// test connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return errors.New("mongodb connection client ping error:" + err.Error())
	}

	// success
	return nil
}
