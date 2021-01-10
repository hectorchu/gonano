// Package rpc enables RPC interfacing with a Nano node
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
}

func (c *Client) send(ctx context.Context, body interface{}) ([]byte, error) {
	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.URL, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	buf.Reset()

	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return nil, err
	}

	_ = resp.Body.Close()

	var v struct{ Error, Message string }

	if err := json.Unmarshal(buf.Bytes(), &v); err != nil {
		return nil, err
	}

	if v.Error != "" {
		err = errors.New(v.Error)
	} else if v.Message != "" {
		err = errors.New(v.Message)
	}

	return buf.Bytes(), err
}
