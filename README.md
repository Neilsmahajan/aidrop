# aidrop

A fast command-line tool that maintains a structured staging area at `~/AIDrop/` for files you intend to share with an AI assistant or chat application.

Files are organized by **project** (auto-detected from the current git repository, or set explicitly) and optionally grouped into named **sessions** with automatic date prefixes. A `clean` command keeps the staging area tidy over time.

---

## Installation

```sh
go install github.com/neilsmahajan/aidrop@latest
```

Or build from source:

```sh
git clone https://github.com/neilsmahajan/aidrop
cd aidrop/cmd
make install
```

---

## Configuration

### `AIDROP_DIR`

By default aidrop stores files at `~/AIDrop/`. Set `AIDROP_DIR` to use a different location:

```sh
export AIDROP_DIR=/tmp/aidrop-scratch
aidrop add myfile.go
# → /tmp/aidrop-scratch/<project>/myfile.go
```

---

## Directory structure

```
~/AIDrop/
├── my-api/                         ← project (inferred from git or --project)
│   ├── handler.go                  ← loose file (added without --session)
│   └── 2026-05-31-auth-bug/        ← session directory (--session auth-bug)
│       ├── middleware.go
│       └── token.go
└── default/                        ← fallback project when outside a git repo
    └── notes.md
```

---

## Commands

### `aidrop add`

Copy (or move) one or more files or directories into the staging area.

```
aidrop add [flags] <file|dir> [files|dirs...]
```

| Flag | Short | Default | Description |
|---|---|---|---|
| `--project` | `-p` | _(git repo name or "default")_ | Project name to place files under |
| `--session` | `-s` | _(none)_ | Session name; creates a `YYYY-MM-DD-<session>` subdirectory |
| `--move` | `-m` | false | Move files instead of copying them |

**Project resolution order:**
1. `--project` value if provided.
2. Name of the current git repository root (automatic inference).
3. `"default"` if no git repository is detected.

**Notes:**
- Directories are copied recursively, preserving their internal structure. Hidden files (dotfiles, `.DS_Store`) are skipped automatically.
- Symbolic links are followed — the resolved target is copied, not the link itself.
- Filename conflicts are resolved automatically by appending a numeric suffix (`file-2.txt`, `file-3.txt`, …).

**Examples:**

```sh
# Copy two files into the current git project
aidrop add README.md internal/models.go

# Copy into an explicit project
aidrop add -p federation-service README.md internal/models.go

# Copy into a named session (creates ~/AIDrop/snake-game/2026-05-31-add-animation/)
aidrop add -p snake-game -s add-animation animate.go

# Move a file (remove from source after staging)
aidrop add -s stack-overflow-issue -m output.log

# Copy an entire directory tree
aidrop add -p my-project src/
```

---

### `aidrop ls`

Print a tree of all staged files, organized by project and session.

```
aidrop ls [flags]
```

| Flag | Short | Default | Description |
|---|---|---|---|
| `--project` | `-p` | _(all)_ | Restrict output to a specific project |

Hidden files (dotfiles, `.DS_Store`) are never shown.

**Examples:**

```sh
# Show everything
aidrop ls

# Show only one project
aidrop ls -p federation-service
```

**Sample output:**

```
AIDrop  /Users/you/AIDrop
├── federation-service/
│   ├── 2026-05-31-auth-bug/
│   │   ├── middleware.go
│   │   └── token.go
│   └── handler.go
└── default/
    └── notes.md
```

---

### `aidrop open`

Open the AIDrop directory, a project, or a session in the system file manager (Finder on macOS, `xdg-open` on Linux).

```
aidrop open [project [session]]
```

**Examples:**

```sh
# Open ~/AIDrop/ in Finder
aidrop open

# Open a project folder
aidrop open federation-service

# Open a specific session
aidrop open federation-service 2026-05-31-auth-bug
```

---

### `aidrop rm`

Remove a project or session directory from the staging area.

```
aidrop rm [flags] <project> [session]
```

| Flag | Short | Default | Description |
|---|---|---|---|
| `--soft` | `-s` | false | Move to the system trash instead of permanently deleting |
| `--dry-run` | | false | Preview what would be removed without making any changes |

**Examples:**

```sh
# Delete an entire project
aidrop rm federation-service

# Delete a single session
aidrop rm federation-service 2026-05-31-auth-bug

# Move to trash instead of deleting
aidrop rm -s federation-service

# Preview without deleting
aidrop rm --dry-run federation-service
```

---

### `aidrop clean`

Remove session directories from the staging area that are older than a given number of days.

```
aidrop clean [flags]
```

| Flag | Short | Default | Description |
|---|---|---|---|
| `--days` | `-d` | 7 | Remove sessions older than this many days |
| `--soft` | `-s` | false | Move items to the system trash instead of permanently deleting them |
| `--loose` | | false | Also remove loose files at the project root (all if `--days` not set explicitly; otherwise only those older than `--days`) |
| `--dry-run` | | false | Preview what would be removed without making any changes |

**Notes:**
- Only dated session directories (names beginning with `YYYY-MM-DD`) are considered for age-based removal. Undated subdirectories are never touched automatically.
- `--soft` moves items to `~/.Trash` on macOS or `~/.local/share/Trash/files` on Linux.
- Combine `--dry-run` with any other flags to safely preview the effect before committing.

**Examples:**

```sh
# Remove sessions older than 7 days (default)
aidrop clean

# Remove sessions older than 30 days
aidrop clean -d 30

# Move old sessions to trash instead of deleting
aidrop clean -s

# Also remove all loose project-root files (no age filter)
aidrop clean --loose

# Remove loose files older than 30 days
aidrop clean --loose -d 30

# Preview without deleting
aidrop clean --dry-run
aidrop clean --loose -d 30 --dry-run
```

---

## Shell completions

Cobra provides built-in tab completion for bash, zsh, fish, and PowerShell.

**Manual setup:**

```sh
# bash
aidrop completion bash > /etc/bash_completion.d/aidrop
# or for user-local install:
aidrop completion bash > ~/.local/share/bash-completion/completions/aidrop

# zsh  (add ~/.zsh/completions to your fpath first)
aidrop completion zsh > ~/.zsh/completions/_aidrop

# fish
aidrop completion fish > ~/.config/fish/completions/aidrop.fish

# PowerShell
aidrop completion powershell | Out-String | Invoke-Expression
```

**Via Makefile:**

```sh
make completion-bash   # Install bash completion
make completion-zsh    # Install zsh completion
make completion-fish   # Install fish completion
make completion        # Install all three
```

---

## Build targets

```sh
make build            # Compile the binary
make install          # Install to $GOPATH/bin
make test             # Run tests
make fmt              # Format source code
make vet              # Run go vet
make lint             # Run golangci-lint
make coverage         # Generate HTML coverage report
make build-all        # Cross-compile for Linux, Windows, and macOS
make clean            # Remove build artifacts
make completion       # Install shell completions (bash, zsh, fish)
make help             # List all targets
```
