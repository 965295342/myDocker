package main

import (
	"fmt"
	"main/container"
	"main/subsystem"

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
		cli.StringFlag{
			Name:  "v",
			Usage: "volume ,: -v /local/file:/docker/file",
		},
		cli.BoolFlag{
			Name:  "d",
			Usage: "detach container",
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

		//tty or detach
		tty := ctx.Bool("it")
		detach := ctx.Bool("d")
		if tty && detach {
			return fmt.Errorf("-it and -d parameter can not both provided")
		}

		config := &subsystem.ResourceConfig{}
		config.MemoryLimit = ctx.String("mem")
		config.CpuCfsQuota = ctx.Int("cpu")

		volume := ctx.String("v")

		Run(tty, cmd, config, volume)

		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",

	Action: func(ctx *cli.Context) error {
		cmd := ctx.Args().Get(0)
		err := container.RunContainerInitProcess(cmd, nil)
		return err
	},
}

var commitCommand = cli.Command{
	Name:  "commit",
	Usage: "commit container to image e.g.:commit redis",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "v",
			Usage: "image save path",
		},
	},
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("commit command missing image name")
		}
		imageName := ctx.Args().Get(0)
		imageSavePath := ctx.String("v")
		container.CommitContainer(imageName, imageSavePath)
		return nil
	},
}
