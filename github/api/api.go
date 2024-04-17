package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type API struct {
	Token func(context.Context) (string, error)
}

type requestOpts struct {
	method string
	path   string
	body   io.Reader
}

type RequestOpt func(*requestOpts)

func withMethod(method string) RequestOpt {
	return func(ro *requestOpts) {
		ro.method = method
	}
}

func withPath(path string) RequestOpt {
	return func(ro *requestOpts) {
		ro.path = strings.TrimPrefix(path, "/")
	}
}

func withBody(body io.Reader) RequestOpt {
	return func(ro *requestOpts) {
		ro.body = body
	}
}

func withBodyData(data any) RequestOpt {
	dt, _ := json.Marshal(data)
	fmt.Println(string(dt))
	return withBody(bytes.NewReader(dt))
}

func (api *API) makeRequest(ctx context.Context, dest any, opts ...RequestOpt) error {
	opt := requestOpts{
		method: http.MethodGet,
	}
	for _, o := range opts {
		o(&opt)
	}

	target := fmt.Sprintf("https://api.github.com/%s", opt.path)
	fmt.Println(target)
	req, err := http.NewRequestWithContext(ctx, opt.method, target, opt.body)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	if api.Token != nil {
		token, err := api.Token(ctx)
		if err != nil {
			return err
		}
		req.Header.Add("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// needs a 2XX status code to be successful
		msg := fmt.Sprintf("bad response %s", resp.Status)
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			msg += "\n" + string(body)
		}
		return fmt.Errorf(msg)
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(dest); err != nil {
		return err
	}

	return nil
}
