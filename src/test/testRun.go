package test

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"os/exec"
)

func Run() {
	cmd := exec.Command("sh")
	cmd.SysProcAttr = &unix.SysProcAttr{
		Cloneflags: unix.CLONE_NEWUTS | unix.CLONE_NEWIPC |
			unix.CLONE_NEWPID | unix.CLONE_NEWNS |
			unix.CLONE_NEWUSER | unix.CLONE_NEWNET,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}



