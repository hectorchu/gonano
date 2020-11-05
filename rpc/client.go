package rpc

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Client is used for connecting to http rpc endpoints.
type Client struct {
	URL string
}

func (c *Client) send(body map[string]string) (result map[string]string, err error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.Encode(body)
	resp, err := http.Post(c.URL, "application/json", &buf)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&result)
	return
}
