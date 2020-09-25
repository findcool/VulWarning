package common

import (
	"github.com/sirupsen/logrus"
)

const (
	// ServiceName -
	ServiceName = "Tiku"
	// Description -
	Description = "Tiku @Venom"
	// Version -
	Version = "0.1.1"
	// ConfigFile -
	ConfigFile = "config.yaml"
	// LogFile -
	LogFile = "vulwarning.log"
)

var (
	// Logger 日志工具
	Logger *logrus.Logger
	// Conf -
	Conf Config
	// DebugMode -
	DebugMode = false
)
