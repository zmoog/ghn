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
	daysAgo          int
	beforeDaysAgo    int
	perPage          int
	reason           string
	repo             string
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
