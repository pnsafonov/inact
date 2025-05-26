package main

import (
	"fmt"
	"github.com/pnsafonov/inact/core/utils/git_utils"
	"github.com/pnsafonov/inact/logic"
	"os"
	"strconv"
)

func main() {
	doMain(os.Args)
}

func doMain(args []string) {

	days0 := "7"
	mins0 := "0"
	verbose := true
	l0 := len(args)
	for i := 1; i < l0; i++ {
		arg := args[i]
		switch arg {
		case "-v", "--version":
			{
				printVersion()
			}
		case "-h", "--help":
			{
				printHelp()
			}
		case "-d", "--days":
			{
				i0 := i + 1
				if i0 < l0 {
					days0 = args[i0]
				}
				i++
				continue
			}
		case "-m", "--mins":
			{
				i0 := i + 1
				if i0 < l0 {
					mins0 = args[i0]
				}
				i++
				continue
			}
		case "--verbose":
			{
				verbose = true
				continue
			}
		case "--no-verbose":
			{
				verbose = false
				continue
			}
		}
	}

	days1, err := strconv.Atoi(days0)
	if err != nil {
		days1 = 7
	}

	mins1, err := strconv.Atoi(mins0)
	if err != nil {
		mins1 = 0
	}

	ctx := &logic.Context{
		Days:    days1,
		Mins:    mins1,
		Verbose: verbose,
	}

	err = logic.DoRun(ctx)

	code := 0
	if err != nil {
		code = 1
	}

	os.Exit(code)
}

func printHelp() {
	helpMsg := `Usage: inact [OPTIONS]
inact checks last logins and do shutdown if no recent logins

    -d, --days <days>     days count before current day to check for logins, default 7
    -m, --mins <mins>     mins count before current time to check for logins
        --verbose               print information about last logins
        --no-verbose      don't print information about last logins

    -h, --help            display this help and exit
    -v, --version         output version information and exit`

	fmt.Println(helpMsg)
	os.Exit(0)
}

func printVersion() {
	gitHash := git_utils.GetGitHashShort()
	fmt.Printf("inact %s %s\n", logic.Version, gitHash)
	os.Exit(0)
}
