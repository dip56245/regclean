package auto

import (
	"github.com/dip56245/regclean/pkg/hub"
	"github.com/dip56245/regclean/pkg/modules/old"
)

type ModuleWorker interface {
	Setup(repo string, setup map[string]string) error
	Tag(tag *hub.TagItem)
	Clean(hub *hub.Hub)
}

func GetModule(name string) ModuleWorker {
	switch name {
	case "old":
		return &old.ModuleOld{}
	default:
		return nil
	}
}
