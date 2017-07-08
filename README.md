Janus is a reusable tool for versioning and deploying builds to Google Cloud Provider (GCP) Storage from the CI
environment.

## Install
> - [ ] TODO: some slick `curl` command to grab the latest os-specific binary from `janus` releases to use in arbitrary CI

In the meantime:

```shell
$ go get github.com/ethereumproject/janus/...
$ cd $GOPATH/src/github.com/ethereumproject/janus
$ go build -o janus main.go
$ mv janus [path/to/project/]
```

## Usage
Janus has two subcomands: `deploy` and `version`.
#### Example
This repo is it's own example! See [.travis.yml](https://github.com/ethereumproject/janus/blob/master/.travis.yml) and [deploy.sh](https://github.com/ethereumproject/janus/blob/master/deploy.sh)


#### Deploy
Janus can use an encrypted _or_ decrypted `.json` GCP service key file. In case of an _encrypted_ JSON key file, Janus will attempt to decrypt it using `openssl`,
and depends on an __environment variable `GCP_PASSWD`__ to be set.


```shell
$ janus deploy -bucket builds.etcdevteam.com -object go-ethereum/v3.5.x/geth-linux-xxx.zip -file geth-linux-xxx.zip -key gcloud-service-encrypted-or-decrypted.json
> Deploying...
```

| flag | use |
| --- | --- |
| `-bucket` | eg `builds.etcdevteam.com`|
| `-object` | location (path) in which to store uploaded file(s) |
| `-file` | file(s) to upload |
| `-key` | encrypted _or_ decrypted JSON GCP service key file |



#### Version
`version` uses some variant of `git describe` or `git rev-list` to produce a
version number, as defined by `-format`
```shell
$ janus version -format v%M.%m.%P+%C-%S
> v3.5.0+55-asdf123
```

where `-format` value takes the interpolated forms:
```txt
%M - major version
%m - minor version
%P - patch version
%C - commit count since last tag
%S - HEAD sha1
```

This makes __previously:__
```yml
- git describe --tags --always > version.txt
- sed -E 's/v([[:digit:]]+\.[[:digit:]]+)\.[[:digit:]]-([[:digit:]]+)-g([a-f0-9]+)/v\1.\2+\3/' version.txt > version-app.txt
- sed -E 's/v([[:digit:]]+\.[[:digit:]]+)\.[[:digit:]]-([[:digit:]]+).+/v\1.\2/' version.txt > version-only.txt
- sed -E 's/v([[:digit:]]+\.[[:digit:]]+)\.[[:digit:]]-([[:digit:]]+).+/v\1.x/' version.txt > version-base.txt

```

__becomes:__

| sed output (.txt) | format syntax |
| --- | --- |
| `version-base.txt` | `--format v%M.%m.x` |
| `version-app.txt` | `--format v%M.%m.%P+%C-%S` |

----

> In ancient Roman religion and myth, Janus (/ˈdʒeɪnəs/; Latin: Iānus, pronounced [ˈjaː.nus]) is the god of beginnings, gates, transitions, time, duality, doorways,[1] passages, and endings.
- https://en.wikipedia.org/wiki/Janus
