package main

import (
	"os/exec"
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

	remote, err := repo.Remote("origin")
	if err != nil {
		return nil, err
	}

	cfg := remote.Config()
	if len(cfg.URLs) == 0 {
		return nil, errors.New("could not find remote URL")
	}

	u := cfg.URLs[0]
	ep, err := transport.NewEndpoint(u)
	if err != nil {
		return nil, err
	}

	spaceKey, domain := extractSpaceKeyAndDomain(ep.Host)
	projectKey, repoName := extractProjectKeyAndRepoName(ep.Path)

	head, err := repo.Head()
	if err != nil {
		return nil, err
	}

	return &BacklogRepository{
		openBrowser: openBrowser,
		repo:        repo,
		head:        head,
		domain:      domain,
		spaceKey:    spaceKey,
		projectKey:  projectKey,
		repoName:    repoName,
	}, nil
}

func extractSpaceKeyAndDomain(host string) (spaceKey, domain string) {
	delimitedHost := strings.Split(host, ".")
	spaceKey = delimitedHost[0]
	domain = strings.Join(delimitedHost[len(delimitedHost)-2:], ".")
	return
}

func extractProjectKeyAndRepoName(path string) (projectKey, repoName string) {
	epPath := strings.TrimPrefix(path, "/git")
	delimitedPath := strings.Split(epPath, "/")
	projectKey = delimitedPath[1]
	repoName = strings.TrimSuffix(delimitedPath[2], ".git")
	return
}

type BacklogRepository struct {
	openBrowser func(url string) error
	repo        *git.Repository
	head        *plumbing.Reference
	domain      string
	spaceKey    string
	projectKey  string
	repoName    string
}

func (b *BacklogRepository) OpenRepositoryList() error {
	return b.openBrowser(NewBacklogURLBuilder(b.domain, b.spaceKey).
		SetProjectKey(b.projectKey).
		SetRepoName(b.repoName).
		GitBaseURL())
}

func (b *BacklogRepository) OpenTree(refOrHash string) error {
	if refOrHash == "" {
		refOrHash = b.head.Name().Short()
	}
	return b.openBrowser(NewBacklogURLBuilder(b.domain, b.spaceKey).
		SetProjectKey(b.projectKey).
		SetRepoName(b.repoName).
		TreeURL(refOrHash))
}

func (b *BacklogRepository) OpenHistory(refOrHash string) error {
	if refOrHash == "" {
		refOrHash = b.head.Name().Short()
	}
	return b.openBrowser(NewBacklogURLBuilder(b.domain, b.spaceKey).
		SetProjectKey(b.projectKey).
		SetRepoName(b.repoName).
		HistoryURL(refOrHash))
}

func (b *BacklogRepository) OpenNetwork(refOrHash string) error {
	if refOrHash == "" {
		refOrHash = b.head.Name().Short()
	}
	return b.openBrowser(NewBacklogURLBuilder(b.domain, b.spaceKey).
		SetProjectKey(b.projectKey).
		SetRepoName(b.repoName).
		NetworkURL(refOrHash))
}

func (b *BacklogRepository) OpenBranchList() error {
	return b.openBrowser(NewBacklogURLBuilder(b.domain, b.spaceKey).
		SetProjectKey(b.projectKey).
		SetRepoName(b.repoName).
		BranchListURL())
}

func (b *BacklogRepository) OpenTagList() error {
	return b.openBrowser(NewBacklogURLBuilder(b.domain, b.spaceKey).
		SetProjectKey(b.projectKey).
		SetRepoName(b.repoName).
		TagListURL())
}

func (b *BacklogRepository) OpenPullRequestList(status string) error {
	s, err := PRStatusFromString(status)
	if err != nil {
		return err
	}
	return b.openBrowser(NewBacklogURLBuilder(b.domain, b.spaceKey).
		SetProjectKey(b.projectKey).
		SetRepoName(b.repoName).
		PullRequestListURL(s.Int()))
}

type PRStatus int

const (
	PRStatusAll PRStatus = iota
	PRStatusOpen
	PRStatusClosed
	PRStatusMerged
)

func (p PRStatus) Int() int {
	return int(p)
}

func PRStatusFromString(s string) (status PRStatus, err error) {
	strToStatus := make(map[string]PRStatus)
	strToStatus["all"] = PRStatusAll
	strToStatus["open"] = PRStatusOpen
	strToStatus["closed"] = PRStatusClosed
	strToStatus["merged"] = PRStatusMerged
	v, ok := strToStatus[s]
	if !ok {
		var specs []string
		for s := range strToStatus {
			specs = append(specs, s)
		}
		err = errors.Errorf("invalid pull request's. choose from %v", specs)
	}
	status = v
	return
}

