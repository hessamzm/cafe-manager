package ui

import (
	"cafe-manager/db"
	"cafe-manager/utils"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"path/filepath"
	"time"
)

// NewReportsPage نمایش صفحه گزارش‌گیری
func NewReportsPage(win fyne.Window) fyne.CanvasObject {
	// ویجت نمایش نتیجه؛ رنگش خودکار بر اساس تم تنظیم می‌شود
	resultLabel := widget.NewLabel("")
	resultLabel.Wrapping = fyne.TextWrapWord

	startDate := widget.NewEntry()
	endDate := widget.NewEntry()
	startDate.SetText(utils.GetDateStr(time.Now()))
	endDate.SetText(utils.GetDateStr(time.Now()))

	// اجرای گزارش با دریافت بازه میلادی
	runReport := func(start, end time.Time) {
		orders, total := db.GetOrdersBetween(
			start.Format("2006-01-02"),
			end.Format("2006-01-02"),
		)

		// جزییات سفارشات و جمع تعداد هر محصول
		var details string
		counts := make(map[string]int)
		for _, o := range orders {
			for _, i := range o.Items {
				details += fmt.Sprintf("- %s × %d = %.0f تومان\n",
					i.ProductName, i.Quantity,
					float64(i.Quantity)*i.UnitPrice,
				)
				counts[i.ProductName] += i.Quantity
			}
		}

		// بخش خلاصه تعداد
		if len(counts) > 0 {
			details += "\n=== خلاصه تعداد ===\n"
			for name, qty := range counts {
				details += fmt.Sprintf("* %s: %d عدد\n", name, qty)
			}
		}

		details += fmt.Sprintf("\nمبلغ کل (بدون تخفیف): %.0f تومان", total)
		resultLabel.SetText(details)
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
			runReport(s, e.Add(24*time.Hour))
		}),
		widget.NewButton("گزارش امروز", func() {
			now := time.Now()
			s := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			runReport(s, s.Add(24*time.Hour))
		}),
		widget.NewButton("گزارش ۷ روز اخیر", func() {
			end := time.Now()
			start := end.AddDate(0, 0, -7)
			runReport(start, end)
		}),
		widget.NewButton("گزارش این ماه", func() {
			now := time.Now()
			start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			runReport(start, now)
		}),
	)
	exportBtn := widget.NewButton("خروجی Excel", func() {
		// بازکردن دیالوگ انتخاب مسیر فایل
		dialog.ShowFileSave(func(uri fyne.URIWriteCloser, err error) {
			if err != nil || uri == nil {
				return
			}
			path := uri.URI().Path()
			// مطمئن شو پسوند xlsx دارد
			if filepath.Ext(path) == "" {
				path += ".xlsx"
			}
			// پارس تاریخ
			s, err1 := utils.ParsePersianDate(startDate.Text)
			e, err2 := utils.ParsePersianDate(endDate.Text)
			if err1 != nil || err2 != nil {
				dialog.ShowError(fmt.Errorf("فرمت تاریخ نادرست است"), win)
				return
			}
			// گرفتن داده‌ها
			orders, total := db.GetOrdersBetween(s.Format("2006-01-02"), e.Format("2006-01-02"))
			// اضافه کردن CreatedAt به مدل (فرض)
			for i := range orders {
				orders[i].CreatedAt = orders[i].CreatedAt // در GetOrdersBetween برگردانده‌اید
			}
			// ذخیره Excel
			if err := utils.SaveReportToExcel(path, orders, total, s, e); err != nil {
				dialog.ShowError(err, win)
			} else {
				dialog.ShowInformation("خروجی اکسل", "فایل با موفقیت در "+path+" ذخیره شد", win)
			}
		}, win)
	})
	toolbar.Add(exportBtn)
	return container.NewBorder(toolbar, nil, nil, nil, container.NewVScroll(resultLabel))
}
