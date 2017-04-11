package pushy

import (
	"context"
	"time"

	"net/http"
	"net/url"

	"strings"

	"github.com/healthimation/go-client/client"
	"github.com/healthimation/go-glitch/glitch"
)

//Error codes
const (
	ErrorAPI = "ERROR_API"
)

// Client can make requests to the pushy api
type Client interface {
	PushToDevices(ctx context.Context, tokens []string, data interface{}, options *PushOptions) glitch.DataError
	PushToTopic(ctx context.Context, topic string, data interface{}, options *PushOptions) glitch.DataError
}

type pushyClient struct {
	c      client.BaseClient
	apiKey string
}

// NewClient returns a new pushy client
func NewClient(apiKey string, timeout time.Duration) Client {
	return &pushyClient{
		c:      client.NewBaseClient(findPushy, "pushy", true, timeout),
		apiKey: apiKey,
	}
}

func (p *pushyClient) PushToDevices(ctx context.Context, tokens []string, data interface{}, options *PushOptions) glitch.DataError {
	bodyObj := pushRequest{
		Tokens: tokens,
		Data:   data,
	}

	return p.push(ctx, bodyObj, options)
}

func (p *pushyClient) PushToTopic(ctx context.Context, topic string, data interface{}, options *PushOptions) glitch.DataError {
	if !strings.HasPrefix(topic, "/topics/") {
		topic = "/topics/" + topic
	}

	bodyObj := pushRequest{
		To:   topic,
		Data: data,
	}

	return p.push(ctx, bodyObj, options)
}

func (p *pushyClient) push(ctx context.Context, bodyObj pushRequest, options *PushOptions) glitch.DataError {
	slug := "push"
	h := http.Header{}
	h.Set("Content-type", "application/json")

	if options != nil {
		bodyObj.ContentAvailable = options.ContentAvailable
		bodyObj.MutableContent = options.MutableContent
		bodyObj.Notification = options.Notification
		bodyObj.TimeToLive = options.TimeToLive

	}

	body, err := client.ObjectToJSONReader(bodyObj)
	if err != nil {
		return err
	}

	query := url.Values{}
	query.Set("api_key", p.apiKey)

	res := &response{}

	err = p.c.Do(ctx, http.MethodPost, slug, query, h, body, res)
	if err != nil {
		return err
	}

	if res.Error != nil {
		return glitch.NewDataError(nil, ErrorAPI, *res.Error)
	}
	return nil
}
