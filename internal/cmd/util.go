package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/objectisundefined/github-cli/internal/config"
)

// SearchOptions is the options for Search Issues
type SearchOptions struct {
	Sort  string
	Order string
	Limit int
}

// TimeFormat is the foramt for time output
const TimeFormat string = "2006-01-02T15:04:05"

// RangeTime is a time range in [start, end]
type RangeTime struct {
	Start time.Time
	End   time.Time
}

func newRangeTime() RangeTime {
	n := time.Now()
	return RangeTime{
		Start: n.Add(-7 * 24 * time.Hour),
		End:   n,
	}
}

func (r *RangeTime) String() string {
	return fmt.Sprintf("%s..%s", r.Start.Format(TimeFormat), r.End.Format(TimeFormat))
}

func (r *RangeTime) adjust(sinceTime string, offsetDur string) {
	if len(sinceTime) > 0 {
		end, err := time.Parse(TimeFormat, sinceTime)

		if err != nil {
			panic(err)
		}

		r.End = end
	}

	d, err := time.ParseDuration(offsetDur)

	if err != nil {
		panic(err)
	}

	r.Start = r.End.Add(d)
	if r.Start.After(r.End) {
		r.Start, r.End = r.End, r.Start
	}
}

func splitUsers(s string) []string {
	if len(s) == 0 {
		return []string{}
	}

	return strings.Split(s, ",")
}

func adjustRepoName(owner string, args []string) (string, string) {
	name := ""
	if len(args) > 0 {
		name = args[0]
	}

	if len(name) == 0 {
		return owner, name
	}

	for i := 0; i < len(name); i++ {
		if name[i] == '/' {
			// The name has already been the format of owner/name
			return name[0:i], name[i+1:]
		}
	}

	return owner, name
}

func findRepo(c *config.Config, args []string) config.Repository {
	owner, name := adjustRepoName("", args)
	if repo := c.FindRepo(owner, name); repo != nil {
		return *repo
	}

	// use specail owner and repo
	return config.Repository{Owner: owner, Name: name}

}

func filterRepo(c *config.Config, owner string, args []string) []config.Repository {
	var name string
	owner, name = adjustRepoName(owner, args)
	if len(name) == 0 && len(owner) == 0 {
		return c.Repos
	} else if len(name) == 0 {
		// only owner, filter repos by owner
		var repos []config.Repository
		for _, repo := range c.Repos {
			if repo.Owner == owner {
				repos = append(repos, repo)
			}
		}
		return repos
	}

	// only name
	if r := c.FindRepo(owner, name); r != nil {
		return []config.Repository{*r}
	}

	// use specail owner and repo
	return []config.Repository{
		{Owner: owner, Name: name},
	}
}

func filterUsers(users []*github.User, names []string) bool {
	if len(names) == 0 {
		return true
	}

	for _, name := range names {
		for _, user := range users {
			if user.GetLogin() == name {
				return true
			}
		}
	}

	return false
}
