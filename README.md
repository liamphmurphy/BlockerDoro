# BlockerDoro

BlockerDoro is a pomodoro timer that, during the "work" periods, will modify `/etc/hosts` to block certain domains the user may find distracting.

## Running

Before running, make sure your `/etc/hosts` file is in a good state before starting. This program, at least for now, will use the original state of `/etc/hosts` as the "source of truth"
for setting that file back to normal.

Set the path to your local config directory at the top of the main func in `main.go`. This is manual until I figure out a workaround to having to run as root.

`sudo go run *.go`

**NOTE:** As of now, root is required since the program directly interfaces with `/etc/hosts`. I'm hoping to figure out a way around this, or at least restricting the privileges this program has.

# Trello Board

I'm testing out using Trello to scope out / gather my thoughts for this project, it's a new thing for me, so I'm curious to see how it goes. If you want to see the board, go [here](https://trello.com/b/27ex4Tzu/blockerdoro).