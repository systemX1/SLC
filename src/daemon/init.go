package daemon

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"os"
	"os/exec"
	"syscall"
)

func RunParentProcess(tty bool)  {
	// re-run itself
	cmd := exec.Command("/proc/self/exe", "daemon", "-i")
	cmd.SysProcAttr = &unix.SysProcAttr{
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
		// if the value of GidMappingsEnableSetgroups is false as default, the child process will have no permission to use setgroups syscall regardless of whether it has root privileges.
		GidMappingsEnableSetgroups: true,
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	readPipe, _, err := os.Pipe()
	if err != nil {
		fmt.Println(err)
	}
	cmd.ExtraFiles = []*os.File{readPipe}



	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}


}

func NewInitProcess(tty bool)  {
	cmd := exec.Command("sh")
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := setUpMount(); err != nil {
		fmt.Println(err)
	}

	if err := syscall.Sethostname([]byte("slc")); err != nil {
		fmt.Println(err)
	}

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}



	//if err := syscall.Exec(cmd, "", os.Environ()); err != nil {
	//	fmt.Println(err)
	//}

}

func setUpMount() error {
	// systemd 加入linux之后, mount namespace 就变成 shared by default, 必须显式声明这个新的mount namespace独立。
	// MS_PRIVATE Make this mount point private. Mount and unmount events do not propagate into or out of this mount point. MS_REC recursive 递归.
	err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	if err != nil {
		return err
	}

	// mount proc
	// MS_NOEXEC Do not allow other programs to be executed from this filesystem.
	// MS_NOSUID Do not honor set-user-ID and set-group-ID bits or file capabilities when executing programs from this filesystem.
	// MS_NODEV Do not allow access to devices (special files) on this filesystem.
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		logrus.Errorf("while mount proc, err: %v", err)
		return err
	}




	return nil
}

// pivot_root must be called from within the new Mount namespace, otherwise we'll end up changing the host's / which is not the intention
func setPivotRoot()  {
	// get working dir
	root, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("current location is %s\n", root)

	syscall.Mount("rootfs", "rootfs", "", syscall.MS_BIND, "")
	os.MkdirAll("rootfs/oldrootfs", 0700)
	syscall.PivotRoot("rootfs", "rootfs/oldrootfs")
	os.Chdir("/")
}


