package compute

import (
	"strconv"
	"time"

	"github.com/DennisJensen95/household-power/communicators/eloverblik"
	"github.com/DennisJensen95/household-power/communicators/energydatahub"
	log "github.com/sirupsen/logrus"
)

func CalculateThePeriodUseOfKWH(data eloverblik.TimeSeriesData) float64 {
	var total_kwh float64

	for _, time_series := range data.Result {
		for _, period := range time_series.MyEnergyDataMarketDocument.TimeSeries[0].Period {
			for _, point := range period.Point {
				kwh, err := strconv.ParseFloat(point.OutQuantityQuantity, 64)
				if err != nil {
					log.Println("Error while converting string to float:", err)
				}
				total_kwh += kwh
			}
		}
	}

	return total_kwh
}

func parseDateFromEnergyDataHub(date string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05Z", date)
}

func parseDateFromElOverblik(date string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05", date)
}

func prependPriceRecord(price energydatahub.Record, prices energydatahub.PriceDataRecords) energydatahub.PriceDataRecords {
	prices.Records = append(prices.Records, price)
	copy(prices.Records[1:], prices.Records)
	prices.Records[0] = price
	return prices
}

func getPricesBetweenDates(start_date string, end_date string, price_data energydatahub.PriceDataRecords) (energydatahub.PriceDataRecords, error) {
	var prices energydatahub.PriceDataRecords
	log.WithFields(log.Fields{
		"start_date": start_date,
		"end_date":   end_date,
	}).Debug("Getting price records between dates")

	start_time, err := parseDateFromEnergyDataHub(start_date)
	if err != nil {
		log.WithFields(log.Fields{
			"start_date": start_date,
			"error":      err,
		}).Error("Error while parsing date")
		return prices, err
	}

	end_time, err := parseDateFromEnergyDataHub(end_date)
	if err != nil {
		log.WithFields(log.Fields{
			"end_time": end_time,
			"error":    err,
		}).Error("Error while parsing date")
		return prices, err
	}

	for _, price := range price_data.Records {
		price_time, err := parseDateFromElOverblik(price.HourDK)
		if err != nil {
			log.WithFields(log.Fields{
				"price_time": price_time,
				"error":      err,
			}).Error("Error while parsing date")
			return prices, err
		}

		if price_time.After(start_time) && price_time.Before(end_time) || price_time.Equal(start_time) {
			prices.Records = prependPriceRecord(price, prices).Records
		}
	}

	return prices, nil
}

func CalculateVariablePrice(data eloverblik.TimeSeriesData, price_data energydatahub.PriceDataRecords, expense float64) (float64, error) {
	var total_price float64

	for _, time_series := range data.Result {
		for _, period := range time_series.MyEnergyDataMarketDocument.TimeSeries[0].Period {
			start_date := period.TimeInterval.Start
			end_date := period.TimeInterval.End

			prices, err := getPricesBetweenDates(start_date, end_date, price_data)
			if err != nil {
				log.WithFields(log.Fields{
					"start_date":        start_date,
					"end_date":          end_date,
					"error":             err,
					"Number of records": len(prices.Records),
				}).Error("Error while getting prices between dates")
				return 0, err
			}

			// The same number of price points is necessary as there are kwh points
			if len(prices.Records) != len(period.Point) {
				log.WithFields(log.Fields{
					"start_date":                   start_date,
					"end_date":                     end_date,
					"error":                        err,
					"Price data points":            len(prices.Records),
					"Number of power usage points": len(period.Point),
					"Total number of price points": len(price_data.Records),
				}).Warning("Number of records from price data does not match number of records from eloverblik")
				continue
			}

			for i, point := range period.Point {
				kwh, err := strconv.ParseFloat(point.OutQuantityQuantity, 64)
				if err != nil {
					log.Println("Error while converting string to float:", err)
				}
				price := prices.Records[i].SpotPriceDKK
				kwh_price := price/1000 + expense
				log.WithFields(log.Fields{
					"price":             kwh_price,
					"kwh":               kwh,
					"Record price":      prices.Records[i].HourDK,
					"Start date period": period.TimeInterval.Start,
					"Increment period":  point.Position,
				}).Debug("Calculating price")
				total_price += kwh * kwh_price
			}
		}
	}

	return total_price, nil
}
