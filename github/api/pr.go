package api

import (
	"context"
	"fmt"
	"net/http"
)

type PR struct {
	URL    string `json:"html_url"`
	Number int    `json:"number"`

	Title string `json:"title"`
	Body  string `json:"body"`
	State string `json:"state"`
	Draft bool   `json:"draft"`

	CreatedAt string `json:"created_at"`
}

func (api *API) MakePR(ctx context.Context, repo string, dt any) (*PR, error) {
	var pr PR
	err := api.makeRequest(ctx, &pr,
		withMethod(http.MethodPost),
		withBodyData(dt),
		withPath(fmt.Sprintf("/repos/%s/pulls", repo)),
	)
	if err != nil {
		return nil, err
	}
	return &pr, err
}
