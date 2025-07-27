package runner

func run(command string, args ...string) {
	switch command {
	case "/state":
		fmt.Println(args)
	}
}