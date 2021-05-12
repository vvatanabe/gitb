package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func Test_toRefToHash(t *testing.T) {
	out := []byte(`e73e35d0a86218a9624167110ff8e7fe42596234	HEAD
e73e35d0a86218a9624167110ff8e7fe42596234	refs/heads/master
2b2b5f9e8508a976096a50bd37c81c17ccdf7fb4	refs/heads/patch-1
2b2b5f9e8508a976096a50bd37c81c17ccdf7fb4	refs/pull/3/head
2674ad54e116b4a05d933aa75c7af0657afd0079	refs/tags/0.0.0
`)
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want RefToHash
	}{
		{
			args: args{out},
			want: RefToHash{
				"HEAD":               "e73e35d0a86218a9624167110ff8e7fe42596234",
				"refs/heads/master":  "e73e35d0a86218a9624167110ff8e7fe42596234",
				"refs/heads/patch-1": "2b2b5f9e8508a976096a50bd37c81c17ccdf7fb4",
				"refs/pull/3/head":   "2b2b5f9e8508a976096a50bd37c81c17ccdf7fb4",
				"refs/tags/0.0.0":    "2674ad54e116b4a05d933aa75c7af0657afd0079",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toRefToHash(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toRefToHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractSpaceKeyAndDomain(t *testing.T) {
	type args struct {
		host string
	}
	tests := []struct {
		name         string
		args         args
		wantSpaceKey string
		wantDomain   string
	}{
		{
			args:         args{"foo.backlog.com"},
			wantSpaceKey: "foo",
			wantDomain:   "backlog.com",
		},
		{
			args:         args{"foo.backlog.jp"},
			wantSpaceKey: "foo",
			wantDomain:   "backlog.jp",
		},
		{
			args:         args{"foo.backlogtool.com"},
			wantSpaceKey: "foo",
			wantDomain:   "backlogtool.com",
		},
		{
			args:         args{"bar.git.backlog.com"},
			wantSpaceKey: "bar",
			wantDomain:   "backlog.com",
		},
		{
			args:         args{"bar.git.backlog.jp"},
			wantSpaceKey: "bar",
			wantDomain:   "backlog.jp",
		},
		{
			args:         args{"bar.git.backlogtool.com"},
			wantSpaceKey: "bar",
			wantDomain:   "backlogtool.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSpaceKey, gotDomain := extractSpaceKeyAndDomain(tt.args.host)
			if gotSpaceKey != tt.wantSpaceKey {
				t.Errorf("extractSpaceKeyAndDomain() gotSpaceKey = %v, want %v", gotSpaceKey, tt.wantSpaceKey)
			}
			if gotDomain != tt.wantDomain {
				t.Errorf("extractSpaceKeyAndDomain() gotDomain = %v, want %v", gotDomain, tt.wantDomain)
			}
		})
	}
}

func Test_extractProjectKeyAndRepoName(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name           string
		args           args
		wantProjectKey string
		wantRepoName   string
	}{
		{
			args:           args{"/FOO/bar.git"},
			wantProjectKey: "FOO",
			wantRepoName:   "bar",
		},
		{
			args:           args{"/git/FOO/bar.git"},
			wantProjectKey: "FOO",
			wantRepoName:   "bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotProjectKey, gotRepoName := extractProjectKeyAndRepoName(tt.args.path)
			if gotProjectKey != tt.wantProjectKey {
				t.Errorf("extractProjectKeyAndRepoName() gotProjectKey = %v, want %v", gotProjectKey, tt.wantProjectKey)
			}
			if gotRepoName != tt.wantRepoName {
				t.Errorf("extractProjectKeyAndRepoName() gotRepoName = %v, want %v", gotRepoName, tt.wantRepoName)
			}
		})
	}
}

func TestBacklogRepository_OpenRepositoryList(t *testing.T) {
	type fields struct {
		openBrowser func(url string) error
		repo        Repository
		domain      string
		spaceKey    string
		projectKey  string
		repoName    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			fields: fields{
				openBrowser: func(url string) error {
					want := "https://foo.backlog.com/git/BAR"
					if url != want {
						return errors.New(fmt.Sprintf("result %v, want %v", url, want))
					}
					return nil
				},
				domain:     "backlog.com",
				spaceKey:   "foo",
				projectKey: "BAR",
				repoName:   "baz",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogRepository{
				openBrowser: tt.fields.openBrowser,
				repo:        tt.fields.repo,
				domain:      tt.fields.domain,
				spaceKey:    tt.fields.spaceKey,
				projectKey:  tt.fields.projectKey,
				repoName:    tt.fields.repoName,
			}
			if err := b.OpenRepositoryList(); (err != nil) != tt.wantErr {
				t.Errorf("BacklogRepository.OpenRepositoryList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBacklogRepository_OpenTree(t *testing.T) {
	type fields struct {
		openBrowser func(url string) error
		repo        Repository
		domain      string
		spaceKey    string
		projectKey  string
		repoName    string
	}
	type args struct {
		refOrHash string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields: fields{
				func(url string) error {
					want := "https://foo.backlog.com/git/BAR/baz/tree/develop"
					if url != want {
						return errors.New(fmt.Sprintf("result %v, want %v", url, want))
					}
					return nil
				},
				&RepositoryMock{
					HeadShortNameFunc: func() string {
						return "develop"
					},
				},
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogRepository{
				openBrowser: tt.fields.openBrowser,
				repo:        tt.fields.repo,
				domain:      tt.fields.domain,
				spaceKey:    tt.fields.spaceKey,
				projectKey:  tt.fields.projectKey,
				repoName:    tt.fields.repoName,
			}
			if err := b.OpenTree(tt.args.refOrHash); (err != nil) != tt.wantErr {
				t.Errorf("BacklogRepository.OpenTree() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBacklogRepository_OpenObject(t *testing.T) {
	type fields struct {
		openBrowser func(url string) error
		repo        Repository
		domain      string
		spaceKey    string
		projectKey  string
		repoName    string
	}
	type args struct {
		refOrHash string
	}
	openedUrl := ""
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields: fields{
				func(url string) error {
					openedUrl = url
					return nil
				},
				&RepositoryMock{
					HeadShortNameFunc: func() string {
						return "develop"
					},
					RootDirectoryFunc: func() string {
						return "/path/to/repo"
					},
				},
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogRepository{
				openBrowser: tt.fields.openBrowser,
				repo:        tt.fields.repo,
				domain:      tt.fields.domain,
				spaceKey:    tt.fields.spaceKey,
				projectKey:  tt.fields.projectKey,
				repoName:    tt.fields.repoName,
			}
			if err := b.OpenObject("/path/to/repo/path/to/dir", true, ""); (err != nil) != tt.wantErr {
				t.Errorf("BacklogRepository.OpenObject() error = %v, wantErr %v", err, tt.wantErr)
			}
			expected := "https://foo.backlog.com/git/BAR/baz/tree/develop/path/to/dir"
			if openedUrl != expected {
				t.Errorf("BacklogRepository.OpenObject() error = expected %s but result was %s , wantErr %v", expected, openedUrl, tt.wantErr)
			}
			if err := b.OpenObject("/path/to/repo/path/to/file", false, "10-20"); (err != nil) != tt.wantErr {
				t.Errorf("BacklogRepository.OpenObject() error = %v, wantErr %v", err, tt.wantErr)
			}
			expected = "https://foo.backlog.com/git/BAR/baz/blob/develop/path/to/file#10-20"
			if openedUrl != expected {
				t.Errorf("BacklogRepository.OpenObject() error = expected %s but result was %s , wantErr %v", expected, openedUrl, tt.wantErr)
			}
			if err := b.OpenObject("/path/to/repo/path/to/file", false, "a10-20"); err == nil {
				t.Errorf("line validation doesn't work properly.")
			}
			if err := b.OpenObject("/path/to/repo/path/to/dir", true, "100"); err == nil {
				t.Errorf("line validation doesn't work properly.")
			}

		})
	}
}

func TestBacklogRepository_OpenHistory(t *testing.T) {
	type fields struct {
		openBrowser func(url string) error
		repo        Repository
		domain      string
		spaceKey    string
		projectKey  string
		repoName    string
	}
	type args struct {
		refOrHash string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields: fields{
				func(url string) error {
					want := "https://foo.backlog.com/git/BAR/baz/history/develop"
					if url != want {
						return errors.New(fmt.Sprintf("result %v, want %v", url, want))
					}
					return nil
				},
				&RepositoryMock{
					HeadShortNameFunc: func() string {
						return "develop"
					},
				},
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogRepository{
				openBrowser: tt.fields.openBrowser,
				repo:        tt.fields.repo,
				domain:      tt.fields.domain,
				spaceKey:    tt.fields.spaceKey,
				projectKey:  tt.fields.projectKey,
				repoName:    tt.fields.repoName,
			}
			if err := b.OpenHistory(tt.args.refOrHash); (err != nil) != tt.wantErr {
				t.Errorf("BacklogRepository.OpenHistory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBacklogRepository_OpenNetwork(t *testing.T) {
	type fields struct {
		openBrowser func(url string) error
		repo        Repository
		domain      string
		spaceKey    string
		projectKey  string
		repoName    string
	}
	type args struct {
		refOrHash string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields: fields{
				func(url string) error {
					want := "https://foo.backlog.com/git/BAR/baz/network/develop"
					if url != want {
						return errors.New(fmt.Sprintf("result %v, want %v", url, want))
					}
					return nil
				},
				&RepositoryMock{
					HeadShortNameFunc: func() string {
						return "develop"
					},
				},
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogRepository{
				openBrowser: tt.fields.openBrowser,
				repo:        tt.fields.repo,
				domain:      tt.fields.domain,
				spaceKey:    tt.fields.spaceKey,
				projectKey:  tt.fields.projectKey,
				repoName:    tt.fields.repoName,
			}
			if err := b.OpenNetwork(tt.args.refOrHash); (err != nil) != tt.wantErr {
				t.Errorf("BacklogRepository.OpenNetwork() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBacklogRepository_OpenBranchList(t *testing.T) {
	type fields struct {
		openBrowser func(url string) error
		repo        Repository
		domain      string
		spaceKey    string
		projectKey  string
		repoName    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			fields: fields{
				func(url string) error {
					want := "https://foo.backlog.com/git/BAR/baz/branches"
					if url != want {
						return errors.New(fmt.Sprintf("result %v, want %v", url, want))
					}
					return nil
				},
				&RepositoryMock{},
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogRepository{
				openBrowser: tt.fields.openBrowser,
				repo:        tt.fields.repo,
				domain:      tt.fields.domain,
				spaceKey:    tt.fields.spaceKey,
				projectKey:  tt.fields.projectKey,
				repoName:    tt.fields.repoName,
			}
			if err := b.OpenBranchList(); (err != nil) != tt.wantErr {
				t.Errorf("BacklogRepository.OpenBranchList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBacklogRepository_OpenTagList(t *testing.T) {
	type fields struct {
		openBrowser func(url string) error
		repo        Repository
		domain      string
		spaceKey    string
		projectKey  string
		repoName    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			fields: fields{
				func(url string) error {
					want := "https://foo.backlog.com/git/BAR/baz/tags"
					if url != want {
						return errors.New(fmt.Sprintf("result %v, want %v", url, want))
					}
					return nil
				},
				&RepositoryMock{},
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogRepository{
				openBrowser: tt.fields.openBrowser,
				repo:        tt.fields.repo,
				domain:      tt.fields.domain,
				spaceKey:    tt.fields.spaceKey,
				projectKey:  tt.fields.projectKey,
				repoName:    tt.fields.repoName,
			}
			if err := b.OpenTagList(); (err != nil) != tt.wantErr {
				t.Errorf("BacklogRepository.OpenTagList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBacklogRepository_OpenPullRequestList(t *testing.T) {
	type fields struct {
		openBrowser func(url string) error
		repo        Repository
		domain      string
		spaceKey    string
		projectKey  string
		repoName    string
	}
	type args struct {
		status string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields: fields{
				func(url string) error {
					want := "https://foo.backlog.com/git/BAR/baz/pullRequests?q.statusId=1"
					if url != want {
						return errors.New(fmt.Sprintf("result %v, want %v", url, want))
					}
					return nil
				},
				&RepositoryMock{},
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			args:    args{"open"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogRepository{
				openBrowser: tt.fields.openBrowser,
				repo:        tt.fields.repo,
				domain:      tt.fields.domain,
				spaceKey:    tt.fields.spaceKey,
				projectKey:  tt.fields.projectKey,
				repoName:    tt.fields.repoName,
			}
			if err := b.OpenPullRequestList(tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("BacklogRepository.OpenPullRequestList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPRStatus_Int(t *testing.T) {
	tests := []struct {
		name string
		p    PRStatus
		want int
	}{
		{
			p: PRStatusAll,
		},
		{
			p:    PRStatusOpen,
			want: 1,
		},
		{
			p:    PRStatusClosed,
			want: 2,
		},
		{
			p:    PRStatusMerged,
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Int(); got != tt.want {
				t.Errorf("PRStatus.Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPRStatusFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus PRStatus
		wantErr    bool
	}{
		{
			args:       args{"all"},
			wantStatus: PRStatusAll,
			wantErr:    false,
		},
		{
			args:       args{"open"},
			wantStatus: PRStatusOpen,
			wantErr:    false,
		},
		{
			args:       args{"closed"},
			wantStatus: PRStatusClosed,
			wantErr:    false,
		},
		{
			args:       args{"merged"},
			wantStatus: PRStatusMerged,
			wantErr:    false,
		},
		{
			args:    args{"test"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStatus, err := PRStatusFromString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("PRStatusFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("PRStatusFromString() = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}

func TestBacklogRepository_OpenPullRequest(t *testing.T) {
	type fields struct {
		openBrowser func(url string) error
		repo        Repository
		domain      string
		spaceKey    string
		projectKey  string
		repoName    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			fields: fields{
				func(url string) error {
					want := "https://foo.backlog.com/git/BAR/baz/pullRequests/3"
					if url != want {
						return errors.New(fmt.Sprintf("result %v, want %v", url, want))
					}
					return nil
				},
				&RepositoryMock{
					HeadNameFunc: func() string {
						return "refs/heads/patch-1"
					},
					LsRemoteFunc: func() (RefToHash, error) {
						refToHash := make(RefToHash)
						refToHash["HEAD"] = "e73e35d0a86218a9624167110ff8e7fe42596234"
						refToHash["refs/heads/master"] = "e73e35d0a86218a9624167110ff8e7fe42596234"
						refToHash["refs/heads/patch-1"] = "2b2b5f9e8508a976096a50bd37c81c17ccdf7fb4"
						refToHash["refs/pull/3/head"] = "2b2b5f9e8508a976096a50bd37c81c17ccdf7fb4"
						refToHash["refs/tags/0.0.0"] = "2674ad54e116b4a05d933aa75c7af0657afd0079"
						return refToHash, nil
					},
				},
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			wantErr: false,
		},
		{
			fields: fields{
				func(url string) error {
					want := "https://foo.backlog.com/git/BAR/baz/pullRequests/3"
					if url != want {
						return errors.New(fmt.Sprintf("result %v, want %v", url, want))
					}
					return nil
				},
				&RepositoryMock{
					HeadNameFunc: func() string {
						return "refs/tags/0.0.0"
					},
					LsRemoteFunc: func() (RefToHash, error) {
						refToHash := make(RefToHash)
						refToHash["HEAD"] = "e73e35d0a86218a9624167110ff8e7fe42596234"
						refToHash["refs/heads/master"] = "e73e35d0a86218a9624167110ff8e7fe42596234"
						refToHash["refs/heads/patch-1"] = "2b2b5f9e8508a976096a50bd37c81c17ccdf7fb4"
						refToHash["refs/pull/3/head"] = "2b2b5f9e8508a976096a50bd37c81c17ccdf7fb4"
						refToHash["refs/tags/0.0.0"] = "2674ad54e116b4a05d933aa75c7af0657afd0079"
						return refToHash, nil
					},
				},
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			wantErr: true,
		},
		{
			fields: fields{
				func(url string) error {
					want := "https://foo.backlog.com/git/BAR/baz/pullRequests/3"
					if url != want {
						return errors.New(fmt.Sprintf("result %v, want %v", url, want))
					}
					return nil
				},
				&RepositoryMock{
					HeadNameFunc: func() string {
						return "refs/heads/patch-2"
					},
					LsRemoteFunc: func() (RefToHash, error) {
						refToHash := make(RefToHash)
						refToHash["HEAD"] = "e73e35d0a86218a9624167110ff8e7fe42596234"
						refToHash["refs/heads/master"] = "e73e35d0a86218a9624167110ff8e7fe42596234"
						refToHash["refs/heads/patch-1"] = "2b2b5f9e8508a976096a50bd37c81c17ccdf7fb4"
						refToHash["refs/pull/3/head"] = "2b2b5f9e8508a976096a50bd37c81c17ccdf7fb4"
						refToHash["refs/tags/0.0.0"] = "2674ad54e116b4a05d933aa75c7af0657afd0079"
						return refToHash, nil
					},
				},
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogRepository{
				openBrowser: tt.fields.openBrowser,
				repo:        tt.fields.repo,
				domain:      tt.fields.domain,
				spaceKey:    tt.fields.spaceKey,
				projectKey:  tt.fields.projectKey,
				repoName:    tt.fields.repoName,
			}
			if err := b.OpenPullRequest(); (err != nil) != tt.wantErr {
				t.Errorf("BacklogRepository.OpenPullRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_isPRRef(t *testing.T) {
	type args struct {
		ref string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{"refs/pull/3/head"},
			want: true,
		},
		{
			args: args{"refs/heads/master"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPRRef(tt.args.ref); got != tt.want {
				t.Errorf("isPRRef() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractPRID(t *testing.T) {
	type args struct {
		ref string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{"refs/pull/3/head"},
			want: "3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractPRID(tt.args.ref); got != tt.want {
				t.Errorf("extractPRID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBacklogRepository_OpenAddPullRequest(t *testing.T) {
	type fields struct {
		openBrowser func(url string) error
		repo        Repository
		domain      string
		spaceKey    string
		projectKey  string
		repoName    string
	}
	type args struct {
		base  string
		topic string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields: fields{
				func(url string) error {
					want := "https://foo.backlog.com/git/BAR/baz/pullRequests/add/master...patch-2"
					if url != want {
						return errors.New(fmt.Sprintf("result %v, want %v", url, want))
					}
					return nil
				},
				&RepositoryMock{
					HeadShortNameFunc: func() string {
						return "patch-2"
					},
					LsRemoteFunc: func() (RefToHash, error) {
						refToHash := make(RefToHash)
						refToHash["HEAD"] = "e73e35d0a86218a9624167110ff8e7fe42596234"
						refToHash["refs/heads/master"] = "e73e35d0a86218a9624167110ff8e7fe42596234"
						refToHash["refs/heads/patch-1"] = "2b2b5f9e8508a976096a50bd37c81c17ccdf7fb4"
						refToHash["refs/heads/patch-2"] = "117591be8e3911e4e34d28d9e4bad26d6aa00460"
						refToHash["refs/pull/3/head"] = "2b2b5f9e8508a976096a50bd37c81c17ccdf7fb4"
						refToHash["refs/tags/0.0.0"] = "2674ad54e116b4a05d933aa75c7af0657afd0079"
						return refToHash, nil
					},
				},
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			args:    args{"master", ""},
			wantErr: false,
		},
		{
			fields: fields{
				func(url string) error {
					want := "https://foo.backlog.com/git/BAR/baz/pullRequests/add/master...patch-1"
					if url != want {
						return errors.New(fmt.Sprintf("result %v, want %v", url, want))
					}
					return nil
				},
				&RepositoryMock{
					LsRemoteFunc: func() (RefToHash, error) {
						refToHash := make(RefToHash)
						refToHash["HEAD"] = "e73e35d0a86218a9624167110ff8e7fe42596234"
						refToHash["refs/heads/master"] = "e73e35d0a86218a9624167110ff8e7fe42596234"
						refToHash["refs/heads/patch-1"] = "2b2b5f9e8508a976096a50bd37c81c17ccdf7fb4"
						refToHash["refs/heads/patch-2"] = "117591be8e3911e4e34d28d9e4bad26d6aa00460"
						refToHash["refs/pull/3/head"] = "2b2b5f9e8508a976096a50bd37c81c17ccdf7fb4"
						refToHash["refs/tags/0.0.0"] = "2674ad54e116b4a05d933aa75c7af0657afd0079"
						return refToHash, nil
					},
				},
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			args:    args{"master", "patch-1"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogRepository{
				openBrowser: tt.fields.openBrowser,
				repo:        tt.fields.repo,
				domain:      tt.fields.domain,
				spaceKey:    tt.fields.spaceKey,
				projectKey:  tt.fields.projectKey,
				repoName:    tt.fields.repoName,
			}
			if err := b.OpenAddPullRequest(tt.args.base, tt.args.topic); (err != nil) != tt.wantErr {
				t.Errorf("BacklogRepository.OpenAddPullRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBacklogRepository_OpenIssue(t *testing.T) {
	type fields struct {
		openBrowser func(url string) error
		repo        Repository
		domain      string
		spaceKey    string
		projectKey  string
		repoName    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			fields: fields{
				func(url string) error {
					want := "https://foo.backlog.com/view/BAR-1234"
					if url != want {
						return errors.New(fmt.Sprintf("result %v, want %v", url, want))
					}
					return nil
				},
				&RepositoryMock{
					HeadShortNameFunc: func() string {
						return "BAR-1234"
					},
				},
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			wantErr: false,
		},
		{
			fields: fields{
				func(url string) error {
					return nil
				},
				&RepositoryMock{
					HeadShortNameFunc: func() string {
						return "patch-1"
					},
				},
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogRepository{
				openBrowser: tt.fields.openBrowser,
				repo:        tt.fields.repo,
				domain:      tt.fields.domain,
				spaceKey:    tt.fields.spaceKey,
				projectKey:  tt.fields.projectKey,
				repoName:    tt.fields.repoName,
			}
			if err := b.OpenIssue(); (err != nil) != tt.wantErr {
				t.Errorf("BacklogRepository.OpenIssue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_extractIssueKey(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{"BLG-1234"},
			want: "BLG-1234",
		},
		{
			args: args{"BLG-1234/patch-1"},
			want: "BLG-1234",
		},
		{
			args: args{"BLG-1234.patch-1"},
			want: "BLG-1234",
		},
		{
			args: args{"BLG-1234-patch-1"},
			want: "BLG-1234",
		},
		{
			args: args{"patch-1"},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractIssueKey(tt.args.s); got != tt.want {
				t.Errorf("extractIssueKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBacklogRepository_OpenAddIssue(t *testing.T) {
	type fields struct {
		openBrowser func(url string) error
		repo        Repository
		domain      string
		spaceKey    string
		projectKey  string
		repoName    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			fields: fields{
				func(url string) error {
					want := "https://foo.backlog.com/add/BAR"
					if url != want {
						return errors.New(fmt.Sprintf("result %v, want %v", url, want))
					}
					return nil
				},
				&RepositoryMock{},
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogRepository{
				openBrowser: tt.fields.openBrowser,
				repo:        tt.fields.repo,
				domain:      tt.fields.domain,
				spaceKey:    tt.fields.spaceKey,
				projectKey:  tt.fields.projectKey,
				repoName:    tt.fields.repoName,
			}
			if err := b.OpenAddIssue(); (err != nil) != tt.wantErr {
				t.Errorf("BacklogRepository.OpenAddIssue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIssueStatus_Int(t *testing.T) {
	tests := []struct {
		name string
		p    IssueStatus
		want int
	}{
		{
			p:    IssueStatusAll,
			want: 0,
		},
		{
			p:    IssueStatusOpen,
			want: 1,
		},
		{
			p:    IssueStatusInProgress,
			want: 2,
		},
		{
			p:    IssueStatusResolved,
			want: 3,
		},
		{
			p:    IssueStatusClosed,
			want: 4,
		},
		{
			p:    IssueStatusNotClosed,
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Int(); got != tt.want {
				t.Errorf("IssueStatus.Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIssueStatusFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus IssueStatus
		wantErr    bool
	}{
		{
			args:       args{"all"},
			wantStatus: IssueStatusAll,
			wantErr:    false,
		},
		{
			args:       args{"open"},
			wantStatus: IssueStatusOpen,
			wantErr:    false,
		},
		{
			args:       args{"in_progress"},
			wantStatus: IssueStatusInProgress,
			wantErr:    false,
		},
		{
			args:       args{"resolved"},
			wantStatus: IssueStatusResolved,
			wantErr:    false,
		},
		{
			args:       args{"closed"},
			wantStatus: IssueStatusClosed,
			wantErr:    false,
		},
		{
			args:       args{"not_closed"},
			wantStatus: IssueStatusNotClosed,
			wantErr:    false,
		},
		{
			args:    args{"test"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStatus, err := IssueStatusFromString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("IssueStatusFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("IssueStatusFromString() = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}

func TestBacklogRepository_OpenIssueList(t *testing.T) {
	type fields struct {
		openBrowser func(url string) error
		repo        Repository
		domain      string
		spaceKey    string
		projectKey  string
		repoName    string
	}
	type args struct {
		state string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields: fields{
				func(url string) error {
					want := "https://foo.backlog.com/find/BAR?condition.simpleSearch=true"
					if url != want {
						return errors.New(fmt.Sprintf("result %v, want %v", url, want))
					}
					return nil
				},
				&RepositoryMock{},
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			args:    args{"all"},
			wantErr: false,
		},
		{
			fields: fields{
				func(url string) error {
					want := "https://foo.backlog.com/find/BAR?condition.simpleSearch=true&condition.statusId=1"
					if url != want {
						return errors.New(fmt.Sprintf("result %v, want %v", url, want))
					}
					return nil
				},
				&RepositoryMock{},
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			args:    args{"open"},
			wantErr: false,
		},
		{
			fields: fields{
				func(url string) error {
					want := "https://foo.backlog.com/find/BAR?condition.simpleSearch=true&condition.statusId=1&condition.statusId=2&condition.statusId=3"
					if url != want {
						return errors.New(fmt.Sprintf("result %v, want %v", url, want))
					}
					return nil
				},
				&RepositoryMock{},
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			args:    args{"not_closed"},
			wantErr: false,
		},
		{
			fields: fields{
				func(url string) error {
					return nil
				},
				&RepositoryMock{},
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			args:    args{"test"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogRepository{
				openBrowser: tt.fields.openBrowser,
				repo:        tt.fields.repo,
				domain:      tt.fields.domain,
				spaceKey:    tt.fields.spaceKey,
				projectKey:  tt.fields.projectKey,
				repoName:    tt.fields.repoName,
			}
			if err := b.OpenIssueList(tt.args.state); (err != nil) != tt.wantErr {
				t.Errorf("BacklogRepository.OpenIssueList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBacklogRepository_OpenCommit(t *testing.T) {
	type fields struct {
		openBrowser func(url string) error
		repo        Repository
		domain      string
		spaceKey    string
		projectKey  string
		repoName    string
		hash        string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			fields: fields{
				openBrowser: func(url string) error {
					want := "https://foo.backlog.com/git/BAR/baz/commit/qux"
					if url != want {
						return errors.New(fmt.Sprintf("result %v, want %v", url, want))
					}
					return nil
				},
				domain:     "backlog.com",
				spaceKey:   "foo",
				projectKey: "BAR",
				repoName:   "baz",
				hash:       "qux",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogRepository{
				openBrowser: tt.fields.openBrowser,
				repo:        tt.fields.repo,
				domain:      tt.fields.domain,
				spaceKey:    tt.fields.spaceKey,
				projectKey:  tt.fields.projectKey,
				repoName:    tt.fields.repoName,
			}
			if err := b.OpenCommit(tt.fields.hash); (err != nil) != tt.wantErr {
				t.Errorf("BacklogRepository.OpenCommit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
