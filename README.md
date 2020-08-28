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

## Generated buildpack

Bootstrapper generates the following buildpack:

```
├── build.go
├── buildpack.toml
├── build_test.go
├── detect.go
├── detect_test.go
├── go.mod
├── init_test.go
├── integration
│   ├── default_test.go
│   ├── init_test.go
│   └── testdata
│       └── default_app
│           └── my-app
├── integration.json
├── LICENSE
├── NOTICE
├── package.toml
├── run
│   └── main.go
└── scripts
    ├── build.sh
    ├── integration.sh
    ├── package.sh
    └── unit.sh
```

A good place to start is making the integration tests fail in a way that demonstrates the correct behavior,
and then working your way through detect and build.

