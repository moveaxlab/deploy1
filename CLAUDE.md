# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

`deploy1` is a Go CLI tool that builds Docker images and deploys them via ArgoCD. It wraps `docker build/push` and the `argocd` CLI, driven by a `deploy1.json` config file in the consumer project root.

## Commands

```bash
go build ./...          # compile
go test ./...           # run all tests (only utils/ has tests currently)
go vet ./...            # static analysis
./lint.sh               # gofmt check (run before committing)
```

`mise install` pins Go to 1.19 (see `.mise.toml`). CI runs `go vet`, `./lint.sh`, and `go test`.

## Architecture

The CLI is built with [cobra](https://github.com/spf13/cobra). Entry point: `main.go` → `cmd.Execute()`.

**Package layout:**

| Package | Responsibility |
|---------|---------------|
| `cmd/` | Cobra command definitions (`build`, `build-all`, `deploy`, `deploy-all`, `reset`); flag parsing in `flags.go` |
| `config/` | Loads `./deploy1.json` at init time into `config.Config`; exposes typed accessors for services, registry, argo, scripts |
| `tag/` | Resolves the Docker image tag: explicit `--tag` → `"dev"` on `dev` branch → task ID extracted from branch name (regex `[A-Z]+-[0-9]+`) |
| `git/` | Reads current branch name via `git` |
| `utils/` | `GetTaskName` regex helper; has the only unit tests |
| `bundle/` | Runs `prepare_bundle`, per-service `bundle`, and `post_bundle` scripts |
| `docker/` | Wraps `docker build`, `docker push`, `docker inspect` (hash), and `docker manifest inspect` (tag existence check) |
| `argo/` | Wraps `argocd` CLI: `app get` (reads current tag + resource kind), `app set` (helm override), `app unset` (remove override), `app actions run restart`, `app wait` |
| `output/` | `io.Writer` adapters that pipe subprocess stdout/stderr through logrus |

**Build flow** (`deploy1 build <svc>`):
1. Run `prepare_bundle` script once
2. For each service: run per-service `bundle` script → `docker build` → `docker push`
3. Run `post_bundle` cleanup script
4. If `--deploy`: run ArgoCD deploys in parallel via `errgroup`

**Deploy logic** (`argo.Deploy`): fetches the service's current Helm image tag from ArgoCD. If the tag is unchanged, it calls `argocd app actions run restart` (rolling restart); if the tag differs, it calls `argocd app set --helm-set-string image.tag=<tag>`. With `--wait`, it additionally calls `argocd app wait --health --sync`.

**Reset logic** (`argo.Reset`): calls `argocd app unset -p image.tag` (or the service's custom `imageTagParameter`) to remove the Helm override, reverting the service to its chart default tag. Accepts one or more services; runs in parallel. Supports `--wait`. Note: ArgoCD v2.7.1 uses `-p <key>` (not `--helm-set-string`) on the `unset` subcommand.

**Flag groups** (`cmd/flags.go`): each command family owns its flag group — `addBaseFlags`/`getBaseFlags` (tag + env, for build/deploy), `addBuildFlags`/`getBuildFlags`, `addDeployFlags`/`getDeployFlags`, `addResetFlags`/`getResetFlags` (env + wait, for reset). `ValidateEnvironment` checks `Config.Registry.Environments`; Argo environments are accessed directly via `Config.Argo.Environments[env]` and are not independently validated.

**ArgoCD auth**: credentials are never hardcoded. The `ARGOCD_AUTH_TOKEN` and `ARGOCD_SERVER` env vars are set per-command from the environment variable names listed in `deploy1.json` (`argo.environments.<env>.auth_token`).

**Config loading**: `config.init()` reads `./deploy1.json` relative to the working directory where `deploy1` is run (the consumer project root, not this repo).
