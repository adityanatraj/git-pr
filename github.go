package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	ErrPrNotLoading = errors.New("prs not loading")
	ErrPrNotFound   = errors.New("pr not found")
)

type GitTransaction struct {
	mergeInto string

	GitInfo
}

func NewTransactionFromInfo(gitInfo GitInfo, mergeInto string) GitTransaction {
	return GitTransaction{
		GitInfo:   gitInfo,
		mergeInto: mergeInto,
	}
}

type GithubClient struct {
	context context.Context
	client  *github.Client
}

func NewGithubClient(token string) *GithubClient {
	client := &GithubClient{
		context: context.Background(),
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tokenClient := oauth2.NewClient(client.context, tokenSource)

	client.client = github.NewClient(tokenClient)

	return client
}

func (gc GithubClient) getPR(gitInfo GitInfo) (*github.PullRequest, error) {

	prs, _, err := gc.client.PullRequests.List(gc.context, gitInfo.namespace, gitInfo.repo, &github.PullRequestListOptions{})
	if err != nil {
		return nil, ErrPrNotLoading
	}

	for _, pr := range prs {
		if *(pr.Head.Ref) == gitInfo.branch {
			return pr, nil
		}
	}

	return nil, ErrPrNotFound
}

func (gc GithubClient) updatePR(gitInfo GitInfo, pr *github.PullRequest, title string, message string) (*github.PullRequest, error) {

	if title != "" {
		pr.Title = &title
	}

	if message != "" {
		pr.Body = &message
	}

	pr, _, err := gc.client.PullRequests.Edit(gc.context, gitInfo.namespace, gitInfo.repo, *pr.Number, pr)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (gc GithubClient) createPR(transaction GitTransaction, title string, message string) (*github.PullRequest, error) {

	if title == "" {
		title = transaction.GitInfo.generateTitle()
	}

	if message == "" {
		message = getTemplateMessage()
	}

	pr, _, err := gc.client.PullRequests.Create(gc.context, transaction.namespace, transaction.repo, &github.NewPullRequest{
		Title:               &title,
		Body:                &message,
		Head:                github.String(fmt.Sprintf("%s:%s", transaction.namespace, transaction.branch)),
		Base:                &transaction.mergeInto,
		MaintainerCanModify: github.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (gc GithubClient) getReviewStatus(pr *github.PullRequest) {
	// TODO: implement this function to print out some details about the PR
	// possibly: {list of changed files}, {assignees + status of their review}, {# current comments} ?
	fmt.Println("Sorry, not yet implemented")
	return
}
