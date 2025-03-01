package issues

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/google/go-github/v69/github"
	"github.com/pterm/pterm"
)

type IssueResult struct {
	Issues []*github.Issue
}

func (r IssueResult) Data() any {
	return r.Issues
}

func (r IssueResult) String() string {
	return r.Table()
}

func (r IssueResult) Table() string {
	table := pterm.TableData{}
	table = append(table, []string{
		"Title",
		"Number",
		"State",
		"User",
		"Created",
		"Updated",
	})

	for _, issue := range r.Issues {
		table = append(table, []string{
			*issue.Title,
			fmt.Sprintf("https://github.com/%s/%s/issues/%d", owner, repo, *issue.Number),
			*issue.State,
			*issue.User.Login,
			humanize.Time(issue.CreatedAt.Time),
			humanize.Time(issue.UpdatedAt.Time),
		})
	}

	render, _ := pterm.DefaultTable.WithHasHeader().WithData(table).Srender()

	return render
}
