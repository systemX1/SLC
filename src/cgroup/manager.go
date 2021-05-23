package cgroups

import (
	"SLC/src/cgroup/subsystem"
	"github.com/sirupsen/logrus"
)

type CGroupManager struct {
	Path string
}

func NewCGroupManager(path string) *CGroupManager {
	return &CGroupManager{Path: path}
}

func (c *CGroupManager) Set(res *subsystem.ResourceConfig) {
	for _, sub := range subsystem.Subsystems {
		err := sub.Set(c.Path, res)
		if err != nil {
			logrus.Errorf("set %s err: %v", sub.Name(), err)
		}
	}
}

func (c *CGroupManager) Apply(pid int) {
	for _, sub := range subsystem.Subsystems {
		err := sub.Apply(c.Path, pid)
		if err != nil {
			logrus.Errorf("apply task, err: %v", err)
		}
	}
}

func (c *CGroupManager) Destroy() {
	for _, sub := range subsystem.Subsystems {
		err := sub.Remove(c.Path)
		if err != nil {
			logrus.Errorf("remove %s err: %v", sub.Name(), err)
		}
	}
}
