<h1 align="center"> dibo </h1>
<p align="center">
  <img src="https://github.com/anton-fuji/dibo/blob/main/img/dibo-icon.png" alt="dibo アイコン" width="500">
</p>

<p align="center">
  <a href="./README.md">English</a> · <strong>日本語</strong>
</p>

<p align="center">
  <img src="https://img.shields.io/github/go-mod/go-version/anton-fuji/dibo?filename=go.mod&style=flat-square" alt="Go バージョン">
  <a href="https://github.com/anton-fuji/dibo/releases/latest"><img src="https://img.shields.io/github/v/release/anton-fuji/dibo?style=flat-square" alt="最新リリース"></a>
  <a href="https://github.com/anton-fuji/dibo/actions/workflows/ci.yml"><img src="https://img.shields.io/github/actions/workflow/status/anton-fuji/dibo/ci.yml?logo=github&style=flat-square" alt="GitHub Actions ステータス"></a>
  <a href="https://goreportcard.com/report/github.com/anton-fuji/dibo"><img src="https://goreportcard.com/badge/github.com/anton-fuji/dibo?style=flat-square" alt="Go Report Card"></a>
  <a href="./LICENSE"><img src="https://img.shields.io/github/license/anton-fuji/dibo?style=flat-square" alt="ライセンス"></a>
</p>

`dibo`（dockerignore boilerplates）は、再利用可能なテンプレートから `.dockerignore` を生成するシンプルな Go 製 CLI ツールです。

