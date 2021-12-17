package old

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dip56245/regclean/pkg/hub"
)

type ModuleOld struct {
	repo       string
	timeItem   time.Time
	needDelete []*hub.TagItem
}

func (m *ModuleOld) Setup(repo string, setup map[string]string) error {
	m.repo = repo
	days, err := strconv.Atoi(setup["days"])
	if err != nil {
		return fmt.Errorf("error parse setup.days: %w", err)
	}
	m.timeItem = time.Now().AddDate(0, 0, -days)
	log.Printf("Try found tags before %s\n", m.timeItem)
	m.needDelete = make([]*hub.TagItem, 0)
	return nil
}

func (m *ModuleOld) Tag(tag *hub.TagItem) {
	if tag.CreateTime.Before(m.timeItem) {
		m.needDelete = append(m.needDelete, tag)
	} else {
		fmt.Printf("%s (%s) - skip\n", tag.Name, tag.CreateTime)
	}
}

func (m *ModuleOld) Clean(hub *hub.Hub) {
	for i := 0; i < len(m.needDelete); i++ {
		err := hub.DeleteManifest(m.repo, m.needDelete[i].Digest)
		if err != nil {
			fmt.Printf("delete %s/%s - %s\n", m.repo, m.needDelete[i].Name, err)
		} else {
			fmt.Printf("delete %s/%s - ok\n", m.repo, m.needDelete[i].Name)
		}
	}
}
