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

Copy (or move) one or more files into the staging area.

```
aidrop add [flags] <file> [files...]
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
- Symbolic links are followed — the resolved target is copied, not the link itself.
- Filename conflicts are resolved automatically by appending a numeric suffix (`file-2.txt`, `file-3.txt`, …).
- Directories are not supported; use a shell glob to expand them (`*.go`, `src/**/*.ts`).

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

## Build targets

```sh
make build        # Compile the binary
make install      # Install to $GOPATH/bin
make test         # Run tests
make fmt          # Format source code
make vet          # Run go vet
make lint         # Run golangci-lint
make coverage     # Generate HTML coverage report
make build-all    # Cross-compile for Linux, Windows, and macOS
make clean        # Remove build artifacts
make help         # List all targets
```
