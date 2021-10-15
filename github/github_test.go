package github

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v39/github"
)

func TestGetIssueCommentEvent(t *testing.T) {
	t.Parallel()

	var cases = []struct {
		title string
		event *github.IssueCommentEvent
		want  *issueCommentEvent
	}{
		{
			title: "Get repository config",
			event: &github.IssueCommentEvent{
				Repo: &github.Repository{
					Owner: &github.User{
						Login: github.String("tyuhara"),
					},
					Name: github.String("repository"),
				},
				Issue: &github.Issue{
					Number: github.Int(1),
				},
				Comment: &github.IssueComment{
					ID: github.Int64(12345678),
				},
			},
			want: &issueCommentEvent{
				owner:       "tyuhara",
				repo:        "repository",
				issueNumber: 1,
				commentID:   12345678,
			},
		},
	}

	for _, test := range cases {
		t.Run(test.title, func(t *testing.T) {

			client := newFakeGithubClient()
			got, err := client.getIssueCommentEvent(context.Background(), test.event)
			if err != nil {
				t.Errorf("failed: %v", err)
			}

			opt := cmp.AllowUnexported(issueCommentEvent{})
			if diff := cmp.Diff(got, test.want, opt); diff != "" {
				t.Errorf("(-got, +want)\n%s", diff)
			}
		})
	}
}

func TestListReviews(t *testing.T) {
	t.Parallel()

	var cases = []struct {
		title string
		event *issueCommentEvent
		want  []*github.PullRequestReview
	}{
		{
			title: "test",
			event: &issueCommentEvent{
				owner:       "tyuhara",
				repo:        "repository",
				issueNumber: 1,
				commentID:   12345678,
			},
			want: []*github.PullRequestReview{
				{
					ID: github.Int64(371748792),
					User: &github.User{
						Login: github.String("user-a"),
					},
				},
				{
					ID: github.Int64(371748793),
					User: &github.User{
						Login: github.String("user-b"),
					},
				},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.title, func(t *testing.T) {
			client := newFakeGithubClient()
			got, _, _ := client.listReviews(context.Background(), test.event)
			if diff := cmp.Diff(got, test.want); diff != "" {
				t.Errorf("(-got, +want)\n%s", diff)
			}
		})
	}
}

func TestListRepositoryTeamIDs(t *testing.T) {
	t.Parallel()

	var cases = []struct {
		title string
		event *issueCommentEvent
		want  []teamID
	}{
		{
			title: "get team id",
			event: &issueCommentEvent{
				owner:       "tyuhara",
				repo:        "repository",
				issueNumber: 1,
				commentID:   12345678,
			},
			want: []teamID{
				{
					id:   1,
					slug: "team-a",
				},
				{
					id:   2,
					slug: "team-b",
				},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.title, func(t *testing.T) {
			client := newFakeGithubClient()
			got, _ := client.listRepositoryTeamIDs(context.Background(), test.event)
			opt := cmp.AllowUnexported(teamID{})
			if diff := cmp.Diff(got, test.want, opt); diff != "" {
				t.Errorf("(-got, +want)\n%s", diff)
			}
		})
	}
}

func TestHasRequiredTeamApproval(t *testing.T) {
	t.Parallel()

	var cases = []struct {
		title        string
		event        *issueCommentEvent
		reviews      []*github.PullRequestReview
		teamIDs      []teamID
		want         bool
		expectedBool bool
	}{
		{
			title: "get team id",
			event: &issueCommentEvent{
				owner:       "tyuhara",
				repo:        "repository",
				issueNumber: 1,
				commentID:   12345678,
			},
			reviews: []*github.PullRequestReview{
				{
					ID: github.Int64(371748792),
					User: &github.User{
						Login: github.String("user-a"),
					},
				},
				{
					ID: github.Int64(371748793),
					User: &github.User{
						Login: github.String("user-b"),
					},
				},
			},
			teamIDs: []teamID{
				{
					id:   1,
					slug: "team-a",
				},
				{
					id:   2,
					slug: "team-b",
				},
			},
			want:         true,
			expectedBool: true,
		},
		{
			title: "get team id",
			event: &issueCommentEvent{
				owner:       "tyuhara",
				repo:        "repository",
				issueNumber: 1,
				commentID:   12345678,
			},
			reviews: []*github.PullRequestReview{
				{
					ID: github.Int64(371748792),
					User: &github.User{
						Login: github.String("user-a"),
					},
				},
			},
			teamIDs: []teamID{
				{
					id:   1,
					slug: "team-a",
				},
				{
					id:   2,
					slug: "team-b",
				},
			},
			want:         false,
			expectedBool: false,
		},
	}

	client := newFakeGithubClient()

	for _, test := range cases {
		if test.expectedBool {
			t.Run(test.title, func(t *testing.T) {
				got := client.hasRequiredTeamApproval(context.Background(), test.event, test.reviews, test.teamIDs)
				opt := cmp.AllowUnexported(approvedTeam{})
				if got != test.expectedBool {
					t.Errorf("Has not required approval - wanted: %t, found %t", test.expectedBool, got)
				}
				if diff := cmp.Diff(got, test.want, opt); diff != "" {
					t.Errorf("(-got, +want)\n%s", diff)
				}
			})
		}
	}
}
