# Knowledge Base

## Vision & Strategy

**Business overview**
- dots is a lightweight dotfile manager focused on speed, clarity, and edit-in-place workflows.
- The product goal is a frictionless CLI that feels like a natural extension of existing shell workflows.

**Target market**
- Developers who want dotfile management without template engines, heavy config layers, or complex onboarding.
- Teams or individuals maintaining multiple machines who value predictability over endless features.

**Competitive positioning vs chezmoi**
- dots prioritizes simplicity and real file names; chezmoi prioritizes templating and configuration power.
- dots is ideal for minimalism and transparency; chezmoi is ideal for heavily customized environments.
- dots differentiates by keeping YAML manifest and symlink-first philosophy front and center.

**Revenue model**
- GitHub Sponsors, targeting $1K–$3K/month in passive income.
- Emphasize community-driven updates and transparent roadmaps to encourage sponsorships.

## Technical Documentation

**Tech stack**
- Go 1.21+.
- CLI framework: Cobra.
- Config format: YAML via `gopkg.in/yaml.v3`.
- Output styling: `github.com/fatih/color`.

**Architecture decisions**
- Symlink approach for edit-in-place: files are copied into `~/.dots/` and symlinked back to their original locations.
- Real filenames are stored in the repo (no `dot_` prefixes) to keep the repository readable.
- YAML manifest (`~/.dots/dots.yaml`) tracks source and target paths.

**Repository structure**
- `cmd/`: Cobra commands (init, add, status, apply).
- `internal/config/`: Manifest loading/saving and helpers.
- `internal/dotfile/`: File operations, status detection, symlink handling.
- `docs/`: GitHub Pages landing page.

## Progress Tracker

**Completed**
- CLI implementation (init/add/status/apply).
- README with features, usage, and comparison table.
- CI/CD with GoReleaser and GitHub Actions.
- GitHub Pages landing site.

**Pending**
- Show HN launch.
- GitHub Sponsors setup and marketing.

## Launch Plan

**Show HN strategy**
- Title: "Show HN: dots – Dotfiles that just work".
- Post timing: US morning, Tuesday–Thursday.
- Lead with one sentence summary and a short usage snippet.

**Marketing channels**
- Reddit: r/unixporn, r/linux, r/commandline.
- Twitter/X: share GIF of status output and the repo link.
- Dev.to: short post on why simple dotfile management matters.

## Session Notes (AI continuity)

**Key links**
- GitHub repo: https://github.com/subcode-labs/dots
- GitHub Pages: https://subcode-labs.github.io/dots

**Technical notes**
- OpenCode runs best with single-line prompts for command execution.
- Go tooling may not be available in all environments; be ready to note missing `go`.
- `dots init` creates `~/.dots/` and initializes a git repo there.
- `dots add` copies the file into `~/.dots/` with real name and symlinks back to original path.
- Status labels: synced, missing, diverged, conflict.

**Decisions made**
- Use YAML manifest at `~/.dots/dots.yaml`.
- Keep file names unchanged (no dot_ prefix).
- GoReleaser for tagged releases and GitHub Actions for CI.
- GitHub Pages site served from `docs/index.html`.
