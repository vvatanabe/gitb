package main

import (
	"reflect"
	"testing"
)

func TestNewBacklogURLBuilder(t *testing.T) {
	type args struct {
		domain   string
		spaceKey string
	}
	tests := []struct {
		name string
		args args
		want *BacklogURLBuilder
	}{
		{
			args: args{domain: "backlog.com", spaceKey: "foo"},
			want: &BacklogURLBuilder{domain: "backlog.com", spaceKey: "foo"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBacklogURLBuilder(tt.args.domain, tt.args.spaceKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBacklogURLBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBacklogURLBuilder_SetProjectKey(t *testing.T) {
	type fields struct {
		domain     string
		spaceKey   string
		projectKey string
		repoName   string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *BacklogURLBuilder
	}{
		{
			fields: fields{
				domain:   "backlog.com",
				spaceKey: "foo",
			},
			args: args{key: "BAR"},
			want: &BacklogURLBuilder{domain: "backlog.com", spaceKey: "foo", projectKey: "BAR"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogURLBuilder{
				domain:     tt.fields.domain,
				spaceKey:   tt.fields.spaceKey,
				projectKey: tt.fields.projectKey,
				repoName:   tt.fields.repoName,
			}
			if got := b.SetProjectKey(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BacklogURLBuilder.SetProjectKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBacklogURLBuilder_SetRepoName(t *testing.T) {
	type fields struct {
		domain     string
		spaceKey   string
		projectKey string
		repoName   string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *BacklogURLBuilder
	}{
		{
			fields: fields{
				domain:   "backlog.com",
				spaceKey: "foo",
			},
			args: args{name: "baz"},
			want: &BacklogURLBuilder{domain: "backlog.com", spaceKey: "foo", repoName: "baz"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogURLBuilder{
				domain:     tt.fields.domain,
				spaceKey:   tt.fields.spaceKey,
				projectKey: tt.fields.projectKey,
				repoName:   tt.fields.repoName,
			}
			if got := b.SetRepoName(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BacklogURLBuilder.SetRepoName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBacklogURLBuilder_Host(t *testing.T) {
	type fields struct {
		domain     string
		spaceKey   string
		projectKey string
		repoName   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				domain:   "backlog.com",
				spaceKey: "foo",
			},
			want: "foo.backlog.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogURLBuilder{
				domain:     tt.fields.domain,
				spaceKey:   tt.fields.spaceKey,
				projectKey: tt.fields.projectKey,
				repoName:   tt.fields.repoName,
			}
			if got := b.Host(); got != tt.want {
				t.Errorf("BacklogURLBuilder.Host() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBacklogURLBuilder_BaseURL(t *testing.T) {
	type fields struct {
		domain     string
		spaceKey   string
		projectKey string
		repoName   string
	}
	tests := []struct {
		fields fields
		want   string
	}{
		{
			fields: fields{
				domain:   "backlog.com",
				spaceKey: "foo",
			},
			want: "https://foo.backlog.com",
		},
	}
	for _, tt := range tests {
		b := &BacklogURLBuilder{
			domain:     tt.fields.domain,
			spaceKey:   tt.fields.spaceKey,
			projectKey: tt.fields.projectKey,
			repoName:   tt.fields.repoName,
		}
		if got := b.BaseURL(); got != tt.want {
			t.Errorf("BacklogURLBuilder.BaseURL() = %v, want %v", got, tt.want)
		}
	}
}

func TestBacklogURLBuilder_GitBaseURL(t *testing.T) {
	type fields struct {
		domain     string
		spaceKey   string
		projectKey string
		repoName   string
	}
	tests := []struct {
		fields fields
		want   string
	}{
		{
			fields: fields{
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			want: "https://foo.backlog.com/git/BAR",
		},
	}
	for _, tt := range tests {
		b := &BacklogURLBuilder{
			domain:     tt.fields.domain,
			spaceKey:   tt.fields.spaceKey,
			projectKey: tt.fields.projectKey,
			repoName:   tt.fields.repoName,
		}
		if got := b.GitBaseURL(); got != tt.want {
			t.Errorf("BacklogURLBuilder.GitBaseURL() = %v, want %v", got, tt.want)
		}
	}
}

func TestBacklogURLBuilder_GitRepoBaseURL(t *testing.T) {
	type fields struct {
		domain     string
		spaceKey   string
		projectKey string
		repoName   string
	}
	tests := []struct {
		fields fields
		want   string
	}{
		{
			fields: fields{
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			want: "https://foo.backlog.com/git/BAR/baz",
		},
	}
	for _, tt := range tests {
		b := &BacklogURLBuilder{
			domain:     tt.fields.domain,
			spaceKey:   tt.fields.spaceKey,
			projectKey: tt.fields.projectKey,
			repoName:   tt.fields.repoName,
		}
		if got := b.GitRepoBaseURL(); got != tt.want {
			t.Errorf("BacklogURLBuilder.GitRepoBaseURL() = %v, want %v", got, tt.want)
		}
	}
}

func TestBacklogURLBuilder_TreeURL(t *testing.T) {
	type fields struct {
		domain     string
		spaceKey   string
		projectKey string
		repoName   string
	}
	type args struct {
		rev string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			fields: fields{
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			args: args{
				"master",
			},
			want: "https://foo.backlog.com/git/BAR/baz/tree/master",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogURLBuilder{
				domain:     tt.fields.domain,
				spaceKey:   tt.fields.spaceKey,
				projectKey: tt.fields.projectKey,
				repoName:   tt.fields.repoName,
			}
			if got := b.TreeURL(tt.args.rev); got != tt.want {
				t.Errorf("BacklogURLBuilder.TreeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBacklogURLBuilder_HistoryURL(t *testing.T) {
	type fields struct {
		domain     string
		spaceKey   string
		projectKey string
		repoName   string
	}
	type args struct {
		rev string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			fields: fields{
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			args: args{
				"master",
			},
			want: "https://foo.backlog.com/git/BAR/baz/history/master",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogURLBuilder{
				domain:     tt.fields.domain,
				spaceKey:   tt.fields.spaceKey,
				projectKey: tt.fields.projectKey,
				repoName:   tt.fields.repoName,
			}
			if got := b.HistoryURL(tt.args.rev); got != tt.want {
				t.Errorf("BacklogURLBuilder.HistoryURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBacklogURLBuilder_BranchListURL(t *testing.T) {
	type fields struct {
		domain     string
		spaceKey   string
		projectKey string
		repoName   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			want: "https://foo.backlog.com/git/BAR/baz/branches",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogURLBuilder{
				domain:     tt.fields.domain,
				spaceKey:   tt.fields.spaceKey,
				projectKey: tt.fields.projectKey,
				repoName:   tt.fields.repoName,
			}
			if got := b.BranchListURL(); got != tt.want {
				t.Errorf("BacklogURLBuilder.BranchListURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBacklogURLBuilder_TagListURL(t *testing.T) {
	type fields struct {
		domain     string
		spaceKey   string
		projectKey string
		repoName   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			want: "https://foo.backlog.com/git/BAR/baz/tags",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogURLBuilder{
				domain:     tt.fields.domain,
				spaceKey:   tt.fields.spaceKey,
				projectKey: tt.fields.projectKey,
				repoName:   tt.fields.repoName,
			}
			if got := b.TagListURL(); got != tt.want {
				t.Errorf("BacklogURLBuilder.TagListURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBacklogURLBuilder_PullRequestListURL(t *testing.T) {
	type fields struct {
		domain     string
		spaceKey   string
		projectKey string
		repoName   string
	}
	type args struct {
		statusID int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			fields: fields{
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			args: args{1},
			want: "https://foo.backlog.com/git/BAR/baz/pullRequests?q.statusId=1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogURLBuilder{
				domain:     tt.fields.domain,
				spaceKey:   tt.fields.spaceKey,
				projectKey: tt.fields.projectKey,
				repoName:   tt.fields.repoName,
			}
			if got := b.PullRequestListURL(tt.args.statusID); got != tt.want {
				t.Errorf("BacklogURLBuilder.PullRequestListURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBacklogURLBuilder_PullRequestURL(t *testing.T) {
	type fields struct {
		domain     string
		spaceKey   string
		projectKey string
		repoName   string
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			fields: fields{
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			args: args{"1"},
			want: "https://foo.backlog.com/git/BAR/baz/pullRequests/1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogURLBuilder{
				domain:     tt.fields.domain,
				spaceKey:   tt.fields.spaceKey,
				projectKey: tt.fields.projectKey,
				repoName:   tt.fields.repoName,
			}
			if got := b.PullRequestURL(tt.args.id); got != tt.want {
				t.Errorf("BacklogURLBuilder.PullRequestURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBacklogURLBuilder_AddPullRequestURL(t *testing.T) {
	type fields struct {
		domain     string
		spaceKey   string
		projectKey string
		repoName   string
	}
	type args struct {
		base  string
		topic string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			fields: fields{
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			args: args{"master", "develop"},
			want: "https://foo.backlog.com/git/BAR/baz/pullRequests/add/master...develop",
		},
		{
			fields: fields{
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			args: args{"", "develop"},
			want: "https://foo.backlog.com/git/BAR/baz/pullRequests/add/...develop",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogURLBuilder{
				domain:     tt.fields.domain,
				spaceKey:   tt.fields.spaceKey,
				projectKey: tt.fields.projectKey,
				repoName:   tt.fields.repoName,
			}
			if got := b.AddPullRequestURL(tt.args.base, tt.args.topic); got != tt.want {
				t.Errorf("BacklogURLBuilder.AddPullRequestURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBacklogURLBuilder_IssueURL(t *testing.T) {
	type fields struct {
		domain     string
		spaceKey   string
		projectKey string
		repoName   string
	}
	type args struct {
		issueKey string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			fields: fields{
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			args: args{"BAR-1234"},
			want: "https://foo.backlog.com/view/BAR-1234",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogURLBuilder{
				domain:     tt.fields.domain,
				spaceKey:   tt.fields.spaceKey,
				projectKey: tt.fields.projectKey,
				repoName:   tt.fields.repoName,
			}
			if got := b.IssueURL(tt.args.issueKey); got != tt.want {
				t.Errorf("BacklogURLBuilder.IssueURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBacklogURLBuilder_AddIssueURL(t *testing.T) {
	type fields struct {
		domain     string
		spaceKey   string
		projectKey string
		repoName   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				"backlog.com",
				"foo",
				"BAR",
				"baz",
			},
			want: "https://foo.backlog.com/add/BAR",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BacklogURLBuilder{
				domain:     tt.fields.domain,
				spaceKey:   tt.fields.spaceKey,
				projectKey: tt.fields.projectKey,
				repoName:   tt.fields.repoName,
			}
			if got := b.AddIssueURL(); got != tt.want {
				t.Errorf("BacklogURLBuilder.AddIssueURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
