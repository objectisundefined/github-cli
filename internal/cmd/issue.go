package cmd

import (
	"bytes"
	"context"
	"fmt"
	"net/url"

	"github.com/google/go-github/github"
	"github.com/objectisundefined/github-cli/internal/config"
)

type IssueClient struct {
	*github.Client
}

// SearchIssues provides a common way to search issues or pull requests.
func (c *IssueClient) SearchIssues(ctx context.Context, repos []config.Repository, opts SearchOptions, queryArgs url.Values) (map[string][]github.Issue, error) {
	m := make(map[string][]github.Issue, len(repos))
	for _, repo := range repos {
		issues, err := c.SearchIssuesByRepo(ctx, repo, opts, queryArgs)
		if err != nil {
			return nil, err
		}
		m[fmt.Sprintf("%s/%s", repo.Owner, repo.Name)] = issues
	}
	return m, nil
}

// SearchIssuesByRepo provides a common way to search issues or pull requests.
func (c *IssueClient) SearchIssuesByRepo(ctx context.Context, repo config.Repository, opts SearchOptions, queryArgs url.Values) ([]github.Issue, error) {
	opt := github.SearchOptions{
		Sort:  opts.Sort,
		Order: opts.Order,
	}

	queryArgs.Del("repo")
	queryArgs.Add("repo", repo.String())

	var (
		query bytes.Buffer
		first = true
	)
	for key, values := range queryArgs {
		for _, value := range values {
			if !first {
				query.WriteByte(' ')
			}
			first = false
			query.WriteString(fmt.Sprintf("%s:%s", key, value))
		}
	}

	var allIssues []github.Issue
	for {
		issues, resp, err := c.Search.Issues(ctx, query.String(), &opt)
		if err != nil {
			return nil, err
		}

		allIssues = append(allIssues, issues.Issues...)

		if opts.Limit > 0 && len(allIssues) >= opts.Limit {
			break
		}

		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}

	return allIssues, nil
}

func (c *IssueClient) GetIssue(ctx context.Context, owner string, repo string, id int) (*github.Issue, error) {
	r, _, err := c.Issues.Get(ctx, owner, repo, id)
	return r, err
}

func (c *IssueClient) ListIssueComments(ctx context.Context, owner string, repo string, number int) ([]*github.IssueComment, error) {
	var allComments []*github.IssueComment

	opts := github.IssueListCommentsOptions{
		Sort:      "updated",
		Direction: "desc",
	}

	for {
		comments, resp, err := c.Issues.ListComments(ctx, owner, repo, number, &opts)
		if err != nil {
			return nil, err
		}
		allComments = append(allComments, comments...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return allComments, nil
}
