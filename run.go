package main

import (
	"main/cgroups"
	"main/container"
	"main/subsystem"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func Run(tty bool, cmd []string, config *subsystem.ResourceConfig, volume string) {
	parent, writePipe := container.NewParentProcess(tty, volume)
	if parent == nil {
		log.Errorf("parent not exist")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	if parent.Process == nil {
		log.Error("parent process is nil")
		return
	}
	log.Infof("parent pid : %v", parent.Process.Pid)
	// 创建cgroup manager, 并通过调用set和apply设置资源限制并使限制在容器上生效
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	cgroupManager.Resource = config
	defer cgroupManager.Destroy()
	_ = cgroupManager.Set(config)
	_ = cgroupManager.Apply(parent.Process.Pid)
	// 在子进程创建后通过管道来发送参数
	sendInitCommand(cmd, writePipe)

	if tty {
		_ = parent.Wait()
		container.DeleteWorkSpace(container.RootURL, container.MntURL, volume)
	}

	//os.Exit(-1)
}

func sendInitCommand(commands []string, writePipe *os.File) {
	command := strings.Join(commands, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
