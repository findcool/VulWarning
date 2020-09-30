package common

import (
	"github.com/sirupsen/logrus"
)

const (
	// ServiceName -
	ServiceName = "VulWarning"
	// Description -
	Description = "VulWarning"
	// Version -
	Version = "0.1.6"
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
