package cmd

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/objectisundefined/github-cli/internal/client"
	"github.com/spf13/cobra"
)

var (
	pullsState     string
	pullsLimit     int
	pullsSinceTime string
	pullsOffsetDur string
	pullsOwner     string
	pullsReviewers string
)

func NewPullsCommand() *cobra.Command {
	m := &cobra.Command{
		Use:   "pulls [repo]",
		Short: "Github CLI for listing pulls",
		Args:  cobra.MinimumNArgs(0),
		Run:   runPullsCommandFunc,
	}
	m.Flags().StringVar(&pullsState, "state", "open", "PR state: open or closed")
	m.Flags().IntVar(&pullsLimit, "limit", 20, "Maximum pull limit for a repository")
	m.Flags().StringVar(&pullsSinceTime, "since", "", fmt.Sprintf("Issue Since Time, format is %s", TimeFormat))
	m.Flags().StringVar(&pullsOffsetDur, "offset", "-48h", "The offset of since time")
	m.Flags().StringVar(&pullsOwner, "owner", "", "The Github account")
	m.Flags().StringVar(&pullsReviewers, "reviewers", "", "Request reviewers, separated by comma")
	return m
}

func runPullsCommandFunc(cmd *cobra.Command, args []string) {
	opts := SearchOptions{
		Order: "desc",
		Sort:  "updated",
		Limit: issuesLimit,
	}

	queryArgs := url.Values{}
	users := splitUsers(pullsReviewers)
	for _, user := range users {
		queryArgs.Add("assignee", user)
	}

	queryArgs.Add("is", "pr")
	rangeTime := newRangeTime()
	rangeTime.adjust(pullsSinceTime, pullsOffsetDur)

	queryArgs.Add("updated", rangeTime.String())
	queryArgs.Add("state", pullsState)

	repos := filterRepo(cfg, pullsOwner, args)

	ic := IssueClient{client.NewClient(ctx, cfg)}

	m, err := ic.SearchIssues(ctx, repos, opts, queryArgs)

	if err != nil {
		panic(err)
	}

	for repo, pulls := range m {
		if len(pulls) == 0 {
			continue
		}

		fmt.Println(repo)
		for _, pull := range pulls {
			fmt.Fprintf(cmd.OutOrStdout(), "%s %s %s\n", pull.GetUpdatedAt().Format(TimeFormat), pull.GetHTMLURL(), pull.GetTitle())
		}
	}
}

var (
	pullCommentLimit int
)

func NewPullCommand() *cobra.Command {
	m := &cobra.Command{
		Use:   "pull [repo] [id]",
		Short: "Github CLI for getting one pull",
		Args:  cobra.MinimumNArgs(2),
		Run:   runPullCommandFunc,
	}

	m.Flags().IntVar(&pullCommentLimit, "comments-limit", 3, "Comments limit")
	return m
}

func runPullCommandFunc(cmd *cobra.Command, args []string) {
	id, err := strconv.Atoi(args[1])

	if err != nil {
		panic(err)
	}

	repo := findRepo(cfg, args)

	pc := PullClient{client.NewClient(ctx, cfg)}

	pull, err := pc.GetPull(ctx, repo.Owner, repo.Name, id)

	if err != nil {
		panic(err)
	}

	comments, err := pc.ListPullComments(ctx, repo.Owner, repo.Name, id)

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Title: %s\n", pull.GetTitle())
	fmt.Fprintf(cmd.OutOrStdout(), "Created at %s\n", pull.GetCreatedAt().Format(TimeFormat))
	fmt.Fprintf(cmd.OutOrStdout(), "Message:\n %s\n", pull.GetBody())
	if len(comments) > pullCommentLimit {
		comments = comments[0:pullCommentLimit]
	}
	for _, comment := range comments {
		fmt.Fprintf(cmd.OutOrStdout(), "Comment:\n %s\n", comment.GetBody())
	}
}
