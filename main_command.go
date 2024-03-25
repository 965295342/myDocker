package main

import (
	"fmt"
	"main/container"
	"main/subsystem"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name:  "run",
	Usage: "Create a container with namespace and cgroups limit : mydocker run -it [command]",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "mem",
			Usage: "memory limit ,e.g.: -mem 100m",
		},
		cli.StringFlag{
			Name:  "cpu",
			Usage: "cpu quota,e.g.:-cpu 100",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit ,e.g.:-cpuset 2,4",
		},
	},
	/*
		这里是run命令执行的真正函数。
		1.判断参数是否包含command
		2.获取用户指定的command
		3.调用Run function去准备启动容器:
	*/
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("missing container command")
		}
		cmd := make([]string, 0)
		for _, arg := range ctx.Args() {
			cmd = append(cmd, arg)
		}

		log.Infof("cmd : %v", cmd)
		tty := ctx.Bool("it")
		config := &subsystem.ResourceConfig{}
		config.MemoryLimit = ctx.String("mem")
		config.CpuCfsQuota = ctx.Int("cpu")
		logrus.Infof("config.MemoryLimi: %v", config.MemoryLimit)
		Run(tty, cmd, config)

		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",

	Action: func(ctx *cli.Context) error {
		log.Infof("init come")
		cmd := ctx.Args().Get(0)
		log.Infof("command %s", cmd)
		err := container.RunContainerInitProcess(cmd, nil)
		return err
	},
}
