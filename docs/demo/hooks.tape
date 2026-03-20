Output "../wt-hooks.gif"

Set Theme "Dracula"
Set FontSize 16
Set Width 1200
Set Height 600
Set TypingSpeed 30ms
Set Padding 20

Hide
Type "docker run --rm -it -e WORKTREE_ROOT=/home/demo/worktrees -e WT_CONFIG=/home/demo/.config/wt/variants/hooks.toml wt-demo bash --rcfile /home/demo/.bashrc"
Enter
Sleep 3s
Type "cd src/main-app"
Enter
Sleep 1s
Type "clear"
Enter
Sleep 1s
Show

Type "cat .env"
Enter
Sleep 4s

Type "wt co fix/login-bug"
Enter
Sleep 4s

Type "cat ~/worktrees/main-app/fix/login-bug/.env"
Enter
Sleep 6s

Sleep 5s
