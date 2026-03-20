#!/usr/bin/env bash
# Shared helpers for shellwright recording scripts
# Source this file: source "$(dirname "$0")/lib.sh"

SHELLWRIGHT_URL="${SHELLWRIGHT_URL:-http://localhost:7498}"
OUTPUT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
ID_COUNTER="${ID_COUNTER:-0}"

next_id() { ID_COUNTER=$((ID_COUNTER + 1)); echo $ID_COUNTER; }

mcp_call() {
  local name="$1" args="$2"
  local id; id=$(next_id)
  curl -s -X POST "$SHELLWRIGHT_URL/mcp" \
    -H "Content-Type: application/json" \
    -H "Accept: application/json, text/event-stream" \
    -H "Mcp-Session-Id: $MCP_SESSION" \
    -d "{\"jsonrpc\":\"2.0\",\"id\":$id,\"method\":\"tools/call\",\"params\":{\"name\":\"$name\",\"arguments\":$args}}"
}

# Fire-and-forget: send without waiting for response
mcp_fire() {
  local name="$1" args="$2"
  local id; id=$(next_id)
  curl -s --max-time 3 -X POST "$SHELLWRIGHT_URL/mcp" \
    -H "Content-Type: application/json" \
    -H "Accept: application/json, text/event-stream" \
    -H "Mcp-Session-Id: $MCP_SESSION" \
    -d "{\"jsonrpc\":\"2.0\",\"id\":$id,\"method\":\"tools/call\",\"params\":{\"name\":\"$name\",\"arguments\":$args}}" > /dev/null 2>&1 &
}

mcp_init() {
  MCP_SESSION=$(curl -s -D - -X POST "$SHELLWRIGHT_URL/mcp" \
    -H "Content-Type: application/json" \
    -H "Accept: application/json, text/event-stream" \
    -d "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"initialize\",\"params\":{\"protocolVersion\":\"2025-03-26\",\"capabilities\":{},\"clientInfo\":{\"name\":\"rec\",\"version\":\"0.1\"}}}" \
    2>&1 | grep -i "mcp-session-id" | awk '{print $2}' | tr -d '\r')
  curl -s -X POST "$SHELLWRIGHT_URL/mcp" \
    -H "Content-Type: application/json" -H "Accept: application/json, text/event-stream" \
    -H "Mcp-Session-Id: $MCP_SESSION" \
    -d '{"jsonrpc":"2.0","method":"notifications/initialized","params":{}}' > /dev/null 2>&1
  echo "MCP session: $MCP_SESSION"
}

send() {
  mcp_call "shell_send" "{\"session_id\":\"$SHELL_SESSION\",\"input\":\"$1\",\"delay_ms\":0}" > /dev/null 2>&1
}

# Type command char by char (fire-and-forget, 30ms between chars) then press enter
run() {
  local text="$1"
  for (( i=0; i<${#text}; i++ )); do
    local char="${text:$i:1}"
    case "$char" in '"') char='\\"' ;; '\\') char='\\\\' ;; esac
    mcp_fire "shell_send" "{\"session_id\":\"$SHELL_SESSION\",\"input\":\"$char\",\"delay_ms\":0}"
    sleep 0.03
  done
  sleep 0.3
  send "\\r"
}

# Pause that defeats frame deduplication by typing/deleting a char
pause() {
  local duration="${1:-6}"
  local end=$((SECONDS + duration))
  while [ $SECONDS -lt $end ]; do
    mcp_fire "shell_send" "{\"session_id\":\"$SHELL_SESSION\",\"input\":\" \",\"delay_ms\":0}"
    sleep 0.5
    mcp_fire "shell_send" "{\"session_id\":\"$SHELL_SESSION\",\"input\":\"\\\\b \\\\b\",\"delay_ms\":0}"
    sleep 0.5
  done
}

shell_start() {
  local docker_args="$1"
  local theme="${2:-dracula}"
  local result
  result=$(mcp_call "shell_start" "{\"command\":\"docker\",\"args\":$docker_args,\"cols\":90,\"rows\":18,\"theme\":\"$theme\"}")
  SHELL_SESSION=$(echo "$result" | grep -o 'shell-session-[a-f0-9]*')
  echo "Shell session: $SHELL_SESSION"
  sleep 2
}

shell_stop() {
  mcp_call "shell_stop" "{\"session_id\":\"$SHELL_SESSION\"}" > /dev/null 2>&1
}

record_start() {
  mcp_call "shell_record_start" "{\"session_id\":\"$SHELL_SESSION\",\"fps\":30}" > /dev/null 2>&1
  sleep 0.3
}

# Extra pause before GIF loops
end_pause() { sleep 5; }

record_save() {
  local name="$1"
  local result
  result=$(mcp_call "shell_record_stop" "{\"session_id\":\"$SHELL_SESSION\",\"name\":\"$name\"}")
  local url
  url=$(echo "$result" | sed 's/\\"/"/g' | grep -o '"download_url":"[^"]*"' | cut -d'"' -f4)
  curl -s -o "$OUTPUT_DIR/$name.gif" "$url"
  echo "Saved $OUTPUT_DIR/$name.gif ($(du -h "$OUTPUT_DIR/$name.gif" | cut -f1))"
}

# Clear screen (starship is auto-loaded via .bashrc)
setup_prompt() {
  send "clear\\r"
  sleep 0.8
}
