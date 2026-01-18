# dots

**Dotfiles that just work**

[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg?style=flat&logo=go&logoColor=white)](https://go.dev/doc/install)
[![License](https://img.shields.io/badge/license-MIT-2ea44f.svg?style=flat)](LICENSE)
[![Build](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat)](#)

A focused, minimal CLI for managing dotfiles with real filenames, readable YAML, and reliable symlinks. No template engines. No magic. Just your dotfiles, under version control.

## Why dots?

Most dotfile managers are powerful, but they often feel like a configuration framework. Tools like chezmoi are feature-rich but require learning templates, scripts, and abstractions before you can move your files. **dots** is intentionally simpler: copy the file, track it in YAML, and link it back. You get a clean repo, predictable behavior, and dotfiles you can edit with any editor.

## Features

- ✅ **Real filenames** stored as-is (e.g. `.bashrc`, `.gitconfig`)
- 🔗 **Symlink-based workflow** for instant edits and rollbacks
- 📄 **Readable YAML manifest** (`~/.dots/dots.yaml`)
- 🎯 **Minimal commands**: init, add, status, apply
- 🎨 **Colorful status output** for quick sync checks

## Install

### Go install

```bash
go install github.com/subcode-labs/dots@latest
```

### Manual

```bash
git clone https://github.com/subcode-labs/dots.git
cd dots
go build -o dots
sudo mv dots /usr/local/bin/
```

## Usage

### Initialize

```bash
$ dots init
Initialized dots at /home/jonty/.dots
```

### Add a dotfile

```bash
$ dots add ~/.bashrc
Tracked /home/jonty/.bashrc -> /home/jonty/.dots/.bashrc
```

### Check status

```bash
$ dots status
synced    /home/jonty/.bashrc
missing   /home/jonty/.vimrc (target missing)
diverged  /home/jonty/.gitconfig
conflict  /home/jonty/.zshrc (not a symlink)
```

### Apply symlinks

```bash
$ dots apply
Linked /home/jonty/.bashrc -> /home/jonty/.dots/.bashrc
Linked /home/jonty/.vimrc -> /home/jonty/.dots/.vimrc
```

## How it works

**dots** keeps a dedicated directory at `~/.dots/` that contains your real dotfiles. When you run `dots add`, it copies the file into that directory and creates a symlink at the original location. The manifest (`~/.dots/dots.yaml`) stores the source and target so you can apply the links on any machine with a single command.

## Comparison

| Feature | dots | chezmoi |
| --- | --- | --- |
| Real filenames in repo | ✅ | ⚠️ (requires templates or config) |
| YAML manifest | ✅ | ❌ |
| Templating / scripting required | ❌ | ✅ |
| Learning curve | Low | Medium / High |
| Symlink-first workflow | ✅ | ✅ |

## Contributing

Contributions are welcome! If you have ideas, open an issue or submit a pull request. Please keep changes focused and consistent with the minimal philosophy of the project.

## License

MIT — see `LICENSE`.
