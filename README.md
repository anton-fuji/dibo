# dibo
`dibo` (dockerignore boilerplates) is a simple CLI tool written in Go for generating `.dockerignore` files from reusable templates.

It is inspired by [**gibo**](https://github.com/simonwhitaker/gibo?tab=readme-ov-file), but focused specifically on Docker.
With `dibo`, you can quickly create a `.dockerignore` tailored to your project by combining multiple templates.

# Installation

Preparing...


# usage
## List Templates
```sh
dibo list
```
Show a list of available templates.

## Dump a template
```sh
dibo dump <template>
```
Output the contents of a specified template to stdout.


## Generate .dockerignore
```sh
dibo init <template> [<template>...]
```
Generate a .dockerignore by combining templates.

## typical usage
if you want to generate .dockerignore for go.
```
dibo dump go >> .dockerignore
```


