package daemon

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func RunParentProcess(cmds []string, tty bool)  {
	// re-run itself
	cmd := exec.Command("/proc/self/exe", "daemon")
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

	getInfo("RunParentProcess")

	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		fmt.Println(err)
	}
	cmd.ExtraFiles = []*os.File{readPipe}
	sendStartDaemonCommand(cmds, writePipe)

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}

func sendStartDaemonCommand(comArray []string, writePipe *os.File) {
	cmd := strings.Join(comArray, " ")
	_, err := writePipe.WriteString(cmd)
	if err != nil {
		fmt.Println(err)
	}

	err = writePipe.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	defer func(pipe *os.File) {
		err := pipe.Close()
		if err != nil {

		}
	}(pipe)
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		fmt.Printf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)
	fmt.Println("msg reci:", msgStr)
	return strings.Split(msgStr, " ")
}

func NewInitProcess(tty bool)  {
	cmds := readUserCommand()

	if err := unix.Sethostname([]byte("slc")); err != nil {
		fmt.Println(err)
	}

	if err := setUpMount(); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("cmd all is %s\n", cmds)

	if err := unix.Exec(cmds[0], cmds[1:], os.Environ()); err != nil {
		fmt.Println(err)
	}
}

func setUpMount() error {
	//PivotRoot()

	// systemd加入linux之后, mount namespace变成shared by default, 必须显式声明这个新的mount namespace独立。
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

// PivotRoot must be called from within the new Mount namespace, otherwise we'll end up changing the host's '/' which is not the intention
func PivotRoot()  {
	getInfo("PivotRoot")


	// get working dir
	newRoot, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	putOld := filepath.Join(newRoot, "/.pivot_root")

	// 声明新的mount namespace独立。
	if err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
		fmt.Println(err)
	}

	// bind mount new_root to itself - this is a slight hack needed to satisfy requirement (2)
	//
	// The following restrictions apply to new_root and put_old:
	// 1.  They must be directories.
	// 2.  new_root and put_old must not be on the same filesystem as the current root.
	// 3.  put_old must be underneath new_root, that is, adding a nonzero
	//     number of /.. to the string pointed to by put_old must yield the same directory as new_root.
	// 4.  No other filesystem may be mounted on put_old.
	if err := syscall.Mount(newRoot, newRoot, "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("syscall.Mount(%s, %s, \"\", syscall.MS_BIND|syscall.MS_REC, \"\") failed\n", newRoot, newRoot))
		fmt.Println(err)
	}

	// umask
	mask := syscall.Umask(0)
	defer syscall.Umask(mask)

	// create put_old directory
	if err := os.MkdirAll(putOld, 0700); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("os.MkdirAll(%s, 0700) failed\n", putOld))
		fmt.Println(err)
	}

	// call pivotRoot
	if err := syscall.PivotRoot(newRoot, putOld); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("syscall.PivotRoot(%s, %s) failed\n", newRoot, putOld))
		fmt.Println(err)
	}

	// Note that this also applies to the calling process: pivotRoot() may
	// or may not affect its current working directory.  It is therefore
	// recommended to call chdir("/") immediately after pivotRoot().
	if err := os.Chdir("/"); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("os.Chdir(\"/\") failed\n"))
		fmt.Println(err)
	}

	// umount put_old, which now lives at /.pivot_root
	putOld = "/.pivot_root"
	if err := syscall.Unmount(putOld, syscall.MNT_DETACH); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("syscall.Unmount(%s, syscall.MNT_DETACH) failed\n", putOld))
		fmt.Println(err)
	}

	// remove put_old
	if err := os.RemoveAll(putOld); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("os.RemoveAll(%s) failed\n", putOld))
		fmt.Println(err)
	}
}

func getInfo(s string) {
	fmt.Printf("%s:\n", s)
	fmt.Printf("Username: %s\n", os.Getenv("USER"))
	pwd, _ := unix.Getwd()
	fmt.Printf("Present working directory: %s\n", pwd)
	fmt.Println("Gid: ", unix.Getegid(), "Uid: ", unix.Geteuid())
}
