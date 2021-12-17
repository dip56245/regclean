package hub

import (
	"context"
	"encoding/json"
	"time"
)

type blobAnswer struct {
	Architecture  string    `json:"architecture"`
	Author        string    `json:"author"`
	Container     string    `json:"container"`
	Created       time.Time `json:"created"`
	DockerVersion string    `json:"docker_version"`
	Os            string    `json:"os"`
}

func (h *Hub) getCreateTime(repo string, configDigest string) (time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), h.Config.RequestTimeOut)
	defer cancel()
	resp, err := h.methodGet(ctx, "/v2/%s/blobs/%s", repo, configDigest)
	if err != nil {
		return time.Now(), err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var blob blobAnswer
	if err = decoder.Decode(&blob); err != nil {
		return time.Now(), err
	}
	return blob.Created, nil
}
