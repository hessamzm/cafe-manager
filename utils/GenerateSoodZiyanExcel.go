package utils

import (
	"cafe-manager/db"
	"fmt"
	"github.com/xuri/excelize/v2"
	"time"
)

func GenerateSoodZiyanExcel(fixedCosts float64) (string, error) {
	// گرفتن تاریخ ابتدای ماه جاری
	now := time.Now()
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	firstDayStr := firstDay.Format("2006-01-02 00:00:00")

	// گرفتن سفارش‌های تسویه‌شده از ابتدای ماه
	rows, err := db.DB.Query(`
		SELECT oi.product_id, p.name, p.cost_price, oi.quantity, oi.unit_price
		FROM order_items oi
		JOIN orders o ON oi.order_id = o.id
		JOIN products p ON oi.product_id = p.id
		WHERE o.settled = 1 AND o.created_at >= ?
	`, firstDayStr)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	type rowData struct {
		ProductName string
		Quantity    int
		SaleTotal   float64
		CostTotal   float64
	}

	summary := map[int]rowData{}
	var totalSales float64
	var totalCost float64

	for rows.Next() {
		var productID int
		var name string
		var costPrice float64
		var qty int
		var unitPrice float64

		rows.Scan(&productID, &name, &costPrice, &qty, &unitPrice)

		s := summary[productID]
		s.ProductName = name
		s.Quantity += qty
		s.SaleTotal += float64(qty) * unitPrice
		s.CostTotal += float64(qty) * costPrice
		summary[productID] = s

		totalSales += float64(qty) * unitPrice
		totalCost += float64(qty) * costPrice
	}

	profit := totalSales - totalCost - fixedCosts

	// ساخت اکسل
	f := excelize.NewFile()
	sheet := f.GetSheetName(0)

	// عنوان
	f.SetCellValue(sheet, "A1", "گزارش سود و زیان ماهانه")
	f.MergeCell(sheet, "A1", "E1")

	// هدر جدول
	headers := []string{"نام محصول", "تعداد فروش", "فروش کل", "قیمت تمام‌شده کل", "سود ناخالص"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 3)
		f.SetCellValue(sheet, cell, h)
	}

	// محتوای جدول
	row := 4
	for _, s := range summary {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), s.ProductName)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), s.Quantity)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), s.SaleTotal)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), s.CostTotal)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), s.SaleTotal-s.CostTotal)
		row++
	}

	// جمع کل و سود نهایی
	row += 1
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "فروش کل:")
	f.SetCellValue(sheet, fmt.Sprintf("C%d", row), totalSales)
	row++
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "قیمت تمام‌شده کل:")
	f.SetCellValue(sheet, fmt.Sprintf("C%d", row), totalCost)
	row++
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "هزینه‌های ثابت:")
	f.SetCellValue(sheet, fmt.Sprintf("C%d", row), fixedCosts)
	row++
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "سود / زیان نهایی:")
	f.SetCellValue(sheet, fmt.Sprintf("C%d", row), profit)

	// ذخیره فایل
	path := fmt.Sprintf("sood_ziyan_%s.xlsx", now.Format("2006_01"))
	err = f.SaveAs(path)
	if err != nil {
		return "", err
	}

	return path, nil
}
