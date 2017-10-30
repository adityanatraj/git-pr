package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrNotGit             = errors.New("fatal: Not a git repository (or any of the parent directories): .git")
	ErrUncommittedChanges = errors.New("fatal: you have uncommitted/unstashed changes. deal with them and try again")

	ErrNoOrigin  = errors.New("fatal: No remote named 'origin'. ")
	ErrBadRemote = errors.New("fatal: couldnt read 'git remote get-url origin' output")

	RemotePattern = regexp.MustCompile(`git@(?P<host>[a-zA-Z0-9._\-]+[.][a-z]+):(?P<namespace>[a-zA-Z0-9\-]+)/(?P<repo>.*)[.]git`)
	TitlePattern  = regexp.MustCompile(`(?:(?:[fF]eature|[fF]ix)\/[pP][lL][tT][oO]-)?(?P<number>[0-9]+)?-?(?P<raw_title>.*)`)
)

type GitInfo struct {
	host      string
	namespace string
	repo      string
	branch    string
}

func getGitInfo() (gitInfo GitInfo, err error) {

	if branchBytes, er := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output(); er != nil {
		err = ErrNotGit
		return
	} else {
		gitInfo.branch = string(bytes.TrimSpace(branchBytes))
	}

	output, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		err = ErrNotGit
		return
	} else if len(bytes.TrimSpace(output)) != 0 {
		err = ErrUncommittedChanges
		return
	}

	output, err = exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		err = ErrNoOrigin
		return
	}

	if matches := RemotePattern.FindAllStringSubmatch(string(output), -1); len(matches) == 0 {
		err = ErrBadRemote
		return
	} else {
		gitInfo.host, gitInfo.namespace, gitInfo.repo = matches[0][1], matches[0][2], matches[0][3]
	}

	return
}

func getTemplateMessage() string {
	output, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return ""
	}

	gitRoot := string(bytes.TrimSpace(output))
	pullRequestTemplatePath := path.Join(gitRoot, ".github/PULL_REQUEST_TEMPLATE.md")

	f, err := os.Open(pullRequestTemplatePath)
	if err != nil {
		return ""
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return ""
	}

	return string(bytes.TrimSpace(data))
}

func (gitInfo GitInfo) generateTitle() string {

	if matches := TitlePattern.FindAllStringSubmatch(gitInfo.branch, -1); len(matches) == 0 {
		return gitInfo.branch
	} else {

		var frontMatter string
		if matches[0][1] == "" {
			frontMatter = ""
		} else {
			frontMatter = fmt.Sprintf("[PLTO-%s] ", matches[0][1])
		}

		var title string
		breakReplacer := strings.NewReplacer("_", " ", "-", " ")
		title = breakReplacer.Replace(matches[0][2])

		firstRune, n := utf8.DecodeRuneInString(title)
		return fmt.Sprintf("%s%c%s", frontMatter, unicode.ToUpper(firstRune), title[n:])
	}
}
