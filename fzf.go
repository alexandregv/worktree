package main

import fzf "github.com/junegunn/fzf/src"

func initFzfOptions(inputs []string, customOptions []string) (options *fzf.Options, err error) {
	options, err = fzf.ParseOptions(true, customOptions)
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
