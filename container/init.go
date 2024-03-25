package container

// RunContainerInitProcess 启动容器的init进程
/*
这里的init函数是在容器内部执行的，也就是说，代码执行到这里后，容器所在的进程其实就已经创建出来了，
这是本容器执行的第一个进程。
使用mount先去挂载proc文件系统，以便后面通过ps等系统命令去查看当前进程资源的情况。
*/
import (
	"errors"
	"os"
	"os/exec"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func RunContainerInitProcess(command string, args []string) error {
	log.Infof("Init command:%s", command)
	setUpMount()

	// 从 pipe 中读取命令
	cmdArray := readUserCommand()
	if len(cmdArray) == 0 {
		return errors.New("run container get user command error, cmdArray is nil")
	}

	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		log.Errorf("Exec loop path error %v", err)
		return err
	}

	log.Infof("Find path %s", path)
	if err = syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		log.Errorf("RunContainerInitProcess exec :" + err.Error())
	}

	return nil
}
