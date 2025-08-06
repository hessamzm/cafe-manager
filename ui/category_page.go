package ui

import (
	"cafe-manager/db"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func NewCategoryPage() fyne.CanvasObject {
	categoriesBox := container.NewVBox()

	var refreshCategories func()

	// ساختن ردیف دسته‌بندی
	categoryItem := func(name string) fyne.CanvasObject {
		label := widget.NewLabel(name)

		editBtn := widget.NewButton("ویرایش", func() {
			input := widget.NewEntry()
			input.SetText(name)

			dialog := container.NewVBox(
				widget.NewLabel("ویرایش دسته‌بندی"),
				input,
				widget.NewButton("ذخیره", func() {
					db.UpdateCategory(name, input.Text) // تابع موردنیاز در db: (oldName, newName)
					refreshCategories()
				}),
			)

			w := fyne.CurrentApp().NewWindow("ویرایش دسته‌بندی")
			w.SetContent(dialog)
			w.Resize(fyne.NewSize(300, 150))
			w.Show()
		})

		deleteBtn := widget.NewButton("حذف", func() {
			db.DeleteCategory(name)
			refreshCategories()
		})

		return container.NewHBox(label, layout.NewSpacer(), editBtn, deleteBtn)
	}

	refreshCategories = func() {
		categoriesBox.Objects = nil
		for _, cat := range db.GetAllCategories() {
			categoriesBox.Add(categoryItem(cat))
		}
		categoriesBox.Refresh()
	}

	input := widget.NewEntry()
	input.SetPlaceHolder("نام دسته‌بندی جدید")
	addBtn := widget.NewButton("افزودن", func() {
		db.AddCategory(input.Text)
		input.SetText("")
		refreshCategories()
	})

	toolbar := container.NewHBox(widget.NewLabel("مدیریت دسته‌بندی‌ها"), layout.NewSpacer())

	refreshCategories()

	return container.NewBorder(toolbar, container.NewVBox(input, addBtn), nil, nil, container.NewVScroll(categoriesBox))
}
