# gitb [![Build Status](https://travis-ci.org/vvatanabe/gitb.svg?branch=master)](https://travis-ci.org/vvatanabe/gitb) [![Coverage Status](https://coveralls.io/repos/github/vvatanabe/gitb/badge.svg?branch=master)](https://coveralls.io/github/vvatanabe/gitb?branch=master)

Backlogのgitを快適に使うためのコマンドラインツール。

## 概要

`gitb`コマンドはBacklogのgitを快適に使うのに役立ちます。たとえば、PR、issue、branch、tagsなどをブラウザで1回のアクションで開くことができます。

`gitb`コマンドはすべての`git`コマンドをラップしているので、 `gitb fetch`、`gitb pull`、 `gitb push`のように` gitb`コマンドを使って`git-command`を実行することができます。

接尾辞Bには複数の意味があります。Backlog、Browser、Bダッシュ。

## インストール

### Homebrew

MacOSで使用可能なパッケージマネージャであるHomebrewでインストールできます。

```
$ brew tap vvatanabe/gitb
$ brew install gitb
```

### Go

Go言語(go1.13+)をインストールしていれば、go getコマンドでもインストールできます。

```
$ go get github.com/vvatanabe/gitb
```

### GitHub Release Page

ビルドされたバイナリはGithubのリリースで利用可能です:  
https://github.com/vvatanabe/gitb/releases

## 使い方

### プルリクエスト

現在のリポジトリに対するBacklogのプルリクエストに関連するコマンドです。

__COMMANDS:__

`gitb pr [-s <STATE>]`

&emsp;現在のリポジトリのプルリクエスト一覧ページを開きます。

`gitb pr show [<PR-ID>]`

&emsp;指定した`<PR-ID>`のプルリクエストのページを開きます。`<PR-ID>`を指定しない時は、現在のブランチに関連したプルリクエストのページを開きます。

`gitb pr add [-b <BASE>]`

&emsp;現在のブランチでプルリクエストを追加するページを開きます。

`gitb pr blame [git blame command options] <PATH>`

&emsp;指定した`<PATH>`の変更に関連するプルリクエストIDを行単位で表示します。`git blame`コマンドのオプションを適用できます。

__OPTIONS:__

`-s, --state <STATE>`

&emsp;STATEでプルリクエストをフィルタリングします。値: "open" (初期値), "closed", "merged", "all".

`-b, --base <BASE>`

&emsp;BASEはプルリクエストのベースとなるブランチ名です。デフォルトは空です。

### 課題

現在のリポジトリに対するBacklogの課題に関連するコマンドです。

__COMMANDS:__

`gitb issue [-s <STATE>]`

&emsp;現在のプロジェクトの課題一覧ページを開きます。

`gitb issue show`

&emsp;現在のブランチに関連する課題ページを開きます。

`gitb issue add`

&emsp;現在のプロジェクトに課題を追加するページを開きます。

__OPTIONS:__

`-s, --state <STATE>`

&emsp;STATEで課題をフィルタリングします。 値: "all", "open", "in_progress", "resolved", "closed", "not_closed" (初期値).

### Browse

現在のリポジトリに関するGitページ（ブランチ、ツリー、タグ等）を開きます。

__COMMANDS:__

`gitb browse branch`

&emsp;現在のリポジトリのブランチ一覧ページを開きます。

`gitb browse tag`

&emsp;現在のリポジトリのタグ一覧ページを開きます。

`gitb browse tree`

&emsp;現在のブランチのツリーページを開きます。

`gitb browse history`

&emsp;現在のブランチの履歴ページを開きます。

`gitb browse network`

&emsp;現在のブランチのネットワークページを開きます。

`gitb browse repo`

&emsp;現在のプロジェクトのリポジトリ一覧ページを開きます。

## エイリアス

`gitb <command>`を`git <command>`として使いたい場合は、.XXXrc（.bashrc、.zshrc、config.fish）に以下のエイリアスを書いてください。

### Bash, Zsh

```
function git(){
  gitb "$@"
}
```

### Fish

```
function git
  gitb $argv
end
```

## 謝辞

- Inspired by [github.com/github/hub](https://github.com/github/hub)
- `gitb pr blame` is a Golang port of [kazuho/git-blame-pr.pl](https://gist.github.com/kazuho/eab551e5527cb465847d6b0796d64a39)


## バグとフィードバック

バグ、質問、ディスカッションについてはGithub Issuesを利用してください。

## ライセンス

[MIT License](http://www.opensource.org/licenses/mit-license.php)

## 著者

[vvatanabe](https://github.com/vvatanabe)