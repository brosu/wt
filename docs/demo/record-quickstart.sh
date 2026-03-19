#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "$0")/lib.sh"
mcp_init

shell_start '["run","--rm","-it","wt-demo"]'
send "cd src/main-app\\r"; sleep 1
setup_prompt

record_start "wt — quickstart"
run "wt co feat/auth";        pause
run "wt co feat/dashboard";   pause
run "wt ls";                  pause 5
run "wt co feat/auth";        pause
run "wt rm feat/dashboard";   pause
run "wt ls";                  pause 5
end_pause
record_save "wt-quickstart"

shell_stop
