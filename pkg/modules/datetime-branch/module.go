package datetimebranch

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/dip56245/regclean/pkg/hub"
)

type internalTagItem struct {
	tag      *hub.TagItem
	branch   string
	datetime string
}

type ModuleDateTimeBranch struct {
	repo      string
	skipCount int
	regexp    *regexp.Regexp
	list      map[string][]*internalTagItem
}

func (m *ModuleDateTimeBranch) Setup(repo string, setup map[string]string) error {
	m.repo = repo
	m.list = make(map[string][]*internalTagItem)
	count, err := strconv.Atoi(setup["count"])
	if err != nil {
		return fmt.Errorf("error parse setup.days: %w", err)
	}
	m.skipCount = count
	m.regexp = regexp.MustCompile("([0-9]+)-([a-z0-9]+)")
	return nil
}

func (m *ModuleDateTimeBranch) Tag(tag *hub.TagItem) {
	if !m.regexp.Match([]byte(tag.Name)) {
		log.Printf("tag: %s:%s - skip", m.repo, tag.Name)
		return
	}
	pos := strings.Index(tag.Name, "-")
	if pos == -1 {
		log.Printf("tag: %s:%s - skip [not found -]", m.repo, tag.Name)
		return
	}
	dateTime := tag.Name[:pos]
	branchName := tag.Name[pos+1:]
	internalItem := &internalTagItem{tag: tag, branch: branchName, datetime: dateTime}
	_, exist := m.list[branchName]
	if !exist {
		m.list[branchName] = []*internalTagItem{internalItem}
	} else {
		m.list[branchName] = append(m.list[branchName], internalItem)
	}
}

func (m *ModuleDateTimeBranch) Clean(hub *hub.Hub) {
	log.Printf("found %d branches\n", len(m.list))
	for branchName := range m.list {
		sort.Slice(m.list[branchName], func(i, j int) bool {
			return m.list[branchName][j].datetime < m.list[branchName][i].datetime
		})
		if len(m.list[branchName]) <= m.skipCount {
			log.Printf("branch-tag: %s - skip count(%d)\n", branchName, len(m.list[branchName]))
		} else {
			log.Printf("branch-tag: %s - cleaning (%d)\n", branchName, len(m.list[branchName]))
			for i := m.skipCount; i < len(m.list[branchName]); i++ {
				err := hub.DeleteManifest(m.repo, m.list[branchName][i].tag.Digest)
				if err != nil {
					log.Printf("branch-tag delete %s - %s", m.list[branchName][i].tag.Name, err)
				} else {
					log.Printf("branch-tag delete %s - ok", m.list[branchName][i].tag.Name)
				}
			}
		}
	}
}
