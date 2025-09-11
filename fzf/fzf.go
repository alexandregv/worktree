package fzf

import fzfLib "github.com/junegunn/fzf/src"

func InitFzfOptions(inputs []string) (options *fzfLib.Options, err error) {
	options, err = fromDefaults(
		inputs,
		[]string{
			"--height=40%",
			"--prompt=worktree: ",
			"--with-nth=2..",
			"--layout=reverse",

			// Jump to a line with Space + <n>
			"--bind=space:jump",
			"--jump-labels=0123456789;:,<.>/?'\"!@#$%^&*",

			// Preview window with git log on the right
			// sh -c '' is used to expand $HOME
			"--preview=sh -c \"git -C {3} log --color=always --oneline --graph --decorate --all -n20\"",
			"--preview-window=right,45%,<90(bottom,45%)",

			// Theme: Catppuccin Macchiato (https://github.com/catppuccin/fzf)
			"--color=bg+:#363a4f,bg:#24273a,spinner:#f4dbd6,hl:#ed8796",
			"--color=fg:#cad3f5,header:#ed8796,info:#c6a0f6,pointer:#f4dbd6",
			"--color=marker:#b7bdf8,fg+:#cad3f5,prompt:#c6a0f6,hl+:#ed8796",
			"--color=selected-bg:#494d64",
		},
	)
	return options, err
}

func fromDefaults(inputs []string, customOptions []string) (options *fzfLib.Options, err error) {
	options, err = fzfLib.ParseOptions(true, customOptions)
	if err != nil {
		return nil, err
	}

	options.Input = make(chan string)
	options.Output = make(chan string)

	go func() {
		defer close(options.Input)

		for _, input := range inputs {
			options.Input <- input
		}
	}()

	return options, err
}
