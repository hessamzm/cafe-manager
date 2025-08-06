package ui

import (
	"cafe-manager/db"
	"cafe-manager/models"
	"cafe-manager/utils"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

func NewTablesPage(win fyne.Window) fyne.CanvasObject {
	tablesBox := container.NewVBox()
	var RefreshTables func()
	RefreshTables = func() {
		tablesBox.Objects = nil
		all := db.GetAllTables()

		for _, t := range all {
			label := widget.NewLabel("میز " + t.Name + " (" + string(t.Status) + ")")
			btn := widget.NewButton("عملیات", func() {
				ShowTableOptions(win, t, RefreshTables)
			})
			tablesBox.Add(container.NewHBox(label, layout.NewSpacer(), btn))
		}
		tablesBox.Refresh()
	}

	toolbar := container.NewHBox(
		widget.NewButtonWithIcon("افزودن میز", utils.LoadIcon("add.png"), func() {
			name := "میز " + strconv.Itoa(len(db.GetAllTables())+1)
			db.AddTable(name)
			RefreshTables()
		}),
		widget.NewButton("سفارش بیرون‌بر", func() {
			ShowOrderForm(win, nil, RefreshTables).Show()
		}),
	)

	RefreshTables()

	return container.NewBorder(toolbar, nil, nil, nil, container.NewVScroll(tablesBox))
}
func ShowTableOptions(win fyne.Window, t models.Table, refresh func()) {
	options := container.NewVBox()

	options.Add(widget.NewButton("تغییر وضعیت", func() {
		db.ToggleTableStatus(t.ID)
		refresh()
		win.Canvas().Overlays().Top().Hide()
	}))

	if t.Status == "مشغول" {
		options.Add(widget.NewButton("افزودن سفارش", func() {
			ShowOrderForm(win, &t, refresh).Show()
		}))
		options.Add(widget.NewButton("تسویه حساب", func() {
			db.CloseTable(t.ID)
			refresh()
		}))
	} else {
		options.Add(widget.NewButton("ثبت سفارش", func() {
			ShowOrderForm(win, &t, refresh).Show()
		}))
	}
	options.Add(widget.NewButton("تسویه حساب", func() {
		ShowCloseTableDialog(win, t.ID, refresh)
	}))

	dialog.ShowCustom("عملیات میز "+t.Name, "بستن", options, win)
}
func ShowCloseTableDialog(win fyne.Window, tableID int, refresh func()) {
	orders, total := db.GetOrdersByTableID(tableID)
	discountEntry := widget.NewEntry()
	discountEntry.SetPlaceHolder("تخفیف (تومان)")

	// نمایش سفارشات
	orderList := widget.NewMultiLineEntry()
	for _, o := range orders {
		for _, i := range o.Items {
			line := fmt.Sprintf("- %s × %d = %.0f\n", i.ProductName, i.Quantity, float64(i.Quantity)*i.UnitPrice)
			orderList.Text += line
		}
	}
	orderList.Text += fmt.Sprintf("\nمبلغ کل: %.0f تومان", total)
	orderList.Disable()

	totalLabel := widget.NewLabel(fmt.Sprintf("مبلغ نهایی: %.0f تومان", total))

	// به‌روزرسانی مبلغ نهایی با اعمال تخفیف
	updateTotal := func() {
		discount, _ := strconv.ParseFloat(discountEntry.Text, 64)
		final := total - discount
		if final < 0 {
			final = 0
		}
		totalLabel.SetText(fmt.Sprintf("مبلغ نهایی: %.0f تومان", final))
	}

	discountEntry.OnChanged = func(s string) {
		updateTotal()
	}

	content := container.NewVBox(
		widget.NewLabel("سفارشات میز"),
		orderList,
		widget.NewButton("اصلاح سفارشات", func() {
			dialog.ShowInformation("درحال توسعه", "ویرایش سفارش در نسخه بعدی فعال می‌شود", win)
		}),
		widget.NewButton("افزودن سفارش", func() {
			ShowOrderForm(win, db.GetTableByID(tableID), refresh).Show()
		}),
		widget.NewLabel("تخفیف:"),
		discountEntry,
		totalLabel,
		widget.NewButton("تسویه سفارش", func() {
			db.CloseTableAndClearOrders(tableID)
			dialog.ShowInformation("تسویه انجام شد", "میز با موفقیت تسویه شد", win)
			refresh()
		}),
	)

	d := dialog.NewCustom("تسویه حساب میز", "بستن", content, win)
	d.Show()
}
