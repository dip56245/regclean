package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"code.cloudfoundry.org/bytefmt"
	"github.com/dip56245/regclean/pkg/hub"
	"github.com/urfave/cli/v2"
)

const (
	FlagRegistry = "registry"
)

func main() {
	app := &cli.App{
		Name:     "regclean",
		Compiled: time.Now(),
		Version:  "0.0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    FlagRegistry,
				Value:   "http://localhost:5000",
				Usage:   "uri docker registry example: http://ip:5000",
				EnvVars: []string{"REGISTRY"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "ping",
				Usage:  "check connection",
				Action: actionPing,
			},
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "show list of repos",
				Action:  actionList,
			},
			{
				Name:    "clist",
				Aliases: []string{"cls"},
				Usage:   "show list of repos sort by tag count",
				Action:  actionCList,
			},
			{
				Name:   "tags",
				Usage:  "regclean tags <repo>",
				Action: actionTags,
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func actionPing(c *cli.Context) error {
	hub := hub.New(c.String(FlagRegistry))
	if hub.Ping() {
		log.Println("connection - ok")
	} else {
		log.Println("connection - error")
	}
	return nil
}

func actionList(c *cli.Context) error {
	h := hub.New(c.String(FlagRegistry))
	list, err := h.ListRepos(hub.SortTypeAbc)
	for _, a := range list {
		fmt.Printf("%d\t%s\n", a.TagCount, a.Name)
	}
	return err
}

func actionCList(c *cli.Context) error {
	h := hub.New(c.String(FlagRegistry))
	list, err := h.ListRepos(hub.SortTypeNum)
	for _, a := range list {
		fmt.Printf("%d\t%s\n", a.TagCount, a.Name)
	}
	return err
}

func actionTags(c *cli.Context) error {
	h := hub.New(c.String(FlagRegistry))
	tags, err := h.GetTags(c.Args().First())
	for _, t := range tags {
		fmt.Printf("%s\t%s\n", bytefmt.ByteSize(t.Size), t.Name)
	}
	return err
}
