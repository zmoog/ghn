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

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// issuesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// issuesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	listIssuesCmd.Flags().BoolVarP(&includeAllRepos, "all", "a", true, "Include issues from all repositories")

	// Add the list command to the issues command
	issuesCmd.AddCommand(listIssuesCmd)
}
