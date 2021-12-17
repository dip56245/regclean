package hub

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"
)

type tagsListAnswer struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type TagItem struct {
	Name         string
	Size         uint64
	Digest       string
	ConfigDigest string
	CreateTime   time.Time
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
			item := TagItem{Name: tname}
			if err == nil {
				item.Size = manifest.Size
				item.Digest = manifest.Digest
				item.ConfigDigest = manifest.ConfigDigest
			}
			createTime, err := h.getCreateTime(repo, item.ConfigDigest)
			if err == nil {
				item.CreateTime = createTime
			}
			outChan <- &item
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
		return out[i].CreateTime.Before(out[j].CreateTime)
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
	return h.DeleteManifest(repo, manifest.Digest)
}

func (h *Hub) DeleteAllTags(repo string) error {
	tags, err := h.GetTags(repo)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, tag := range tags {
		wg.Add(1)
		go func(wg *sync.WaitGroup, repo string, tag *TagItem) {
			status := "ok"
			err := h.DeleteManifest(repo, tag.Digest)
			if err != nil {
				status = err.Error()
			}
			fmt.Printf("delete %s - %s\n", tag.Name, status)
			wg.Done()
		}(&wg, repo, tag)
	}
	wg.Wait()
	return nil
}
