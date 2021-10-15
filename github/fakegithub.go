package github

import (
	"context"

	"github.com/google/go-github/v39/github"
)

type fakeClient struct {
	FakeGetIssueCommentEvent    func(ctx context.Context, event *github.IssueCommentEvent) (*issueCommentEvent, error)
	FakeListReviews             func(ctx context.Context, owner, repo string, number int, opt *github.ListOptions) ([]*github.PullRequestReview, *github.Response, error)
	FakeListRepositoryTeamIDs   func(ctx context.Context, event *issueCommentEvent) ([]teamID, error)
	FakeHasRequiredTeamApproval func(ctx context.Context, event *issueCommentEvent, reviews []*github.PullRequestReview, teamIDs []teamID) bool
}

func (f *fakeClient) getIssueCommentEvent(ctx context.Context, event *github.IssueCommentEvent) (*issueCommentEvent, error) {
	return f.FakeGetIssueCommentEvent(ctx, event)
}

func (f *fakeClient) listReviews(ctx context.Context, e *issueCommentEvent) ([]*github.PullRequestReview, *github.Response, error) {
	return f.FakeListReviews(ctx, e.owner, e.repo, e.issueNumber, &github.ListOptions{})
}

func (f *fakeClient) listRepositoryTeamIDs(ctx context.Context, event *issueCommentEvent) ([]teamID, error) {
	return f.FakeListRepositoryTeamIDs(ctx, event)
}

func (f *fakeClient) hasRequiredTeamApproval(ctx context.Context, event *issueCommentEvent, reviews []*github.PullRequestReview, teamIDs []teamID) bool {
	return f.FakeHasRequiredTeamApproval(ctx, event, reviews, teamIDs)
}

func newFakeGithubClient() *fakeClient {
	return &fakeClient{
		FakeGetIssueCommentEvent: func(ctx context.Context, event *github.IssueCommentEvent) (*issueCommentEvent, error) {
			return &issueCommentEvent{
				owner:       event.Repo.Owner.GetLogin(),
				repo:        event.Repo.GetName(),
				issueNumber: event.Issue.GetNumber(),
				commentID:   event.Comment.GetID(),
			}, nil
		},
		FakeListReviews: func(ctx context.Context, owner, repo string, number int, opt *github.ListOptions) ([]*github.PullRequestReview, *github.Response, error) {
			return []*github.PullRequestReview{
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
			}, nil, nil
		},
		FakeListRepositoryTeamIDs: func(ctx context.Context, event *issueCommentEvent) ([]teamID, error) {
			return []teamID{
				{
					id:   1,
					slug: "team-a",
				},
				{
					id:   2,
					slug: "team-b",
				},
			}, nil
		},
		FakeHasRequiredTeamApproval: func(ctx context.Context, event *issueCommentEvent, reviews []*github.PullRequestReview, teamIDs []teamID) bool {
			return true
		},
	}
}
