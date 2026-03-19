#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "$0")/lib.sh"
mcp_init

shell_start '["run","--rm","-it","-e","WORKTREE_ROOT=/home/demo/worktrees","-e","WT_CONFIG=/home/demo/.config/wt/variants/hooks.toml","wt-demo"]'
send "cd ~/src/main-app\\r"; sleep 1
setup_prompt

record_start "wt — hooks"
# Show .env in main checkout
run "cat .env";                pause
# Create worktree — hook copies .env automatically
run "wt co fix/login-bug";    pause
# Verify .env was copied to the worktree
run "cat ~/worktrees/main-app/fix/login-bug/.env";  pause
end_pause
record_save "wt-hooks"

shell_stop
