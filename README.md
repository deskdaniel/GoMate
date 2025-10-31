# GoMate
GoMate is a terminal-based chess game written in Go.
It features a TUI (text-based user interface) for local two-player matches.
No AI opponent is included — both players play locally on the same machine.

## Motivation
This project started as a way to practice Go while building something familiar and self-contained.
Chess provides a well-defined set of rules and logic, making it a good fit for exploring concepts like turn-based gameplay, move validation and tracking game state.
It also serves as a lightweight way to play locally with a friend.

## Features
- Register, log in, log out, and view player statistics
- Core chess rules implemented, including:
    - Check and checkmate detection
    - Draw by stalemate
    - Draw by insufficient material
    - Draw by the fifty-move rule
    - Castling
    - En passant
- Ability to offer or accept draws
- Option to forfeit a game
- Player statistics automatically update after each game

## Requirements
- Go: version 1.25.1 was used during development (recommended).
Older versions (1.22+) should work, but behavior is not guaranteed.
- SQLite driver: this project uses [go-sqlite3](https://github.com/mattn/go-sqlite3), which requires working GCC compiler in your system PATH.
You can find setup instructions [here](https://github.com/mattn/go-sqlite3?tab=readme-ov-file#compiling)
- Other dependencies will be automaticaly downloaded when building the executable/running app for first time

## Quick Start
1. Make sure [Go](https://go.dev/dl/) is installed.
2. Clone the repository:
```
git clone <repo-url>
cd GoMate
```
3. Build the executable
```
go build
```
This will produce an executable in the current directory.
To specify a custom executable name:
```
go build -o <desired-name>
```

## Usage
Run the executable to launch the game in terminal mode.
Use `arrow keys` to navigate the main menu and `Enter` to select options.
You can also type the number next to an option for quick access.

You can play as a registered user or as a guest (stats are not tracked for guests).

Exit any screen (except during a game) using `Esc` or `Ctrl + C`.
During a game, type your move into the input field using the format:
```
a2 a4
```
This moves the piece from A2 to A4 (if the move is legal).

### Ending a Game Early
You can end the game before checkmate by:
- Offering a draw: type `draw`
Your opponent must accept for it to take effect.
- Forfeiting: type one of the following:
    - `ff`
    - `forfeit`
    - `surrender`
    - `surr`
The forfeiting player records a loss, while the opponent records a win.
- Closing the terminal window.

## Contributing
If you want to contribute you can fork the repository and open pull request.
Please add tests to test your suggested changes, and make sure you pass already existing tests.

## Notes
- Some terminal fonts may not display chess pieces or board symbols correctly. For best results, use DejaVu Sans Mono (tested and confirmed working).
- Detailed explanations of chess rules, draw conditions, and piece movements are available in the `Help` section of the main menu.
