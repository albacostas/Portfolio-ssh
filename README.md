# SSH Portfolio

An interactive terminal-based portfolio accessible via SSH. Built with Go, showcasing a full-featured TUI (Terminal User Interface) with navigation, projects, skills, and experience.

## Features

- **Interactive SSH Server** - Connect directly via SSH to explore the portfolio
- **Beautiful Terminal UI** - Built with Charmbracelet libraries (Bubbletea, Lipgloss, Wish)
- **Multiple Sections**:
  - About - Personal introduction
  - Skills - Programming languages, frameworks, tools, and platforms
  - Projects - Showcase of featured projects with descriptions and GitHub links
  - Experience - Education and professional background
  - Contact - Email, GitHub, LinkedIn, and social media
- **Real-time Navigation** - Smooth keyboard-based navigation with vim keybindings
- **ASCII Art Profile** - Custom ASCII art avatar displayed alongside information
- **Clean, Modern Design** - Purple and cyan color scheme for an elegant appearance

## Tech Stack

- **Language**: Go
- **UI Framework**: Charmbracelet (Bubbletea, Lipgloss)
- **SSH Server**: Charmbracelet Wish
- **Build System**: Go Modules

## Installation

### Prerequisites
- Go 1.19 or higher
- SSH client (to connect to the portfolio)
- An SSH key

## Usage

### Connect via SSH

```bash
ssh -p 2222 localhost
```

### Navigation

| Key | Action |
|-----|--------|
| `→` / `l` / `Tab` | Next tab |
| `←` / `h` / `Shift+Tab` | Previous tab |
| `↑` / `k` | Navigate up/previous item |
| `↓` / `j` | Navigate down/next item |
| `q` / `Ctrl+C` | Quit |

## Customization

Edit `main.go` to customize:

- **Colors**: Modify the color variables at the top of the file
- **Personal Info**: Update the `projects`, `skills`, and `experience` arrays
- **ASCII Art**: Replace the `avatarASCII` constant with your own design
- **Port**: Change the SSH server port in the `main()` function

