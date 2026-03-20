Output "../wt-interactive.gif"

Set Theme "Dracula"
Set FontSize 16
Set Width 1200
Set Height 600
Set TypingSpeed 30ms
Set Padding 20

Hide
Type "docker run --rm -it -e WORKTREE_ROOT=/home/demo/worktrees wt-demo bash --rcfile /home/demo/.bashrc"
Enter
Sleep 3s
Type "cd src/main-app"
Enter
Sleep 1s
Type "clear"
Enter
Sleep 1s
Show

# Tab completion for subcommands
Type "wt "
Tab
Tab
Sleep 5s

# Clear and show interactive branch selection
Ctrl+U
Type "clear"
Enter
Sleep 1s

Type "wt co"
Enter
Sleep 3s
Down
Sleep 1s
Down
Sleep 1s
Up
Sleep 1s
Enter
Sleep 5s

Sleep 5s
