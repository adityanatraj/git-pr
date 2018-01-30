package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"strings"
)

func check(err error, msg ...string) {
	if err != nil {
		if len(msg) == 0 {
			fmt.Println(err)
		} else {
			fullMsg := strings.Join(msg, "\n\t")
			fmt.Printf("%s\n", fullMsg)
		}
		os.Exit(1)
	}
}

const (
	DefaultConfigPath = ".config/git-pr"
	DefaultBranchInto = "master"
	DefaultTitle      = ""
	DefaultMessage    = ""
	DefaultIsDetailed = false
)

var commandArgs struct {
	configPath string
	branchName string
	title      string
	message    string
	isDetailed bool
}

func init() {
	pflag.StringVarP(&commandArgs.configPath, "config", "c", DefaultConfigPath, "$HOME-relative path to config file")
	pflag.StringVarP(&commandArgs.branchName, "branchInto", "b", DefaultBranchInto, "branch to merge into")
	pflag.StringVarP(&commandArgs.title, "title", "t", DefaultTitle, "override title generation with this")
	pflag.StringVarP(&commandArgs.message, "message", "m", DefaultMessage, "message body for the PR")
	pflag.BoolVarP(&commandArgs.isDetailed, "details", "d", DefaultIsDetailed, "output details of PR (if already exists)")

	pflag.Parse()
}

func main() {
	gitInfo, err := getGitInfo()
	check(err)

	config, err := readConfiguration(commandArgs.configPath)
	check(err)

	// only read config for mergeInto when the user DOES NOT command line set the value
	// note: because the unset default looks the same as the setting the default value
	// it will use the config value when manually using -b <DefaultValue>
	if commandArgs.branchName == DefaultBranchInto {
		commandArgs.branchName = config.getRepoMergeInto(gitInfo.namespace, gitInfo.repo, DefaultBranchInto)
	}

	githubClient := NewGithubClient(config.GithubToken)

	pr, err := githubClient.getPR(gitInfo)
	if err != nil {
		if err == ErrPrNotLoading {
			check(err)
		} else if err == ErrPrNotFound {
			transaction := NewTransactionFromInfo(gitInfo, commandArgs.branchName)
			pr, err := githubClient.createPR(transaction, commandArgs.title, commandArgs.message)
			check(err)

			fmt.Printf("updated pr is here: %s\n", *pr.HTMLURL)
		}
	} else {
		if commandArgs.title == "" && commandArgs.message == "" {
			if commandArgs.isDetailed {
				githubClient.getReviewStatus(pr)
			} else {
				fmt.Printf("pr already exists here: %s\n", *pr.HTMLURL)
			}
		} else {
			pr, err := githubClient.updatePR(gitInfo, pr, commandArgs.title, commandArgs.message)
			check(err)

			fmt.Printf("created pr here: %s\n", *pr.HTMLURL)
		}
	}
}
