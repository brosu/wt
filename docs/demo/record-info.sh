#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "$0")/lib.sh"
mcp_init

shell_start '["run","--rm","-it","-e","WORKTREE_ROOT=/home/demo/worktrees","wt-demo"]'
send "cd src/main-app\\r"; sleep 1
setup_prompt

record_start "wt — info"
run "wt info";          pause 5
run "wt config show";   pause 5
end_pause
record_save "wt-info"

shell_stop
