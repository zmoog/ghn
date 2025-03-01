package issues

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/google/go-github/v69/github"
	"github.com/spf13/cobra"
	"github.com/zmoog/ogh/internal/client"
	"github.com/zmoog/ws/feedback"
)

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
		owner := args[0]
		repo := args[1]

		client, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create GitHub client: %w", err)
		}

		opts := &github.IssueListByRepoOptions{
			State:  state,
			Since:  time.Now().Add(-time.Hour * 24 * time.Duration(daysAgo)),
			Sort:   sort,
			Labels: labels,
		}

		if creator != "" {
			opts.Creator = creator
		}

		// if assignee != "*" {
		// 	opts.Assignee = assignee
		// }

		issues, _, err := client.Issues.ListByRepo(context.Background(), owner, repo, opts)
		if err != nil {
			return fmt.Errorf("failed to list issues: %w", err)
		}

		feedback.Println(fmt.Sprintf("Found %d issues", len(issues)))
		feedback.PrintResult(IssueResult{
			Issues: issues,
			Owner:  owner,
			Repo:   repo,
		})

		return nil
	},
}

func init() {
	listByRepoCmd.Flags().StringVarP(&assignee, "assignee", "a", "*", "Assignee of the issues")
	listByRepoCmd.Flags().StringVarP(&creator, "creator", "c", "", "Creator of the issues")
	listByRepoCmd.Flags().StringVarP(&state, "state", "s", "open", "State of the issues")
	listByRepoCmd.Flags().StringVarP(&sort, "sort", "o", "updated", "Sort by (updated or created)")
	listByRepoCmd.Flags().IntVarP(&daysAgo, "days-ago", "d", 30, "Days ago to search for issues")
	listByRepoCmd.Flags().StringSliceVarP(&labels, "labels", "l", []string{}, "Labels to filter issues by")

	// Add the list command to the issues command
	issuesCmd.AddCommand(listByRepoCmd)
}
