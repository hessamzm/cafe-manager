package models

type OrderItem struct {
	ProductID   int
	ProductName string
	Quantity    int
	UnitPrice   float64
}

type Order struct {
	ID          int
	TableID     *int
	Items       []OrderItem
	TotalAmount float64
	Note        string // به جای Description
	CreatedAt   string
}
