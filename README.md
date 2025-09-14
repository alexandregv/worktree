# 🗂️ worktree

[![Go Reference](https://pkg.go.dev/badge/github.com/alexandregv/worktree.svg)](https://pkg.go.dev/github.com/alexandregv/worktree)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexandregv/worktree)](https://goreportcard.com/report/github.com/alexandregv/worktree)

Familar with Git Worktrees? No? Well, `worktree` is here to make it easy!

`worktree` is a CLI utility to easily navigate between [Git Worktrees](https://git-scm.com/docs/git-worktree), list them and clone a multiple-worktrees-enabled repo.  
That's it. No overkill, no you-will-never-be-able-to-work-without-it-anymore feature. Just a bit of handiness.

![Demo GIF](/assets/demo.gif)

<!--TOC-->

- [Dependencies](#dependencies)
- [Download](#download)
- [Setup](#setup)
- [Usage](#usage)
- [Building](#building)
- [FAQ](#faq)
- [Similar projects](#similar-projects)
- [Contributing](#contributing)
- [Contributors](#contributors)
- [Stargazers over time](#stargazers-over-time)

<!--TOC-->

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
  if test $status -eq 0; and string match -q '/*' $output
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
wt help                 ## Help
wt                      ## Chose a worktree via TUI and switch to it
wt list                 ## List worktrees (alias ls)
wt switch <name>        ## Quickly switch to a named worktree (alias s, cd)
wt clone  <link> [path] ## Clone a repository, with worktrees at the root folder
```

Typical flow:
1. Clone a repo with `wt clone -a <url>`, creating a directory for each branch
2. Navigate between worktrees with `wt` or `wt cd <worktree>`
3. If needed, create a new worktree with `git worktree add <path> -B <branch> [base-branch]`

### Building

Run `make build` or just `go build .` at the root of the directory.  
See `make help` for more commands.

### FAQ

1. What is a _multiple-worktrees-enabled repo_?
   => By default, a git repository has only one worktree ("active" branch). You can have multiple worktrees, and `wt clone --all <url>` will clone every branch in its own directory, at the root of your repository.
2. Do I need `fzf` installed?  
   => No, it's built in `worktree`.
3. What if I don't want to create the shell function?  
   => You will be able to use all commands except switching to a worktree. As a workaround, use `cd $(worktree)` and `cd $(worktree switch <path>)` for this matter.
4. Why is there no `wt add` (or similar) command?  
   => The goal for `worktree` is to help you use git worktrees, while not fully replacing basic git commands. You should be able to use git worktrees even if `worktree` is not available. Creating a new worktree is simple: `git worktree add <path> -B <branch> [base-branch]`. On the other hand, `wt clone <url>` replaces `git clone --no-checkout <url>; cd <path>; git config core.bare true; for each branch -> git worktree add $(basename <branch>) -B <branch>`. A bit too cumbersome, so `worktree` implements it.

### Similar projects

- https://github.com/3rd/work/
- https://github.com/yankeexe/git-worktree-switcher
- https://github.com/davvid/gcd
- https://github.com/egyptianbman/zsh-git-worktrees

### Contributing

Contributions are welcome!  
Please make sure you use [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/#summary) and branch names like `feat/<my-feature>`, `fix/<my-fix>`, etc. **`make commit` will help you with that**.  
This repository also provides a [pre-commit](https://pre-commit.com/) hook, **which you can setup with `make pre-commit-install`**. You still need to install the pre-commit software on your own.

1. Fork the repository (<https://github.com/alexandregv/worktree/fork>)
2. Create your branch (`git checkout -b feat/my-feature`, `git checkout -b fix/list-width`, etc)
3. Commit your changes (`git add -p && git commit -m 'fix(list): use correct width'`)
4. Push to the branch (`git push origin fix/list-width`)
5. Create a new Pull Request on GitHub

### Contributors

- [alexandregv](https://github.com/alexandregv) - Creator and maintainer

### Stargazers over time

[![Stargazers over time](https://starchart.cc/alexandregv/worktree.svg?variant=adaptive)](https://starchart.cc/alexandregv/worktree)
