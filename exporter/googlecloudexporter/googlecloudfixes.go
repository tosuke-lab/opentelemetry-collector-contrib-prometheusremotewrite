package googlecloudexporter

import "go.opentelemetry.io/collector/model/pdata"

/* Fixes empty StartTimestamp field in input metrics datapoints by updating StartTimestamp = Timestamp - 1 micro sec
   Some metrics from Hostmetric receiver has Datapoints with empty StartTimestamp = 0
   This causes that later on StartTimestamp is assigned as Timestamp and such datapoints are rejected by google monitoring api */
func fixEmptyStartTime(metrics pdata.Metrics) {
	resourceMetrics := metrics.ResourceMetrics()
	for i := 0; i < resourceMetrics.Len(); i++ {
		rm := resourceMetrics.At(i)
		ilms := rm.InstrumentationLibraryMetrics()

		for j := 0; j < ilms.Len(); j++ {
			ilm := ilms.At(j)
			metrics := ilm.Metrics()

			for k := 0; k < metrics.Len(); k++ {
				metric := metrics.At(k)

				switch metric.DataType() {
				case pdata.MetricDataTypeSum:
					fixEmptyStartTimestampInIntData(metric.Sum().DataPoints())
					break
					//case pdata.MetricDataTypeSum:
					//	fixEmptyStartTimestampInDoubleData(metric.Sum().DataPoints())
					//	break
				}
			}
		}
	}
}

const microSecAsNano = 1000

func fixEmptyStartTimestampInIntData(dps pdata.NumberDataPointSlice) {
	for i := 0; i < dps.Len(); i++ {
		dp := dps.At(i)
		if dp.StartTimestamp() == 0 && dp.Timestamp() > microSecAsNano {
			dp.SetStartTimestamp(dp.Timestamp() - microSecAsNano)
		}
	}
}

//func fixEmptyStartTimestampInDoubleData(dps pdata.NumberDataPointSlice) {
//	for i := 0; i < dps.Len(); i++ {
//		dp := dps.At(i)
//		if dp.StartTimestamp() == 0 && dp.Timestamp() > microSecAsNano {
//			dp.SetStartTimestamp(dp.Timestamp() - microSecAsNano)
//		}
//	}
//}
