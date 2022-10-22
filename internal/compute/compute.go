package compute

import (
	"log"
	"strconv"

	"github.com/DennisJensen95/household-power/communicators/eloverblik"
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
