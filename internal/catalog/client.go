package catalog

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client interface {
	Exists(articleID string, authHeader string) (bool, error)
}

type client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string) Client {
	return &client{
		baseURL: strings.TrimRight(baseURL, "/"),
		http:    &http.Client{Timeout: 3 * time.Second},
	}
}

func (c *client) Exists(articleID string, authHeader string) (bool, error) {
	if articleID == "" {
		return false, nil
	}

	url := fmt.Sprintf("%s/articles/%s", c.baseURL, articleID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	// Forward auth tal cual viene ("Bearer ...")
	if strings.TrimSpace(authHeader) != "" {
		req.Header.Set("Authorization", authHeader)
	}

	res, err := c.http.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return false, fmt.Errorf("catalog auth error: status %d", res.StatusCode)
	default:
		return false, fmt.Errorf("catalog returned status %d", res.StatusCode)
	}
}
