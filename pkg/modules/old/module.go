package old

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dip56245/regclean/pkg/hub"
)

const dtFormat string = "2006/01/02 15:04"

type ModuleOld struct {
	repo       string
	timeItem   time.Time
	needDelete []*hub.TagItem
	needSkip   []*hub.TagItem
}

func (m *ModuleOld) Setup(repo string, setup map[string]string) error {
	m.repo = repo
	days, err := strconv.Atoi(setup["days"])
	if err != nil {
		return fmt.Errorf("error parse setup.days: %w", err)
	}
	m.timeItem = time.Now().AddDate(0, 0, -days)
	fmt.Printf("  + %s older %s\n", m.repo, m.timeItem.Format(dtFormat))
	m.needDelete = make([]*hub.TagItem, 0)
	m.needSkip = make([]*hub.TagItem, 0)
	return nil
}

func (m *ModuleOld) Tag(tag *hub.TagItem) {
	if tag.CreateTime.Before(m.timeItem) {
		m.needDelete = append(m.needDelete, tag)
	} else {
		m.needSkip = append(m.needSkip, tag)
		// fmt.Printf("%s (%s) - skip\n", tag.Name, tag.CreateTime.Format(dtFormat))
	}
}

func (m *ModuleOld) Clean(hub *hub.Hub, realDelete bool) {
	if len(m.needDelete) == 0 {
		return
	}
	fmt.Printf("    deleted tags:\n")
	for i := 0; i < len(m.needDelete); i++ {
		fmt.Printf("      %s [%s] - ", m.needDelete[i].Name, m.needDelete[i].CreateTime.Format(dtFormat))
		if realDelete {
			err := hub.DeleteManifest(m.repo, m.needDelete[i].Digest)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("ok")
			}
		} else {
			fmt.Println("skip")
		}
	}
}

func (m *ModuleOld) PrintSkipped() {
	if len(m.needSkip) == 0 {
		return
	}
	fmt.Printf("    leaved tags:\n")
	for i := 0; i < len(m.needSkip); i++ {
		fmt.Printf("      %s [%s]\n", m.needSkip[i].Name, m.needSkip[i].CreateTime.Format(dtFormat))
	}
}
