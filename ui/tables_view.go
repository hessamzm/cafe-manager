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

	hasOrders := db.HasOrdersForTable(t.ID) // تابعی که بررسی کند آیا میز سفارش دارد یا نه

	if t.Status == "مشغول" {
		options.Add(widget.NewButton("افزودن سفارش", func() {
			ShowOrderForm(win, &t, refresh).Show()
		}))

		options.Add(widget.NewButton("اصلاح سفارشات", func() {
			ShowEditOrdersDialog(t.ID, win)
		}))

		if hasOrders {
			options.Add(widget.NewButton("تسویه حساب", func() {
				ShowCloseTableDialog(win, t.ID, refresh)
			}))
		}
	} else {
		options.Add(widget.NewButton("ثبت سفارش", func() {
			ShowOrderForm(win, &t, refresh).Show()
		}))
	}

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
func ShowEditOrdersDialog(tableID int, win fyne.Window) {
	orders, _ := db.GetOrdersByTableID(tableID)
	if len(orders) == 0 {
		dialog.ShowInformation("بدون سفارش", "هیچ سفارشی برای این میز ثبت نشده است", win)
		return
	}

	orderItems := db.GetOrderItemsForTable(tableID)
	if len(orderItems) == 0 {
		dialog.ShowInformation("بدون آیتم", "سفارش‌های این میز آیتمی ندارند", win)
		return
	}

	type ItemRow struct {
		ID            int
		OrderID       int
		QuantityEntry *widget.Entry
		RowContainer  *fyne.Container
	}

	var itemRows []ItemRow
	listContainer := container.NewVBox()

	// برای حذف یا آپدیت بهتر است از لیستی از آیتم‌ها استفاده کنیم
	for _, item := range orderItems {
		quantityEntry := widget.NewEntry()
		quantityEntry.SetText(fmt.Sprintf("%d", item.Quantity))

		rowLabel := widget.NewLabel(fmt.Sprintf("%s - قیمت واحد: %.0f تومان", item.ProductName, item.UnitPrice))

		deleteBtn := widget.NewButton("حذف", func(id int, rowCont *fyne.Container) func() {
			return func() {
				db.DeleteOrderItemByID(id)
				rowCont.Hide() // پنهان کردن ردیف از UI
			}
		}(item.ID, listContainer))

		rowContainer := container.NewHBox(
			rowLabel,
			widget.NewLabel("تعداد:"),
			quantityEntry,
			deleteBtn,
		)

		itemRows = append(itemRows, ItemRow{
			ID:            item.ID,
			OrderID:       item.OrderID,
			QuantityEntry: quantityEntry,
			RowContainer:  rowContainer,
		})

		listContainer.Add(rowContainer)
	}

	saveBtn := widget.NewButton("ذخیره تغییرات", func() {
		affectedOrders := make(map[int]bool)

		for _, row := range itemRows {
			text := row.QuantityEntry.Text
			if text == "" {
				continue
			}
			newQty, err := strconv.Atoi(text)
			if err == nil {
				if newQty <= 0 {
					db.DeleteOrderItemByID(row.ID)
					row.RowContainer.Hide()
				} else {
					db.UpdateOrderItemQuantity(row.ID, newQty)
				}
				affectedOrders[row.OrderID] = true
			}
		}

		for orderID := range affectedOrders {
			db.RecalculateOrderTotalAmount(orderID)
		}

		dialog.ShowInformation("ثبت شد", "تغییرات ذخیره شدند", win)
	})

	content := container.NewVBox(
		widget.NewLabel("ویرایش آیتم‌های سفارش میز"),
		widget.NewSeparator(),
		listContainer,
		widget.NewSeparator(),
		saveBtn,
	)

	editDialog := dialog.NewCustom("ویرایش سفارشات", "بستن", content, win)
	editDialog.Show()
}
