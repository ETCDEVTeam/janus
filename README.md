Janus aims to be a reusable tool for versioning and deploying builds from the CI
environment.

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

