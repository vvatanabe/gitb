package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"os"

	"strings"

	"path"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
)

const (
	Name    = "gitb"
	version = "0.0.0"
)

var (
	commit string
	date   string
)

func FmtVersion() string {
	if commit == "" || date == "" {
		return version
	}
	return fmt.Sprintf("%s, build %s, date %s", version, commit, date)
}

func main() {
	NewCLI().Run(os.Args)
}

func NewCLI() *CLI {
	app := cli.NewApp()
	app.Name = Name
	app.Usage = "A command line tool for using Backlog's git comfortably."
	app.UsageText = Name + " [sub] [options]"
	app.Version = FmtVersion()
	app.Commands = []cli.Command{
		{
			Name:   "commit",
			Usage:  "open commits page",
			Action: commitCommand,
			Flags: []cli.Flag{
				cli.IntFlag{
					EnvVar: "CONCURRENT_LIMIT",
					Name:   "concurrent-limit",
					Usage:  "size of concurrent limit",
					Value:  4,
				},
				cli.IntFlag{
					EnvVar: "TIMEOUT_EACH_TASK",
					Name:   "timeout-each-task",
					Usage:  "timeout for each task(sec)",
					Value:  180,
				},
			},
		},
	}
	return &CLI{app}
}

type CLI struct {
	app *cli.App
}

func (cli *CLI) Run(argv []string) error {
	return cli.app.Run(argv)
}

func commitCommand(c *cli.Context) error {

	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		// TODO
		return err
	}

	remo, err := repo.Remote("origin")
	if err != nil {
		// TODO
		return err
	}

	cfg := remo.Config()
	if len(cfg.URLs) == 0 {
		return errors.New("could not find remote URL")
	}
	u := cfg.URLs[0]

	ep, err := transport.NewEndpoint(u)
	if err != nil {
		return err
	}

	spaceKey := strings.Split(ep.Host, ".")[0]

	var domain string
	if ep.Protocol == "ssh" {
		domain = strings.Replace(ep.Host, spaceKey+".git.", "", 1)
	} else {
		domain = strings.Replace(ep.Host, spaceKey+".", "", 1)
	}

	head, err := repo.Head()
	if err != nil {
		return err
	}

	head.Type()

	fmt.Println(head)

	base := fmt.Sprintf("https://%s.%s/git", spaceKey, domain)
	gitURL := path.Join(base, strings.TrimSuffix(ep.Path, ".git"), "history", "master")
	openBrowser(gitURL)
	fmt.Println(gitURL)
	return nil
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
