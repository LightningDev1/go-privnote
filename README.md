# go-privnote

go-privnote is a Go client library for creating and reading notes on [Privnote](https://privnote.com/). It bypasses Cloudflare bot detection by using a [TLS client](https://github.com/bogdanfinn/tls-client) for TLS fingerprinting.

[![Reference](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/LightningDev1/go-privnote)
[![Linter](https://goreportcard.com/badge/github.com/LightningDev1/go-privnote?style=flat-square)](https://goreportcard.com/report/github.com/LightningDev1/go-privnote)
[![Build status](https://github.com/LightningDev1/go-privnote/actions/workflows/ci.yml/badge.svg)](https://github.com/LightningDev1/go-privnote/actions)

```go
client := privnote.NewClient()

noteLink, err := client.CreateNote(privnote.CreateNoteData{
    Data: "Hello, World!",
})

noteContent, err := client.ReadNoteFromLink("https://privnote.com/note-id#password")

noteContent, err := client.ReadNoteFromID("note-id", "password")
```

## Installation

```bash
go get github.com/LightningDev1/go-privnote
```

## Usage

See [example/main.go](./example/main.go) for a complete example program.

