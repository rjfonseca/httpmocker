# HTTP Mocker

A simple HTTP mocking tool. It allows for serving simple text or structured
files as responses of a HTTP request.

## What is this project about?

- A toy project;
- A simple HTTP mock tool for helping me with my daily coding;
- No external imports, only plain and simple go packages;
- Idiomatic (I try);

## How to install

You can `go get` this package:

``` shell
go get github.com/rjfonseca/httpmocker
```


## How to use it

Use `-h` to see all available options.

It serves a directory as mocked responses (it iterates the files on each
request). The received request and returned responses are saved in a
`output-<timestamp>` directory.

