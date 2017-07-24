Janus is a tool for versioning and deploying builds to Google Cloud Provider (GCP) Storage from the CI
environment.

## Install

#### CI System Requirements:
- [ ] __JSON GCP Service Account Key__, with access to GCP _Storage_ enabled.
- [ ] __CI environment variable `GCP_PASSWD`__ to be set if the key is encrypted.
- [ ] __openssl__ is required for key decryption. This is standard on Travis. AppVeyor may require that you add some extra things to your `PATH`, but you may not have to install anything extra.
- [ ] __gpg__ is required to verify the Janus binary. This is standard on Travis and AppVeyor.
- [ ] __rev__, __curl__, and a few other basic bash commands are required for the installer script. Standard on Travis, can be added to PATH for AppVeyor as per example below

#### Install Janus:

##### Travis
```shell
- curl -sL https://raw.githubusercontent.com/ethereumproject/janus/master/get.sh | bash
- export PATH=./janusbin:$PATH
```

##### AppVeyor
```shell
- set PATH=C:\msys64\mingw64\bin;C:\msys64\usr\bin\;%PATH%
- curl -sL https://raw.githubusercontent.com/ethereumproject/janus/master/get-windows.sh | bash
- set PATH=./janusbin;%PATH%
```

__Security note:__ The installer scripts `get.sh` and `get-windows.sh` will use GPG to verify the latest Janus release binary against
the signing GPG key downloaded from a [specific commit at ethereumproject/volunteer](https://raw.githubusercontent.com/ethereumproject/volunteer/7a78a94307d67a0b20e418568b7bccac83c3d143/Volunteer-Public-Keys/isaac.ardis%40gmail.com).
For an additional layer of security, you may use the provided installer script signatures (`./*.sig`) to verify the installer script itself before using Janus
to deploy from your CI build. For maximum security, use a locally tracked version of [the signing key](https://raw.githubusercontent.com/ethereumproject/volunteer/7a78a94307d67a0b20e418568b7bccac83c3d143/Volunteer-Public-Keys/isaac.ardis%40gmail.com)
in your own repo. Alternatively, you can mimic the installer script itself, and use `curl` to download the key from the specific commit as mentioned previously. The link is:

> https://raw.githubusercontent.com/ethereumproject/volunteer/7a78a94307d67a0b20e418568b7bccac83c3d143/Volunteer-Public-Keys/isaac.ardis%40gmail.com

In practice, this would look like:
```yml
 - curl -sLO https://raw.githubusercontent.com/ethereumproject/volunteer/7a78a94307d67a0b20e418568b7bccac83c3d143/Volunteer-Public-Keys/isaac.ardis%40gmail.com
 - gpg --import isaac.ardis@gmail.com
 - curl -sLO https://raw.githubusercontent.com/ethereumproject/janus/master/get.sh
 - curl -sLO https://raw.githubusercontent.com/ethereumproject/janus/master/get.sh.sig
 - gpg --verify get.sh.sig get.sh
 - chmod +x get.sh
 - bash get.sh
```

Note that if you implement this additional layer and the signing key changes, you'll need to update either your tracked version of the key or download link accordingly.

## Usage
Janus has two subcommands: `deploy` and `version`.

#### Deploy
Janus can use an encrypted _or_ decrypted `.json` GCP service key file. In case of an _encrypted_ JSON key file, Janus will attempt to decrypt it using `openssl`,
and depends on an __environment variable `GCP_PASSWD`__ to be set. After successfully decrypting the key and deploying the files, Janus will automatically destroy (`rm`) the decrypted key from the CI.

| flag | example | description |
| --- | --- | --- |
| `-to` | `builds.etcdevteam.com/go-ethereum/v3.5.x/`| bucket, followed by 'directory' in which to hold the uploaded archive |
| `-files` | `./dist/*.zip` | file(s) to upload, can use relative or absolute path and/or wildcard globbing |
| `-key` | `./gcloud-travis.enc.json` | encrypted or decrypted JSON GCP service key file |

```shell
$ janus deploy -to builds.etcdevteam.com/go-ethereum/v3.5.x/ -files ./dist/*.zip -key gcloud-service-encrypted-or-decrypted.json
> Deploying...
```

#### Version
`version` uses `git` subcommands to produce a
version number, as defined by `-format`

```shell
$ janus version -format='v%M.%m.%P+%C-%S'
> v3.5.0+55-asdf123
```

`-format=value` takes the interpolated forms:
```txt
%M, _M - major version
%m, _m - minor version
%P, _P - patch version
%B, _B - hybrid patch version: `(%P * 100) + %C`
%C, _C - commit count since last tag
%S, _S - HEAD sha1 (first 7 characters)
```
_Note_: you may use either `%M` or `_M` syntax to interpolate version variables, since escaping `%` in batch scripts is rather tricky.

So this:

| sed output (.txt) | format syntax |
| --- | --- |
| `version-base.txt` | `-format v%M.%m.x` |
| `version-app.txt` | `-format v%M.%m.%P+%C-%S` |

replaces this:
```yml
- git describe --tags --always > version.txt
- sed -E 's/v([[:digit:]]+\.[[:digit:]]+)\.[[:digit:]]-([[:digit:]]+)-g([a-f0-9]+)/v\1.\2+\3/' version.txt > version-app.txt
- sed -E 's/v([[:digit:]]+\.[[:digit:]]+)\.[[:digit:]]-([[:digit:]]+).+/v\1.\2/' version.txt > version-only.txt
- sed -E 's/v([[:digit:]]+\.[[:digit:]]+)\.[[:digit:]]-([[:digit:]]+).+/v\1.x/' version.txt > version-base.txt
```

## Examples
Please visit the [/examples directory](./examples) to find example Travis and AppVeyor configuration files, deploy script, and service key.


----

> In ancient Roman religion and myth, Janus (/ˈdʒeɪnəs/; Latin: Iānus, pronounced [ˈjaː.nus]) is the god of beginnings, gates, transitions, time, duality, doorways,[1] passages, and endings.
- https://en.wikipedia.org/wiki/Janus
