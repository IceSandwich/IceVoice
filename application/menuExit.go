package application

import "github.com/sirupsen/logrus"

func clickedExit() {
	logrus.Info("clicked exit")

	Destroy()
}

func init() {
	RegisterAction("Exit", clickedExit)
}
