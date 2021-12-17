package hub

import (
	"context"
	"encoding/json"
)

type tagsListAnswer struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

func (h *Hub) getTags(repo string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), h.Config.RequestTimeOut)
	defer cancel()
	resp, err := h.methodGet(ctx, "/v2/%s/tags/list", repo)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var tList tagsListAnswer
	if err = decoder.Decode(&tList); err != nil {
		return nil, err
	}
	return tList.Tags, nil
}
