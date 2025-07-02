package main

import (
	"cafe-manager/db"
	"cafe-manager/ui"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow("مدیریت کافه")

	// اتصال به دیتابیس
	db.Connect()
	db.Migrate() // ایجاد جداول در صورت عدم وجود

	// بارگذاری UI
	ui.LoadMainUI(w)

	w.Resize(fyne.NewSize(900, 600))
	w.ShowAndRun()
}
