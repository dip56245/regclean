package hub

import (
	"context"
	"net/http"
	"time"
)

type Hub struct {
	Config Config
	Path   string
	client *http.Client
}

type Config struct {
	DryRun         bool
	MaxReposCount  int
	ThreadCount    int
	RequestTimeOut time.Duration
	Login          string
	Password       string
}

func New(path string, login string, password string) *Hub {
	return &Hub{
		Path: path,
		Config: Config{
			DryRun:         true,
			MaxReposCount:  1000,
			RequestTimeOut: time.Second * 3,
			ThreadCount:    25,
			Login:          login,
			Password:       password,
		},
		client: &http.Client{
			Transport: http.DefaultTransport,
		},
	}
}

func (h *Hub) Ping() bool {
	ctx, cancel := context.WithTimeout(context.Background(), h.Config.RequestTimeOut)
	defer cancel()
	resp, err := h.methodGet(ctx, "/v2/")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}
