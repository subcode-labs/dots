# Show HN: dots - Dotfiles that just work (no templates, edit-in-place)

Hi HN! I built **dots** because I wanted a dotfile manager that felt as simple as `git add` for config files. Tools like chezmoi are powerful, but I found the templates and layers of config too heavy for a small set of dotfiles.

**dots** keeps things straightforward: it copies your file into `~/.dots/`, tracks it in a simple YAML manifest, and symlinks it back to the original location so you can edit in place.

**Key differentiators**
- Real filenames in the repo (`.bashrc`, `.gitconfig`) with no `dot_` prefixes
- Minimal YAML manifest (`~/.dots/dots.yaml`)
- Symlink-first workflow so edits update the tracked file instantly
- Simple commands: `init`, `add`, `status`, `apply`

**Quick usage**

```bash
$ dots init
Initialized dots at /home/jonty/.dots

$ dots add ~/.bashrc
Tracked /home/jonty/.bashrc -> /home/jonty/.dots/.bashrc

$ dots status
synced    /home/jonty/.bashrc
```

GitHub: https://github.com/subcode-labs/dots
Website: https://subcode-labs.github.io/dots

## Why I built this

I wanted a dotfile tool that was easy to explain in one sentence and easy to understand by looking at the repo. I kept the core workflow to copying files, linking them back, and tracking paths in YAML. It’s intentionally minimal and feels closer to native file management than a configuration framework.
