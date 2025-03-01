package issues

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v69/github"
	"github.com/spf13/cobra"
	"github.com/zmoog/ogh/internal/client"
	"github.com/zmoog/ws/feedback"
)

// listIssuesCmd represents the list command
var listIssuesCmd = &cobra.Command{
	Use:   "list",
	Short: "List issues",
	Long:  `List issues from GitHub.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create GitHub client: %w", err)
		}

		opts := &github.IssueListOptions{
			Sort:   sort,
			State:  state,
			Since:  time.Now().Add(-time.Hour * 24 * time.Duration(daysAgo)),
			Filter: filter,
			Labels: labels,
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

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// issuesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// issuesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	listIssuesCmd.Flags().BoolVarP(&includeAllRepos, "all", "a", true, "Include issues from all repositories")
	listIssuesCmd.Flags().IntVarP(&daysAgo, "days-ago", "d", 30, "Days ago to search for issues")
	listIssuesCmd.Flags().StringVarP(&sort, "sort", "o", "updated", "Sort by (updated or created)")
	listIssuesCmd.Flags().StringVarP(&state, "state", "s", "open", "State of the issues")
	listIssuesCmd.Flags().StringVarP(&filter, "filter", "f", "assigned", "Filter issues by (assignee, creator, mentioned, subscribed, all)")
	listIssuesCmd.Flags().StringSliceVarP(&labels, "labels", "l", []string{}, "Labels to filter issues by")
	// Add the list command to the issues command
	issuesCmd.AddCommand(listIssuesCmd)
}
