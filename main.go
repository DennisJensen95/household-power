package main

import (
	"flag"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/DennisJensen95/household-power/communicators/eloverblik"
	"github.com/DennisJensen95/household-power/internal/power_usage"
	log "github.com/sirupsen/logrus"
)

var (
	start_date      = flag.String("start-date", "", "Start date for the time series")
	end_date        = flag.String("end-date", "", "End date for the time series")
	token           = flag.String("token", "", "access-token")
	meteringPointId = flag.String("meteringPointId", "", "meter-point-id")
	log_level       = flag.String("log-level", "info", "Set log level (debug, info, warn, error, fatal, panic)")
	fixed_price     = flag.Float64("fixed-price", 0, "Fixed price per Kr/kWh")
	transport_cost  = flag.Float64("transport-cost", -1, "Transport cost per Kr/kWh")
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
		log.Fatal("Please provide a start-date in the format YYYY-MM-DD")
	}

	if _, err := time.Parse("2006-01-02", *start_date); err != nil {
		flag.Usage()
		log.Fatal("Please provide a start-date in the format YYYY-MM-DD")
	}

	if *end_date == "" {
		flag.Usage()
		log.Fatal("Please provide a end-date")
	}

	if _, err := time.Parse("2006-01-02", *end_date); err != nil {
		flag.Usage()
		log.Fatal("Please provide a end-date in the format YYYY-MM-DD")
	}

	if *fixed_price == 0 {
		flag.Usage()
		log.Fatal("Please provide a fixed-price")
	}

	if *transport_cost == -1 {
		flag.Usage()
		log.Fatal("Please provide a transport-cost")
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
	tokenElOverblik := eloverblik.ElOverblikCommunicator{
		Client: &http.Client{},
	}
	eloverblik_access_token, err := tokenElOverblik.GetAccessToken(*token)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error while getting access token")
	}

	power, err := power_usage.PeriodUseKWH(eloverblik_access_token, *meteringPointId, *start_date, *end_date)
	if err != nil {
		log.Fatal(err)
	}
	log.WithField("PowerKWH", math.Round(power)).Info("Power used")

	price, err := power_usage.PeriodUsePriceVariable(eloverblik_access_token, *meteringPointId, *start_date, *end_date, *transport_cost)

	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(log.Fields{
		"Money spend with variable agreement [DKK]": math.Round(price*10) / 10,
		"Money spend with fixed agreement [DKK]":    math.Round(power**fixed_price*10) / 10,
	}).Info("Power used")
}
