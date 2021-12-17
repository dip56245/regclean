package hub

import (
	"context"
	"fmt"
	"net/http"
)

func (h *Hub) url(template string, args ...interface{}) string {
	suffix := fmt.Sprintf(template, args...)
	return fmt.Sprintf("%s%s", h.Path, suffix)
}

func (h *Hub) methodGet(ctx context.Context, url string, args ...interface{}) (*http.Response, error) {
	reqURL := h.url(url, args...)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	return h.client.Do(req)
}
