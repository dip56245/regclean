package hub

import (
	"context"
	"encoding/json"
	"log"
	"sort"
)

type repositoryList struct {
	Repository []string `json:"repositories"`
}

type RepoItem struct {
	Name     string
	TagCount int
}

type SortType int

const (
	SortTypeAbc SortType = iota
	SortTypeNum
)

func (h *Hub) getListRepos() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), h.Config.RequestTimeOut)
	defer cancel()
	resp, err := h.methodGet(ctx, "/v2/_catalog?n=%d", h.Config.MaxReposCount)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var rList repositoryList
	if err = decoder.Decode(&rList); err != nil {
		return nil, err
	}
	sort.Strings(rList.Repository)
	return rList.Repository, nil
}

func (h *Hub) ListRepos(sortType SortType) ([]*RepoItem, error) {
	reposList, err := h.getListRepos()
	if err != nil {
		return nil, err
	}
	semaphoreChan := make(chan struct{}, h.Config.ThreadCount)
	resultChan := make(chan *RepoItem)
	defer func() {
		close(semaphoreChan)
		close(resultChan)
	}()

	for _, name := range reposList {
		go func(name string) {
			semaphoreChan <- struct{}{}
			tags, err := h.getTags(name)
			if err != nil {
				resultChan <- &RepoItem{Name: name, TagCount: -1}
			}
			resultChan <- &RepoItem{Name: name, TagCount: len(tags)}
			<-semaphoreChan
		}(name)
	}
	result := make([]*RepoItem, 0, len(reposList))
	for {
		res := <-resultChan
		result = append(result, res)
		if len(result) == len(reposList) {
			break
		}
	}
	switch sortType {
	case SortTypeAbc:
		sort.Slice(result, func(i, j int) bool {
			return result[i].Name < result[j].Name
		})
	case SortTypeNum:
		sort.Slice(result, func(i, j int) bool {
			return result[i].TagCount < result[j].TagCount
		})
	default:
		log.Panicf("sortType %d - not found", sortType)
	}
	return result, nil
}
