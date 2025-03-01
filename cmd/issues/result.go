package issues

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/google/go-github/v69/github"
	"github.com/pterm/pterm"
)

type IssueResult struct {
	Issues []*github.Issue
	Repo   string
	Owner  string
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
		"URL",
		"State",
		"User",
		"Created",
		"Updated",
	})

	for _, issue := range r.Issues {
		table = append(table, []string{
			*issue.Title,
			r.asLink(issue),
			*issue.State,
			*issue.User.Login,
			humanize.Time(issue.CreatedAt.Time),
			humanize.Time(issue.UpdatedAt.Time),
		})
	}

	render, err := pterm.DefaultTable.WithHasHeader().WithData(table).Srender()
	if err != nil {
		return fmt.Sprintf("Error rendering table: %s", err)
	}

	return render
}

func (r IssueResult) asLink(issue *github.Issue) string {
	if issue.Repository != nil {
		return fmt.Sprintf("https://github.com/%s/%s/issues/%d", *issue.Repository.Owner.Login, *issue.Repository.Name, *issue.Number)
	}

	if r.Owner != "" && r.Repo != "" {
		return fmt.Sprintf("https://github.com/%s/%s/issues/%d", r.Owner, r.Repo, *issue.Number)
	}

	return "-"
}
