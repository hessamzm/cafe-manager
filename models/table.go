package models

type TableStatus string

const (
	TableAvailable TableStatus = "آماده استفاده"
	TableBusy      TableStatus = "مشغول"
	TableReserved  TableStatus = "رزرو"
)

type Table struct {
	ID     int
	Name   string
	Status TableStatus
}
