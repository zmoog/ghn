/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
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

// notificationsCmd represents the notifications command
var notificationsCmd = &cobra.Command{
	Use:   "notifications",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("notification called")
	},
}

// listCmd represents the list command
var listNotificationsCmd = &cobra.Command{
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
	rootCmd.AddCommand(notificationsCmd)
	notificationsCmd.AddCommand(listNotificationsCmd)

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
