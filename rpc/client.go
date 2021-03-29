package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Client is used for connecting to http rpc endpoints.
type Client struct {
	URL string
	AuthHeader string
	Ctx context.Context
}

func (c *Client) send(body interface{}) (result []byte, err error) {
	var buf bytes.Buffer
	if err = json.NewEncoder(&buf).Encode(body); err != nil {
		return
	}
	if c.Ctx == nil {
		c.Ctx = context.Background()
	}
	req, err := http.NewRequestWithContext(c.Ctx, http.MethodPost, c.URL, &buf)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	if c.AuthHeader != "" {
		req.Header.Set("Authorization", c.AuthHeader)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	buf.Reset()
	if _, err = io.Copy(&buf, resp.Body); err != nil {
		return
	}
	if err = resp.Body.Close(); err != nil {
		return
	}
	var v struct{ Error, Message string }
	if err = json.Unmarshal(buf.Bytes(), &v); err != nil {
		return
	}
	if v.Error != "" {
		err = errors.New(v.Error)
	} else if v.Message != "" {
		err = errors.New(v.Message)
	}
	return buf.Bytes(), err
}
