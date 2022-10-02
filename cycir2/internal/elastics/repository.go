package elastics

type Repository interface {
	InsertHostStatusReport(name, statusCode, date string) error
	YesterdayUptimeReport(name string) ([]Report, error)
	YesterdayReport(name string) ([]Report, error)
	RangeUptimeReport(name, startDate, endDate string) ([]Report, error)
	RangeReport(name, startDate, endDate string) ([]Report, error)
}