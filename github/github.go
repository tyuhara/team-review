package github

import (
	"context"
	"net/http"
	"time"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v39/github"
	"go.uber.org/zap"

	"github.com/tyuhara/team-review/config"
)

type client struct {
	client *github.Client
	logger *zap.SugaredLogger
}

type issueCommentEvent struct {
	owner       string
	repo        string
	issueNumber int
	commentID   int64
}

type teamID struct {
	id   int64
	slug string
}

type approvedTeam struct {
	slug string
}

func newGithubClient(httpClient *http.Client, logger *zap.SugaredLogger) (*client, error) {
	return &client{
		client: github.NewClient(httpClient),
		logger: logger,
	}, nil
}

func HandleMergeIssueComment(ctx context.Context, event *github.IssueCommentEvent, conf *config.Config, logger *zap.SugaredLogger) error {
	if event.GetComment().GetBody() != "/merge" {
		return nil
	}

	// Create github client
	installationID := event.GetInstallation().GetID()
	tr := http.DefaultTransport
	itr, err := ghinstallation.New(tr, conf.GithubAppID, installationID, []byte(conf.GithubAppPrivateKey))
	if err != nil {
		return err
	}

	c, err := newGithubClient(&http.Client{
		Transport: itr,
		Timeout:   5 * time.Second,
	},
		logger)
	if err != nil {
		logger.Error(err)
	}

	// get necessary event from github webhook event
	e, err := c.getIssueCommentEvent(ctx, event)
	if err != nil {
		logger.Error(err)
	}

	// Leave "eyes" to issue comment
	content := "eyes"
	if _, _, err := c.client.Reactions.CreateIssueCommentReaction(ctx, e.owner, e.repo, e.commentID, content); err != nil {
		return err
	}

	// List reviews (list approved member)
	reviews, err := c.listReviews(ctx, e)
	if err != nil {
		c.logger.Errorf("err")
	}

	// List repositories teams
	teamIDs, err := c.listRepositoryTeamIDs(ctx, e)
	if err != nil {
		c.logger.Errorf("err")
	}

	// If don't have the required team approval
	if !c.hasRequiredTeamApproval(ctx, e, reviews, teamIDs) {
		body := "Need to approve from all reviewer team"
		c.logger.Infof("%s", body)
		comment := &github.IssueComment{
			Body: &body,
		}
		if _, _, err := c.client.Issues.CreateComment(ctx, e.owner, e.repo, e.issueNumber, comment); err != nil {
			return nil
		}
	}

	//　If the required team approval is sufficient
	if c.hasRequiredTeamApproval(ctx, e, reviews, teamIDs) {
		comment := "merged"
		c.logger.Infof("%s", comment)
		options := &github.PullRequestOptions{
			MergeMethod: "merge",
		}
		if _, _, err := c.client.PullRequests.Merge(ctx, e.owner, e.repo, e.issueNumber, comment, options); err != nil {
			return err
		}
	}

	return nil
}

// Get repository config
func (c *client) getIssueCommentEvent(ctx context.Context, event *github.IssueCommentEvent) (*issueCommentEvent, error) {
	owner := event.Repo.GetOwner().GetLogin()
	repo := event.Repo.GetName()
	issueNumber := event.Issue.GetNumber()
	commentID := event.Comment.GetID()

	return &issueCommentEvent{
		owner:       owner,
		repo:        repo,
		issueNumber: issueNumber,
		commentID:   commentID,
	}, nil
}

// List repositories teams
func (c *client) listRepositoryTeamIDs(ctx context.Context, event *issueCommentEvent) ([]teamID, error) {
	opt := &github.ListOptions{}
	teams, _, err := c.client.Repositories.ListTeams(ctx, event.owner, event.repo, opt)
	if err != nil {
		c.logger.Errorf("[ERROR] teams: %s", err)
	}

	teamIDs := make([]teamID, 0)
	for _, v := range teams {
		teamIDs = append(teamIDs, teamID{v.GetID(), v.GetSlug()})
	}

	return teamIDs, nil
}

func (c *client) listReviews(ctx context.Context, e *issueCommentEvent) ([]*github.PullRequestReview, error) {
	opt := &github.ListOptions{}
	reviews, _, err := c.client.PullRequests.ListReviews(ctx, e.owner, e.repo, e.issueNumber, opt)
	if err != nil {
		c.logger.Errorf("[ERROR] reviews: %s", err)
		return nil, err
	}
	return reviews, nil
}

func (c *client) hasRequiredTeamApproval(ctx context.Context, event *issueCommentEvent, reviews []*github.PullRequestReview, teamIDs []teamID) bool {
	opt := &github.TeamListTeamMembersOptions{}
	approvedTeams := make([]approvedTeam, 0)

	for _, v := range reviews {
		if v.GetState() != "APPROVED" {
			continue
		}

		for _, x := range teamIDs {
			// List team members
			// https://github.com/google/go-github/blob/33ae6f3d80a32cc768320aad783f70c58a38b45a/github/teams_members.go#L52
			users, _, err := c.client.Teams.ListTeamMembersBySlug(ctx, event.owner, x.slug, opt)
			if err != nil {
				c.logger.Errorf("[ERROR] reviews: %s", err)
			}

			for _, y := range users {
				// Approved user == team member
				if v.User.GetLogin() != y.GetLogin() {
					continue
				}
				approvedTeams = append(approvedTeams, approvedTeam{x.slug})
			}
		}
	}

	return len(teamIDs) == len(approvedTeams)
}
