package stathat

import "encoding/json"
import "net/http"
import "bytes"
import "time"
import "log"

// EZ endpoint.
var Endpoint = "http://api.stathat.com/ez"

// Count.
type count struct {
	Stat      string `json:"stat"`
	Count     int64  `json:"count"`
	Timestamp int64  `json:"t,emitempty"`
}

// Value.
type value struct {
	Stat      string  `json:"stat"`
	Value     float64 `json:"value"`
	Timestamp int64   `json:"t,emitempty"`
}

// Request body.
type body struct {
	Key  string        `json:"ezkey"`
	Data []interface{} `json:"data"`
}

// Client which batches stats and flushes at the given Interval or
// when the Size limit is exceeded. Set Verbose to true to enable
// logging output.
type Client struct {
	Verbose  bool
	Interval time.Duration
	Size     int
	key      string
	stats    chan interface{}
	quit     chan bool
}

// New return a client with "EZ key" and default buffer size of 1000 and interval of 1s.
func New(key string) *Client {
	c := &Client{
		stats:    make(chan interface{}, 1000),
		quit:     make(chan bool),
		Interval: time.Second,
		Size:     200,
		key:      key,
	}

	go c.loop()

	return c
}

// Count for `name`.
func (c *Client) Count(name string, n int64) error {
	return c.CountTime(name, n, time.Now())
}

// Count for `name` with explicit time.
func (c *Client) CountTime(name string, n int64, t time.Time) error {
	c.stats <- &count{
		Stat:      name,
		Count:     n,
		Timestamp: t.Unix(),
	}
	return nil
}

// Value for `name`.
func (c *Client) Value(name string, n float64) error {
	return c.ValueTime(name, n, time.Now())
}

// Value for `name` with explicit time.
func (c *Client) ValueTime(name string, n float64, t time.Time) error {
	c.stats <- &value{
		Stat:      name,
		Value:     n,
		Timestamp: t.Unix(),
	}
	return nil
}

// Close and flush metrics.
func (c *Client) Close() error {
	c.quit <- true
	close(c.stats)
	<-c.quit
	return nil
}

// Log when everbose.
func (c *Client) log(msg string, args ...interface{}) {
	if c.Verbose {
		log.Printf("stathat: "+msg, args...)
	}
}

// Send batch request to Stathat.
func (c *Client) send(stats []interface{}) {
	body := &body{
		Key:  c.key,
		Data: stats,
	}

	b, err := json.Marshal(body)
	if err != nil {
		c.log("error marshalling stats: %s", err)
		return
	}

	res, err := http.Post(Endpoint, "application/json", bytes.NewReader(b))
	if err != nil {
		c.log("error posting stats: %s", err)
		return
	}

	c.log("response -> %s", res.Status)
}

// Batch loop.
func (c *Client) loop() {
	var stats []interface{}
	tick := time.NewTicker(c.Interval)

	for {
		select {
		case stat := <-c.stats:
			c.log("buffer (%d/%d) %#v", len(stats), c.Size, stat)
			stats = append(stats, stat)
			if len(stats) == c.Size {
				c.log("exceeded %d messages – flushing", c.Size)
				c.send(stats)
				stats = nil
			}
		case <-tick.C:
			if len(stats) > 0 {
				c.log("interval reached - flushing")
				c.send(stats)
				stats = nil
			} else {
				c.log("interval reached – nothing to send")
			}
		case <-c.quit:
			c.log("exit requested – flushing")
			c.send(stats)
			c.log("exit")
			c.quit <- true
			return
		}
	}
}
