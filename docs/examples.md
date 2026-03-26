# Examples

## Claude Code + tmux integration

Launch a [Claude Code](https://docs.anthropic.com/en/docs/claude-code) session in tmux for every worktree automatically:

![wt claude-tmux](wt-claude-tmux.gif)

```toml
[hooks]
post_create = [
  "tmux new-session -d -s \"$WT_REPO_NAME/$WT_BRANCH\" -c \"$WT_PATH\" \"claude -n '$WT_REPO_NAME/$WT_BRANCH'\" 2>/dev/null; echo \"tmux session: $WT_REPO_NAME/$WT_BRANCH\""
]
post_checkout = [
  "tmux new-session -d -s \"$WT_REPO_NAME/$WT_BRANCH\" -c \"$WT_PATH\" \"claude -n '$WT_REPO_NAME/$WT_BRANCH'\" 2>/dev/null; echo \"tmux session: $WT_REPO_NAME/$WT_BRANCH\""
]
pre_remove = [
  "tmux kill-session -t \"$WT_REPO_NAME/$WT_BRANCH\" 2>/dev/null || true"
]
```

Claude Code also has built-in worktree support via `claude --worktree --tmux`, but it uses a fixed directory layout. The hooks approach lets you keep `wt`'s configurable strategies and naming conventions.

## Task spanning multiple repositories

When a task or story requires changes across multiple repositories (e.g. a shared library and a main application), you can organize worktrees by feature instead of by repo using a custom pattern:

![wt multi-repo](wt-multi-repo.gif)

```toml
# ~/.config/wt/config.toml
strategy = "custom"
pattern = "{.worktreeRoot}/{.branch}/{.repo.Name}"
```

Use the same branch name in each repository:

```bash
cd ~/src/shared-lib
wt create feat/PROJ-123

cd ~/src/main-app
wt create feat/PROJ-123
```

This groups all repositories for a task together:

```
~/dev/worktrees/
  feat/PROJ-123/
    shared-lib/
    main-app/
```

## Using environment variables in patterns

You can reference any environment variable in your pattern with `{.env.VARNAME}`. This lets you group worktrees by an external value such as a feature name, ticket ID, or sprint.

```toml
# ~/.config/wt/config.toml
strategy = "custom"
pattern = "{.worktreeRoot}/{.env.FEATURE}/{.repo.Name}"
```

```bash
export FEATURE=PROJ-42-new-checkout

cd ~/src/frontend
wt create main

cd ~/src/backend
wt create main
```

```
~/dev/worktrees/
  PROJ-42-new-checkout/
    frontend/
    backend/
```

If the referenced environment variable is not set, `wt` will return an error.

## Deterministic dev server port per worktree

When running multiple worktrees simultaneously, each dev server needs a unique port. Use a post-checkout hook to compute a deterministic port offset from the branch name:

```toml
[hooks]
post_create = [
  "printf 'PORT=%d\\n' $(( 3000 + $(printf '%s' \"$WT_BRANCH\" | cksum | cut -d' ' -f1) % 997 )) > $WT_PATH/.env.port"
]
post_checkout = [
  "printf 'PORT=%d\\n' $(( 3000 + $(printf '%s' \"$WT_BRANCH\" | cksum | cut -d' ' -f1) % 997 )) > $WT_PATH/.env.port"
]
```

This hashes the branch name with `cksum` and maps it to a port in the range 3000–3996. The result is stable — the same branch always gets the same port.

Then read it in your dev server config (e.g. `vite.config.ts`):

```js
import { readFileSync } from 'fs';

const port = (() => {
  try {
    const env = readFileSync('.env.port', 'utf8');
    const match = env.match(/PORT=(\d+)/);
    return match ? parseInt(match[1]) : 3000;
  } catch { return 3000; }
})();

export default { server: { port } };
```

Or source it in a shell script:

```bash
source .env.port 2>/dev/null || PORT=3000
echo "Starting dev server on port $PORT"
```

## Shared build cache across worktrees

Each worktree normally gets its own `node_modules`, `.venv`, or build output directory. This wastes disk space and install time. Use a hook to symlink these from a shared per-repo cache:

```toml
[hooks]
# Node.js — share node_modules via a central cache per repo
post_create = [
  "mkdir -p $HOME/.cache/wt/$WT_REPO_NAME/node_modules && ln -sf $HOME/.cache/wt/$WT_REPO_NAME/node_modules $WT_PATH/node_modules && cd $WT_PATH && npm install"
]

# Python — share a single venv per repo
# post_create = [
#   "mkdir -p $HOME/.cache/wt/$WT_REPO_NAME/venv && ln -sf $HOME/.cache/wt/$WT_REPO_NAME/venv $WT_PATH/.venv && cd $WT_PATH && uv sync"
# ]

# Rust — share the target directory
# post_create = [
#   "mkdir -p $HOME/.cache/wt/$WT_REPO_NAME/target && ln -sf $HOME/.cache/wt/$WT_REPO_NAME/target $WT_PATH/target"
# ]

# Terraform — share the .terraform directory (providers, modules)
# post_create = [
#   "mkdir -p $HOME/.cache/wt/$WT_REPO_NAME/.terraform && ln -sf $HOME/.cache/wt/$WT_REPO_NAME/.terraform $WT_PATH/.terraform && cd $WT_PATH && terraform init"
# ]
```

All worktrees for the same repo point to one `node_modules` (or `.venv`, or `target/`). The first `npm install` populates the cache; subsequent worktrees reuse it instantly.

> **Note:** Shared mutable caches work well for dependencies that are branch-independent. If branches use different dependency versions, consider per-branch cache keys instead:
> ```bash
> mkdir -p $HOME/.cache/wt/$WT_REPO_NAME/$WT_BRANCH/node_modules
> ```
