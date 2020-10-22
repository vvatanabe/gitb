package main

import (
	"fmt"
	"path"
)

func NewBacklogURLBuilder(domain, spaceKey string) *BacklogURLBuilder {
	return &BacklogURLBuilder{
		domain:   domain,
		spaceKey: spaceKey,
	}
}

type BacklogURLBuilder struct {
	domain     string
	spaceKey   string
	projectKey string
	repoName   string
}

func (b *BacklogURLBuilder) SetProjectKey(key string) *BacklogURLBuilder {
	b.projectKey = key
	return b
}

func (b *BacklogURLBuilder) SetRepoName(name string) *BacklogURLBuilder {
	b.repoName = name
	return b
}

func (b *BacklogURLBuilder) Host() string {
	return fmt.Sprintf("%s.%s", b.spaceKey, b.domain)
}

func (b *BacklogURLBuilder) BaseURL() string {
	return "https://" + b.Host()
}

func (b *BacklogURLBuilder) GitBaseURL() string {
	return b.BaseURL() + path.Join("/", "git", b.projectKey)
}

func (b *BacklogURLBuilder) GitRepoBaseURL() string {
	return b.GitBaseURL() + path.Join("/", b.repoName)
}

func (b *BacklogURLBuilder) ObjectURL(refOrHash string, relPath string, isDirectory bool, line string) string {
	base := "blob"
	if isDirectory {
		base = "tree"
	}
	hash := ""
	if line != "" {
		hash = "#" + line
	}
	return b.GitRepoBaseURL() + path.Join("/", base, refOrHash, relPath) + hash
}

func (b *BacklogURLBuilder) TreeURL(refOrHash string) string {
	return b.GitRepoBaseURL() + path.Join("/", "tree", refOrHash)
}

func (b *BacklogURLBuilder) HistoryURL(refOrHash string) string {
	return b.GitRepoBaseURL() + path.Join("/", "history", refOrHash)
}

func (b *BacklogURLBuilder) NetworkURL(refOrHash string) string {
	return b.GitRepoBaseURL() + path.Join("/", "network", refOrHash)
}

func (b *BacklogURLBuilder) BranchListURL() string {
	return b.GitRepoBaseURL() + path.Join("/", "branches")
}

func (b *BacklogURLBuilder) TagListURL() string {
	return b.GitRepoBaseURL() + path.Join("/", "tags")
}

func (b *BacklogURLBuilder) PullRequestListURL(statusID int) string {
	var q string
	if statusID > 0 {
		q += fmt.Sprintf("?q.statusId=%d", statusID)
	}
	return b.GitRepoBaseURL() + path.Join("/", "pullRequests") + q
}

func (b *BacklogURLBuilder) PullRequestURL(id string) string {
	return b.GitRepoBaseURL() + path.Join("/", "pullRequests", id)
}

func (b *BacklogURLBuilder) AddPullRequestURL(base, topic string) string {
	s := fmt.Sprintf("%s...%s", base, topic)
	return b.GitRepoBaseURL() + path.Join("/", "pullRequests", "add", s)
}

func (b *BacklogURLBuilder) IssueListURL(statusIDs []int) string {
	q := "?condition.simpleSearch=true"
	for _, v := range statusIDs {
		q += fmt.Sprintf("&condition.statusId=%d", v)
	}
	return b.BaseURL() + path.Join("/", "find", b.projectKey) + q
}

func (b *BacklogURLBuilder) IssueURL(issueKey string) string {
	return b.BaseURL() + path.Join("/", "view", issueKey)
}

func (b *BacklogURLBuilder) AddIssueURL() string {
	return b.BaseURL() + path.Join("/", "add", b.projectKey)
}
