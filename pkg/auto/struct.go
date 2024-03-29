package auto

import (
	"github.com/dip56245/regclean/pkg/hub"
	datetimebranch "github.com/dip56245/regclean/pkg/modules/datetime-branch"
	"github.com/dip56245/regclean/pkg/modules/old"
)

type ModuleWorker interface {
	Setup(repo string, setup map[string]string) error
	Tag(tag *hub.TagItem)
	Clean(hub *hub.Hub, realDelete bool)
	PrintSkipped()
}

func GetModule(name string) ModuleWorker {
	switch name {
	case "old":
		return &old.ModuleOld{}
	case "datetime-branch":
		return &datetimebranch.ModuleDateTimeBranch{}
	default:
		return nil
	}
}
