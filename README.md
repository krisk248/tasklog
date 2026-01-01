# tasklog

A beautiful terminal-based task logger and time tracker with an **Ultraviolet** theme. Built with Go and [Bubbletea](https://github.com/charmbracelet/bubbletea).

![Ultraviolet Theme](https://via.placeholder.com/800x400/0d0d14/a855f7?text=tasklog+%E2%80%A2+Ultraviolet+Theme)

## Features

- **Three-pane layout**: Calendar | Tasks | Timeline
- **Nested tasks**: Infinite subtask hierarchy with expand/collapse
- **Task states**: Todo, Completed, Delegated, Delayed
- **Task priorities**: P1 (Critical), P2 (Important), P3 (Normal)
- **Time tracking**: Start/stop timer on tasks
- **Activity timeline**: Automatic logging of all task state changes
- **Search & Filter**: Find tasks quickly, filter by state or priority
- **Export**: Markdown, JSON, or Plain Text
- **Beautiful themes**: Ultraviolet (default), Terminal, Minimal, Nord
- **Undo/Redo**: 50-state history
- **Vim-style navigation**: hjkl + arrow keys
- **Single binary**: No dependencies, runs anywhere

## Installation

### From Releases (Recommended)

Download the latest release for your platform from the [Releases](https://github.com/krisk248/tasklog/releases) page.

### From Source

```bash
go install github.com/krisk248/tasklog/cmd/tasklog@latest
```

Or clone and build:

```bash
git clone https://github.com/krisk248/tasklog.git
cd tasklog
go build -o tasklog ./cmd/tasklog
./tasklog
```

## Keyboard Shortcuts

### Global

| Key | Action |
|-----|--------|
| `Ctrl+C` | Exit (press twice) |
| `Ctrl+U` | Undo |
| `Ctrl+T` | Theme selector |
| `Ctrl+E` | Export dialog |
| `?` | Help |
| `:` | Month overview |
| `/` | Search tasks |
| `Esc` | Clear search/filter |
| `1/2/3` | Switch panes |
| `Tab` | Next pane |

### Calendar Pane

| Key | Action |
|-----|--------|
| `h/l` or `←/→` | Previous/next day |
| `j/k` or `↑/↓` | Previous/next week |
| `n/p` | Next/previous month |
| `T` | Jump to today |

### Tasks Pane

| Key | Action |
|-----|--------|
| `j/k` or `↑/↓` | Navigate up/down |
| `a` | Add task |
| `e` | Edit task |
| `d` | Delete task |
| `Space` | Toggle complete |
| `D` | Delegate task |
| `x` | Toggle delayed |
| `s` | Start/stop timer |
| `1/2/3` | Set priority P1/P2/P3 |
| `0` | Clear priority |
| `Enter` or `→` | Expand/collapse |
| `←` | Collapse |

### Timeline Pane

| Key | Action |
|-----|--------|
| `j/k` or `↑/↓` | Scroll |
| `Ctrl+D` | Page down |
| `C` | Clear timeline |

## Themes

tasklog comes with 4 beautiful themes:

### Ultraviolet (Default)
Deep space black with electric violet accents. A modern, eye-catching theme.

### Terminal
Classic green/amber phosphor on black. Retro CRT vibes.

### Minimal
Clean grayscale. Distraction-free productivity.

### Nord
Arctic blue palette. Calm and professional.

Switch themes with `Ctrl+T`.

## Data Storage

Your data is stored locally:

- **macOS**: `~/Library/Application Support/tasklog/data.json`
- **Linux**: `~/.config/tasklog/data.json`
- **Windows**: `%APPDATA%\tasklog\data.json`

## Color Palette (Ultraviolet)

```
Background:     #0d0d14  (deep space black)
Surface:        #1a1625  (dark purple-gray)
Border:         #2d2640  (muted violet)
Primary:        #a855f7  (electric violet)
Secondary:      #c084fc  (soft lavender)
Accent:         #e879f9  (hot pink/magenta)
Success:        #22d3ee  (cyan glow)
Warning:        #f59e0b  (amber)
Error:          #f43f5e  (rose)
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

Built with these amazing Go libraries:

- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components

---

Made with 💜 and Go
