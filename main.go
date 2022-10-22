package main

import (
	"flag"
	"math"
	"strings"

	"github.com/DennisJensen95/household-power/internal/power_usage"
	log "github.com/sirupsen/logrus"
)

var (
	start_date      = flag.String("start-date", "", "Start date for the time series")
	end_date        = flag.String("end-date", "", "End date for the time series")
	token           = flag.String("token", "", "access-token")
	meteringPointId = flag.String("meteringPointId", "", "meter-point-id")
	log_level       = flag.String("log-level", "info", "Set log level (debug, info, warn, error, fatal, panic)")
)

func init_parser() {
	flag.Parse()

	if *token == "" {
		flag.Usage()
		log.Fatal("Please provide an access-token")
	}

	if *meteringPointId == "" {
		flag.Usage()
		log.Fatal("Please provide a metering-point-id")
	}

	if *start_date == "" {
		flag.Usage()
		log.Fatal("Please provide a start-date")
	}

	if *end_date == "" {
		flag.Usage()
		log.Fatal("Please provide a end-date")
	}

}

var logLevel = map[string]log.Level{
	"info":  log.InfoLevel,
	"debug": log.DebugLevel,
	"warn":  log.WarnLevel,
	"error": log.ErrorLevel,
	"fatal": log.FatalLevel,
	"panic": log.PanicLevel,
	"i":     log.InfoLevel,
	"d":     log.DebugLevel,
	"w":     log.WarnLevel,
	"e":     log.ErrorLevel,
	"f":     log.FatalLevel,
	"p":     log.PanicLevel,
}

func main() {
	init_parser()

	*log_level = strings.ToLower(*log_level)
	log.SetLevel(logLevel[*log_level])

	// log.SetFormatter(&log.JSONFormatter{})

	power, err := power_usage.PeriodUseKWH(*token, *meteringPointId, *start_date, *end_date)
	if err != nil {
		log.Fatal(err)
	}
	log.WithField("PowerKWH", math.Round(power)).Info("Power used")
}
