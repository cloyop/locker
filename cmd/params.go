package cmd

func validAction(action *string) bool {
	if *action != "get" &&
		*action != "set" &&
		*action != "rm" &&
		*action != "ls" &&
		*action != "exit" &&
		*action != "save" &&
		*action != "clear" {
		return false
	}
	return true
}
