package notifications

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

// listCmd represents the list command
var listNotificationsCmd = &cobra.Command{
	Use:   "list",
	Short: "List notifications",
	Long:  `List notifications from GitHub.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// avoid non-overlapping time ranges
		if sinceDaysAgo > 0 && beforeDaysAgo > 0 && sinceDaysAgo <= beforeDaysAgo {
			return fmt.Errorf("since-days-ago must be greater than before-days-ago")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create GitHub client: %w", err)
		}

		// https://docs.github.com/en/rest/activity/notifications?apiVersion=2022-11-28#list-notifications-for-the-authenticated-user
		now := time.Now()
		opts := &github.NotificationListOptions{
			All:           all,
			Participating: participating,
		}

		if sinceDaysAgo > 0 {
			opts.Since = now.Add(-time.Hour * 24 * time.Duration(sinceDaysAgo))
		}

		if beforeDaysAgo > 0 {
			opts.Before = now.Add(-time.Hour * 24 * time.Duration(beforeDaysAgo))
		}

		notifications, _, err := client.Activity.ListNotifications(context.Background(), opts)
		if err != nil {
			return fmt.Errorf("failed to list notifications: %w", err)
		}

		// Filter out notifications that don't match
		// criteria not available in the API
		filtered := filter(notifications)

		// highlight the number of notifications
		// filtered out and print the result
		feedback.Println(
			fmt.Sprintf("Found %d notifications (%d filtered out)",
				len(notifications),
				len(notifications)-len(filtered),
			),
		)
		feedback.PrintResult(Result{
			Notifications: filtered,
		})

		return nil
	},
}

func filter(notifications []*github.Notification) []*github.Notification {
	filtered := []*github.Notification{}
	for _, notification := range notifications {
		if repo != "" && !slices.Contains(strings.Split(repo, ","), *notification.Repository.Name) {
			continue
		}

		if reason != "" && !slices.Contains(strings.Split(reason, ","), *notification.Reason) {
			continue
		}

		if owner != "" && !slices.Contains(strings.Split(owner, ","), *notification.Repository.Owner.Login) {
			continue
		}

		if unseen && notification.LastReadAt != nil {
			continue
		}

		if notificationType != "" && !slices.Contains(strings.Split(notificationType, ","), *notification.Subject.Type) {
			continue
		}

		filtered = append(filtered, notification)
	}

	return filtered
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	listNotificationsCmd.Flags().StringVarP(&repo, "repo", "r", "", "Show notifications for a specific repository")
	listNotificationsCmd.Flags().StringVarP(&reason, "reason", "R", "", "Show notifications for a specific reason")
	listNotificationsCmd.Flags().StringVarP(&owner, "owner", "o", "", "Show notifications for a specific owner")
	listNotificationsCmd.Flags().StringVarP(&notificationType, "type", "t", "", "Show notifications for a specific subject type")
	listNotificationsCmd.Flags().BoolVarP(&unseen, "unseen", "u", false, "Show only unseen notifications")

	listNotificationsCmd.Flags().IntVarP(&sinceDaysAgo, "since-days-ago", "s", 30, "Show notifications from the last n days")
	listNotificationsCmd.Flags().IntVarP(&beforeDaysAgo, "before-days-ago", "b", 0, "Show notifications before the last n days")
	listNotificationsCmd.Flags().BoolVarP(&participating, "participating", "p", false, "If true, only shows notifications in which the user is directly participating or mentioned.")
	listNotificationsCmd.Flags().BoolVarP(&all, "all", "a", false, "If true, show notifications marked as read.")

	// Add the list command to the notifications command
	notificationsCmd.AddCommand(listNotificationsCmd)
}
