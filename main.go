package main

import (
	"cafe-manager/db"
	"cafe-manager/ui"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

type customTheme struct {
	baseTheme fyne.Theme
	fontRes   fyne.Resource
}

// Color برمی‌گرداند (از تم پایه استفاده می‌کنیم)
func (t *customTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return t.baseTheme.Color(name, variant)
}

// Icon برمی‌گرداند
func (t *customTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return t.baseTheme.Icon(name)
}

// Font برمی‌گرداند (همیشه قلم ما)
func (t *customTheme) Font(style fyne.TextStyle) fyne.Resource {
	return t.fontRes
}

// Size برمی‌گرداند
func (t *customTheme) Size(name fyne.ThemeSizeName) float32 {
	return t.baseTheme.Size(name)
}

func main() {
	// دیگر نیازی به بارگذاری فایل از مسیر نیست چون از منابع embed شده استفاده می‌کنیم
	// fontRes, err := fyne.LoadResourceFromPath("assets/fonts/IRANSansXFaNum-Regular.ttf")
	// ...

	a := app.NewWithID("com.hessamzm.cafemanager")

	// تنظیم تم سفارشی با استفاده از فونت embed شده
	custom := &customTheme{
		baseTheme: theme.DefaultTheme(),
		fontRes:   FontIRANSansRes, // <<-- استفاده از متغیر تعریف‌شده در assets.go
	}
	a.Settings().SetTheme(custom)

	w := a.NewWindow("مدیریت کافه")

	// تنظیم آیکون با استفاده از منبع embed شده
	// دیگر نیازی به استفاده از storage نیست
	w.SetIcon(IconRes) // <<-- استفاده از متغیر تعریف‌شده در assets.go

	db.InitDB()

	tabs := container.NewAppTabs(
		container.NewTabItem("مدیریت میزها", ui.NewTablesPage(w)),
		container.NewTabItem("محصولات", ui.NewProductPage(w)),
		container.NewTabItem("دسته‌بندی‌ها", ui.NewCategoryPage()),
		container.NewTabItem("گزارش‌ها", ui.NewReportsPage(w)),
		container.NewTabItem("سود و زیان", ui.NewSoodZiyanPage(w)),
	)

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(1200, 900))
	w.ShowAndRun()
}
