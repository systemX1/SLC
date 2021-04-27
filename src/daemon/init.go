package daemon

import (
	"SLC/src/reexec"
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"syscall"
)

func Init(cmds []string, tty bool) {
	fmt.Println("func init")

	reexec.Register("nsInit", nsInit)
	if reexec.Init(os.Args[0]) {
		os.Exit(0)
	}

	Run(cmds, tty)
}

func nsInit() {
	fmt.Println("func nsInit")


	cmds := readCommand()
	fmt.Println("cmds reci: ", cmds)

	if err := pivotRoot(); err != nil {
		fmt.Printf("Error running pivot_root - %s\n", err)
		os.Exit(1)
	}

	if err := mountProc(); err != nil {
		fmt.Printf("Error mounting - %s\n", err)
		os.Exit(1)
	}

	if err := unix.Sethostname([]byte("slc")); err != nil {
		fmt.Printf("Error setting hostname - %s\n", err)
		os.Exit(1)
	}

	//if err := waitForNetwork(); err != nil {
	//	fmt.Printf("Error waiting for network - %s\n", err)
	//	os.Exit(1)
	//}

	nsRun(cmds)
}

func nsRun(cmds []string) {
	fmt.Println("func nsRun")

	if err := unix.Exec("/bin/sh", nil, os.Environ()); err != nil {
		fmt.Println("Execve", err)
	}
}

func Run(cmds []string, tty bool) {
	newRoot := "/tmp/ubuntu"
	cmd := reexec.Command("nsInit", "daemon", "-i")

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

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

	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		fmt.Println(err)
	}
	cmd.ExtraFiles = []*os.File{readPipe}
	sendCommand(cmds, writePipe)

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








