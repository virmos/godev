package elastics

func (es *TestElasticRepository) CreateIndex(index string) error {
	return nil
}

// testing
func (es *TestElasticRepository) GetAllReports(index string) (map[string]interface{}, error) {
	var r map[string]interface{}
	return r, nil
}

func (es *TestElasticRepository) InsertHostStatusReport(index, hostName, statusCode, date string) error {
	return nil
}

// startDate and Endate are in UTC format
func (es *TestElasticRepository) GetYesterdayUptimeReport(index string) (map[string]Report, error) {
	var r map[string]Report
	return r, nil
}

// startDate and Endate are in UTC format
func (es *TestElasticRepository) GetYesterdayReport(index string) (map[string]Report, error) {
	var r map[string]Report
	return r, nil
}

// startDate and Endate are in UTC format
func (es *TestElasticRepository) GetRangeUptimeReport(index, hostName, startDate, endDate string) (map[string]Report, error) {
	var r map[string]Report
	return r, nil
}

// startDate and Endate are in UTC format
func (es *TestElasticRepository) GetRangeReport(index, hostName, startDate, endDate string) (map[string]Report, error) {
	var r map[string]Report
	return r, nil
}
