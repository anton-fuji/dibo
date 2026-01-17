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
🔧 Flexible — Combine multiple templates as needed
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


# usage
## Generate .dockerignore
Generate .dockerignore for your project
```sh
# Single template
dibo init Go

# Multiple templates
dibo init Go Node
```

## List Templates
Show a list of available templates.
```sh
dibo list
```

## Dump a template
Output the contents of a specified template to stdout.
```sh
dibo dump <template> [<template>...]
```

## typical usage
if you want to generate .dockerignore for Go.
```
dibo dump Go >> .dockerignore
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
