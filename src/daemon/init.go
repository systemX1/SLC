package daemon

import (
	"SLC/src/reexec"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func init() {
	fmt.Println("func init")
	fmt.Println("0: ", os.Args[0])
	reexec.Register("nsInit", nsInit)
	if reexec.Init(os.Args[0]) {
		fmt.Println("if reexec init")
		os.Exit(0)
	}
}

func nsInit() {
	fmt.Println("func nsInit")
	fmt.Println("0: ", os.Args[0], "1: ", os.Args[1])
	//newrootPath := os.Args[1]

	//if err := mountProc(newrootPath); err != nil {
	//	fmt.Printf("Error mounting /proc - %s\n", err)
	//	os.Exit(1)
	//}
	//
	//if err := pivotRoot(newrootPath); err != nil {
	//	fmt.Printf("Error running pivot_root - %s\n", err)
	//	os.Exit(1)
	//}

	if err := syscall.Sethostname([]byte("ns-process")); err != nil {
		fmt.Printf("Error setting hostname - %s\n", err)
		os.Exit(1)
	}

	//if err := waitForNetwork(); err != nil {
	//	fmt.Printf("Error waiting for network - %s\n", err)
	//	os.Exit(1)
	//}

	nsRun()
}

func nsRun() {
	fmt.Println("func nsRun")
	cmd := exec.Command("/bin/sh")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = []string{"slc # "}

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running the /bin/sh command - %s\n", err)
		os.Exit(1)
	}
}

func Run() {
	cmd := reexec.Command("nsInit", "/tmp/ns-process/rootfs")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,
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
	fmt.Println("after cmd wait")
}








