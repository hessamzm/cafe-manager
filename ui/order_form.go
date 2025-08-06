package ui

import (
	"cafe-manager/db"
	"cafe-manager/models"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

func ShowOrderForm(win fyne.Window, table *models.Table, refresh func()) fyne.Window {
	w := fyne.CurrentApp().NewWindow("ثبت سفارش")

	products := db.GetAllProducts()

	selectedItems := []models.OrderItem{}

	productSelect := widget.NewSelect([]string{}, nil)
	productMap := make(map[string]models.Product)
	for _, p := range products {
		productSelect.Options = append(productSelect.Options, p.Name)
		productMap[p.Name] = p
	}
	productSelect.Refresh()

	quantity := widget.NewEntry()
	quantity.SetPlaceHolder("تعداد")

	list := widget.NewMultiLineEntry()
	list.SetPlaceHolder("محصولات انتخاب‌شده")

	addBtn := widget.NewButton("افزودن به سفارش", func() {
		qty, _ := strconv.Atoi(quantity.Text)
		p := productMap[productSelect.Selected]
		item := models.OrderItem{
			ProductID:   p.ID,
			ProductName: p.Name,
			Quantity:    qty,
			UnitPrice:   p.SalePrice,
		}
		selectedItems = append(selectedItems, item)

		list.Text += fmt.Sprintf("%s × %d = %.0f\n", item.ProductName, item.Quantity, float64(item.Quantity)*item.UnitPrice)
		list.Refresh()
	})

	description := widget.NewMultiLineEntry()
	description.SetPlaceHolder("توضیحات سفارش")

	submitBtn := widget.NewButton("ثبت سفارش", func() {
		var total float64
		for _, i := range selectedItems {
			total += float64(i.Quantity) * i.UnitPrice
		}

		order := models.Order{
			Items:       selectedItems,
			Note:        description.Text,
			TotalAmount: total,
		}
		if table != nil {
			order.TableID = &table.ID
		}

		db.SubmitOrder(order)
		w.Close()
		refresh()
	})

	form := container.NewVBox(
		widget.NewLabel("محصول"), productSelect,
		widget.NewLabel("تعداد"), quantity,
		addBtn,
		list,
		widget.NewLabel("توضیحات"), description,
		submitBtn,
	)

	w.SetContent(form)
	w.Resize(fyne.NewSize(400, 400))
	if table != nil {
		db.UpdateTableStatus(table.ID, "مشغول") // به وضعیت مشغول برو
	}
	return w
}
