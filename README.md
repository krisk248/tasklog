# seyal

A beautiful terminal-based task manager with an **Ultraviolet** theme. Built with Go and [Bubbletea](https://github.com/charmbracelet/bubbletea).

*"Seyal" (செயல்) means "action" or "deed" in Tamil - take action on your tasks!*

## Features

- **Three-pane layout**: Calendar | Tasks | Timeline
- **Nested tasks**: Infinite subtask hierarchy with expand/collapse
- **Task states**: Todo, Completed, Delegated, Delayed
- **Task priorities**: P1 (Critical), P2 (Important), P3 (Normal)
- **Time tracking**: Start/stop timer on tasks
- **Push to next day**: Move tasks forward with pushed count tracking
- **Activity timeline**: Automatic logging of all task state changes
- **Search & Filter**: Find tasks quickly, filter by state or priority
- **Export**: Markdown, JSON, or Plain Text to ~/Documents/seyal-exports/
- **Month Overview**: See all tasks in a month grid (`:`)
- **Undo**: 50-state history
- **Vim-style navigation**: hjkl + arrow keys
- **Single binary**: No dependencies, runs anywhere

## Installation

### From Source (Recommended)

```bash
go install github.com/krisk248/seyal/cmd/seyal@latest
```

Or clone and build:

```bash
git clone https://github.com/krisk248/seyal.git
cd seyal
go build -o seyal ./cmd/seyal
./seyal
```

## Keyboard Shortcuts

### Global

| Key | Action |
|-----|--------|
| `Ctrl+C` | Exit (press twice) |
| `Ctrl+U` | Undo |
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
| `v` | View full task details |
| `Space` | Toggle complete |
| `D` | Delegate task |
| `x` | Toggle delayed |
| `s` | Start/stop timer |
| `n` | Push to next day |
| `1/2/3` | Set priority P1/P2/P3 |
| `0` | Clear priority |
| `Enter` or `→` | Expand/collapse |
| `←` | Collapse |

### Timeline Pane

| Key | Action |
|-----|--------|
| `j/k` or `↑/↓` | Scroll |
| `Shift+C` | Clear timeline |

## Data Storage

Your data is stored locally in a human-readable JSON file. This allows for easy backups or manual editing if necessary.

- **macOS**: `~/Library/Application Support/seyal/data.json`
- **Linux**: `~/.local/share/seyal/data.json` (or `$XDG_DATA_HOME`)
- **Windows**: `%APPDATA%\seyal\data.json`

## Export

Exports are saved to a common folder for easy access:

- **All platforms**: `~/Documents/seyal-exports/`

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

Made with 💜 by krisk248
