package runner

func Run(command string, args ...string) {
	switch command {
	case "/state":
		fmt.Println(args)
	}
}