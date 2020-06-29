# CNB Bootstrapper

A tool to bootstrap [packit](https://github.com/paketo-buildpacks/packit) compliant CNBs.

## Usage

- Edit the `config.yml` file included in this directory where `organization` is the name of your github org,
and `buildpack` is the name of the buildpack.

Your github repo will be: `github.com/<organization>/<buildpack>`

- Run the following command:

```bash
$ go run executer/main.go
```

- You will find your packit compliant buildpack template in `/tmp/<buildpack>`
