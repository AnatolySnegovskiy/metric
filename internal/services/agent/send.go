package agent

func SendMetricsPeriodically(addr string, s Storage) error {
	for storageType, storage := range s.GetList() {
		for metricName, metric := range storage.GetList() {
			err := sendMetric(addr, storageType, metricName, metric)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
