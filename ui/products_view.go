package ui

import (
	"cafe-manager/db"
	"cafe-manager/models"
	"cafe-manager/utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

var showBoxed = false

func NewProductPage(win fyne.Window) fyne.CanvasObject {
	productsBox := container.NewVBox()
	showBoxed := false
	var RefreshProducts func()
	RefreshProducts = func() {
		productsBox.Objects = nil
		all := db.GetAllProducts() // باید []models.Product برگرداند

		grouped := map[string][]models.Product{}
		for _, p := range all {
			grouped[p.Category] = append(grouped[p.Category], p)
		}

		// ساختن کارت‌های محصول (داخل تابع برای دسترسی به refresh)
		productCard := func(prod models.Product) fyne.CanvasObject {
			name := widget.NewLabel(prod.Name + " (" + strconv.Itoa(int(prod.SalePrice)) + " تومان)")

			editBtn := widget.NewButtonWithIcon("", utils.LoadIcon("edite.png"), func() {
				ShowEditProductForm(prod, RefreshProducts).Show()
			})

			deleteBtn := widget.NewButtonWithIcon("", utils.LoadIcon("trash.png"), func() {
				db.DeleteProduct(prod.ID)

			})

			return container.NewHBox(name, layout.NewSpacer(), editBtn, deleteBtn)
		}

		for cat, items := range grouped {
			productsBox.Add(widget.NewLabelWithStyle(cat, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))

			if showBoxed {
				grid := container.NewGridWithColumns(3)
				for _, prod := range items {
					grid.Add(productCard(prod))
				}
				productsBox.Add(grid)
			} else {
				box := container.NewVBox()
				for _, prod := range items {
					box.Add(productCard(prod))
				}
				productsBox.Add(box)
			}
		}
		productsBox.Refresh()
	}

	toolbar := container.NewHBox(
		widget.NewButtonWithIcon("اضافه کردن محصول", utils.LoadIcon("addprodoct.png"), func() {
			ShowAddProductForm(RefreshProducts).Show()
		}),
		widget.NewButtonWithIcon("بروزرسانی", utils.LoadIcon("refresh.png"), func() {
			RefreshProducts()
		}),
		layout.NewSpacer(),
		widget.NewButtonWithIcon("ردیف", utils.LoadIcon("liner.png"), func() {
			showBoxed = false
			RefreshProducts()
		}),
		widget.NewButtonWithIcon("جعبه ای", utils.LoadIcon("boxed.png"), func() {
			showBoxed = true
			RefreshProducts()
		}),
	)

	RefreshProducts()

	return container.NewBorder(toolbar, nil, nil, nil, container.NewVScroll(productsBox))
}
