package utils

import (
	"helper/log"
	"os/exec"
)

func RunCmd(cmd string, wd string, shell bool) []byte {
	log.Info(">>> CMD:", cmd)
	var c *exec.Cmd
	if shell {
		c = exec.Command("bash", "-c", cmd)
	} else {
		c = exec.Command(cmd)
	}

	if len(wd) > 0 {
		c.Dir = wd
	}
	out, err := c.Output()
	if err != nil {
		log.Error(err)
		panic("some error found")
	}
	log.Info("Output:" + string(out) + " <<<")
	return out
}
