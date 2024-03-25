package subsystem

import (
	"fmt"
	"main/constant"
	"os"
	"path"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	PeriodDefault = 100000
	Percent       = 100
)

type CpuSubSystem struct {
}

func (s *CpuSubSystem) Name() string {
	return "cpu"
}

// 具体来说,
// cpu.shares 参数用于设置 CPU 分配的相对权重,
// 而不是绝对值。它是一个相对于默认值1024的相对权重。例如,
// 如果一个 cgroup 中的进程的 cpu.shares 值为 2048,
// 而另一个 cgroup 中的进程的值为 1024,
// 则前者将获得相对于后者两倍的 CPU 时间片。
func (s *CpuSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	if res.CpuCfsQuota == 0 && res.CpuShare == "" {
		return nil
	}
	logrus.Infof("CpuSubSystem:%v ", res.CpuCfsQuota, res.CpuShare)
	subsysCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	logrus.Infof("subsysCgroupPath: %v", subsysCgroupPath)

	if res.CpuShare != "" {
		err := os.WriteFile(path.Join(subsysCgroupPath, "cpu.shares"), []byte(res.CpuShare), constant.Perm0644)
		if err != nil {
			return fmt.Errorf("set cgroup cpuShare fail %v", err)
		}
	}
	// 举个例子，假设设置了 cpu.cfs_period_us 为 100000us
	// （100毫秒），cpu.cfs_quota_us 为 50000us
	// （50毫秒）。那么在每个 100 毫秒的周期内，
	// 进程最多可以使用 50 毫秒的 CPU 时间。
	if res.CpuCfsQuota != 0 {
		err := os.WriteFile(path.Join(subsysCgroupPath, "cpu.cfs_period_us"), []byte("100000"), constant.Perm0644)
		if err != nil {
			return fmt.Errorf("set cgroup cfs_period_us fail %v", err)
		}
		err = os.WriteFile(path.Join(subsysCgroupPath, "cpu.cfs_quota_us"), []byte(strconv.Itoa(PeriodDefault/Percent*res.CpuCfsQuota)), constant.Perm0644)
		if err != nil {
			return fmt.Errorf("set cgroup cfs_quota_us fail %v", err)
		}
	}
	return nil
}

func (s *CpuSubSystem) Apply(cgroupPath string, pid int, res *ResourceConfig) error {
	if res.CpuShare == "" && res.CpuCfsQuota == 0 {
		return nil
	}
	subsysCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, false)
	if err != nil {
		return errors.Wrapf(err, "get cgroup %s", cgroupPath)
	}
	if err := os.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), constant.Perm0644); err != nil {
		return fmt.Errorf("set cgroup proc fail %v", err)
	}
	return nil
}

func (s *CpuSubSystem) Remove(cgroupPath string) error {
	subsysCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, false)
	if err != nil {
		return err
	}
	return os.RemoveAll(subsysCgroupPath)
}
