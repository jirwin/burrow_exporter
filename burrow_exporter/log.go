package burrow_exporter

import (
	log "github.com/Sirupsen/logrus"
)

func SetLogLevel(mode int) {
	switch mode {
	case 1:
		log.SetLevel(log.PanicLevel)
	case 2:
		log.SetLevel(log.FatalLevel)
	case 3:
		log.SetLevel(log.ErrorLevel)
	case 4:
		log.SetLevel(log.WarnLevel)
	case 5:
		log.SetLevel(log.InfoLevel)
	case 6:
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}
