package main

import (
	"encoding/json"
	"errors"
	"github.com/mitchellh/go-homedir"
	"os"
	"path"
)

type Config struct {
	GithubToken string `json:"githubToken"`

	Repos []struct {
		Name      string `json:"name"`
		MergeInto string `json:"mergeInto"`
	} `json:"repos"`
}

var (
	ErrMissingConfigFile  = errors.New("couldnt find config file")
	ErrMissingGithubToken = errors.New("missing githubToken in config")
)

func readConfiguration(providedPath string) (conf Config, err error) {

	if dir, e := homedir.Dir(); e == nil {
		expandedPath := path.Join(dir, providedPath)
		if fConf, e := os.Open(expandedPath); e == nil {
			defer fConf.Close()
			err = json.NewDecoder(fConf).Decode(&conf)
		}
	} else {
		err = ErrMissingConfigFile
		return
	}

	if conf.GithubToken == "" {
		err = ErrMissingGithubToken
		return
	}

	return
}

func (conf Config) getRepoMergeInto(namespace, repo, defaulted string) string {

	qualifier := namespace + ":" + repo

	for _, repo := range conf.Repos {
		if repo.Name == qualifier && repo.MergeInto != "" {
			return repo.MergeInto
		}
	}

	return defaulted
}
