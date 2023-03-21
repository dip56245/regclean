package auto

import (
	"fmt"
	"log"
	"os"

	"github.com/dip56245/regclean/pkg/hub"
	"gopkg.in/yaml.v2"
)

type Worker struct {
	Hub   *hub.Hub
	Tasks []Item
}

type Item struct {
	Name   string            `yaml:"name"`
	Type   string            `yaml:"type"`
	Delete bool              `yaml:"delete"`
	Setup  map[string]string `yaml:"setup"`
	Repos  []string          `yaml:"repos"`
}

func New(hub *hub.Hub, fileName string) (*Worker, error) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	tasks := make([]Item, 0)
	err = yaml.Unmarshal(file, &tasks)
	if err != nil {
		return nil, err
	}
	return &Worker{
			Hub:   hub,
			Tasks: tasks,
		},
		nil
}

func (w *Worker) Work() error {
	for i := 0; i < len(w.Tasks); i++ {
		w.oneWorker(&w.Tasks[i])
	}
	return nil
}

func (w *Worker) oneWorker(task *Item) {
	fmt.Printf("+ %s (%s)\n", task.Name, task.Type)
	for i := 0; i < len(task.Repos); i++ {
		module := GetModule(task.Type)
		if module == nil {
			log.Printf("not found module type: %s - SKIP\n", task.Type)
			break
		}
		if err := module.Setup(task.Repos[i], task.Setup); err != nil {
			log.Printf("error setup module type: %s - %s\n", task.Type, err)
			break
		}
		tags, err := w.Hub.GetTags(task.Repos[i])
		if err != nil {
			log.Printf("Error: %s\n", err)
		}
		for i := 0; i < len(tags); i++ {
			module.Tag(tags[i])
		}
		module.Clean(w.Hub, task.Delete)
	}
}
