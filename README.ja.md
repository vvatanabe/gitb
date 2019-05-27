# gitb

Backlogのgitを快適に使うためのコマンドラインツール。

## 概要

`gitb`コマンドはBacklogのgitを快適に使うのに役立ちます。たとえば、PR、issue、branch、tagsなどをブラウザで1回のアクションで開くことができます。

接尾辞Bには複数の意味があります。Backlog、Browser、Bダッシュ（
素早くスーパーマリオを移動するような）。

## インストール

ビルドされたバイナリはGithubのリリースで利用可能です:  
https://github.com/vvatanabe/gitb/releases

このパッケージはgo getコマンドでもインストールできます:

`$ go get github.com/vvatanabe/gitb`

## 使い方

```
使い方:
   gitb <command>

コマンド:
     pr         現在のブランチに関連するプルリクエストページを開く
     ls-pr      プルリクエスト一覧ページを開く
     add-pr     現在のブランチでプルリクエストを作成するページを開きます
     issue      現在のブランチに関連したissueページを開く
     add-issue  現在のプロジェクトにissueを作成するためのページを開く
     ls-branch  現在のリポジトリのブランチ一覧ページを開く
     ls-tag     現在のリポジトリのタグリストページを開く
     tree       現在のブランチのツリーページを開く
     log        現在のブランチのコミットログページを開く
     ls-repo    現在のプロジェクトのリポジトリ一覧ページを開く
     help, h    コマンドのヘルプを表示する
```

## バグとフィードバック

バグ、質問、ディスカッションについてはGithub Issuesを利用してください。

## ライセンス

[MIT License](http://www.opensource.org/licenses/mit-license.php)

## 著者

[vvatanabe](https://github.com/vvatanabe)