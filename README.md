📌 Overview | معرفی
Cafe Manager
A desktop application for managing café operations: table reservations, orders (dine-in & takeaway), inventory, reporting, and accounting.
یک برنامه دسکتاپ برای مدیریت کافه: رزرو میز، سفارشات (درون‌کافه و بیرون‌بر)، موجودی کالا، گزارش‌گیری و حسابداری.

🌟 Features | امکانات
Table Management | مدیریت میزها

Add/remove tables

Change status: Available, Reserved, Busy

Open order form per table

Product Catalog | کاتالوگ محصولات

Categories & products

Add, edit, delete categories & products

Order Processing | پردازش سفارش

Dine-in & takeaway

Multi-item orders with quantities & notes

Settle orders, calculate totals & apply discounts

Reporting | گزارش‌گیری

Daily, weekly, monthly & custom-range sales reports (Persian dates)

Profit & loss closing with fixed cost entries

Excel export: detailed sales, item summary, order-type breakdown

Accounting | حسابداری

Track fixed expenses (rent, utilities…)

Calculate cost of goods sold vs. revenue

Generate profit/loss statements

🚀 Installation | نصب
Prerequisites

Go ≥1.18

SQLite3 (database)

Fyne toolkit (go get fyne.io/fyne/v2)

Clone

bash
Copy
Edit
git clone https://github.com/yourusername/cafe-manager.git
cd cafe-manager
Build & Run

bash
Copy
Edit
go mod tidy
go run main.go
— or —

bash
Copy
Edit
go build -o CafeManager
./CafeManager
📋 Usage | راهنمای استفاده
First launch: creates cafe.db in working directory.

Tables: Add tables → click “Operations” → take or close orders.

Products: Define categories & products → use toolbar to add/edit/delete.

Orders: From tables page or “Takeaway” → select items & quantities → save.

Reports: Go to Reports tab → choose date range or presets → view & export to Excel.

Profit & Loss: In P&L tab → enter fixed monthly costs → generate closing & export Excel.

🤝 Contributing | مشارکت
Feel free to submit issues or pull requests!
برای باگ‌ها، درخواست ویژگی یا درگیری کد خوشحال می‌شویم Issue یا PR ارسال کنید.

📫 Contact | تماس
For customization, upgrades or support, email:

✉️ your.email@example.com

⚖️ License | مجوز
This project is licensed under the MIT License.
این پروژه تحت مجوز MIT عرضه می‌شود.

Thank you for using Cafe Manager!
از همراهی شما سپاس‌گزاریم!
