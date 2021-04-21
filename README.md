# CNB Bootstrapper

A tool to bootstrap [packit](https://github.com/paketo-buildpacks/packit) compliant CNBs.

## Usage

- Run the following command:

```bash
$ go run cmd/bootstrapper/main.go --buildpack <organization>/<buildpack>
```

Your github repo will be: `github.com/<organization>/<buildpack>`


- You will find your packit compliant buildpack template in `/tmp/<buildpack>`

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
├── .packit
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

