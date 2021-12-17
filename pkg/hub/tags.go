package hub

import (
	"context"
	"encoding/json"
	"sort"
)

type tagsListAnswer struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type TagItem struct {
	Name   string
	Size   uint64
	Digest string
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

func (h *Hub) GetTags(repo string) ([]*TagItem, error) {
	tags, err := h.getTags(repo)
	if err != nil {
		return []*TagItem{}, err
	}
	semaphoreChan := make(chan struct{}, h.Config.ThreadCount)
	outChan := make(chan *TagItem)
	defer func() {
		defer close(semaphoreChan)
		defer close(outChan)
	}()
	for _, tname := range tags {
		go func(tname string) {
			semaphoreChan <- struct{}{}
			manifest, err := h.Manifest(repo, tname)
			if err != nil {
				outChan <- &TagItem{
					Name:   tname,
					Size:   0,
					Digest: "",
				}
			} else {
				outChan <- &TagItem{
					Name:   tname,
					Size:   manifest.Size,
					Digest: manifest.Digest,
				}
			}
			<-semaphoreChan
		}(tname)
	}
	out := make([]*TagItem, 0, len(tags))
	for {
		res := <-outChan
		out = append(out, res)
		if len(tags) == len(out) {
			break
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out, nil
}

func (h *Hub) DeleteTag(repo string, tag string) error {
	manifest, err := h.Manifest(repo, tag)
	if err != nil {
		return err
	}
	if manifest.Digest == "" {
		return nil
	}
	h.DeleteManifest(repo, manifest.Digest)
	return nil
}
