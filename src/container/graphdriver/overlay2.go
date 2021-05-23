package graphdriver

import (
	"SLC/src/reexec"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"os"
	"os/exec"
	"syscall"
)

var (
	readonlyDir = "/tmp/image/ubuntu/lower"
	writeDir    = "/tmp/image/ubuntu/upper"
	mountDir    = "/tmp/image/ubuntu/merged"
	workdir     = "/tmp/image/ubuntu/work"
	//readonlyDir = "/lower"
	//writeDir = "/upper"
	//mountDir = "/merged"
	//workdir = "/work"
)

func init() {
	reexec.Register("slc-mountfrom", NewWorkSpace)
}

func ReexecMount() {
	cmd := reexec.Command("slc-mountfrom")

	cmd.SysProcAttr = &unix.SysProcAttr{
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

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error reexec the monut command - %s\n", err)
		os.Exit(1)
	}
}

func NewWorkSpace() {
	//createReadOnlyLayer()
	//createWritableLayer()
	createMountPoint()
}

func createReadOnlyLayer() {

}

func createWritableLayer() {
	if err := os.Mkdir(writeDir, 0777); err != nil {
		log.Errorf("creating writable layer in %s : %v", writeDir, err)
	}
}

func createMountPoint() {
	if err := os.Mkdir(mountDir, 0777); err != nil {
		log.Errorf("creating mount point in %s : %v", mountDir, err)
	}
	if err := os.Mkdir(workdir, 0777); err != nil {
		log.Errorf("creating workdir in %s : %v", workdir, err)
	}

	/*	if err := unix.Mount("", "/", "overlay", unix.MS_REC, ""); err != nil {
		log.Errorf("while mounting overlay2 in %s : %v", mountDir, err)
	}*/

	dir := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", readonlyDir, writeDir, workdir)

	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", dir, mountDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("overlay2 mounting : %v", err)
		os.Exit(1)
	}

	if err := os.Chdir("/merged"); err != nil {
		log.Errorf("while Chdir %v", err)
	}

	log.Infoln("createMountPoint finished")
}

func DeleteWorkSpace() {

}
