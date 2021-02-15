package main

// /opt/X11/bin/xterm -hold -e %C
//TODO: add lock file

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

const SessionName string = "BURP_SEND_TO"

func command(cmd string) (string, error) {
	fmt.Printf("$ %s\n", cmd)
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", errors.WithMessagef(err, "command('%s')", cmd)
	}
	log.Printf("%s", string(out))
	return string(out), nil
}

func E(s string) string {
	return fmt.Sprintf("%#v", s)
}

func windowNum(s string) (int, error) {
	out, err := command(fmt.Sprintf("tmux list-windows -t %s | awk -F: '{print $1}'", E(s)))
	if err != nil {
		return -1, errors.WithStack(err)
	}
	winnumbrs := []int{}
	for _, v := range funk.FilterString(strings.Split(out, "\n"), func(x string) bool { return x != "" }) {
		num, err := strconv.Atoi(v)
		if err != nil {
			return 0, errors.WithStack(err)
		}
		winnumbrs = append(winnumbrs, num)
	}
	num := 1
	winnumbrs = append([]int{0}, winnumbrs...)
	winnumbrs = append(winnumbrs, math.MaxInt32)
	log.Printf("winnumbrs = %#v\n", winnumbrs)
	for i := 0; i < len(winnumbrs)-1; i++ {
		j := i + 1
		if winnumbrs[j]-winnumbrs[i] > 1 {
			num = winnumbrs[i] + 1
			break
		}
	}
	return num, nil
}

func main() {
	log.SetFlags(0)

	var err error
	var paramCommand string
	var killMe bool
	var PATH string
	flag.StringVar(&paramCommand, "c", "", "bash command")
	flag.BoolVar(&killMe, "kill", false, "kill tmux session")
	flag.StringVar(&PATH, "path", "/usr/local/opt/tmux/bin/", "PATH env")
	flag.Parse()

	os.Setenv("PATH", PATH+":"+os.Getenv("PATH"))

	if killMe {
		if _, err := command(fmt.Sprintf("tmux kill-session -t %s", E(SessionName))); err != nil {
			log.Print(err)
		}
		return
	}

	if paramCommand == "" {
		log.Fatal("usage: burpsendto2tmux -c \"ls\" ")
	}

	isNewSession := true
	if _, err := command(fmt.Sprintf("tmux new-session -d -s %s -n tmp", E(SessionName))); err != nil {
		isNewSession = false
		log.Print(err)
	}

	numWin := 1
	if !isNewSession {
		numWin, err = windowNum(SessionName)
		if err != nil {
			log.Fatalf("%+v\n", err)
		}
		log.Printf("numWin = %#v\n", numWin)
	}

	if _, err := command(
		fmt.Sprintf(
			"tmux new-window -t %s -n \"new\"",
			E(fmt.Sprintf("%s:%d", SessionName, numWin)),
		),
	); err != nil {
		log.Print(err)
	}

	if _, err := command(
		fmt.Sprintf(
			"tmux rename-window -t %s %s",
			E(fmt.Sprintf("%s:%d", SessionName, numWin)),
			E(strings.Fields(paramCommand)[0]),
		),
	); err != nil {
		log.Print(err)
	}

	if _, err := command(
		fmt.Sprintf(
			"tmux select-window -t %s ",
			E(fmt.Sprintf("%s:%d.0", SessionName, numWin)),
		),
	); err != nil {
		log.Print(err)
	}

	if _, err := command(
		fmt.Sprintf(
			"tmux send-keys -t %s %s C-m",
			E(fmt.Sprintf("%s:%d.0", SessionName, numWin)),
			E("cd /tmp"),
		),
	); err != nil {
		log.Fatalf("%+v\n", err)
	}

	if _, err := command(
		fmt.Sprintf(
			"tmux send-keys -t %s %s C-m",
			E(fmt.Sprintf("%s:%d.0", SessionName, numWin)),
			E(paramCommand),
		),
	); err != nil {
		log.Print(err)
	}

	//NOTE: tmux display-message -p '#S'
	currentSession, err := command(
		fmt.Sprintf("tmux display-message -p '#S'"),
	)
	currentSession = strings.TrimSuffix(currentSession, "\n")
	if err != nil {
		log.Print(err)
	}
	log.Printf("currentSession = %#v\n", currentSession)

	if _, err := command(
		fmt.Sprintf(
			"tmux switch -t %s ",
			E(fmt.Sprintf("%s", SessionName)),
		),
	); err != nil {
		log.Print(err)
	}

	//comeback
	if currentSession != SessionName {
		if _, err := command(
			fmt.Sprintf(
				"tmux send-keys -t %s %s",
				E(fmt.Sprintf("%s:%d.0", SessionName, numWin)),
				E(fmt.Sprintf("tmux switch -t %s", E(currentSession))),
			),
		); err != nil {
			log.Print(err)
		}
	}
}
