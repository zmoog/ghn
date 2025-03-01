package pulls

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/google/go-github/v69/github"
	"github.com/spf13/cobra"
	"github.com/zmoog/ogh/internal/client"
	"github.com/zmoog/ws/feedback"
)

// listByRepoCmd represents the listByRepo command
var listByRepoCmd = &cobra.Command{
	Use:   "list-by-repo [owner] [repo]",
	Short: "List pull requests by repository",
	Long:  `List pull requests by repository from GitHub.`,
	Args:  cobra.ExactArgs(2),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		validSortValues := []string{"updated", "created", "popularity", "long-running"}
		if !slices.Contains(validSortValues, sort) {
			return fmt.Errorf("sort must be either [%s], got '%s'", strings.Join(validSortValues, ", "), sort)
		}

		validStateValues := []string{"open", "closed", "all"}
		if !slices.Contains(validStateValues, state) {
			return fmt.Errorf("state must be either [%s], got '%s'", strings.Join(validStateValues, ", "), state)
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

		opts := &github.PullRequestListOptions{
			State: state,
			Sort:  sort,
		}

		pulls, _, err := client.PullRequests.List(context.Background(), owner, repo, opts)
		if err != nil {
			return fmt.Errorf("failed to list pull requests: %w", err)
		}

		feedback.Println(fmt.Sprintf("Found %d pull requests", len(pulls)))
		feedback.PrintResult(PullRequestResult{
			Pulls: pulls,
			Owner: owner,
			Repo:  repo,
		})

		return nil
	},
}

func init() {
	listByRepoCmd.Flags().StringVarP(&state, "state", "s", "open", "State of the issues")
	listByRepoCmd.Flags().StringVarP(&sort, "sort", "o", "updated", "Sort by (updated or created)")

	// Add the list command to the pulls command
	pullsCmd.AddCommand(listByRepoCmd)
}
