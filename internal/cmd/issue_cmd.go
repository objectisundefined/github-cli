package cmd

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/objectisundefined/github-cli/internal/client"
	"github.com/spf13/cobra"
)

var (
	issuesState     string
	issuesLimit     int
	issuesSinceTime string
	issuesOffsetDur string
	issuesOwner     string
	issuesAssignees string
)

func NewIssuesCommand() *cobra.Command {
	m := &cobra.Command{
		Use:   "issues [repo]",
		Short: "Github CLI for listing issues",
		Args:  cobra.MinimumNArgs(0),
		Run:   runIssuesCommandFunc,
	}
	m.Flags().StringVar(&issuesState, "state", "open", "Issue state: open or closed")
	m.Flags().IntVar(&issuesLimit, "limit", 20, "Maximum issues limit for a repository")
	m.Flags().StringVar(&issuesSinceTime, "since", "", fmt.Sprintf("Issue Since Time, format is %s", TimeFormat))
	m.Flags().StringVar(&issuesOffsetDur, "offset", "-48h", "The offset of since time")
	m.Flags().StringVar(&issuesOwner, "owner", "", "The Github account")
	m.Flags().StringVar(&issuesAssignees, "assignees", "", "Assignees for the issue, separated by comma")
	return m
}

func runIssuesCommandFunc(cmd *cobra.Command, args []string) {
	opts := SearchOptions{
		Order: "desc",
		Sort:  "updated",
		Limit: issuesLimit,
	}

	queryArgs := url.Values{}
	users := splitUsers(issuesAssignees)
	for _, user := range users {
		queryArgs.Add("assignee", user)
	}

	queryArgs.Add("is", "issue")
	rangeTime := newRangeTime()
	rangeTime.adjust(issuesSinceTime, issuesOffsetDur)

	queryArgs.Add("updated", rangeTime.String())
	queryArgs.Add("state", issuesState)

	repos := filterRepo(cfg, issuesOwner, args)

	ic := IssueClient{client.NewClient(ctx, cfg)}

	m, err := ic.SearchIssues(ctx, repos, opts, queryArgs)

	if err != nil {
		panic(err)
	}

	for repo, issues := range m {
		if len(issues) == 0 {
			continue
		}

		fmt.Fprintln(cmd.OutOrStdout(), repo)
		for _, issue := range issues {
			fmt.Printf("%s %s %s\n", issue.GetUpdatedAt().Format(TimeFormat), issue.GetHTMLURL(), issue.GetTitle())
		}
	}
}

var (
	issueCommentLimit int
)

func NewIssueCommand() *cobra.Command {
	m := &cobra.Command{
		Use:   "issue [repo] [id]",
		Short: "Github CLI for getting one pull",
		Args:  cobra.MinimumNArgs(2),
		Run:   runIssueCommandFunc,
	}

	m.Flags().IntVar(&pullCommentLimit, "comments-limit", 3, "Comments limit")
	return m
}

func runIssueCommandFunc(cmd *cobra.Command, args []string) {
	id, err := strconv.Atoi(args[1])

	if err != nil {
		panic(err)
	}

	repo := findRepo(cfg, args)

	ic := IssueClient{client.NewClient(ctx, cfg)}

	issue, err := ic.GetIssue(ctx, repo.Owner, repo.Name, id)

	if err != nil {
		panic(err)
	}

	comments, err := ic.ListIssueComments(ctx, repo.Owner, repo.Name, id)

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Title: %s\n", issue.GetTitle())
	fmt.Fprintf(cmd.OutOrStdout(), "Created at %s\n", issue.GetCreatedAt().Format(TimeFormat))
	fmt.Fprintf(cmd.OutOrStdout(), "Message:\n %s\n", issue.GetBody())
	if len(comments) > issueCommentLimit {
		comments = comments[0:issueCommentLimit]
	}
	for _, comment := range comments {
		fmt.Fprintf(cmd.OutOrStdout(), "Comment:\n %s\n", comment.GetBody())
	}
}
