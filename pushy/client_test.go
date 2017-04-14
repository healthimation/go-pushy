package pushy

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"reflect"

	"github.com/healthimation/go-client/client"
)

func testClient(handler http.HandlerFunc, timeout time.Duration) (Client, *httptest.Server) {
	ts := httptest.NewServer(handler)
	finder := func(serviceName string, useTLS bool) (url.URL, error) {
		ret, err := url.Parse(ts.URL)
		if err != nil || ret == nil {
			return url.URL{}, err
		}
		return *ret, err
	}
	c := &pushyClient{
		c:      client.NewBaseClient(finder, "pushy", true, timeout),
		apiKey: "",
	}
	return c, ts
}

func makeBoolPtr(v bool) *bool {
	return &v
}
func makeInt64Ptr(v int64) *int64 {
	return &v
}
func makeStrPtr(v string) *string {
	return &v
}

func TestUnit_PushToDevices(t *testing.T) {

	type testcase struct {
		name             string
		handler          http.HandlerFunc
		timeout          time.Duration
		ctx              context.Context
		tokens           []string
		data             interface{}
		options          *PushOptions
		expectedErrCode  *string
		expectedResponse *string
	}

	testcases := []testcase{
		{
			name: "base path",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"success":true, "id":"5742ea5dacf3a92e17ba7126"}`)
			}),
			timeout:          5 * time.Second,
			ctx:              context.Background(),
			tokens:           []string{"1", "2", "3"},
			data:             map[string]string{"foo": "bar"},
			expectedResponse: makeStrPtr("5742ea5dacf3a92e17ba7126"),
		},
		{
			name: "alternate path - with options",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"success":true, "id":"5742ea5dacf3a92e17ba7126"}`)
			}),
			timeout:          5 * time.Second,
			ctx:              context.Background(),
			tokens:           []string{"1", "2", "3"},
			data:             map[string]string{"foo": "bar"},
			expectedResponse: makeStrPtr("5742ea5dacf3a92e17ba7126"),
			options: &PushOptions{
				ContentAvailable: makeBoolPtr(true),
				MutableContent:   makeBoolPtr(true),
				TimeToLive:       makeInt64Ptr(100),
				Notification: &NotificationOptions{
					Body:  makeStrPtr("test body"),
					Badge: makeInt64Ptr(1),
				},
			},
		},
		{
			name: "exceptional path",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"error":"test error"}`)
			}),
			timeout:         5 * time.Second,
			ctx:             context.Background(),
			tokens:          []string{"1", "2", "3"},
			data:            map[string]string{"foo": "bar"},
			expectedErrCode: makeStrPtr(ErrorAPI),
		},
		{
			name: "exceptional path - timeout",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(2 * time.Millisecond)
				fmt.Fprint(w, `{"success":true, "id":"5742ea5dacf3a92e17ba7126"}`)
			}),
			timeout:         1 * time.Millisecond,
			ctx:             context.Background(),
			tokens:          []string{"1", "2", "3"},
			data:            map[string]string{"foo": "bar"},
			expectedErrCode: makeStrPtr(client.ErrorRequestError),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c, ts := testClient(tc.handler, tc.timeout)
			defer ts.Close()
			ret, err := c.PushToDevices(tc.ctx, tc.tokens, tc.data, tc.options)
			if tc.expectedErrCode != nil || err != nil {
				if tc.expectedErrCode == nil {
					t.Fatalf("Unexpected error occurred (%#v)", err)
				}
				if err == nil {
					t.Fatalf("Expected error did not occur")
				}
				if err.Code() != *tc.expectedErrCode {
					t.Fatalf("Actual error (%#v) did not match expected (%#v)", err.Code(), tc.expectedErrCode)
				}
				if !reflect.DeepEqual(tc.expectedResponse, ret) {
					t.Fatalf("Actual response (%#v) did not match expected (%#v)", ret, tc.expectedResponse)
				}
			}
		})
	}
}

func TestUnit_PushToTopic(t *testing.T) {

	type testcase struct {
		name             string
		handler          http.HandlerFunc
		timeout          time.Duration
		ctx              context.Context
		topic            string
		data             interface{}
		options          *PushOptions
		expectedErrCode  *string
		expectedResponse *string
	}

	testcases := []testcase{
		{
			name: "base path",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"success":true, "id":"5742ea5dacf3a92e17ba7126"}`)
			}),
			timeout:          5 * time.Second,
			ctx:              context.Background(),
			topic:            "foobar",
			data:             map[string]string{"foo": "bar"},
			expectedResponse: makeStrPtr("5742ea5dacf3a92e17ba7126"),
		},
		{
			name: "alternate path - with options",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"success":true, "id":"5742ea5dacf3a92e17ba7126"}`)
			}),
			timeout:          5 * time.Second,
			ctx:              context.Background(),
			topic:            "foobar",
			data:             map[string]string{"foo": "bar"},
			expectedResponse: makeStrPtr("5742ea5dacf3a92e17ba7126"),
			options: &PushOptions{
				ContentAvailable: makeBoolPtr(true),
				MutableContent:   makeBoolPtr(true),
				TimeToLive:       makeInt64Ptr(100),
				Notification: &NotificationOptions{
					Body:  makeStrPtr("test body"),
					Badge: makeInt64Ptr(1),
				},
			},
		},
		{
			name: "alternate path - with /topics/",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"success":true, "id":"5742ea5dacf3a92e17ba7126"}`)
			}),
			timeout:          5 * time.Second,
			ctx:              context.Background(),
			topic:            "/topics/foobar",
			data:             map[string]string{"foo": "bar"},
			expectedResponse: makeStrPtr("5742ea5dacf3a92e17ba7126"),
			options: &PushOptions{
				ContentAvailable: makeBoolPtr(true),
				MutableContent:   makeBoolPtr(true),
				TimeToLive:       makeInt64Ptr(100),
				Notification: &NotificationOptions{
					Body:  makeStrPtr("test body"),
					Badge: makeInt64Ptr(1),
				},
			},
		},
		{
			name: "exceptional path",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"error":"test error"}`)
			}),
			timeout:         5 * time.Second,
			ctx:             context.Background(),
			topic:           "foobar",
			data:            map[string]string{"foo": "bar"},
			expectedErrCode: makeStrPtr(ErrorAPI),
		},
		{
			name: "exceptional path - timeout",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(2 * time.Millisecond)
				fmt.Fprint(w, `{"success":true, "id":"5742ea5dacf3a92e17ba7126"}`)
			}),
			timeout:         1 * time.Millisecond,
			ctx:             context.Background(),
			topic:           "foobar",
			data:            map[string]string{"foo": "bar"},
			expectedErrCode: makeStrPtr(client.ErrorRequestError),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c, ts := testClient(tc.handler, tc.timeout)
			defer ts.Close()
			ret, err := c.PushToTopic(tc.ctx, tc.topic, tc.data, tc.options)
			if tc.expectedErrCode != nil || err != nil {
				if tc.expectedErrCode == nil {
					t.Fatalf("Unexpected error occurred (%#v)", err)
				}
				if err == nil {
					t.Fatalf("Expected error did not occur")
				}
				if err.Code() != *tc.expectedErrCode {
					t.Fatalf("Actual error (%#v) did not match expected (%#v)", err.Code(), tc.expectedErrCode)
				}
				if !reflect.DeepEqual(tc.expectedResponse, ret) {
					t.Fatalf("Actual response (%#v) did not match expected (%#v)", ret, tc.expectedResponse)
				}
			}
		})
	}
}
