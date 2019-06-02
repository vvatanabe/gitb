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

### Pull Request

Related to Backlog Pull Requests for the current repository

__COMMANDS:__

```
gitb pr [-s <STATE>]     Open the pull request list page in the current repository
gitb pr show             Open the pull request page related to the current branch
gitb pr add [-b <BASE>]  Open the page to create pull request with the current branch
```

__OPTIONS:__

```
-s, --state <STATE>

     Filter pull requests by STATE. Values: "open" (default), "closed", "merged", "all".

-b, --base <BASE>

    BASE is base branch name. Default is empty.
```

### Issue

Related to Backlog Issues for the current repository

__COMMANDS:__

```
gitb issue [-s <STATE>]  Open the issue list page in the current project
gitb issue show          Open the issue page related to the current branch
gitb issue add           Open the page to create issue in the current project
```

__OPTIONS:__

```
-s, --state <STATE>

     Filter issues by STATE. Values: "all", "open", "in_progress", "resolved", "closed", "not_closed" (default)
```

### Other Commands

`gitb branch`

Open the branch list page in current repository

`gitb tag`

Open the tag list page in current repository

`gitb tree`

Open the tree page in current branch

`gitb history`

Open the history page in current branch

`gitb network`

Open the network page in current branch

`gitb repo`

Open the repository list page in current project

`gitb help, h`

Shows a list of commands or help for one command


## Bugs and Feedback

For bugs, questions and discussions please use the Github Issues.

## License

[MIT License](http://www.opensource.org/licenses/mit-license.php)

## Author

[vvatanabe](https://github.com/vvatanabe)