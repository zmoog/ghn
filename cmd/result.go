package cmd

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v69/github"
	"github.com/pterm/pterm"
)

type Result struct {
	Notifications []*github.Notification
}

func (r Result) Data() any {
	return r.Notifications
}

func (r Result) String() string {
	return r.Table()
}

func (r Result) Table() string {

	table := pterm.TableData{}
	table = append(table, []string{
		"Type",
		"Repository",
		"Reason",
		"URL",
		"Title",
	})

	for _, notification := range r.Notifications {
		table = append(table, []string{
			*notification.Subject.Type,
			*notification.Repository.Name,
			*notification.Reason,
			asLink(notification),
			*notification.Subject.Title,
		})
	}

	render, _ := pterm.DefaultTable.WithHasHeader().WithData(table).Srender()

	return render
}

func asLink(notification *github.Notification) string {
	// https://github.com/zmoog/ws/issues/1
	// https://github.com/zmoog/ws/pull/2
	baseUrl := "https://github.com"

	parts := strings.Split(*notification.Subject.URL, "/")
	number := parts[len(parts)-1]

	switch *notification.Subject.Type {
	case "Issue":
		return fmt.Sprintf("%s/%s/%s/issues/%s", baseUrl, *notification.Repository.Owner.Login, *notification.Repository.Name, number)
	case "PullRequest":
		return fmt.Sprintf("%s/%s/%s/pull/%s", baseUrl, *notification.Repository.Owner.Login, *notification.Repository.Name, number)
	}

	return "unknown"
}
