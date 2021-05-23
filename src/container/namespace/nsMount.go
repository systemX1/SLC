package namespace

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"path/filepath"
)

// PivotRoot must be called from within the new Mount namespace, otherwise we'll end up changing the host's '/' which is not the intention
func pivotRoot() error {
	newRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Get pwd error %v\n", err)
	}

	// 声明新的mount namespace独立
	if err := unix.Mount("", "/", "", unix.MS_PRIVATE|unix.MS_REC, ""); err != nil {
		return err
	}

	// bind mount new_root to itself - this is a slight hack needed to satisfy requirement (2)
	if err := unix.Mount(newRoot, newRoot, "bind", unix.MS_BIND|unix.MS_REC, ""); err != nil {
		return fmt.Errorf("mount newRoot %s to itself error: %v", newRoot, err)
	}

	// create putOld directory
	putOld := filepath.Join(newRoot, "/.pivot_root")
	if err := os.MkdirAll(putOld, 0777); err != nil {
		return fmt.Errorf("creating putOld directory %v", err)
	}

	// The following restrictions apply to new_root and put_old:
	// 1.  They must be directories.
	// 2.  new_root and put_old must not be on the same filesystem as the current root.
	// 3.  put_old must be underneath new_root, that is, adding a nonzero number of /.. to the string pointed to by put_old must yield the same directory as new_root.
	// 4.  No other filesystem may be mounted on put_old.
	// see https://man7.org/linux/man-pages/man2/pivot_root.2.html

	if err := unix.PivotRoot(newRoot, putOld); err != nil {
		return fmt.Errorf("syscalling PivotRoot %v", err)
	}

	// Note that this also applies to the calling process: pivotRoot() may
	// or may not affect its current working directory.  It is therefore
	// recommended to call chdir("/") immediately after pivotRoot().
	if err := os.Chdir("/"); err != nil {
		return fmt.Errorf("while Chdir %v", err)
	}

	// umount putOld, which now lives at .pivot_root
	putOld = "/.pivot_root"
	if err := unix.Unmount(putOld, unix.MNT_DETACH); err != nil {
		return fmt.Errorf("while unmount putOld %v", err)
	}

	// remove put_old
	if err := os.RemoveAll(putOld); err != nil {
		return fmt.Errorf("while remove putOld %v", err)
	}

	return nil
}

func mountProc() error {
	// systemd加入linux之后, mount namespace变成shared by default, 必须显式声明这个新的mount namespace独立。
	// MS_PRIVATE Make this mount point private. Mount and unmount events do not propagate into or out of this mount point. MS_REC recursive递归
	err := unix.Mount("", "/", "", unix.MS_PRIVATE|unix.MS_REC, "")
	if err != nil {
		return fmt.Errorf("while making mount namespace private: %v", err)
	}

	// MS_NOEXEC Do not allow other programs to be executed from this filesystem.
	// MS_NOSUID Do not honor set-user-ID and set-group-ID bits or file capabilities when executing programs from this filesystem.
	// MS_NODEV Do not allow access to devices (special files) on this filesystem.

	if err := unix.Mount("proc", "/proc", "proc", unix.MS_NOSUID|unix.MS_NODEV|unix.MS_NOEXEC|unix.MS_RELATIME, ""); err != nil {
		return fmt.Errorf("while mount proc error: %v", err)
	}

	if err := unix.Mount("tmpfs", "/dev", "tmpfs", unix.MS_NOSUID|unix.MS_STRICTATIME, "mode=755"); err != nil {
		return fmt.Errorf("while mount tmpfs error: %v", err)
	}

	if err := unix.Mount("sysfs", "/sys", "sysfs", unix.MS_NOSUID|unix.MS_NODEV|unix.MS_RELATIME|unix.MS_NOEXEC|unix.MS_RDONLY, "mode=755"); err != nil {
		return fmt.Errorf("while mount tmpfs error: %v", err)
	}

	if err := unix.Mount("udev", "/dev", "tmpfs", unix.MS_NOSUID|unix.MS_RELATIME|unix.MS_NOEXEC, "mode=755"); err != nil {
		return fmt.Errorf("while mount tmpfs error: %v", err)
	}

	if err := unix.Mount("devpts", "/dev/pts", "devpts", unix.MS_NOSUID|unix.MS_RELATIME|unix.MS_NOEXEC, "mode=620, ptmxmode=666"); err != nil {
		return fmt.Errorf("while mount tmpfs error: %v", err)
	}

	if err := unix.Mount("devpts", "/dev/console", "devpts", unix.MS_NOSUID|unix.MS_RELATIME|unix.MS_NOEXEC, "mode=620, ptmxmode=666"); err != nil {
		return fmt.Errorf("while mount tmpfs error: %v", err)
	}

	return nil
}
