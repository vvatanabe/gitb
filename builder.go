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

func (b *BacklogURLBuilder) TreeURL(rev string) string {
	return b.GitRepoBaseURL() + path.Join("/", "tree", rev)
}

func (b *BacklogURLBuilder) HistoryURL(rev string) string {
	return b.GitRepoBaseURL() + path.Join("/", "history", rev)
}

func (b *BacklogURLBuilder) BranchListURL() string {
	return b.GitRepoBaseURL() + path.Join("/", "branches")
}

func (b *BacklogURLBuilder) TagListURL() string {
	return b.GitRepoBaseURL() + path.Join("/", "tags")
}

func (b *BacklogURLBuilder) PullRequestListURL(statusID int) string {
	return b.GitRepoBaseURL() + path.Join("/", fmt.Sprintf("pullRequests?q.statusId=%d", statusID))
}

func (b *BacklogURLBuilder) PullRequestURL(id string) string {
	return b.GitRepoBaseURL() + path.Join("/", "pullRequests", id)
}

func (b *BacklogURLBuilder) AddPullRequestURL(base, topic string) string {
	s := fmt.Sprintf("%s...%s", base, topic)
	return b.GitRepoBaseURL() + path.Join("/", "pullRequests", "add", s)
}

func (b *BacklogURLBuilder) IssueListURL() string {
	return b.BaseURL() + path.Join("/", "find", b.projectKey)
}

func (b *BacklogURLBuilder) IssueURL(issueKey string) string {
	return b.BaseURL() + path.Join("/", "view", issueKey)
}

func (b *BacklogURLBuilder) AddIssueURL() string {
	return b.BaseURL() + path.Join("/", "add", b.projectKey)
}
