package ui

import (
	"cafe-manager/db"
	"cafe-manager/models"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

func ShowEditProductForm(p models.Product, refresh func()) fyne.Window {
	win := fyne.CurrentApp().NewWindow("ویرایش محصول")

	name := widget.NewEntry()
	name.SetText(p.Name)

	category := widget.NewSelect(db.GetAllCategories(), nil)
	category.SetSelected(p.Category)

	cost := widget.NewEntry()
	cost.SetText(strconv.FormatFloat(p.CostPrice, 'f', 0, 64))

	sale := widget.NewEntry()
	sale.SetText(strconv.FormatFloat(p.SalePrice, 'f', 0, 64))

	form := container.NewVBox(
		widget.NewLabel("نام محصول"), name,
		widget.NewLabel("دسته‌بندی"), category,
		widget.NewLabel("قیمت تمام‌شده"), cost,
		widget.NewLabel("قیمت فروش"), sale,
		widget.NewButton("ذخیره", func() {
			costVal, _ := strconv.ParseFloat(cost.Text, 64)
			saleVal, _ := strconv.ParseFloat(sale.Text, 64)
			db.UpdateProduct(p.ID, name.Text, category.Selected, costVal, saleVal)
			win.Close()
			refresh()
		}),
	)
	refresh()
	win.SetContent(form)
	win.Resize(fyne.NewSize(300, 300))
	return win
}

func ShowAddProductForm(refresh func()) fyne.Window {
	win := fyne.CurrentApp().NewWindow("افزودن محصول")

	name := widget.NewEntry()
	category := widget.NewSelect(db.GetAllCategories(), nil)
	cost := widget.NewEntry()
	sale := widget.NewEntry()

	form := container.NewVBox(
		widget.NewLabel("نام محصول"), name,
		widget.NewLabel("دسته‌بندی"), category,
		widget.NewLabel("قیمت تمام‌شده"), cost,
		widget.NewLabel("قیمت فروش"), sale,
		widget.NewButton("ذخیره", func() {
			costVal, _ := strconv.ParseFloat(cost.Text, 64)
			saleVal, _ := strconv.ParseFloat(sale.Text, 64)
			db.AddProduct(name.Text, category.Selected, costVal, saleVal)
			win.Close()
			refresh()
		}),
	)
	refresh()
	win.SetContent(form)
	win.Resize(fyne.NewSize(300, 300))
	return win
}
