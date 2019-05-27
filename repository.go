package main

import (
	"sort"
	"strings"

	"regexp"

	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
)

func NewBacklogRepository(path string) (*BacklogRepository, error) {
	repo, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, err
	}

	r, err := repo.Remote("origin")
	if err != nil {
		return nil, err
	}

	cfg := r.Config()
	if len(cfg.URLs) == 0 {
		return nil, errors.New("could not find remote URL")
	}
	u := cfg.URLs[0]

	ep, err := transport.NewEndpoint(u)
	if err != nil {
		return nil, err
	}

	delimitedHost := strings.Split(ep.Host, ".")
	spaceKey := delimitedHost[0]
	domain := strings.Join(delimitedHost[len(delimitedHost)-2:], ".")

	delimitedPath := strings.Split(strings.TrimPrefix(ep.Path, "/"), "/")
	projectKey := delimitedPath[0]
	repoName := strings.TrimSuffix(delimitedPath[1], ".git")

	head, err := repo.Head()
	if err != nil {
		return nil, err
	}

	return &BacklogRepository{
		repo:       repo,
		head:       head,
		domain:     domain,
		spaceKey:   spaceKey,
		projectKey: projectKey,
		repoName:   repoName,
	}, nil
}

type BacklogRepository struct {
	repo       *git.Repository
	head       *plumbing.Reference
	domain     string
	spaceKey   string
	projectKey string
	repoName   string
}

func (b *BacklogRepository) OpenRepositoryList() error {
	return OpenBrowser(NewBacklogURLBuilder(b.domain, b.spaceKey).
		SetProjectKey(b.projectKey).
		SetRepoName(b.repoName).
		GitBaseURL())
}

func (b *BacklogRepository) OpenTree(spec string) error {
	if spec == "" {
		spec = b.head.Name().Short()
	}
	return OpenBrowser(NewBacklogURLBuilder(b.domain, b.spaceKey).
		SetProjectKey(b.projectKey).
		SetRepoName(b.repoName).
		TreeURL(spec))
}

func (b *BacklogRepository) OpenHistory(spec string) error {
	if spec == "" {
		spec = b.head.Name().Short()
	}
	return OpenBrowser(NewBacklogURLBuilder(b.domain, b.spaceKey).
		SetProjectKey(b.projectKey).
		SetRepoName(b.repoName).
		TreeURL(spec))
}

func (b *BacklogRepository) OpenBranchList() error {
	return OpenBrowser(NewBacklogURLBuilder(b.domain, b.spaceKey).
		SetProjectKey(b.projectKey).
		SetRepoName(b.repoName).
		BranchListURL())
}

func (b *BacklogRepository) OpenTagList() error {
	return OpenBrowser(NewBacklogURLBuilder(b.domain, b.spaceKey).
		SetProjectKey(b.projectKey).
		SetRepoName(b.repoName).
		TagListURL())
}

func (b *BacklogRepository) OpenPullRequestList() error {
	return OpenBrowser(NewBacklogURLBuilder(b.domain, b.spaceKey).
		SetProjectKey(b.projectKey).
		SetRepoName(b.repoName).
		PullRequestListURL(1))
}

func (b *BacklogRepository) OpenPullRequest() error {
	id, err := b.FindPullRequestIDFromRemote(b.head.Name().Short())
	if err != nil {
		return err
	}
	return OpenBrowser(NewBacklogURLBuilder(b.domain, b.spaceKey).
		SetProjectKey(b.projectKey).
		SetRepoName(b.repoName).
		PullRequestURL(id))
}

func (b *BacklogRepository) FindPullRequestIDFromRemote(branchName string) (string, error) {
	r, err := b.repo.Remote("origin")
	if err != nil {
		return "", err
	}

	rfs, err := r.List(&git.ListOptions{})
	if err != nil {
		return "", err
	}

	var target *plumbing.Reference
	for _, rf := range rfs {
		if rf.Name().Short() == branchName {
			target = rf
			break
		}
	}
	if target == nil {
		return "", errors.New("not found a current branch in remote")
	}

	var prIDs []string
	for _, rf := range rfs {
		sn := rf.Name().Short()
		if !strings.HasPrefix(sn, "pull/") || !strings.HasSuffix(sn, "/head") {
			continue
		}
		if rf.Hash() != target.Hash() {
			continue
		}
		prID := strings.TrimPrefix(sn, "pull/")
		prID = strings.TrimSuffix(prID, "/head")
		prIDs = append(prIDs, prID)
	}

	if len(prIDs) == 0 {
		return "", errors.New("not found a pull request")
	}

	sort.Sort(sort.Reverse(sort.StringSlice(prIDs)))

	return prIDs[0], nil
}

func (b *BacklogRepository) OpenAddPullRequest(base, topic string) error {
	if topic == "" {
		topic = b.head.Name().Short()
	}
	return OpenBrowser(NewBacklogURLBuilder(b.domain, b.spaceKey).
		SetProjectKey(b.projectKey).
		SetRepoName(b.repoName).
		AddPullRequestURL(base, topic))
}

func (b *BacklogRepository) OpenIssue() error {
	key := extractIssueKey(b.head.Name().Short())
	if key == "" {
		return errors.New("could not find issue key in current branch name")
	}
	return OpenBrowser(NewBacklogURLBuilder(b.domain, b.spaceKey).
		SetProjectKey(b.projectKey).
		SetRepoName(b.repoName).
		IssueURL(key))
}

func extractIssueKey(s string) string {
	matches := regexp.MustCompile("([A-Z0-9]+(?:_[A-Z0-9]+)*-[0-9]+)").FindStringSubmatch(s)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

func (b *BacklogRepository) OpenAddIssue() error {
	return OpenBrowser(NewBacklogURLBuilder(b.domain, b.spaceKey).
		SetProjectKey(b.projectKey).
		SetRepoName(b.repoName).
		AddIssueURL())
}
