package db

import (
	"cafe-manager/models"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

var DB *sql.DB

var tables []models.Table

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./cafe.db")
	if err != nil {
		log.Fatal(err)
	}

	createTables()
}

func createTables() {
	categoryTable := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE
	);`
	DB.Exec(categoryTable)

	productTable := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		category TEXT,
		cost_price REAL,
		sale_price REAL
	);`
	DB.Exec(productTable)
	_, err := DB.Exec(productTable)
	if err != nil {
		log.Fatal("خطا در ساخت جدول محصول:", err)
	}

	tableTable := `
	CREATE TABLE IF NOT EXISTS tables (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		status TEXT
	);`
	DB.Exec(tableTable)

	orderTable := `
CREATE TABLE IF NOT EXISTS orders (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	table_id INTEGER,
	description TEXT,
	total_price REAL,
	created_at TEXT,
	settled INTEGER DEFAULT 0
);`
	DB.Exec(orderTable)

	orderItemsTable := `
CREATE TABLE IF NOT EXISTS order_items (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	order_id INTEGER,
	product_id INTEGER,
	product_name TEXT,
	quantity INTEGER,
	unit_price REAL
);`
	DB.Exec(orderItemsTable)
	MonthlyClosings := `	
CREATE TABLE IF NOT EXISTS monthly_closings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		start_date TEXT,
		end_date TEXT,
		total_sales REAL,
		total_cost REAL,
		fixed_expenses REAL,
		profit REAL
	);`
	DB.Exec(MonthlyClosings)

	FixedExpenses := `
CREATE TABLE IF NOT EXISTS fixed_expenses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    closing_id INTEGER,
    title TEXT,
    amount REAL
);

`
	DB.Exec(FixedExpenses)
}

func AddProduct(name, category string, costPrice, salePrice float64) {
	_, err := DB.Exec("INSERT INTO products (name, category, cost_price, sale_price) VALUES (?, ?, ?, ?)", name, category, costPrice, salePrice)
	if err != nil {
		log.Println("خطا در افزودن محصول:", err)
	}
}

func GetAllProducts() []models.Product {
	rows, err := DB.Query("SELECT id, name, category, cost_price, sale_price FROM products ORDER BY category, name")
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.CostPrice, &p.SalePrice)
		if err != nil {
			log.Println(err)
			continue
		}
		products = append(products, p)
	}

	return products
}

func AddCategory(name string) {
	_, err := DB.Exec("INSERT INTO categories (name) VALUES (?)", name)
	if err != nil {
		log.Println("خطا در افزودن دسته:", err)
	}
}

func GetAllCategories() []string {
	rows, _ := DB.Query("SELECT name FROM categories ORDER BY name")
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		categories = append(categories, name)
	}
	return categories
}
func DeleteProduct(id int) {
	_, err := DB.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		log.Println("خطا در حذف محصول:", err)
	}
}

func UpdateProduct(id int, name, category string, costPrice, salePrice float64) {
	_, err := DB.Exec("UPDATE products SET name = ?, category = ?, cost_price = ?, sale_price = ? WHERE id = ?", name, category, costPrice, salePrice, id)
	if err != nil {
		log.Println("خطا در بروزرسانی محصول:", err)
	}
}

func UpdateCategory(oldName, newName string) {
	_, err := DB.Exec("UPDATE categories SET name = ? WHERE name = ?", newName, oldName)
	if err != nil {
		log.Println("خطا در بروزرسانی دسته‌بندی:", err)
	}

	// همچنین باید دسته‌بندی محصولات مرتبط هم آپدیت شود
	_, err = DB.Exec("UPDATE products SET category = ? WHERE category = ?", newName, oldName)
	if err != nil {
		log.Println("خطا در بروزرسانی دسته‌بندی محصولات:", err)
	}
}
func DeleteCategory(name string) {
	_, err := DB.Exec("DELETE FROM categories WHERE name = ?", name)
	if err != nil {
		log.Println("خطا در حذف دسته‌بندی:", err)
	}

	// محصولات مرتبط همچنان باقی می‌مانند. می‌توان به‌دلخواه پاک یا نال کرد.
	_, err = DB.Exec("DELETE FROM products WHERE category = ?", name)
	if err != nil {
		log.Println("خطا در خالی‌کردن دسته‌بندی محصولات:", err)
	}
}

