package main

import (
	"os"

	"github.com/urfave/cli"
)

const help = `

`

func main() {
	app := cli.NewApp()
	app.Name = name
	app.Usage = usage
	app.UsageText = usageText
	app.Version = FmtVersion()
	app.Before = func(c *cli.Context) error {
		b, err := NewBacklogRepository(".")
		if err != nil {
			return err
		}
		SetBacklogRepositoryToContext(c, b)
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:  "pr",
			Usage: "Open the pull request list page in current repository",
			Action: func(c *cli.Context) error {
				return exit(GetBacklogRepositoryFromContext(c).OpenPullRequestList())
			},
			Subcommands: []cli.Command{
				{
					Name:  "show",
					Usage: "Open the pull request page related to current branch",
					Action: func(c *cli.Context) error {
						return exit(GetBacklogRepositoryFromContext(c).OpenPullRequest())
					},
				},
				{
					Name:  "create",
					Usage: "Open the page to create pull request with current branch",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "b, base",
							Usage: "BASE is base branch name. Default is empty",
						},
					},
					Action: func(c *cli.Context) error {
						base := c.Args().Get(0)
						return exit(GetBacklogRepositoryFromContext(c).OpenAddPullRequest(base, ""))
					},
				},
			},
		},
		{
			Name:  "issue",
			Usage: "Open the issue page related to current branch",
			Action: func(c *cli.Context) error {
				return exit(GetBacklogRepositoryFromContext(c).OpenIssue())
			},
			Subcommands: []cli.Command{},
		},
		{
			Name:  "add-issue",
			Usage: "Open the page to create issue in current repository's project",
			Action: func(c *cli.Context) error {
				return GetBacklogRepositoryFromContext(c).OpenAddIssue()
			},
		},
		{
			Name:  "ls-repo",
			Usage: "Open the repository list page of current repository's project",
			Action: func(c *cli.Context) error {
				return exit(GetBacklogRepositoryFromContext(c).OpenRepositoryList())
			},
		},
		{
			Name:  "ls-branch",
			Usage: "Open the branch list page of current repository",
			Action: func(c *cli.Context) error {
				return exit(GetBacklogRepositoryFromContext(c).OpenBranchList())
			},
		},
		{
			Name:  "ls-tag",
			Usage: "Open the tag list page of current repository",
			Action: func(c *cli.Context) error {
				return exit(GetBacklogRepositoryFromContext(c).OpenTagList())
			},
		},
		{
			Name:  "tree",
			Usage: "Open the tree page of current branch",
			Action: func(c *cli.Context) error {
				return exit(GetBacklogRepositoryFromContext(c).OpenTree(""))
			},
		},
		{
			Name:  "log",
			Usage: "Open the commit log page of current branch",
			Action: func(c *cli.Context) error {
				return exit(GetBacklogRepositoryFromContext(c).OpenHistory(""))
			},
		},
	}
	app.Run(os.Args)
}

func exit(err error) error {
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}
