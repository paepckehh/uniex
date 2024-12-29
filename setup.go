package uniex

import (
	"errors"
	"net"
	"net/url"
)

// Setup defaults and Sanitize Config
func (c *Config) setup() (*Config, error) {

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
	case "":
		c.Scope = "client"
	default:
		return c, errors.New("invalid export scope, need: [client|infra], got: " + c.Scope)
	}

	// success
	return c, nil
}
