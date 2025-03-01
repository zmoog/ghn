/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package pulls

import (
	"github.com/spf13/cobra"
)

var (
	// includeAllRepos bool
	state string
	// assignee        string
	// creator         string
	// daysAgo         int
	sort string
	// filter          string
	// labels          []string
)

// pullsCmd represents the pulls command
var pullsCmd = &cobra.Command{
	Use:   "pulls",
	Short: "List and search pull requests",
	Long:  `List and search pull requests from GitHub.`,
}

func Cmd() *cobra.Command {
	return pullsCmd
}
