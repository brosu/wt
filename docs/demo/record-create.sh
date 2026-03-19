#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "$0")/lib.sh"
mcp_init

shell_start '["run","--rm","-it","wt-demo"]'
send "cd src/main-app\\r"; sleep 1
setup_prompt

record_start "wt — create & remove"
run "wt create feat/new-api";  pause
run "wt ls";                   pause 5
run "wt rm feat/new-api";      pause
end_pause
record_save "wt-create"

shell_stop
