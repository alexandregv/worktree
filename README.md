# 🗂️ worktree

Familar with Git Worktrees? No? Well, `worktree` is here to make it easy!

`worktree` is a CLI utility to easily navigate between [Git Worktrees](https://git-scm.com/docs/git-worktree), list them and clone a multiple-worktrees-enabled repo.  
That's it. No overkill feature, no you-will-never-be-able-to-work-without-it-anymore feature. Just a bit of handiness.

### Dependencies

Only `git` is required.

> [!NOTE]  
> `worktree` uses `fzf` under the hood, but it's built in (Go library), so you don't need it installed.

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
  output=$(worktree "$@")
  if [[ $? == 0 ]] && [[ "$output" == /* ]]; then
    cd "$output"
  fi
  printf "$output\n"
}
```

Then source the file or run `exec bash` / `exec zsh` to restart your shell.

</details>

<details>
  <summary>Fish</summary>

Add this function in your `~/.config/fish/config.fish`:

```sh
# https://github.com/alexandregv/worktree#setup
function wt
  set output (worktree $argv)
  if string match -q '/*' $output
    cd $output
  end
  printf "$output\n"
end
```

Then source the file or run `exec fish` to restart your shell.

</details>

<details>
  <summary>Nu Shell</summary>

Add this function in your `~/.config/nushell/config.nu`:

```sh
# https://github.com/alexandregv/worktree#setup
def --env --wrapped wt [...args] {
  let cmd = (worktree ...$args | complete)
  if $cmd.exit_code == 0 and ($cmd.stdout | str starts-with "/") {
    cd $cmd.stdout
  }
  printf $cmd.stdout
}
```

Then source the file or run `exec nu` to restart your shell.

</details>

### Usage

```sh
wt help          ## Help
wt               ## Chose a worktree via TUI and switch to it
wt list          ## List worktrees
wt switch <name> ## Quickly switch to a named worktree
wt clone  <link> ## Clone a repository, following the `.bare` + worktrees convention
```

### Building

Run `make build` or just `go build .` at the root of the directory.  
See `make help` for more commands.

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
