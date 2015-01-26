
# go-stathat

 Stathat client with batch support.

## Usage

```go
var Endpoint = "http://api.stathat.com/ez"
```
EZ endpoint.

#### type Client

```go
type Client struct {
  Verbose  bool
  Interval time.Duration
  Size     int
}
```

Client which batches stats and flushes at the given Interval or
when the Size limit is exceeded. Set Verbose to true to enable
logging output.

#### func  New

```go
func New(key string) *Client
```
New return a client with "EZ key" and default buffer size of 1000 and interval
of 1s.

#### func (*Client) Close

```go
func (c *Client) Close() error
```
Close and flush metrics.

#### func (*Client) Count

```go
func (c *Client) Count(name string, n int64) error
```
Count for `name`.

#### func (*Client) CountTime

```go
func (c *Client) CountTime(name string, n int64, t time.Time) error
```
Count for `name` with explicit time.

#### func (*Client) Value

```go
func (c *Client) Value(name string, n float64) error
```
Value for `name`.

#### func (*Client) ValueTime

```go
func (c *Client) ValueTime(name string, n float64, t time.Time) error
```
Value for `name` with explicit time.


## License

 MIT