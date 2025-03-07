/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package notifications

import (
	"github.com/spf13/cobra"
)

var (
	all              bool
	participating    bool
	sinceDaysAgo     int
	beforeDaysAgo    int
	reason           string
	repo             string
	excludeRepo      string
	owner            string
	unseen           bool
	notificationType string
)

// notificationsCmd represents the notifications command
var notificationsCmd = &cobra.Command{
	Use:   "notifications",
	Short: "List notifications",
	Long:  `List notifications from GitHub.`,
}

func Cmd() *cobra.Command {
	return notificationsCmd
}