func GetAllTables() []models.Table {
	rows, err := DB.Query("SELECT id, name, status FROM tables")
	if err != nil {
		log.Println("خطا در دریافت لیست میزها:", err)
		return nil
	}
	defer rows.Close()

	var tables []models.Table
	for rows.Next() {
		var t models.Table
		err := rows.Scan(&t.ID, &t.Name, &t.Status)
		if err == nil {
			tables = append(tables, t)
		}
	}
	return tables
}
func AddTable(name string) {
	_, err := DB.Exec("INSERT INTO tables (name, status) VALUES (?, ?)", name, models.TableAvailable)
	if err != nil {
		log.Println("خطا در افزودن میز:", err)
	}
}

func UpdateTableStatus(id int, status models.TableStatus) {
	_, err := DB.Exec("UPDATE tables SET status = ? WHERE id = ?", status, id)
	if err != nil {
		log.Println("خطا در بروزرسانی وضعیت میز:", err)
	}
}

func GetTableByID(id int) *models.Table {
	row := DB.QueryRow("SELECT id, name, status FROM tables WHERE id = ?", id)
	var t models.Table
	err := row.Scan(&t.ID, &t.Name, &t.Status)
	if err != nil {
		return nil
	}
	return &t
}

func ToggleTableStatus(id int) {
	t := GetTableByID(id)
	if t == nil {
		return
	}
	var newStatus models.TableStatus
	switch t.Status {
	case models.TableAvailable:
		newStatus = models.TableBusy
	case models.TableBusy:
		newStatus = models.TableReserved
	case models.TableReserved:
		newStatus = models.TableAvailable
	default:
		newStatus = models.TableAvailable
	}
	UpdateTableStatus(id, newStatus)
}
func CloseTable(id int) {
	for i := range tables {
		if tables[i].ID == id {
			tables[i].Status = models.TableAvailable
			break
		}
	}
}

func SubmitOrder(order models.Order) {
	tx, err := DB.Begin()
	if err != nil {
		log.Println("خطا در شروع تراکنش سفارش:", err)
		return
	}
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	settled := 0
	if order.TableID == nil {
		settled = 1
	}
	res, err := tx.Exec(`
        INSERT INTO orders (table_id, description, total_price, created_at, settled)
        VALUES (?, ?, ?, ?, ?)
    `, order.TableID, order.Note, order.TotalAmount, createdAt, settled)
	if err != nil {
		log.Println("خطا در درج سفارش:", err)
		tx.Rollback()
		return
	}

	orderID, _ := res.LastInsertId()
	for _, item := range order.Items {
		_, err := tx.Exec(`
            INSERT INTO order_items (order_id, product_id, product_name, quantity, unit_price)
            VALUES (?, ?, ?, ?, ?)
        `, orderID, item.ProductID, item.ProductName, item.Quantity, item.UnitPrice)
		if err != nil {
			log.Println("خطا در درج آیتم سفارش:", err)
			tx.Rollback()
			return
		}
	}

	tx.Commit()

	if order.TableID != nil {
		UpdateTableStatus(*order.TableID, models.TableBusy)
	}
}

func GetOrdersByTableID(tableID int) ([]models.Order, float64) {
	rows, _ := DB.Query("SELECT id, description, total_price FROM orders WHERE table_id = ? AND settled = 0", tableID)
	defer rows.Close()

	var orders []models.Order
	var total float64

	for rows.Next() {
		var o models.Order
		rows.Scan(&o.ID, &o.Note, &o.TotalAmount)
		o.TableID = &tableID

		o.Items = GetOrderItems(o.ID)
		orders = append(orders, o)
		total += o.TotalAmount
	}
	return orders, total
}

