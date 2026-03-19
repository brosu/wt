#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "$0")/lib.sh"
mcp_init

shell_start '["run","--rm","-it","-e","WORKTREE_ROOT=/home/demo/worktrees","-e","WORKTREE_STRATEGY=custom","-e","WORKTREE_PATTERN={.worktreeRoot}/{.branch}/{.repo.Name}","wt-demo"]'
setup_prompt

record_start "wt — multi-repo"
run "wt info";                                       pause 5
run "cd ~/src/shared-lib && wt co feat/auth";        pause
run "cd ~/src/main-app && wt co feat/auth";          pause
run "tree -C ~/worktrees";                           pause 5
end_pause
record_save "wt-multi-repo"

shell_stop
