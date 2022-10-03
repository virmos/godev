package elastics

type Repository interface {
	InsertHostStatusReport(index, hostName, statusCode, date string) error
	GetAllReports(index string) (map[string]interface{}, error)
	GetYesterdayUptimeReport(index string) (map[string]Report, error)
	GetYesterdayReport(index string) (map[string]Report, error)
	GetRangeUptimeReport(index, startDate, endDate string) (map[string]Report, error)
	GetRangeReport(index, startDate, endDate string) (map[string]Report, error)
}