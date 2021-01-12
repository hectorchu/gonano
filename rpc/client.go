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
	Ctx context.Context
}

func (c *Client) send(body interface{}) (result []byte, err error) {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)
	ctx := c.Ctx
	if ctx == nil {
		ctx = context.Background()
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.URL, &buf)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	buf.Reset()
	io.Copy(&buf, resp.Body)
	resp.Body.Close()
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
