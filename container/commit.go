package container

import (
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func CommitContainer(imageName, imageSavePath string) {
	mntPath := "/root/merged"
	imageFileName := imageSavePath + "/" + imageName + ".tar"
	fmt.Println("commitContainer imageFileName:", imageFileName)
	_, err := exec.Command("tar", "-czf", imageFileName, "-C", mntPath, ".").CombinedOutput()
	if err != nil {
		log.Errorf("commit container error %s", err)
	}
}
