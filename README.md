<h1 align="center"> dibo </h1>
<p align="center">
  <img src="https://github.com/anton-fuji/dibo/blob/main/img/dibo-icon.png" alt="dibo icon" width="500">
</p>

<p align="center">
  <img src="https://img.shields.io/github/go-mod/go-version/anton-fuji/dibo?filename=go.mod&style=flat-square" alt="Go Version">
  <a href="https://github.com/anton-fuji/dibo/releases/latest"><img src="https://img.shields.io/github/v/release/anton-fuji/dibo?style=flat-square" alt="GitHub release (latest by date)"></a>
  <a href="https://github.com/anton-fuji/dibo/actions/workflows/ci.yml"><img src="https://img.shields.io/github/actions/workflow/status/anton-fuji/dibo/ci.yml?logo=github&style=flat-square" alt="GitHub Workflow Status"></a>
  <a href="https://goreportcard.com/report/github.com/anton-fuji/dibo"><img src="https://goreportcard.com/badge/github.com/anton-fuji/dibo?style=flat-square" alt="Go Report Card"></a>
  <a href="./LICENSE"><img src="https://img.shields.io/github/license/anton-fuji/dibo?style=flat-square" alt="LICENSE"></a>
</p>


`dibo` (dockerignore boilerplates) is a simple CLI tool written in Go for generating `.dockerignore` files from reusable templates.

It is inspired by [**gibo**](https://github.com/simonwhitaker/gibo?tab=readme-ov-file), but focused specifically on Docker.
With `dibo`, you can quickly create a `.dockerignore` tailored to your project by combining multiple templates.

# Features
🚀 Fast and simple — Just one command to generate .dockerignore
<br>
📦 Built-in templates — Common languages and frameworks included
<br>
🔧 Flexible — Combine multiple templates as needed (duplicate patterns are removed automatically)
<br>
🔒 Secrets template — Keep credentials, keys, and `.env` files out of your image
<br>
🔠 Case-insensitive — `dibo init go` and `dibo init Go` both work
<br>
🛡️ Safe by default — Won't overwrite an existing `.dockerignore` unless you ask
<br>
⌨️ Shell completion — bash / zsh / fish / powershell
<br>
💡 [`fzf`](https://github.com/junegunn/fzf) integration — Interactive template selection (recommended)
 
# Installation
### Using Go 
```sh
go install github.com/anton-fuji/dibo@latest
```
 
### Using Homebrew
```sh
brew tap anton-fuji/tap
brew install dibo
```
 
### Using Nix (Flakes)
You can run `dibo` directly or install it using Nix Flakes.
 
#### Run without installing
```sh
nix run github:anton-fuji/dibo -- list
```
 
#### Install to your profile
```
nix profile install github:anton-fuji/dibo
```
 
#### Use in a temporary shell
```
nix shell github:anton-fuji/dibo
```
 
 
# usage
Template names are **case-insensitive** — `go`, `Go`, and `GO` all resolve to the same template.
 
## Generate .dockerignore
Generate `.dockerignore` for your project. Multiple templates are merged and duplicate patterns are removed automatically.
```sh
# Single template
dibo init Go
 
# Multiple templates
dibo init Go Node
 
# Keep secrets out of your image
dibo init Go Secrets
```
 
By default `init` will **not** overwrite an existing `.dockerignore`. Use a flag to change that:
 
| Flag | Description |
| --- | --- |
| `-f`, `--force` | Overwrite the file if it already exists |
| `-a`, `--append` | Append to the file instead of overwriting |
| `-o`, `--output <path>` | Write to a different path (default: `.dockerignore`) |
 
```sh
dibo init Go --force
dibo init Node --append
dibo init Go -o build/.dockerignore
```
 
## List Templates
Show a list of available templates.
```sh
dibo list
```
 
## Search Templates
Filter templates by keyword (case-insensitive substring match).
```sh
dibo search ru     # -> Ruby, Rust
```
 
## Dump a template
Output the contents of a specified template to stdout.
```sh
dibo dump <template> [<template>...]
```
 
## Shell completion
`dibo` ships with completion for bash, zsh, fish, and powershell. Template names are completed automatically.
```sh
# zsh (current shell)
source <(dibo completion zsh)
 
# bash (persist)
dibo completion bash > /etc/bash_completion.d/dibo
```
Run `dibo completion --help` for instructions for your shell.
 
## typical usage
If you want to generate `.dockerignore` for Go:
```sh
dibo init Go
```
 
# Recommend : Use with fzf
For an interactive experience, combine dibo with fzf.
```sh
dibo list | fzf
 
# Interactive single selection
dibo dump $(dibo list | fzf) >> .dockerignore
```

# LICENSE
[MIT](https://github.com/anton-fuji/dibo/blob/main/LICENSE)
