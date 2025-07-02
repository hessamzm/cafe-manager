package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func LoadMainUI(w fyne.Window) {
	tabs := container.NewAppTabs(
		container.NewTabItem("میزها", widget.NewLabel("در حال بارگذاری میزها...")),
		container.NewTabItem("محصولات", widget.NewLabel("در حال بارگذاری محصولات...")),
		container.NewTabItem("سفارشات", widget.NewLabel("در حال بارگذاری سفارشات...")),
		container.NewTabItem("گزارشات", widget.NewLabel("در حال بارگذاری گزارشات...")),
		container.NewTabItem("حسابداری", widget.NewLabel("در حال بارگذاری حسابداری...")),
	)
	w.SetContent(tabs)
}