func (b *BacklogRepository) OpenPullRequest() error {
	id, err := b.findPullRequestIDFromRemote(b.head.Name().String())
	if err != nil {
		return err
	}
	return b.openBrowser(NewBacklogURLBuilder(b.domain, b.spaceKey).
		SetProjectKey(b.projectKey).
		SetRepoName(b.repoName).
		PullRequestURL(id))
}

const (
	refPrefix            = "refs/"
	refPullRequestPrefix = refPrefix + "pull/"
	refPullRequestSuffix = "/head"
)

func (b *BacklogRepository) findPullRequestIDFromRemote(branchName string) (string, error) {

	cmd := exec.Command("git", "ls-remote", "-q")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	rfs := make(map[string]string)
	remotes := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	for _, v := range remotes {
		delimited := strings.Split(v, "\t")
		hash := delimited[0]
		ref := delimited[1]
		rfs[ref] = hash
	}

	targetHash, ok := rfs[branchName]
	if !ok {
		return "", errors.New("not found a current branch in remote")
	}

	var prIDs []string
	for refName, hash := range rfs {
		if !strings.HasPrefix(refName, refPullRequestPrefix) || !strings.HasSuffix(refName, refPullRequestSuffix) {
			continue
		}
		if hash != targetHash {
			continue
		}
		prID := strings.TrimPrefix(refName, refPullRequestPrefix)
		prID = strings.TrimSuffix(prID, refPullRequestSuffix)
		prIDs = append(prIDs, prID)
	}

	if len(prIDs) == 0 {
		return "", errors.New("not found a pull request related to current branch")
	}

	sort.Sort(sort.Reverse(sort.StringSlice(prIDs)))

	return prIDs[0], nil
}

func (b *BacklogRepository) OpenAddPullRequest(base, topic string) error {
	if topic == "" {
		topic = b.head.Name().Short()
	}
	return b.openBrowser(NewBacklogURLBuilder(b.domain, b.spaceKey).
		SetProjectKey(b.projectKey).
		SetRepoName(b.repoName).
		AddPullRequestURL(base, topic))
}

func (b *BacklogRepository) OpenIssue() error {
	key := extractIssueKey(b.head.Name().Short())
	if key == "" {
		return errors.New("could not find issue key in current branch name")
	}
	return b.openBrowser(NewBacklogURLBuilder(b.domain, b.spaceKey).
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
	return b.openBrowser(NewBacklogURLBuilder(b.domain, b.spaceKey).
		SetProjectKey(b.projectKey).
		SetRepoName(b.repoName).
		AddIssueURL())
}

type IssueStatus int

const (
	IssueStatusAll IssueStatus = iota
	IssueStatusOpen
	IssueStatusInProgress
	IssueStatusResolved
	IssueStatusClosed
	IssueStatusNotClosed
)

func (p IssueStatus) Int() int {
	return int(p)
}

func IssueStatusFromString(s string) (status IssueStatus, err error) {
	strToStatus := make(map[string]IssueStatus)
	strToStatus["all"] = IssueStatusAll
	strToStatus["open"] = IssueStatusOpen
	strToStatus["in_progress"] = IssueStatusInProgress
	strToStatus["resolved"] = IssueStatusResolved
	strToStatus["closed"] = IssueStatusClosed
	strToStatus["not_closed"] = IssueStatusNotClosed
	v, ok := strToStatus[s]
	if !ok {
		var specs []string
		for s := range strToStatus {
			specs = append(specs, s)
		}
		err = errors.Errorf("invalid issue's status. choose from %v", specs)
	}
	status = v
	return
}

func (b *BacklogRepository) OpenIssueList(state string) error {
	s, err := IssueStatusFromString(state)
	if err != nil {
		return err
	}
	var statusIds []int
	switch s {
	case IssueStatusAll:
		// Don't specify the issue status
	case IssueStatusNotClosed:
		statusIds = append(statusIds, IssueStatusOpen.Int())
		statusIds = append(statusIds, IssueStatusInProgress.Int())
		statusIds = append(statusIds, IssueStatusResolved.Int())
	default:
		statusIds = append(statusIds, s.Int())
	}
	return b.openBrowser(NewBacklogURLBuilder(b.domain, b.spaceKey).
		SetProjectKey(b.projectKey).
		SetRepoName(b.repoName).
		IssueListURL(statusIds))
}
