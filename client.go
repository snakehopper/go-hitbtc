package spiral

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type client struct {
	apiKey      string
	apiSecret   string
	httpClient  *http.Client
	httpTimeout time.Duration
	debug       bool
}

// NewClient return a new Spiral HTTP client
func NewClient(apiKey, apiSecret string) (c *client) {
	return &client{apiKey, apiSecret, &http.Client{}, 30 * time.Second, false}
}

// NewClientWithCustomHttpConfig returns a new Spiral HTTP client using the predefined http client
func NewClientWithCustomHttpConfig(apiKey, apiSecret string, httpClient *http.Client) (c *client) {
	timeout := httpClient.Timeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &client{apiKey, apiSecret, httpClient, timeout, false}
}

// NewClient returns a new Spiral HTTP client with custom timeout
func NewClientWithCustomTimeout(apiKey, apiSecret string, timeout time.Duration) (c *client) {
	return &client{apiKey, apiSecret, &http.Client{}, timeout, false}
}

func (c client) dumpRequest(r *http.Request) {
	if r == nil {
		log.Print("dumpReq ok: <nil>")
		return
	}
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Print("dumpReq err:", err)
	} else {
		log.Print("dumpReq ok:", string(dump))
	}
}

func (c client) dumpResponse(r *http.Response) {
	if r == nil {
		log.Print("dumpResponse ok: <nil>")
		return
	}
	dump, err := httputil.DumpResponse(r, true)
	if err != nil {
		log.Print("dumpResponse err:", err)
	} else {
		log.Print("dumpResponse ok:", string(dump))
	}
}

// doTimeoutRequest do a HTTP request with timeout
func (c *client) doTimeoutRequest(timer *time.Timer, req *http.Request) (*http.Response, error) {
	// Do the request in the background so we can check the timeout
	type result struct {
		resp *http.Response
		err  error
	}
	done := make(chan result, 1)
	go func() {
		if c.debug {
			c.dumpRequest(req)
		}
		resp, err := c.httpClient.Do(req)
		if c.debug {
			c.dumpResponse(resp)
		}
		done <- result{resp, err}
	}()
	// Wait for the read or the timeout
	select {
	case r := <-done:
		return r.resp, r.err
	case <-timer.C:
		return nil, errors.New("timeout on reading data from Spiral API")
	}
}

// do prepare and process HTTP request to Spiral API
func (c *client) do(method string, resource string, params map[string]string, authNeeded bool) (response []byte, err error) {
	connectTimer := time.NewTimer(c.httpTimeout)

	var rawurl string
	if strings.HasPrefix(resource, "http") {
		rawurl = resource
	} else {
		rawurl = fmt.Sprintf("%s/%s", API_BASE, resource)
	}
	var payload string
	if method == "GET" {
		var URL *url.URL
		URL, err = url.Parse(rawurl)
		if err != nil {
			return
		}
		q := URL.Query()
		for key, value := range params {
			q.Set(key, value)
		}
		payload = q.Encode()
		URL.RawQuery = payload
		rawurl = URL.String()
	} else {
		bs, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		payload = string(bs)
	}
	req, err := http.NewRequest(method, rawurl, strings.NewReader(payload))
	if err != nil {
		return
	}
	req.Header.Add("Accept", "application/json")

	// Auth
	if authNeeded {
		if len(c.apiKey) == 0 || len(c.apiSecret) == 0 {
			err = errors.New("you need to set API Key and API Secret to call this method")
			return
		}
		req.Header.Set("api-key", c.apiKey)

		expired := fmt.Sprint(time.Now().Add(5 * time.Second).Unix())
		req.Header.Set("api-expires", expired)

		switch method {
		case "POST":
			sign := c.signature(method, resource, map[string]string{}, expired, payload)
			req.Header.Set("api-signature", sign)
		default:
			sign := c.signature(method, resource, params, expired, "")
			req.Header.Set("api-signature", sign)
		}
	}

	resp, err := c.doTimeoutRequest(connectTimer, req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 401 {
		return response, errors.New(resp.Status)
	}
	return response, err
}

func (c *client) signature(verb, path string, params map[string]string, expired, body string) string {
	ul, err := url.Parse(path)
	if err != nil {
		return err.Error()
	}
	var val = url.Values{}
	for k, v := range params {
		val.Set(k, v)
	}
	ul.RawQuery = val.Encode()

	txt := fmt.Sprintf("%v%v%v%v", verb, ul.String(), expired, body)
	return computeHmac256(txt, c.apiSecret)
}

func computeHmac256(strMessage string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(strMessage))

	return hex.EncodeToString(h.Sum(nil))
}