func GetOrderItems(orderID int) []models.OrderItem {
	rows, _ := DB.Query("SELECT product_id, product_name, quantity, unit_price FROM order_items WHERE order_id = ?", orderID)
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		rows.Scan(&item.ProductID, &item.ProductName, &item.Quantity, &item.UnitPrice)
		items = append(items, item)
	}
	return items
}

func CloseTableAndClearOrders(tableID int) {
	DB.Exec("UPDATE orders SET settled = 1 WHERE table_id = ?", tableID)
	UpdateTableStatus(tableID, models.TableAvailable)
}

func GetOrdersBetween(startDate, endDate string) ([]models.Order, float64) {
	rows, err := DB.Query(`
        SELECT id, table_id, description, total_price, created_at
        FROM orders
        WHERE settled = 1
          AND date(created_at) BETWEEN date(?) AND date(?)
        ORDER BY created_at ASC
    `, startDate, endDate)
	if err != nil {
		log.Println("خطا در دریافت سفارشات در بازه زمانی:", err)
		return nil, 0
	}
	defer rows.Close()

	var orders []models.Order
	var total float64

	for rows.Next() {
		var o models.Order
		var tableID sql.NullInt64
		var createdAt string

		err := rows.Scan(&o.ID, &tableID, &o.Note, &o.TotalAmount, &createdAt)
		if err != nil {
			log.Println("خطا در اسکن سفارش:", err)
			continue
		}
		if tableID.Valid {
			tid := int(tableID.Int64)
			o.TableID = &tid
		}
		o.CreatedAt = createdAt // 🔻 این خط باید قبل از append باشد
		o.Items = GetOrderItems(o.ID)
		orders = append(orders, o)
		total += o.TotalAmount
	}

	return orders, total
}
func GetOrderItemsForTable(tableID int) []models.OrderItem {
	rows, _ := DB.Query(`
		SELECT oi.id, oi.order_id, oi.product_id, p.name, oi.quantity, oi.unit_price
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		JOIN products p ON p.id = oi.product_id
		WHERE o.table_id = ? AND o.settled = 0
	`, tableID)
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.ProductName, &item.Quantity, &item.UnitPrice)
		items = append(items, item)
	}
	return items
}

//func DeleteOrderItemByID(id int) {
//	DB.Exec("DELETE FROM order_items WHERE id = ?", id)
//}

func UpdateOrderItemQuantity(id int, newQty int) {
	DB.Exec("UPDATE order_items SET quantity = ? WHERE id = ?", newQty, id)
}

// RecalculateOrderTotalAmount محاسبه و ذخیره‌ی مجموع مبلغ سفارش
func RecalculateOrderTotalAmount(orderID int) error {
	rows, err := DB.Query("SELECT quantity, unit_price FROM order_items WHERE order_id = ?", orderID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var total float64
	for rows.Next() {
		var quantity int
		var price float64
		rows.Scan(&quantity, &price)
		total += float64(quantity) * price
	}

	_, err = DB.Exec("UPDATE orders SET total_price = ? WHERE id = ?", total, orderID)
	return err
}
func DeleteOrderItemByID(id int) {
	var orderID int
	_ = DB.QueryRow("SELECT order_id FROM order_items WHERE id = ?", id).Scan(&orderID)

	DB.Exec("DELETE FROM order_items WHERE id = ?", id)

	// بعد از حذف، مبلغ کل سفارش را به‌روز کن
	RecalculateOrderTotalAmount(orderID)
}
func HasOrdersForTable(tableID int) bool {
	orders, _ := GetOrdersByTableID(tableID)
	return len(orders) > 0
}
