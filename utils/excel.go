package utils

import (
	"cafe-manager/models"
	"fmt"
	"github.com/xuri/excelize/v2"
	"time"
)

// SaveReportToExcel فایل اکسل گزارش را در مسیر path می‌سازد
func SaveReportToExcel(path string, orders []models.Order, total float64, start, end time.Time) error {
	f := excelize.NewFile()

	// -------------------------------
	// Sheet 1: گزارش فروش (Sales)
	// -------------------------------
	salesSheet := "Sales"
	f.NewSheet(salesSheet)
	headers := []string{"تاریخ", "نام محصول", "تعداد", "قیمت واحد", "جمع"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(salesSheet, cell, h)
	}

	row := 2
	var totalSales float64
	for _, o := range orders {
		t, _ := time.Parse("2006-01-02 15:04:05", o.CreatedAt)
		shamsi := GetDateStr(t)
		for _, item := range o.Items {
			f.SetCellValue(salesSheet, fmt.Sprintf("A%d", row), shamsi)
			f.SetCellValue(salesSheet, fmt.Sprintf("B%d", row), item.ProductName)
			f.SetCellValue(salesSheet, fmt.Sprintf("C%d", row), item.Quantity)
			f.SetCellValue(salesSheet, fmt.Sprintf("D%d", row), item.UnitPrice)
			sum := float64(item.Quantity) * item.UnitPrice
			f.SetCellValue(salesSheet, fmt.Sprintf("E%d", row), sum)
			totalSales += sum
			row++
		}
	}
	f.SetCellValue(salesSheet, fmt.Sprintf("D%d", row), "جمع کل:")
	f.SetCellValue(salesSheet, fmt.Sprintf("E%d", row), totalSales)

	// -------------------------------
	// Sheet 2: تعداد فروش هر آیتم (Item Summary)
	// -------------------------------
	summarySheet := "ItemSummary"
	f.NewSheet(summarySheet)
	f.SetCellValue(summarySheet, "A1", "نام محصول")
	f.SetCellValue(summarySheet, "B1", "تعداد کل فروش")

	summaryMap := make(map[string]int)
	for _, o := range orders {
		for _, item := range o.Items {
			summaryMap[item.ProductName] += item.Quantity
		}
	}

	row = 2
	for name, count := range summaryMap {
		f.SetCellValue(summarySheet, fmt.Sprintf("A%d", row), name)
		f.SetCellValue(summarySheet, fmt.Sprintf("B%d", row), count)
		row++
	}

	// -------------------------------
	// Sheet 3: تفکیک نوع سفارش (Order Type Breakdown)
	// -------------------------------
	typeSheet := "OrderTypes"
	f.NewSheet(typeSheet)
	f.SetCellValue(typeSheet, "A1", "نوع سفارش")
	f.SetCellValue(typeSheet, "B1", "تعداد")
	f.SetCellValue(typeSheet, "C1", "جمع مبلغ")

	var cafeCount, takeCount int
	var cafeTotal, takeTotal float64

	for _, o := range orders {
		var sum float64
		var count int
		for _, item := range o.Items {
			sum += float64(item.Quantity) * item.UnitPrice
			count += item.Quantity
		}
		if o.TableID != nil {
			cafeCount += count
			cafeTotal += sum
		} else {
			takeCount += count
			takeTotal += sum
		}
	}

	f.SetCellValue(typeSheet, "A2", "فروش درون‌کافه")
	f.SetCellValue(typeSheet, "B2", cafeCount)
	f.SetCellValue(typeSheet, "C2", cafeTotal)

	f.SetCellValue(typeSheet, "A3", "فروش بیرون‌بر")
	f.SetCellValue(typeSheet, "B3", takeCount)
	f.SetCellValue(typeSheet, "C3", takeTotal)

	// -------------------------------
	// ذخیره فایل
	// -------------------------------
	f.DeleteSheet("Sheet1") // حذف شیت پیش‌فرض
	if err := f.SaveAs(path); err != nil {
		return err
	}
	return nil
}
