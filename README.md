# CNB Bootstrapper

A tool to bootstrap [packit](https://github.com/paketo-buildpacks/packit) compliant CNBs.

## Usage

```bash
$ export BUILDPACK_NAME=myCNB
$ echo "buildpack: ${BUILDPACK_NAME}" > config.yml
$ go run bootstrapper/main.go
$ ls -la /tmp/myCNB 
```
