/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package issues

import (
	"github.com/spf13/cobra"
)

var (
	includeAllRepos bool
	owner           string
	repo            string
	state           string
	assignee        string
	creator         string
	daysAgo         int
	sort            string
)

// issuesCmd represents the issues command
var issuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "List and search issues",
	Long:  `List and search issues from GitHub.`,
}

func Cmd() *cobra.Command {
	return issuesCmd
}
