package schema

import (
	"github.com/GuiaBolso/darwin"
	"github.com/jmoiron/sqlx"
)

var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Add products",
		Script: `
	CREATE TABLE products (
		id	UUID,
		name	TEXT,
		cost	INT,
		quantity	INT,
		date_created	TIMESTAMP,
		date_updated	TIMESTAMP,

		PRIMARY KEY (id)
	);`,
	},
	{
		Version:     2,
		Description: "Add sales",
		Script: `
	CREATE TABLE sales (
		sale_id	UUID,
		product_id	UUID,
		quantity	INT,
		paid	INT,
		date_created	TIMESTAMP,

		PRIMARY KEY (sale_id),
		FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
	
	);`,
	},
	{
		Version:     3,
		Description: "Add roles",
		Script: `
	CREATE TABLE users(
		user_id	UUID,
		name	TEXT,
		email	TEXT UNIQUE,
		password_hash	TEXT,
		roles	TEXT[],
		date_created	TIMESTAMP,
		date_updated	TIMESTAMP,

		PRIMARY KEY (user_id)
	);`,
	},
	{
		Version:     4,
		Description: "Add user column to products",
		Script: `
	ALTER TABLE products
		ADD COLUMN user_id UUID DEFAULT '00000000-0000-0000-0000-000000000000'
		`,
	},
}

func Migrate(db *sqlx.DB) error {
	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})
	d := darwin.New(driver, migrations, nil)
	return d.Migrate()
}
