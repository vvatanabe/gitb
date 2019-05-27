# gitb

A command line tool for using Backlog's git comfortably.

## Description

`gitb` command helps to use Backlog's git comfortably. For example, can open PR, issue, branches, tags, etc in the browser with one action.

Suffix B has multiple meanings. Backlog, Browser, B-Dash (
Like a Move Super Mario quickly).

## Installation

Built binaries are available on Github releases:  
https://github.com/vvatanabe/gitb/releases

This package can be installed with the go get command too:

`$ go get github.com/vvatanabe/gitb`

## Usage

```
USAGE:
   gitb <command>

COMMANDS:
     pr         Open the pull request page related to current branch
     ls-pr      Open the pull request list page
     add-pr     Open the page to create pull request with current branch
     issue      Open the issue page related to current branch
     add-issue  Open the page to create issue in current repository's project
     ls-branch  Open the branch list page of current repository
     ls-tag     Open the tag list page of current repository
     tree       Open the tree page of current branch
     log        Open the commit log page of current branch
     ls-repo    Open the repository list page of current repository's project
     help, h    Shows a list of commands or help for one command
```

## Bugs and Feedback

For bugs, questions and discussions please use the Github Issues.

## License

[MIT License](http://www.opensource.org/licenses/mit-license.php)

## Author

[vvatanabe](https://github.com/vvatanabe)