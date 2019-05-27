package main

import "github.com/urfave/cli"

const contextKeyGitURLBuilder = "ctx-key-git-url-builder"

func GetBacklogRepositoryFromContext(c *cli.Context) *BacklogRepository {
	v, ok := c.App.Metadata[contextKeyGitURLBuilder]
	if !ok {
		return nil
	}
	if repo, ok := v.(*BacklogRepository); !ok {
		return nil
	} else {
		return repo
	}
}

func SetBacklogRepositoryToContext(c *cli.Context, b *BacklogRepository) {
	c.App.Metadata[contextKeyGitURLBuilder] = b
}
