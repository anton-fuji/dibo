<h1 align="center"> dibo </h1>
<p align="center">
  <img src="https://github.com/anton-fuji/dibo/blob/main/img/dibo-icon.png" alt="dibo icon" width="600">
</p>

<p align="center">
<a href="./LICENSE"><img src="https://img.shields.io/github/license/anton-fuji/dibo?style=flat-square" alt="LICENSE"></a>
</p>


`dibo` (dockerignore boilerplates) is a simple CLI tool written in Go for generating `.dockerignore` files from reusable templates.

It is inspired by [**gibo**](https://github.com/simonwhitaker/gibo?tab=readme-ov-file), but focused specifically on Docker.
With `dibo`, you can quickly create a `.dockerignore` tailored to your project by combining multiple templates.

# Features
🚀 Fast and simple — Just one command to generate .dockerignore
📦 Built-in templates — Common languages and frameworks included
🔧 Flexible — Combine multiple templates as needed
💡 [`fzf`](https://github.com/junegunn/fzf) integration — Interactive template selection (recommended)

# Installation

Preparing...


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
