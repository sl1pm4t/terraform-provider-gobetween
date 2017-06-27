package gobetween

import (
	"errors"
	"fmt"
	"log"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/sl1pm4t/snooze"
	gbconfig "github.com/yyyar/gobetween/src/config"
)

// GbClient is the GoBetween API client
type GbClient struct {
	Api *api
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

type api struct {
	GetSystemInfo  func() error                                     `method:"GET"    contentType:"application/json" path:"/"`
	DumpConfig     func(format string) (string, error)              `method:"GET"    contentType:"application/json" path:"/dump?format={0}"`
	GetServer      func(name string) (*gbconfig.Server, error)      `method:"GET"    contentType:"application/json" path:"/servers/{0}"`
	GetServers     func() ([]*gbconfig.Server, error)               `method:"GET"    contentType:"application/json" path:"/servers"`
	AddServer      func(name string, server *gbconfig.Server) error `method:"POST"   contentType:"application/json" path:"/servers/{0}"`
	DeleteServer   func(name string) error                          `method:"DELETE" contentType:"application/json" path:"/servers/{0}"`
	GetServerStats func(name string) (*ServerStats, error)          `method:"PUT"    contentType:"application/json" path:"/servers/{0}/stats"`
}

type ServerStats struct {
	ActiveConnections int        `json:"active_connections"`
	RxTotal           int        `json:"rx_total"`
	TxTotal           int        `json:"tx_total"`
	RxSecond          int        `json:"rx_second"`
	TxSecond          int        `json:"tx_second"`
	Backends          []*Backend `json:"backends"`
}

type Backend struct {
	Host     string       `json:"host"`
	Port     string       `json:"port"`
	Priority int          `json:"priority"`
	Weight   int          `json:"weight"`
	Stats    BackendStats `json:"stats"`
}

type BackendStats struct {
	Live               bool `json:"live"`
	TotalConnections   int  `json:"total_connections"`
	ActiveConnections  int  `json:"active_connections"`
	RefusedConnections int  `json:"refused_connections"`
	RxBytes            int  `json:"rx"`
	TxBytes            int  `json:"tx"`
	RxPerSecond        int  `json:"rx_second"`
	TxPerSecond        int  `json:"tx_second"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func (c *GbClient) Init(addr, username, password string) error {
	// Create cookie jar
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		log.Fatal(err)
	}

	client := snooze.Client{
		Root: addr,
		Before: func(r *retryablehttp.Request, c *retryablehttp.Client) {
			r.SetBasicAuth(username, password)

			r.Header.Set("User-Agent", "terraform-provider-gobetween")
			r.Header.Add("Accept", `application/json`)

			c.HTTPClient.Jar = jar
			timeout, _ := time.ParseDuration("180s")
			c.HTTPClient.Timeout = timeout
		},
		HandleError: handleApiError,
		// Logger:      logging.Logger(),
	}

	c.Api = new(api)
	client.Create(c.Api)

	return nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func handleApiError(err *snooze.ErrorResponse) error {
	switch {
	case err.StatusCode > 399:
		return fmt.Errorf("GoBetween API Error Response: %s", err.Status)
	case strings.Contains(err.ResponseContentType, "xml"):
		if len(err.ResponseBody) > 0 {
			return errors.New(string(err.ResponseBody))
		}
	default:
		// attempt to transate to ErrorInfo struct
		return fmt.Errorf("got unknown error from API: %s", spew.Sdump(err))
	}

	return nil
}
