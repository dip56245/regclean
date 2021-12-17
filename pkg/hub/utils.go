package hub

import (
	"context"
	"fmt"
	"net/http"
)

type modifyRequest func(req *http.Request)

func (h *Hub) url(template string, args ...interface{}) string {
	suffix := fmt.Sprintf(template, args...)
	return fmt.Sprintf("%s%s", h.Path, suffix)
}

func (h *Hub) methodGet(ctx context.Context, url string, args ...interface{}) (*http.Response, error) {
	reqURL := h.url(url, args...)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	return h.client.Do(req)
}

func (h *Hub) methodGetEx(ctx context.Context, modify modifyRequest, url string, args ...interface{}) (*http.Response, error) {
	reqURL := h.url(url, args...)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if modify != nil {
		modify(req)
	}
	return h.client.Do(req)
}

func (h *Hub) methodDelete(ctx context.Context, url string, args ...interface{}) (*http.Response, error) {
	reqURL := h.url(url, args...)
	req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, reqURL, nil)
	return h.client.Do(req)
}