[**gibo**](https://github.com/simonwhitaker/gibo) に着想を得ていますが、Docker に特化しています。複数のテンプレートを組み合わせて、プロジェクトに合った `.dockerignore` をすばやく作成できます。

# 特徴
🚀 高速でシンプル — 1つのコマンドで `.dockerignore` を生成
<br>
📦 組み込みテンプレート — よく使われる言語やフレームワークに対応
<br>
🔧 柔軟 — 必要なテンプレートを組み合わせて利用可能（重複パターンは自動削除）
<br>
🔒 Secrets テンプレート — 認証情報、鍵、`.env` ファイルをイメージから除外
<br>
🔠 大文字・小文字を区別しない — `dibo init go` と `dibo init Go` は同じように動作
<br>
🛡️ 安全なデフォルト — 明示的に指定しない限り既存の `.dockerignore` を上書きしない
<br>
⌨️ シェル補完 — bash / zsh / fish / powershell
<br>
💡 [`fzf`](https://github.com/junegunn/fzf) 連携 — インタラクティブなテンプレート選択（おすすめ）

# インストール
### Go を使う
```sh
go install github.com/anton-fuji/dibo@latest
```

### Homebrew を使う
```sh
brew install anton-fuji/tap/dibo
```

Homebrew 6 以降では、サードパーティ tap の利用に明示的な信頼設定が必要です。上記の完全修飾コマンドは `dibo` Formula のみを信頼します。

### Nix（Flakes）を使う
Nix Flakes を使って `dibo` を直接実行したり、インストールしたりできます。

#### インストールせずに実行
```sh
nix run github:anton-fuji/dibo -- list
```

#### プロファイルにインストール
```sh
nix profile install github:anton-fuji/dibo
```

#### 一時的なシェルで使う
```sh
nix shell github:anton-fuji/dibo
```

# 使い方
テンプレート名は **大文字・小文字を区別しません**。`go`、`Go`、`GO` はすべて同じテンプレートとして解決されます。

## `.dockerignore` を生成
プロジェクト用の `.dockerignore` を生成します。複数のテンプレートは自動的に結合され、重複パターンは削除されます。
```sh
# 単一テンプレート
dibo init Go

# 複数テンプレート
dibo init Go Node

# Secrets を除外
dibo init Go Secrets
```

デフォルトでは、`init` は既存の `.dockerignore` を上書きしません。動作を変更するにはフラグを使用します。

| フラグ | 説明 |
| --- | --- |
| `-f`, `--force` | 既存のファイルを上書きする |
| `-a`, `--append` | ファイル末尾に追加する |
| `-o`, `--output <path>` | 別のパスに出力する（デフォルト: `.dockerignore`） |

```sh
dibo init Go --force
dibo init Node --append
dibo init Go -o build/.dockerignore
```

### インタラクティブ選択

表示された番号またはテンプレート名を、カンマ区切りで入力して1つ以上のテンプレートを選択できます。

```sh
dibo init --interactive
# Available templates:
#   1. Common
#   2. Go
#   ...
# Select templates by number or name (comma-separated): 1, Go, Secrets
```

## プロジェクトを検出

ディレクトリ内のファイルを調べ、推奨するテンプレートの組み合わせを表示します。`--write` を付けると、そのディレクトリにファイルを生成します。

```sh
dibo detect
# Detected: Go
# Recommended templates: Common, Go, Secrets
# Create it with: dibo init Common Go Secrets

dibo detect ./api
# Create it with: dibo init Common Go Secrets --output api/.dockerignore

dibo detect ./api --write
```

`detect` は Go、Node.js、Python、Ruby、Rust、Java/Kotlin、PHP、.NET のプロジェクトを検出します。検出した言語のテンプレートに加えて、`Common` と `Secrets` も推奨します。

## `.dockerignore` をチェック

検出したプロジェクト種別に必要なビルド成果物と、基本的なシークレットファイルのパターンが現在の `.dockerignore` から除外されているか確認します。不足がある場合は終了コードが 0 以外になるため、CI にも利用できます。

```sh
dibo check
dibo check ./api
dibo check --file docker/.dockerignore
```

## テンプレートを表示
指定したテンプレートの内容を標準出力に表示します。
```sh
dibo dump <template> [<template>...]
```

## シェル補完
`dibo` は bash、zsh、fish、powershell の補完に対応しています。テンプレート名も自動補完されます。
```sh
# zsh（現在のシェル）
source <(dibo completion zsh)

# bash（永続化）
dibo completion bash > /etc/bash_completion.d/dibo
```
詳しくは `dibo completion --help` を実行してください。

## よくある使い方
Go プロジェクト向けの `.dockerignore` を生成する場合:
```sh
dibo init Go
```

# おすすめ: fzf と組み合わせる
インタラクティブに使うには、`dibo` と fzf を組み合わせます。
```sh
dibo list | fzf

# 1つ選択して生成
dibo dump $(dibo list | fzf) >> .dockerignore
```

## テンプレート一覧

| テンプレート | 説明 |
| --- | --- |
| `Common` | 共通の除外設定（VCS ディレクトリ、OS の不要ファイル、エディタのファイル、ログなど） |
| `Secrets` | 認証情報、鍵、証明書、`.env` ファイル |
| `Go` | Go プロジェクト |
| `Node` | Node.js / JavaScript |
| `Python` | Python |
| `Ruby` | Ruby / Rails |
| `Rust` | Rust |
| `Java` | Java / Kotlin |
| `PHP` | PHP |
| `dotNet` | .NET |

## 開発タスク

このリポジトリではタスクランナーとして [just](https://just.systems/) を使用しています。Nix 開発シェルに入り、`just` を実行すると利用可能なタスクを一覧表示できます。

```sh
nix develop
just check    # フォーマット、lint、テスト、ビルド
just run detect
```

## リリース

Release Please は `main` へのマージ後に実行されます。Conventional Commits のプレフィックスによって次のバージョンが決まり、`feat:` はマイナーリリース、`fix:` はパッチリリース、`chore:` はリリースを作成しません。Release PR をマージするとバージョンタグが作成され、GoReleaser が GitHub Release と Homebrew Formula を公開します。

GitHub Actions の secret `HOMEBREW_TAP_GITHUB_TOKEN` には、このリポジトリと `anton-fuji/homebrew-tap` の両方への書き込み権限が必要です。

# ライセンス
[MIT](https://github.com/anton-fuji/dibo/blob/main/LICENSE)
