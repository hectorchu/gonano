package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Client is used for connecting to http rpc endpoints.
type Client struct {
	URL string
}

func (c *Client) send(body interface{}) (result []byte, err error) {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)
	resp, err := http.Post(c.URL, "application/json", &buf)
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
