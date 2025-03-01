package pulls

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/google/go-github/v69/github"
	"github.com/pterm/pterm"
)

type PullRequestResult struct {
	Pulls []*github.PullRequest
	Repo  string
	Owner string
}

func (r PullRequestResult) Data() any {
	return r.Pulls
}

func (r PullRequestResult) String() string {
	return r.Table()
}

func (r PullRequestResult) Table() string {
	table := pterm.TableData{}
	table = append(table, []string{
		"Title",
		"URL",
		"State",
		"User",
		"Created",
		"Updated",
	})

	for _, pull := range r.Pulls {
		table = append(table, []string{
			*pull.Title,
			r.asLink(pull),
			*pull.State,
			*pull.User.Login,
			humanize.Time(pull.CreatedAt.Time),
			humanize.Time(pull.UpdatedAt.Time),
		})
	}

	render, err := pterm.DefaultTable.WithHasHeader().WithData(table).Srender()
	if err != nil {
		return fmt.Sprintf("Error rendering table: %s", err)
	}

	return render
}

func (r PullRequestResult) asLink(pull *github.PullRequest) string {
	return fmt.Sprintf("https://github.com/%s/%s/pull/%d", r.Owner, r.Repo, *pull.Number)
}
