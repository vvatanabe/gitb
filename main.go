package main

import (
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/urfave/cli"
)

func help() string {
	var sb strings.Builder
	out, err := exec.Command("git", "help").Output()
	if err != nil {
		log.Fatalln(err)
	}
	sb.Write(out)
	sb.WriteString(`
These Backlog's git commands are provided by gitb:
     pr       Open the pull request list page in current repository
     issue    Open the issue list page in current project
     browse   Open other git page (e.g. branch, tree, tag, and more...) in current repository
     help, h  Shows a list of commands or help for one command

`)
	return sb.String()
}

func main() {
	app := cli.NewApp()
	app.Name = name
	app.Usage = usage
	app.UsageText = usageText
	app.CustomAppHelpTemplate = help()
	app.Version = FmtVersion()
	app.Commands = []cli.Command{
		{
			Name:  "pr",
			Usage: "Open the pull request list page in current repository",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "s, state",
					Value: "open",
				},
			},
			Action: func(c *cli.Context) error {
				s := c.String("state")
				repo, err := open(".")
				if err != nil {
					return exit(err)
				}
				return exit(repo.OpenPullRequestList(s))
			},
			Subcommands: []cli.Command{
				{
					Name:  "show",
					Usage: "Open the pull request page. When no specify <PR-ID>, open the PR page related to the current branch",
					Action: func(c *cli.Context) error {
						repo, err := open(".")
						if err != nil {
							return exit(err)
						}
						if c.Args().Present() {
							return exit(repo.OpenPullRequestByID(c.Args().First()))
						}
						return exit(repo.OpenPullRequest())
					},
				},
				{
					Name:  "add",
					Usage: "Open the page to add pull request with current branch",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "b, base",
						},
					},
					Action: func(c *cli.Context) error {
						base := c.String("base")
						repo, err := open(".")
						if err != nil {
							return exit(err)
						}
						return exit(repo.OpenAddPullRequest(base, ""))
					},
				},
				{
					Name:            "blame",
					Usage:           "Show pull request id with git blame",
					SkipFlagParsing: true,
					Action: func(c *cli.Context) error {
						repo, err := open(".")
						if err != nil {
							return exit(err)
						}
						if c.Args().Present() {
							return exit(repo.BlamePR(c.Args()))
						}
						return exit(nil)
					},
				},
			},
		},
		{
			Name:  "issue",
			Usage: "Open the issue list page in current project",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "s, state",
					Value: "not_closed",
				},
			},
			Action: func(c *cli.Context) error {
				s := c.String("state")
				repo, err := open(".")
				if err != nil {
					return exit(err)
				}
				return exit(repo.OpenIssueList(s))
			},
			Subcommands: []cli.Command{
				{
					Name:  "show",
					Usage: "Open the issue page related to current branch",
					Action: func(c *cli.Context) error {
						repo, err := open(".")
						if err != nil {
							return exit(err)
						}
						return exit(repo.OpenIssue())
					},
				},
				{
					Name:  "add",
					Usage: "Open the page to add issue in current repository's project",
					Action: func(c *cli.Context) error {
						repo, err := open(".")
						if err != nil {
							return exit(err)
						}
						return exit(repo.OpenAddIssue())
					},
				},
			},
		},
		{
			Name:  "browse",
			Usage: "Open other git page (e.g. branch, tree, tag, and more...) in current repository",
			Subcommands: []cli.Command{
				{
					Name:  "branch",
					Usage: "Open the branch list page in current repository",
					Action: func(c *cli.Context) error {
						repo, err := open(".")
						if err != nil {
							return exit(err)
						}
						return exit(repo.OpenBranchList())
					},
				},
				{
					Name:  "tag",
					Usage: "Open the tag list page in the current repository.",
					Action: func(c *cli.Context) error {
						repo, err := open(".")
						if err != nil {
							return exit(err)
						}
						return exit(repo.OpenTagList())
					},
				},
				{
					Name:  "tree",
					Usage: "Open the tree page in current branch",
					Action: func(c *cli.Context) error {
						repo, err := open(".")
						if err != nil {
							return exit(err)
						}
						return exit(repo.OpenTree(""))
					},
				},
				{
					Name:  "history",
					Usage: "Open the history page in current branch",
					Action: func(c *cli.Context) error {
						repo, err := open(".")
						if err != nil {
							return exit(err)
						}
						return exit(repo.OpenHistory(""))
					},
				},
				{
					Name:  "network",
					Usage: "Open the network page in current branch",
					Action: func(c *cli.Context) error {
						repo, err := open(".")
						if err != nil {
							return exit(err)
						}
						return exit(repo.OpenNetwork(""))
					},
				},
				{
					Name:  "repo",
					Usage: "Open the repository list page in current project",
					Action: func(c *cli.Context) error {
						repo, err := open(".")
						if err != nil {
							return exit(err)
						}
						return exit(repo.OpenRepositoryList())
					},
				},
				{
					Name:  "show",
					Usage: "Open the corresponding page to given file or directory in current project",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "l, line",
						},
					},
					Action: func(c *cli.Context) error {
						repo, err := open(".")
						if err != nil {
							return exit(err)
						}
						filePath := c.Args().First()
						if !path.IsAbs(filePath) {
							wd, err := os.Getwd()
							if err != nil {
								return exit(err)
							}
							filePath = path.Join(wd, filePath)
						}
						fileUrl, err := url.Parse(filePath)
						if err != nil {
							return exit(err)
						}
						fs, err := os.Stat(fileUrl.Path)
						if err != nil {
							return exit(err)
						}
						line := c.String("line")
						if line == "" {
							line = fileUrl.Fragment
						}
						return exit(repo.OpenObject(fileUrl.Path, fs.IsDir(), line))
					},
				},
				{
					Name:  "commit",
					Usage: "Open the network page in current branch",
					Action: func(c *cli.Context) error {
						repo, err := open(".")
						if err != nil {
							return exit(err)
						}
						hash := c.Args().First()
						if len(hash) == 0 {
							return exit(err)
						}
						return exit(repo.OpenCommit(hash))
					},
				},
			},
		},
	}
	app.OnUsageError = func(context *cli.Context, err error, isSubcommand bool) error {
		if isSubcommand {
			return err
		}
		if err := NewGitCmd(os.Args[1:]).Run(); err != nil {
			return err
		}
		return nil
	}
	app.CommandNotFound = func(c *cli.Context, name string) {
		args := append([]string{name}, os.Args[2:]...)
		if err := NewGitCmd(args).Run(); err != nil {
			log.Fatalln(err)
		}
	}
	_ = app.Run(os.Args)
}

func exit(err error) error {
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

func open(path string) (*BacklogRepository, error) {
	repo, err := OpenRepository(path)
	if err != nil {
		return nil, err
	}
	return NewBacklogRepository(repo), nil
}
