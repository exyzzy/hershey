# Hershey Fonts

Decoding and using Hershey vector fonts in Go

Current code matches this Medium article: Hershey Fonts in Go

To install:

    go get github.com/exyzzy/hershey

To test & generate decoded files:

    go test

For examples of:

    Decoding Hershey: see generate.go

    Package use: see examples.go

Quick Start

    package main

    import "github.com/exyzzy/hershey"

    func main() {
        hershey.DrawAllFontImage()
    }

