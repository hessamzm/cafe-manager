package db

func Migrate() {
	createTablesTable := `
	CREATE TABLE IF NOT EXISTS tables (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		status TEXT NOT NULL
	);`

	_, err := DB.Exec(createTablesTable)
	if err != nil {
		panic("خطا در ساخت جدول tables: " + err.Error())
	}
}
