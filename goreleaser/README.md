This is a static archive of [goreleaser
0.26.1](https://github.com/goreleaser/goreleaser/releases).

## Why?

Security. Since `go-ethereum` relies on `goreleaser` to build distrubtion
releases, it's important to mitigate any potential breach of security
on the `goreleaser` side, and refrain from `$ curl`ing the latest
release directly from a repository not directly connected with
`ethereumproject`.
