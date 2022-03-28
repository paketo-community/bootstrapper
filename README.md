# CNB Bootstrapper

A command-line tool to bootstrap
[packit](https://github.com/paketo-buildpacks/packit/v2)-compliant CNBs.

`bootstrapper` comes with the following commands:
* help                    : Help about any command
* run                     : Bootstrap a packit-compliant buildpack
* version                 : Get the version of the bootstrapper

## Usage

#### As a CLI:
- Package the tool for use:
```
$ ./scripts/package.sh -v <version>
```

- Run the CLI:
```
$ ./build/bootstrapper-<os> run --buildpack-name <organization/buildpack> --output <output directory>
```
This will use the default `template-cnb` template to create a packit-compliant
buildpack complete with scripts, Github workflows, and runnable integration/unit test suites.

If you'd like to use your own template, you can pass the optional `--template`
argument to `bootstrapper run` with the path to the template directory. Please
note that the template directory must contain a `go.mod` file.

Your output buildpack repo will be: `github.com/<organization>/<buildpack>`,
discoverable at the given output directory.

## Generated buildpack

Bootstrapper generates the following buildpack:

```
.
├── .github
│   ├── .syncignore
│   ├── dependabot.yml
│   ├── labels.yml
│   └── workflows
│       ├── auto-merge.yml
│       ├── codeql-analysis.yml
│       ├── create-draft-release.yml
│       ├── lint.yml
│       ├── push-buildpackage.yml
│       ├── synchronize-labels.yml
│       ├── test-pull-request.yml
│       └── update-github-config.yml
├── .gitignore
├── LICENSE
├── NOTICE
├── build.go
├── build_test.go
├── buildpack.toml
├── detect.go
├── detect_test.go
├── go.mod
├── go.sum
├── init_test.go
├── integration
│   ├── default_test.go
│   ├── init_test.go
│   └── testdata
│       └── default_app
│           └── my-app
├── integration.json
├── run
│   └── main.go
└── scripts
    ├── .util
    │   ├── git.sh
    │   ├── print.sh
    │   ├── tools.json
    │   └── tools.sh
    ├── build.sh
    ├── integration.sh
    ├── package.sh
    └── unit.sh
```

A good place to start is making the integration tests fail in a way that demonstrates the correct behavior,
and then working your way through detect and build.
