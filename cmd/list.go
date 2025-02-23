/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v69/github"
	"github.com/spf13/cobra"
	"github.com/zmoog/ws/feedback"
)

var (
	reason string
	repo   string
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List notifications",
	Long:  `List notifications from GitHub.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")

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
}

func filter(notifications []*github.Notification) []*github.Notification {
	if repo == "" && reason == "" {
		return notifications
	}

	filtered := []*github.Notification{}
	for _, notification := range notifications {
		if *notification.Repository.Name == repo || *notification.Reason == reason {
			filtered = append(filtered, notification)
		}
	}

	return filtered
}
