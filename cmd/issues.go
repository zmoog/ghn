/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/google/go-github/v69/github"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/zmoog/ws/feedback"
)

var (
	includeAllRepos bool
	state           string
	assignee        string
	creator         string
	daysAgo         int
	sort            string
)

// issuesCmd represents the issues command
var issuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "List and search issues",
	Long:  `List and search issues from GitHub.`,
}

// listIssuesCmd represents the list command
var listIssuesCmd = &cobra.Command{
	Use:   "list",
	Short: "List issues",
	Long:  `List issues from GitHub.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := &github.IssueListOptions{
			Sort:  "updated",
			State: "open",
			Since: time.Now().Add(-time.Hour * 24 * time.Duration(daysAgo)),
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		}

		issues, _, err := client.Issues.List(context.Background(), includeAllRepos, opts)
		if err != nil {
			return fmt.Errorf("failed to list issues: %w", err)
		}

		feedback.Println(fmt.Sprintf("Found %d issues", len(issues)))
		feedback.PrintResult(IssueResult{
			Issues: issues,
		})

		return nil
	},
}

// listByRepoCmd represents the listByRepo command
var listByRepoCmd = &cobra.Command{
	Use:   "list-by-repo [owner] [repo]",
	Short: "List issues by repository",
	Long:  `List issues by repository from GitHub.`,
	Args:  cobra.ExactArgs(2),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		validSortValues := []string{"updated", "created"}

		if !slices.Contains(validSortValues, sort) {
			return fmt.Errorf("sort must be either [%s], got '%s'", strings.Join(validSortValues, " or "), sort)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		owner = args[0]
		repo = args[1]

		opts := &github.IssueListByRepoOptions{
			State: state,
			Since: time.Now().Add(-time.Hour * 24 * time.Duration(daysAgo)),
			Sort:  sort,
		}

		if creator != "" {
			opts.Creator = creator
		}

		if assignee != "*" {
			opts.Assignee = assignee
		}

		issues, _, err := client.Issues.ListByRepo(context.Background(), owner, repo, opts)
		if err != nil {
			return fmt.Errorf("failed to list issues: %w", err)
		}

		feedback.Println(fmt.Sprintf("Found %d issues", len(issues)))
		feedback.PrintResult(IssueResult{
			Issues: issues,
		})

		return nil
	},
}

func init() {
	rootCmd.AddCommand(issuesCmd)

	issuesCmd.AddCommand(listIssuesCmd)
	issuesCmd.AddCommand(listByRepoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// issuesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// issuesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	listIssuesCmd.Flags().BoolVarP(&includeAllRepos, "all", "a", true, "Include issues from all repositories")

	listByRepoCmd.Flags().StringVarP(&assignee, "assignee", "a", "*", "Assignee of the issues")
	listByRepoCmd.Flags().StringVarP(&creator, "creator", "c", "", "Creator of the issues")
	listByRepoCmd.Flags().IntVarP(&daysAgo, "days-ago", "d", 30, "Days ago to search for issues")
	listByRepoCmd.Flags().StringVarP(&state, "state", "s", "open", "State of the issues")
	listByRepoCmd.Flags().StringVarP(&sort, "sort", "o", "updated", "Sort by (updated or created)")
}

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
