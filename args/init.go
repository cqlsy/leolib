package args

import (
	"fmt"
	"os"
)

// value:
// bool： if only this
type OnCatchOrder func(key string, value string) (bool, bool)

// Values without a key will be ignored
func GetArgs(f OnCatchOrder) bool {
	println(os.Args)
	var isOnlyGetArgs = false
	for index := 1; index < len(os.Args); index++ {
		value := os.Args[index]
		if value == "--help" {
			ShowArgsHelpStr()
			return true
		}
		var next = ""
		if len(os.Args) > index+1 {
			next = os.Args[index+1]
			if _, ok := argsList[next]; ok {
				// is a key to command
				if _, ok := argsList[next]; ok {
					// next is also a command
					next = ""
				} else {
					index++
				}
			} else {
				fmt.Println("UnKnown Args; you can learn more by --help")
				ShowArgsHelpStr()
				// jump out
				return true
			}
		}
		end, only := f(value, next)
		if !isOnlyGetArgs && only {
			isOnlyGetArgs = only
		}
		if end {
			return isOnlyGetArgs
		}
	}
	return isOnlyGetArgs
}

func ShowArgsHelpStr() {
	var str = ""
	for key, value := range argsList {
		str = fmt.Sprintf("%s[%s]: %s\n", str, key, value)
	}
	fmt.Println(str)
}

func AddArgsOne(key string, value string) {
	AddArgs(map[string]string{key: value})
}

func AddArgs(newArgs map[string]string) {
	for key, value := range newArgs {
		argsList[key] = value
	}
}

var argsList = map[string]string{
	"--help":       "Command list",
	"--init":       "Init program",
	"--configFile": "Configuration of program operation",
}
