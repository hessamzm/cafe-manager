package ui

import (
	"cafe-manager/utils"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

func NewSoodZiyanPage(win fyne.Window) fyne.CanvasObject {
	btn := widget.NewButton("بستن حساب ماهانه", func() {
		ShowSoodZiyanDialog(win)
	})

	return container.NewVBox(
		widget.NewLabel("محاسبه سود و زیان ماهانه"),
		btn,
	)
}
func ShowSoodZiyanDialog(win fyne.Window) {
	var fixedCosts []struct {
		Title  string
		Amount float64
	}

	// ورودی برای هزینه ثابت
	costTitle := widget.NewEntry()
	costTitle.SetPlaceHolder("عنوان هزینه (مثلاً اجاره)")

	costAmount := widget.NewEntry()
	costAmount.SetPlaceHolder("مبلغ هزینه (تومان)")

	// لیبل جمع هزینه‌ها
	totalCostsLabel := widget.NewLabel("جمع هزینه‌های ثابت: 0 تومان")

	// جدول هزینه‌ها
	costList := widget.NewMultiLineEntry()
	costList.Disable()

	addCost := func() {
		title := costTitle.Text
		amount, err := strconv.ParseFloat(costAmount.Text, 64)
		if title == "" || err != nil {
			dialog.ShowError(fmt.Errorf("عنوان یا مبلغ نامعتبر است"), win)
			return
		}
		fixedCosts = append(fixedCosts, struct {
			Title  string
			Amount float64
		}{title, amount})

		costList.Text += fmt.Sprintf("- %s: %.0f تومان\n", title, amount)
		costList.Refresh()
		costTitle.SetText("")
		costAmount.SetText("")

		// به‌روزرسانی جمع هزینه‌ها
		var total float64
		for _, c := range fixedCosts {
			total += c.Amount
		}
		totalCostsLabel.SetText(fmt.Sprintf("جمع هزینه‌های ثابت: %.0f تومان", total))
	}

	btnCalc := widget.NewButton("محاسبه سود و زیان و ساخت فایل اکسل", func() {
		// محاسبه مجموع هزینه‌ها
		var totalCosts float64
		for _, c := range fixedCosts {
			totalCosts += c.Amount
		}

		// ساخت فایل اکسل
		path, err := utils.GenerateSoodZiyanExcel(totalCosts)
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		dialog.ShowInformation("موفق", fmt.Sprintf("فایل سود و زیان ساخته شد:\n%s", path), win)
	})

	content := container.NewVBox(
		widget.NewLabel("هزینه‌های ثابت ماهانه"),
		costTitle,
		costAmount,
		widget.NewButton("افزودن هزینه", addCost),
		costList,
		totalCostsLabel,
		btnCalc,
	)

	d := dialog.NewCustom("بستن حساب ماهانه", "بستن", content, win)
	d.Show()
}
