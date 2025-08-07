package ui

import (
	"cafe-manager/db"
	"cafe-manager/utils"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// NewReportsPage نمایش صفحه گزارش‌گیری
func NewReportsPage(win fyne.Window) fyne.CanvasObject {
	resultLabel := widget.NewLabel("")
	resultLabel.Wrapping = fyne.TextWrapWord

	startDate := widget.NewEntry()
	endDate := widget.NewEntry()
	startDate.SetText(utils.GetDateStr(time.Now()))
	endDate.SetText(utils.GetDateStr(time.Now()))

	runReport := func(start, end time.Time, label string) {
		orders, total := db.GetOrdersBetween(
			start.Format("2006-01-02"),
			end.Format("2006-01-02"),
		)

		// ساخت جزییات متن نمایش داده‌شده
		var details string
		counts := make(map[string]int)
		for _, o := range orders {
			for _, i := range o.Items {
				line := fmt.Sprintf("- %s × %d = %.0f تومان\n", i.ProductName, i.Quantity, float64(i.Quantity)*i.UnitPrice)
				details += line
				counts[i.ProductName] += i.Quantity
			}
		}
		if len(counts) > 0 {
			details += "\n=== خلاصه تعداد ===\n"
			for name, qty := range counts {
				details += fmt.Sprintf("* %s: %d عدد\n", name, qty)
			}
		}
		details += fmt.Sprintf("\nمبلغ کل (بدون تخفیف): %.0f تومان", total)
		resultLabel.SetText(details)

		// ساخت مسیر ذخیره اکسل روی دسکتاپ
		homeDir, _ := os.UserHomeDir()
		outputDir := filepath.Join(homeDir, "Desktop", "CafeReports")
		os.MkdirAll(outputDir, os.ModePerm)
		dateForFilename := utils.GetDateStr(time.Now())
		safeDateStr := strings.ReplaceAll(dateForFilename, "/", "-")
		safeDateStr = strings.ReplaceAll(safeDateStr, "\\", "-") // برای اطمینان از حذف بک‌اسلش

		// 2. از تاریخ امن‌شده برای ساخت نام فایل استفاده کنید
		filename := fmt.Sprintf("%s-%s-گزارش.xlsx", safeDateStr, label)
		// --- پایان تغییرات ---

		excelPath := filepath.Join(outputDir, filename)
		//filename := fmt.Sprintf("%s-%s-گزارش.xlsx", utils.GetDateStr(time.Now()), label)
		//excelPath := filepath.Join(outputDir, filename)

		if err := utils.SaveReportToExcel(excelPath, orders, total, start, end); err != nil {
			dialog.ShowError(err, win)
		} else {
			dialog.ShowInformation("خروجی اکسل", "فایل گزارش در مسیر زیر ذخیره شد:\n"+excelPath, win)
		}
	}

	toolbar := container.NewVBox(
		widget.NewLabel("تاریخ شروع (yyyy/MM/dd):"), startDate,
		widget.NewLabel("تاریخ پایان (yyyy/MM/dd):"), endDate,
		widget.NewButton("گزارش بازه دلخواه", func() {
			s, err1 := utils.ParsePersianDate(startDate.Text)
			e, err2 := utils.ParsePersianDate(endDate.Text)
			if err1 != nil || err2 != nil {
				resultLabel.SetText("فرمت تاریخ نادرست است")
				return
			}
			runReport(s, e.Add(24*time.Hour), "دلخواه")
		}),
		widget.NewButton("گزارش امروز", func() {
			now := time.Now()
			start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			runReport(start, start.Add(24*time.Hour), "روزانه")
		}),
		widget.NewButton("گزارش ۷ روز اخیر", func() {
			end := time.Now()
			start := end.AddDate(0, 0, -7)
			runReport(start, end, "هفتگی")
		}),
		widget.NewButton("گزارش این ماه", func() {
			now := time.Now()
			start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			runReport(start, now, "ماهانه")
		}),
	)

	return container.NewBorder(toolbar, nil, nil, nil, container.NewVScroll(resultLabel))
}
