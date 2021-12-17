package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"code.cloudfoundry.org/bytefmt"
	"github.com/dip56245/regclean/pkg/auto"
	"github.com/dip56245/regclean/pkg/hub"
	"github.com/urfave/cli/v2"
	"github.com/xeonx/timeago"
)

const (
	FlagRegistry = "registry"
	DryRun       = "dry"
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
			&cli.BoolFlag{
				Name:  DryRun,
				Value: false,
				Usage: "nothing delete, only output",
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
			{
				Name:   "rm",
				Usage:  "regclean rm <repo> <tag>",
				Action: actionRm,
			},
			{
				Name:   "rmall",
				Usage:  "regclean rmall <repo>",
				Action: actionRmAll,
			},
			{
				Name:   "auto",
				Usage:  "regclean auto <file.yaml>",
				Action: actionAuto,
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func getHub(c *cli.Context) *hub.Hub {
	hub := hub.New(c.String(FlagRegistry))
	hub.Config.DryRun = c.Bool(DryRun)
	return hub
}

func actionPing(c *cli.Context) error {
	hub := getHub(c)
	if hub.Ping() {
		log.Println("connection - ok")
	} else {
		log.Println("connection - error")
	}
	return nil
}

func actionList(c *cli.Context) error {
	h := getHub(c)
	list, err := h.ListRepos(hub.SortTypeAbc)
	for _, a := range list {
		fmt.Printf("%d\t%s\n", a.TagCount, a.Name)
	}
	return err
}

func actionCList(c *cli.Context) error {
	h := getHub(c)
	list, err := h.ListRepos(hub.SortTypeNum)
	for _, a := range list {
		fmt.Printf("%d\t%s\n", a.TagCount, a.Name)
	}
	return err
}

func actionTags(c *cli.Context) error {
	h := getHub(c)
	tags, err := h.GetTags(c.Args().First())
	for _, t := range tags {
		fmt.Printf("%s\t%s\t%s\n", bytefmt.ByteSize(t.Size), timeago.English.Format(t.CreateTime), t.Name)
	}
	return err
}

func actionRm(c *cli.Context) error {
	if c.Args().Len() != 2 {
		fmt.Printf("Usage:\n regclean rm <repo> <tag>\n")
		return nil
	}
	h := getHub(c)
	return h.DeleteTag(c.Args().Get(0), c.Args().Get(1))
}

func actionRmAll(c *cli.Context) error {
	if c.Args().Len() != 1 {
		fmt.Printf("Usage:\n regclean rmall <repo>\n")
		return nil
	}
	h := getHub(c)
	return h.DeleteAllTags(c.Args().First())
}

func actionAuto(c *cli.Context) error {
	if c.Args().Len() != 1 {
		fmt.Println("regclean auto <file.yaml>")
		return nil
	}
	h := getHub(c)
	worker, err := auto.New(h, c.Args().First())
	if err != nil {
		return err
	}
	return worker.Work()
}
