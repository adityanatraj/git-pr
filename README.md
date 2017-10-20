# Git PR

`git-pr` is a custom command for the `git` tool that helps creating and editing
pull requests. It is meant to be extremely simple to use. If you'd like to see
a more powerful tool, though limited to github, you should definitely check out
the very awesome [hub tool](github.com/github/hub).

Since it's named `git-pr`, after installation you'll interface with the command
directly through `git` itself. 

```
$ git pr -t "My New PR" -m "this pr is to show off stuff"`
pr created at: https://github.com/adityanatraj/git-pr/pulls/0
```

# Basics

The are really only two things of particular note:

1. when not provided, PR titles are set to "cleansed" branch names
    - ex. this-new-feature => "This new feature"
    - ex. fix/PLTO-1049-do-some-real-work => "[PLTO-1049] Do some real work"
    - ex. i_did_cool-stuff => "I did cool stuff"
2. currently the tool works for a branching (not forking) model 
    - you CAN pick which branch you're merging into though (`-b` flag or via config)

# Installation

## Building it

You'll need to have `go` installed and and an [environment set up](https://golang.org/doc/code.html).
Once you've done that, you can `go get` it and it will be built in `$GOPATH/bin`. 

```
go get github.com/adityanatraj/git-pr
```

As long as you've followed the environment instructions and added `$GOPATH/bin` to your `$PATH`,
you're all set for configuration and use!

## Configuration

`git-pr` expects a `JSON`-encoded configuration file that, by default, is located at `$HOME/.config/git-pr`. You
can change this expectation by using the `-c --config [path]` flag on the command line.

### Github Token

Since it interacts with Github via the API, you'll need to get an API authentication token and put it
in your configuration file.

```
{
  "githubToken": "your github token"
}
```

### Per-repo config

Consider a situation where you don't maintain a `master` branch but rather a `develop`
branch that everything gets merged into. The default branch to merge into is set internally in `git-pr` as `master`. Oops.
Well, you could use the `-b --branch [branchname]` command line flag but this can quickly become tiresome when it's
all the time. Enter per-repo configuration.

```
{
 ...
 "repos": [
   {
      "name": "adityanatraj:git-pr",
      "mergeInto": "develop"
   }
 ]
}
```

Fields:
  - `name`: the `<owner>:<repo-name>` of the project you'd like to observe these settings
  - `mergeInto`: the branch you'd like PRs to merge into
  
_Note_: command line arguments over-ride any default you've set.

# Usage 

Here is the direct output of `git pr -h`. 

```
Usage of git-pr:
  -b, --branchInto string   branch to merge into (default "master")
  -c, --config string       $HOME-relative path to config file (default ".config/git-pr")
  -d, --details             output details of PR (if already exists)
  -m, --message string      message body for the PR
  -t, --title string        override title generation with this
```

# To Do
- [ ] make it work for bitbucket
- [ ] implement showing details on a pr once made
- [ ] allow more flexible titling
- [ ] tests

# MIT License
Copyright 2017 adityanatraj

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
