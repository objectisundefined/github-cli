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

// TimeFormat is the format for time output
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

	// use special owner and repo
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

	// use special owner and repo
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

type style struct{}

func (_ *style) Black(s string) string         { return fmt.Sprintf("\u001b[30m%s\u001b[39m", s) }
func (_ *style) Red(s string) string           { return fmt.Sprintf("\u001b[31m%s\u001b[39m", s) }
func (_ *style) Green(s string) string         { return fmt.Sprintf("\u001b[32m%s\u001b[39m", s) }
func (_ *style) Yellow(s string) string        { return fmt.Sprintf("\u001b[33m%s\u001b[39m", s) }
func (_ *style) Blue(s string) string          { return fmt.Sprintf("\u001b[34m%s\u001b[39m", s) }
func (_ *style) Magenta(s string) string       { return fmt.Sprintf("\u001b[35m%s\u001b[39m", s) }
func (_ *style) Cyan(s string) string          { return fmt.Sprintf("\u001b[36m%s\u001b[39m", s) }
func (_ *style) White(s string) string         { return fmt.Sprintf("\u001b[37m%s\u001b[39m", s) }
func (_ *style) Gray(s string) string          { return fmt.Sprintf("\u001b[90m%s\u001b[39m", s) }
func (_ *style) Grey(s string) string          { return fmt.Sprintf("\u001b[90m%s\u001b[39m", s) }
func (_ *style) BrightRed(s string) string     { return fmt.Sprintf("\u001b[91m%s\u001b[39m", s) }
func (_ *style) BrightGreen(s string) string   { return fmt.Sprintf("\u001b[92m%s\u001b[39m", s) }
func (_ *style) BrightYellow(s string) string  { return fmt.Sprintf("\u001b[93m%s\u001b[39m", s) }
func (_ *style) BrightBlue(s string) string    { return fmt.Sprintf("\u001b[94m%s\u001b[39m", s) }
func (_ *style) BrightMagenta(s string) string { return fmt.Sprintf("\u001b[95m%s\u001b[39m", s) }
func (_ *style) BrightCyan(s string) string    { return fmt.Sprintf("\u001b[96m%s\u001b[39m", s) }
func (_ *style) BrightWhite(s string) string   { return fmt.Sprintf("\u001b[97m%s\u001b[39m", s) }

var Styles style
