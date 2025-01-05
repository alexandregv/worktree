# worktree

Git worktrees utility

### Dependencies

There are no dependencies, not either `git` or `fzf`.

### Download

Download from the [Releases page](https://github.com/alexandregv/worktree/releases/latest) and place it in your `$PATH`.

Or, if you have `go` installed:

```sh
go install github.com/alexandregv/worktree@latest
```

This will install `worktree` in `$GOBIN`, make sure this value is in your `$PATH`.

### Setup

Creating a shell function is **required** to allow `worktree` changing the current directory (only your shell can do so).

<details>
  <summary>Bash, Zsh</summary>

Add this function in your `~/.bashrc` or `~/.zshrc`:

```sh
# https://github.com/alexandregv/worktree#setup
function wt() {
  output=$(worktree)
  printf "$output\n"
  if [[ $? == 0 ]] && [[ "$output" == /* ]]; then
    cd "$output"
  fi
}
```

</details>

<details>
  <summary>Fish</summary>

Add this function in your `~/.config/fish/config.fish`:

```sh
# https://github.com/alexandregv/worktree#setup
function wt
  set output (worktree)
  printf "$output\n"
  if string match -q '/*' $output
    cd $output
  end
end
```

</details>

<details>
  <summary>Nu Shell</summary>

Add this function in your `~/.config/nushell/config.nu`:

```sh
# https://github.com/alexandregv/worktree#setup
def --env wt [] {
  let cmd = (worktree | complete)
  printf $cmd.stdout
  if $cmd.exit_code == 0 and ($cmd.stdout | str starts-with "/") {
    cd $cmd.stdout
  }
}
```

</details>

Then source the file or run `exec $SHELL` to restart your shell.

### Usage

```sh
wt help    ## Help
wt clone   ## Clone a repository, following the `.bare` + worktrees convention
wt list    ## List worktrees
wt <name>  ## Quickly switch to a named worktree
wt         ## Chose a worktree via TUI and switch to it
```

If a local worktree has the same name than a wt command, use `wt -- <name>` instead of `wt <name>` to quickly switch to it.

### FAQ

1. Do I need `fzf` installed?  
   => No, it's built in `worktree`.
2. What if I don't want to create the shell function?  
   => You will be able to use all commands except switching to a worktree.

### Similar projects

- https://github.com/3rd/work/
- https://github.com/yankeexe/git-worktree-switcher
- https://github.com/davvid/gcd
- https://github.com/egyptianbman/zsh-git-worktrees
