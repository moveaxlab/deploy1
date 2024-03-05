# deploy1

`deploy1` (pronounced /dəplɔɪɒnɛ/) is a tool to build and deploy docker images, with strong opinions on image tags.

## Requirements

- `go 1.19` or newer
- `docker`
- `argocd 1.4`

Make sure that your `GOPATH` is set correctly, and that `$GOPATH/bin` is included in your `PATH`.
You can add this to your `.bashrc` or `.zshrc`:

```bash
export GOPATH=$(go env GOPATH)
export PATH="${PATH}:$(go env GOPATH)/bin"
```

More info [here](https://golang.org/doc/gopath_code.html).

## Installation

Download and install `deploy1`:

```bash
go install github.com/moveaxlab/deploy1@latest
```

## Configuration

All configuration happens inside a `deploy1.json` in the root of your project.

### Environments

You can have multiple environments for your registry and deployment configuration.
Environment specific configuration goes into the `environment` key of the `argo` and `registry` configuration.

You can specify a default environment with the `default_environment` key.

### Docker image registry

The configuration for your docker image registry is stored in the `registry` key.

You must provide a `base_path` for your docker registry.

You must provide a `directory` for your docker registry for each environment.

The final image tag for `service` and `tag` inside `env` will be:
- the `base_path` configured for your registry, followed by a `/`
- the `directory` for `env`, followed by a `/`
- the `service_name` of `service`, followed by a `:`
- the `tag` you provided

If `base_path` is `myregistry.com`, `directory` for `env` is `my_dir`, and the service name is `service`,
the final image tag will be:
```
myregistry.com/my_dir/service:tag
```

### Deployment configuration

Deployment happens using `argocd` and image tag override.
All deployment configuration is stored inside the `argo` key.

For each environment, you should provide the following values:
- `auth_token`: this is the name of the environment variable containing the argocd auth token for the given environment
- `server`: this is the url of the server for the given environment

### Bundling configuration

You can run additional steps before and after building the docker images of your service.
The bundle configuration is stored inside the `scripts` key.

You can specify three different scripts:
- `prepare_bundle` is ran only once before any service is built
- `bundle` is ran for each service right before the docker build, and receives in input the service name
- `post_bundle` is ran only once after all services have been built

All keys are treated as paths to a script relative to the root directory.

### Service configuration

The service configuration is stored inside the `services` key.

The `services` key is a map, where each key is the name of the service.
For each service you must provide the following values:
- `directory`: the directory where `docker build` will run
- `service_name`: the name of the service on argo
- `image_name`: the name of the image on your registry
- `dockerfile`: (optional) the path to dockerfile to use for the service, relative to the `directory` of the service
- `scripts`: (optional) an object containing global scripts override. It is possibile to override only the bundle script
  - `bundle` is ran right before the docker build, and receives in input the service name. It is executed instead of global bundle script. The path is relative to the service directory and it is executed in the service directory

If all your services use the same dockerfile, you can specify it inside the `docker` key, as `dockerfile`.
This must be a path relative to the root directory of your project.

### Complete configuration example

A complete configuration looks something like this:
```json
{
  "default_environment": "dev",
  "argo": {
    "retries": 3,
    "environments": {
      "dev": {
        "auth_token": "ARGOCD_AUTH_TOKEN_DEV",
        "server": "argo.myproject.it"
      }
    }
  },
  "registry": {
    "base_path": "my.dkr.ecr.eu-central-1.amazonaws.com",
    "environments": {
      "dev": {
        "directory": "dev/myproject"
      }
    }
  },
  "docker": {
    "dockerfile": "./nest.Dockerfile"
  },
  "scripts": {
    "prepare_bundle": "./scripts/prepare.sh",
    "bundle": "./scripts/build.sh",
    "post_bundle": "./scripts/cleanup.sh"
  },
  "services": {
    "street-corners": {
      "directory": "./services/street-corners/",
      "service_name": "street-corners",
      "image_name": "street-corners"
    }
  }
}
```

## Usage

All commands accept the `--debug` flag.
This will overflow you with debug messages.

You can get help with commands with `deploy1 help` and `deploy1 help <command>`.

### Building a service

You can build one or more services with:

```bash
deploy1 build <service 1> <service 2> ...
```

The name of the service is the key of the `services` config map.

The docker image will be tagged based on the current branch with this policy:

- if you are on branch `dev`, the image will be tagged with `dev`
- if you are on a feature, bugfix, or task branch, the image will be tagged with the task ID.
  The task ID must match this regex: `[A-Z]+-[0-9]+`, e.g. `MX-1234`

If you are not on `dev` or on a branch associated with a task ID,
you must pass the tag manually, with:

```bash
deploy1 build --tag <tag> <service 1> <service 2> ...
```

You can specify which environment to build for using the `--env` flag.
If you don't specify an environment, the `default_environment` will be used.

You can also build all services with `deploy1 build-all`.
The flags are the same as `deploy1 build`.

#### Additional docker build arguments

You can pass some extra docker build arguments with `--build-args`:

```bash
deploy1 build --build-args "MY_VAR=abc OTHER_VAR=def" <service 1> ...
```

Additional build arguments should be enclosed in double quotes, and space-separated.

#### Skipping all bundle scripts

If you are sure your images are ready to be bundled, you can skip the bundle phase with:

```bash
deploy1 build --no-bundle <service 1> ...
```

#### Build and deploy in one shot

If you want to deploy your services right away, add the `--deploy` flag to the build command:

```bash
deploy1 build --deploy <service 1> ...
```

#### Disabling docker cache

The `--no-cache` flag adds the `--no-cache` flag to the `docker build` command:

```bash
deploy1 build --no-cache <service 1> ...
```

### Deploying a service

You can deploy one or more services with:

```bash
deploy1 deploy <service 1> <service 2> ...
```

The name of the service is the key of the `services` config map.

You can specify which environment to deploy to using the `--env` flag.
If you don't specify an environment, the `default_environment` will be used.

You can also deploy all services with `deploy1 deploy-all`.

## Extras

There's an attempt at support for shell completion.
To setup shell completion, run this command and follow the instructions for your shell:

```bash
deploy1 help completion
```

Shell completion is provided by [cobra](https://github.com/spf13/cobra/blob/master/shell_completions.md).
