Janus is a reusable tool for versioning and deploying builds to Google Cloud Provider (GCP) Storage from the CI
environment.

## Install
Janus is designed for use in the CI environment.

#### Requirements:
- [ ] encrypted JSON GCP service account key, with access to GCP _Storage_ feature
- [ ] CI environment variable `GCP_PASSWD` to be set, either via secure global (as below), or via CI GUI interface
- [ ] having Janus installed (as below), via `curl -sL https://raw.githubusercontent.com/ethereumproject/janus/master/get.sh | bash`

_eg_
```yml
# .travis.yml
env:
  global:
    # This value should hold at least environment variable GCP_PASSWD=xxx in order to decrypt the GCP service account key that Janus relies on.
    # eg.
    # $ travis encrypt -r ethereumproject/emerald-rs GCP_PASSWD=asdfasdfasdfasdfasdf
    - secure: "MjvfqrKakMa+z+6LFxaL30n+BtjxUm2BnJ6/+S5cbxo"

before_deploy:
  # Install Janus
  # This should work for Travis (OSX and Linux) and AppVeyor (Windows) environments.
  - curl -sL https://raw.githubusercontent.com/ethereumproject/janus/master/get.sh | bash

  # Prepare file(s) to deploy.
  - zip emerald-"$TRAVIS_OS_NAME"-$(janus version -format v%M.%m.%P+%C-%S).zip emerald-cli

  # Ship it.
  - janus deploy -to builds.etcdevteam.com/emerald/$(janus version -format v%M.%m.x) -files ./*.zip -key gcloud-travis.enc.json
```

## Usage
Janus has two subcommands: `deploy` and `version`.

#### Deploy
Janus can use an encrypted _or_ decrypted `.json` GCP service key file. In case of an _encrypted_ JSON key file, Janus will attempt to decrypt it using `openssl`,
and depends on an __environment variable `GCP_PASSWD`__ to be set. After successfully decrypting the key and deploying the files, Janus will automatically destroy (`rm`) the decrypted key from the CI.

| flag | example | description |
| --- | --- | --- |
| `-to` | `builds.etcdevteam.com/go-ethereum/3.5.x`| __<bucket>__/<path/to/hold/file(s)> |
| `-files` | `./dist/*.zip` | file(s) to upload, can use relative or absolute path and/or wildcard globbing |
| `-key` | `./gcloud-travis.enc.json` |encrypted _or_ decrypted JSON GCP service key file |

```shell
$ janus deploy -to builds.etcdevteam.com/go-ethereum/v3.5.x/ -files ./dist/*.zip -key gcloud-service-encrypted-or-decrypted.json
> Deploying...
```

#### Version
`version` uses some variant of `git describe` or `git rev-list` to produce a
version number, as defined by `-format`

```shell
$ janus version -format v%M.%m.%P+%C-%S
> v3.5.0+55-asdf123
```

`-format=value` takes the interpolated forms:
```txt
%M - major version
%m - minor version
%P - patch version
%C - commit count since last tag
%S - HEAD sha1
```

So this:

| sed output (.txt) | format syntax |
| --- | --- |
| `version-base.txt` | `--format v%M.%m.x` |
| `version-app.txt` | `--format v%M.%m.%P+%C-%S` |

replaces this:
```yml
- git describe --tags --always > version.txt
- sed -E 's/v([[:digit:]]+\.[[:digit:]]+)\.[[:digit:]]-([[:digit:]]+)-g([a-f0-9]+)/v\1.\2+\3/' version.txt > version-app.txt
- sed -E 's/v([[:digit:]]+\.[[:digit:]]+)\.[[:digit:]]-([[:digit:]]+).+/v\1.\2/' version.txt > version-only.txt
- sed -E 's/v([[:digit:]]+\.[[:digit:]]+)\.[[:digit:]]-([[:digit:]]+).+/v\1.x/' version.txt > version-base.txt
```

#### Example
This repo is it's own example! See [.travis.yml](https://github.com/ethereumproject/janus/blob/master/.travis.yml) and [deploy.sh](https://github.com/ethereumproject/janus/blob/master/deploy.sh)

----

> In ancient Roman religion and myth, Janus (/ˈdʒeɪnəs/; Latin: Iānus, pronounced [ˈjaː.nus]) is the god of beginnings, gates, transitions, time, duality, doorways,[1] passages, and endings.
- https://en.wikipedia.org/wiki/Janus
