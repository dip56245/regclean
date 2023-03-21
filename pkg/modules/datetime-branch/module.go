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
	repo         string
	skipCount    int
	regexp       *regexp.Regexp
	list         map[string][]*internalTagItem
	skipedByName []string
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
	m.skipedByName = make([]string, 0)
	fmt.Printf("  + %s\n", repo)
	return nil
}

func (m *ModuleDateTimeBranch) Tag(tag *hub.TagItem) {
	if !m.regexp.Match([]byte(tag.Name)) {
		m.skipedByName = append(m.skipedByName, tag.Name)
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
	if _, exist := m.list[branchName]; !exist {
		m.list[branchName] = []*internalTagItem{internalItem}
	} else {
		m.list[branchName] = append(m.list[branchName], internalItem)
	}
}

func (m *ModuleDateTimeBranch) Clean(hub *hub.Hub, realDelete bool) {
	if len(m.list) == 0 {
		return
	}
	fmt.Printf("    found %d branches:\n", len(m.list))
	for branchName := range m.list {
		fmt.Printf("      %s - %d\n", branchName, len(m.list[branchName]))
		sort.Slice(m.list[branchName], func(i, j int) bool {
			return m.list[branchName][j].datetime < m.list[branchName][i].datetime
		})
		if (len(m.list[branchName]) > m.skipCount) && (realDelete) {
			for i := m.skipCount; i < len(m.list[branchName]); i++ {
				fmt.Printf("        %s - ", m.list[branchName][i].tag.Name)
				err := hub.DeleteManifest(m.repo, m.list[branchName][i].tag.Digest)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("ok")
				}
			}
		}
	}
}

func (m *ModuleDateTimeBranch) PrintSkipped() {
	if len(m.skipedByName) == 0 {
		return
	}
	fmt.Println("    skiped tags (unsupport format):")
	for i := 0; i < len(m.skipedByName); i++ {
		fmt.Printf("      %s\n", m.skipedByName[i])
	}
}
