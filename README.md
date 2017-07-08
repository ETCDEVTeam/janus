Janus aims to be a reusable tool for versioning and deploying builds from the CI
environment.

## Install
> - [ ] TODO: some slick `curl` command to grab the latest os-specific binary from `janus` releases to use in arbitrary CI

> In the meantime:

```shell
$ go get github.com/ethereumproject/janus/...
$ cd $GOPATH/src/github.com/ethereumproject/janus
$ go build -o janus main.go
$ mv janus [path/to/project/]
```

## Usage
Janus has two subcomands: `deploy` and `version`.
#### Deploy
```shell
$ janus deploy -bucket builds.etcdevteam.com -object go-ethereum/v3.5.x/geth-linux-xxx.zip -file geth-linux-xxx.zip -key gcloud-service-encrypted-or-decrypted.json
> Deploying...
```

#### Version
`version` uses some variant of `git describe` or `git rev-list` to produce an applicable
version number, as defined by `-format`
```shell
$ janus version -format v%M.%m.%P+%C-%S
> v3.5.0+55-asdf123
```

`-format` value takes the interpolated forms:
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

### Prospective commands

```shell
$ janus deploy --from ./dist/ --mask ".+\.(dmg|zip|deb)" --bucket builds.etcdevteam.com --project emerald-wallet
> Deployed!

$ janus app-version
> 0.2.1+13

...

```

### Janus should replace
_.travis.yml / appveyor.yml_
```yml
- git describe --tags --always > version.txt
- sed -E 's/v([[:digit:]]+\.[[:digit:]]+)\.[[:digit:]]-([[:digit:]]+)-g([a-f0-9]+)/v\1.\2+\3/' version.txt > version-app.txt
- sed -E 's/v([[:digit:]]+\.[[:digit:]]+)\.[[:digit:]]-([[:digit:]]+).+/v\1.\2/' version.txt > version-only.txt
- sed -E 's/v([[:digit:]]+\.[[:digit:]]+)\.[[:digit:]]-([[:digit:]]+).+/v\1.x/' version.txt > version-base.txt

```

_deploy-gcs.sh_
```shell
./gcs-deploy-$TRAVIS_OS_NAME -bucket builds.etcdevteam.com -object go-ethereum/$(cat version-base.txt)/geth-classic-$TRAVIS_OS_NAME-$(cat version-app.txt).zip -file geth-classic-$TRAVIS_OS_NAME-$(cat version-app.txt).zip -key .gcloud.json
```

### Specs
- Janus will deploy to Google Cloud Provider.
- Janus should be (re)usable in OSX, Linux, and Windows operating systems.

----

> In ancient Roman religion and myth, Janus (/ˈdʒeɪnəs/; Latin: Iānus, pronounced [ˈjaː.nus]) is the god of beginnings, gates, transitions, time, duality, doorways,[1] passages, and endings.
- https://en.wikipedia.org/wiki/Janus
