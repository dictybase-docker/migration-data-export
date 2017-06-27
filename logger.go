package main

import (
	"os"

	"github.com/johntdyer/slackrus"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
)

func getLogger(c *cli.Context) *logrus.Logger {
	log := logrus.New()
	switch c.GlobalString("log-format") {
	case "text":
		log.Formatter = &logrus.TextFormatter{
			TimestampFormat: "02/Jan/2006:15:04:05",
		}
	case "json":
		log.Formatter = &logrus.JSONFormatter{
			TimestampFormat: "02/Jan/2006:15:04:05",
		}
	}
	log.Out = os.Stderr
	l := c.GlobalString("log-level")
	switch l {
	case "debug":
		log.Level = logrus.DebugLevel
	case "info":
		log.Level = logrus.InfoLevel
	case "warn":
		log.Level = logrus.WarnLevel
	case "error":
		log.Level = logrus.ErrorLevel
	case "fatal":
		log.Level = logrus.FatalLevel
	case "panic":
		log.Level = logrus.PanicLevel
	}
	// Set up hook
	lh := make(logrus.LevelHooks)
	for _, h := range c.GlobalStringSlice("hooks") {
		switch h {
		case "slack":
			lh.Add(&slackrus.SlackrusHook{
				HookURL:        c.GlobalString("slack-url"),
				AcceptedLevels: slackrus.LevelThreshold(log.Level),
				IconEmoji:      ":skull:",
			})
		default:
			continue
		}
	}
	log.Hooks = lh
	return log
}
