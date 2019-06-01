package main

import "github.com/urfave/cli"

const contextKeyBacklogRepository = "ctx-key-vacklog-repository"

func GetBacklogRepositoryFromContext(c *cli.Context) *BacklogRepository {
	v, ok := c.App.Metadata[contextKeyBacklogRepository]
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
	c.App.Metadata[contextKeyBacklogRepository] = b
}
