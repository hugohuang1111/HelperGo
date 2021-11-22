package command

import (
	"errors"
	"helper/log"
	"helper/utils"
)

func PrefabProcess(src, dst string) error {
	if !utils.FileExist(src) {
		return errors.New("path is not exist:" + src)
	}

	log.Info(src)
	log.Info(dst)
	return nil
}
