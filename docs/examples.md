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
