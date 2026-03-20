Output "../wt-multi-repo.gif"

Set Theme "Dracula"
Set FontSize 16
Set Width 1200
Set Height 600
Set TypingSpeed 30ms
Set Padding 20

Hide
Type "docker run --rm -it -e WORKTREE_ROOT=/home/demo/worktrees -e WT_CONFIG=/home/demo/.config/wt/variants/multi-repo.toml wt-demo bash --rcfile /home/demo/.bashrc"
Enter
Sleep 3s
Type "clear"
Enter
Sleep 1s
Show

Type "cd ~/src/main-app"
Enter
Sleep 1s

Type "wt co feat/auth"
Enter
Sleep 4s

Type "cd ~/src/shared-lib"
Enter
Sleep 1s

Type "wt co feat/auth"
Enter
Sleep 4s

Type "tree ~/worktrees"
Enter
Sleep 6s

Sleep 5s
