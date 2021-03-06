package hcm

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	DefaultEndpointFmt               = "https://push-api.cloud.huawei.com/v1/%d/messages:send"
	DefaultTimeout     time.Duration = 30 * time.Second
)

const (
	minBackoff = 100 * time.Millisecond
	maxBackoff = 1 * time.Minute
)

var (
	// ErrInvalidAPIKey occurs if API key is not set.
	ErrInvalidAppID = errors.New("app ID cannot be empty")
)

// Client abstracts the interaction between the application server and the
// FCM server via HTTP protocol. The developer must obtain an API key from the
// Google APIs Console page and pass it to the `Client` so that it can
// perform authorized requests on the application server's behalf.
// To send a message to one or more devices use the Client's Send.
//
// If the `HTTP` field is nil, a zeroed http.Client will be allocated and used
// to send messages.
type Client struct {
	apiKey   string
	client   *fasthttp.Client
	endpoint string
	timeout  time.Duration
}

// NewClient creates new Firebase Cloud Messaging Client based on API key and
// with default endpoint and http client.
func NewClient(appID int) (*Client, error) {
	if appID == 0 {
		return nil, ErrInvalidAppID
	}
	c := &Client{
		endpoint: fmt.Sprintf(DefaultEndpointFmt, appID),
		client:   &fasthttp.Client{},
		timeout:  DefaultTimeout,
	}
	return c, nil
}

// Send sends a message to the FCM server without retrying in case of service
// unavailability. A non-nil error is returned if a non-recoverable error
// occurs (i.e. if the response status is not "200 OK").
func (c *Client) Send(msg *Message, accessToken string) (*Response, error) {
	// validate
	if err := msg.Validate(); err != nil {
		return nil, err
	}

	// marshal message
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return c.send(data, accessToken)
}

// SendWithRetry sends a message to the FCM server with defined number of
// retrying in case of temporary error.
func (c *Client) SendWithRetry(msg *Message, accessToken string, retryAttempts int) (*Response, error) {
	// validate
	if err := msg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid msg: %v", err)
	}
	// marshal message
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("cannot create msg json: %v", err)
	}

	resp := new(Response)
	err = retry(func() error {
		var er error
		resp, er = c.send(data, accessToken)
		return er
	}, retryAttempts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// send sends a request.
func (c *Client) send(data []byte, accessToken string) (*Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.SetBody(data)
	req.SetRequestURI(c.endpoint)

	err := c.client.DoTimeout(req, resp, c.timeout)
	if err != nil {
		return nil, err
	}

	sc := resp.StatusCode()
	if sc != http.StatusOK {
		return nil, fmt.Errorf("%d error: %s", sc, resp.String())
	}
	response := new(Response)
	body := resp.Body()
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("cannot parse resp body: %+v", err)
	}

	return response, nil
}

func retry(fn func() error, attempts int) error {
	var attempt int
	for {
		err := fn()
		if err == nil {
			return nil
		}

		if tErr, ok := err.(net.Error); !ok || !tErr.Temporary() {
			return err
		}

		attempt++
		backoff := minBackoff * time.Duration(attempt*attempt)
		if attempt > attempts || backoff > maxBackoff {
			return err
		}

		time.Sleep(backoff)
	}
}
