package test

import (
	"SLC/src/reexec"
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"syscall"
)

func overlay2() {
	newRoot := "/tmp/ubuntu"
	cmd := reexec.Command("nsInit", "container", "-i")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &unix.SysProcAttr{
		Cloneflags: unix.CLONE_NEWNS |
			unix.CLONE_NEWUTS |
			unix.CLONE_NEWIPC |
			unix.CLONE_NEWPID |
			unix.CLONE_NEWNET |
			unix.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	cmd.Dir = newRoot

	fmt.Println("cmd Path: ", cmd.Path, "Args: ", cmd.Args)
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting the reexec.Command - %s\n", err)
		os.Exit(1)
	}
	fmt.Println("after cmd start")

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error waiting for the reexec.Command - %s\n", err)
		os.Exit(1)
	}
}
