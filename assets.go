package main

import (
	_ "embed" // این ایمپورت برای فعال‌سازی کامنت //go:embed ضروری است
	"fyne.io/fyne/v2"
)

// //go:embed دستوری است که به کامپایلر Go می‌گوید محتوای فایل مشخص‌شده را
// در یک متغیر از نوع بایت (`[]byte`) قرار دهد.

// Embed کردن فونت
//
//go:embed assets/fonts/IRANSansXFaNum-Regular.ttf
var fontIRANSans []byte

// Embed کردن آیکون برنامه
//
//go:embed icon.png
var appIcon []byte

// حالا متغیرهای بایت را به نوع Fyne Resource تبدیل می‌کنیم تا در برنامه قابل استفاده باشند.
// این متغیرها به صورت عمومی (با حرف بزرگ) تعریف می‌شوند تا از main.go قابل دسترسی باشند.

var FontIRANSansRes = &fyne.StaticResource{
	StaticName:    "IRANSansXFaNum-Regular.ttf",
	StaticContent: fontIRANSans,
}

var IconRes = &fyne.StaticResource{
	StaticName:    "icon.png",
	StaticContent: appIcon,
}
