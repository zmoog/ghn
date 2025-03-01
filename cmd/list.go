/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"
	"slices"
	"strings"

	"github.com/google/go-github/v69/github"
	"github.com/spf13/cobra"
	"github.com/zmoog/ws/feedback"
)

var (
	reason           string
	repo             string
	owner            string
	unseen           bool
	notificationType string
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List notifications",
	Long:  `List notifications from GitHub.`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := &github.NotificationListOptions{}

		notifications, _, err := client.Activity.ListNotifications(context.Background(), opts)
		if err != nil {
			log.Fatal(err)
		}

		feedback.PrintResult(Result{
			Notifications: filter(notifications),
		})

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	listCmd.Flags().StringVarP(&repo, "repo", "r", "", "Show notifications for a specific repository")
	listCmd.Flags().StringVarP(&reason, "reason", "R", "", "Show notifications for a specific reason")
	listCmd.Flags().StringVarP(&owner, "owner", "o", "", "Show notifications for a specific owner")
	listCmd.Flags().StringVarP(&notificationType, "type", "t", "", "Show notifications for a specific subject type")
	listCmd.Flags().BoolVarP(&unseen, "unseen", "u", false, "Show only unseen notifications")
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
