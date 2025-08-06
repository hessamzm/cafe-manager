<<<<<<< HEAD
# Cafe Manager

**Desktop application for café management**  
یک برنامه دسکتاپ برای مدیریت کافه

---

## 🌟 Features | امکانات  
- **Table Management | مدیریت میزها**  
  - Add/remove tables | افزودن/حذف میز  
  - Change status: Available, Reserved, Busy | تغییر وضعیت: آماده، رزرو، مشغول  
  - Open order form per table | باز کردن فرم ثبت سفارش برای هر میز  
- **Product Catalog | کاتالوگ محصولات**  
  - Categories & products | دسته‌بندی و محصولات  
  - Add, edit, delete categories & products | افزودن، ویرایش، حذف  
- **Order Processing | پردازش سفارش**  
  - Dine-in & takeaway | سفارش درون‌کافه و بیرون‌بر  
  - Multi-item orders with quantities & notes | سفارش چندمحصولی با تعداد و یادداشت  
  - Settle orders, calculate totals & apply discounts | تسویه، محاسبه مجموع، اعمال تخفیف  
- **Reporting | گزارش‌گیری**  
  - Daily, weekly, monthly & custom-range sales reports (Persian dates) | گزارش روزانه، هفتگی، ماهانه و بازه‌ای (تاریخ شمسی)  
  - Profit & Loss closing with fixed cost entries | بستن حساب سود و زیان با ورود هزینه‌های ثابت  
  - Excel export: detailed sales, item summary, order-type breakdown | خروجی اکسل: فروش جزئی، جمع اقلام، تفکیک نوع سفارش  
- **Accounting | حسابداری**  
  - Track fixed expenses (rent, utilities…) | ثبت هزینه‌های ثابت (اجاره، قبوض...)  
  - Calculate cost of goods sold vs. revenue | محاسبه بهای تمام‌شده و درآمد  
  - Generate profit/loss statements | تولید صورت سود و زیان  

---

## 🚀 Installation | نصب  
1. **Prerequisites | پیش‌نیازها**  
   - Go ≥1.18  
   - SQLite3  
   - Fyne toolkit  
2. **Clone the repository | کلون کردن مخزن**  
   ```bash
   git clone https://github.com/yourusername/cafe-manager.git
   cd cafe-manager
   ```  
3. **Build & Run | ساخت و اجرا**  
   ```bash
   go mod tidy
   go run main.go
   ```  
   Or build a binary:  
   ```bash
   go build -o CafeManager
   ./CafeManager
   ```  

---

## 📋 Usage | راهنمای استفاده  
1. **First launch**: creates `cafe.db` | اولین اجرا، ایجاد فایل دیتابیس  
2. **Tables**: Add tables → Operations | مدیریت میزها → عملیات  
3. **Products**: Categories & products → CRUD | مدیریت محصولات  
4. **Orders**: From tables or “Takeaway” | ثبت سفارش درون‌کافه یا بیرون‌بر  
5. **Reports**: Reports tab → choose date range → view/export | گزارش‌ها  
6. **Profit & Loss**: P&L tab → enter fixed costs → close month | سود و زیان  

---

## 🤝 Contributing | مشارکت  
Issues and Pull Requests are welcome!  
برای باگ‌ها و درخواست ویژگی، PR ارسال کنید.

---

## 📫 Contact | تماس  
For customization or support, email:  
✉️ **zaremahmoodih@gmail.com**

---

## ⚖️ License | مجوز  
MIT License  
