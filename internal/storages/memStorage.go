package storages

type MetricStorage struct {
	Metrics map[string]storageInterface
}

func NewMetricStorage() *MetricStorage {
	storage := &MetricStorage{}
	storage.Metrics["gauge"] = NewMetricGauge()
	storage.Metrics["counter"] = NewMetricCounter()
	return storage
}

type storageInterface interface {
	Process(data ...interface{})
}

type MetricGauge struct {
	value float64
}

func NewMetricGauge() *MetricGauge {
	return &MetricGauge{}
}

func (fip MetricGauge) Process(data ...interface{}) {
	fip.value = data[0].(float64)
}

type MetricCounter struct {
	value int
}

func NewMetricCounter() *MetricCounter {
	return &MetricCounter{}
}

func (fip MetricCounter) Process(data ...interface{}) {
	fip.value += data[0].(int)
}
