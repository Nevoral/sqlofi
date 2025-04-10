package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Nevoral/sqlofi/sqlite"
	_ "github.com/tursodatabase/go-libsql"
)

// Product model for testing
type Product struct {
	Id          int64          `sqlofi:"PRIMARY KEY AUTOINCREMENT"`
	Name        string         `sqlofi:"NOT NULL"`
	Description string         `sqlofi:"NOT NULL"`
	Price       float64        `sqlofi:"NOT NULL CHECK(Price >= 0)"`
	Quantity    int            `sqlofi:"NOT NULL DEFAULT 0"`
	CategoryId  sql.NullInt64  `sqlofi:"REFERENCES Category (Id)"`
	Active      sql.NullBool   `sqlofi:"DEFAULT 1"`
	Created     sql.NullString `sqlofi:"DEFAULT (datetime('now'))"`
}

// Category model for testing
type Category struct {
	Id   int64  `sqlofi:"PRIMARY KEY AUTOINCREMENT"`
	Name string `sqlofi:"NOT NULL UNIQUE"`
}

// User model for testing relationships
type User struct {
	Id       int64          `sqlofi:"PRIMARY KEY AUTOINCREMENT"`
	Username string         `sqlofi:"NOT NULL UNIQUE"`
	Email    string         `sqlofi:"NOT NULL UNIQUE"`
	Age      sql.NullInt64  `sqlofi:"CHECK(Age IS NULL OR Age >= 18)"`
	Active   sql.NullBool   `sqlofi:"DEFAULT 1"`
	Created  sql.NullString `sqlofi:"DEFAULT (datetime('now'))"`
}

// Order model for testing relationships
type Order struct {
	Id      int64          `sqlofi:"PRIMARY KEY AUTOINCREMENT"`
	UserId  sql.NullInt64  `sqlofi:"REFERENCES User (Id)"`
	Status  string         `sqlofi:"NOT NULL DEFAULT 'pending'"`
	Total   float64        `sqlofi:"NOT NULL DEFAULT 0"`
	Created sql.NullString `sqlofi:"DEFAULT (datetime('now'))"`
}

// OrderItem model for testing relationships
type OrderItem struct {
	Id        int64         `sqlofi:"PRIMARY KEY AUTOINCREMENT"`
	OrderId   sql.NullInt64 `sqlofi:"REFERENCES Order (Id)"`
	ProductId sql.NullInt64 `sqlofi:"REFERENCES Product (Id)"`
	Quantity  int           `sqlofi:"NOT NULL DEFAULT 1"`
	Price     float64       `sqlofi:"NOT NULL"`
	Subtotal  float64       `sqlofi:"NOT NULL"`
}

// Test expression building
func testExpressions() {
	fmt.Println("===== Testing Expressions =====")

	// Simple expressions
	expr1 := sqlite.Expr("Price > 100")
	fmt.Println("Expression 1:", expr1.Build())

	expr2 := sqlite.Expr("Quantity BETWEEN 10 AND 100")
	fmt.Println("Expression 2:", expr2.Build())

	// Compound expressions
	expr3 := sqlite.Expr("Price > 100 AND Quantity > 0")
	fmt.Println("Expression 3:", expr3.Build())

	expr4 := sqlite.Expr("CategoryId IS NOT NULL OR Price < 10")
	fmt.Println("Expression 4:", expr4.Build())

	// Function expressions
	expr5 := sqlite.Expr("length(Name) > 5")
	fmt.Println("Expression 5:", expr5.Build())

	expr6 := sqlite.Expr("datetime('now') > Created")
	fmt.Println("Expression 6:", expr6.Build())

	// Subquery expressions
	expr7 := sqlite.Expr("Price > (SELECT AVG(Price) FROM product)")
	fmt.Println("Expression 7:", expr7.Build())

	expr8 := sqlite.Expr("CategoryId IN (SELECT Id FROM Category WHERE Name LIKE 'Electronics%')")
	fmt.Println("Expression 8:", expr8.Build())

	fmt.Println()
}

// Test SELECT statement building
func testSelectStatements() {
	fmt.Println("===== Testing SELECT Statements =====")

	// Simple SELECT
	select1 := sqlite.SELECT(sqlite.ALL,
		sqlite.NewExpressionColumn(sqlite.NewExpression("Id")),
		sqlite.NewExpressionColumn(sqlite.NewExpression("Name")),
		sqlite.NewExpressionColumn(sqlite.NewExpression("Price")),
	).FROM(sqlite.NewTableFrom("product"))

	fmt.Println("Simple SELECT:")
	fmt.Println(select1.Build())
	fmt.Println()

	// SELECT with WHERE clause
	select2 := sqlite.SELECT(sqlite.ALL,
		sqlite.NewWildcardColumn(),
	).FROM(
		sqlite.NewTableFrom("product"),
	).WHERE(
		sqlite.NewExpression("Price > 100"),
	)

	fmt.Println("SELECT with WHERE:")
	fmt.Println(select2.Build())
	fmt.Println()

	// SELECT with column alias
	select3 := sqlite.SELECT(sqlite.ALL,
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("Id"), "ProductID"),
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("Name"), "ProductName"),
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("Price * 1.2"), "PriceWithTax"),
	).FROM(
		sqlite.NewTableFrom("product"),
	)

	fmt.Println("SELECT with column aliases:")
	fmt.Println(select3.Build())
	fmt.Println()

	// SELECT with JOIN
	select4 := sqlite.SELECT(sqlite.ALL,
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("p.Id"), "ProductID"),
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("p.Name"), "ProductName"),
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("c.Name"), "CategoryName"),
	).FROM(
		sqlite.NewTableFrom("product").Alias("p").Join(
			sqlite.NewTableJoin(sqlite.LEFT_JOIN, "category").Alias("c").On(
				sqlite.NewExpression("p.CategoryId = c.Id"),
			),
		),
	)

	fmt.Println("SELECT with JOIN:")
	fmt.Println(select4.Build())
	fmt.Println()

	// SELECT with multiple JOINs
	select5 := sqlite.SELECT(sqlite.ALL,
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("o.Id"), "OrderID"),
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("u.Username"), "Customer"),
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("p.Name"), "ProductName"),
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("oi.Quantity"), "Quantity"),
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("oi.Subtotal"), "Subtotal"),
	).FROM(
		sqlite.NewTableFrom("order").Alias("o").
			Join(sqlite.NewTableJoin(sqlite.INNER_JOIN, "user").Alias("u").On(
				sqlite.NewExpression("o.UserId = u.Id"),
			)).
			Join(sqlite.NewTableJoin(sqlite.INNER_JOIN, "orderitem").Alias("oi").On(
				sqlite.NewExpression("oi.OrderId = o.Id"),
			)).
			Join(sqlite.NewTableJoin(sqlite.INNER_JOIN, "product").Alias("p").On(
				sqlite.NewExpression("oi.ProductId = p.Id"),
			)),
	)

	fmt.Println("SELECT with multiple JOINs:")
	fmt.Println(select5.Build())
	fmt.Println()

	// SELECT with GROUP BY and HAVING
	select6 := sqlite.SELECT(sqlite.ALL,
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("c.Name"), "CategoryName"),
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("COUNT(p.Id)"), "ProductCount"),
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("AVG(p.Price)"), "AveragePrice"),
	).FROM(
		sqlite.NewTableFrom("product").Alias("p").
			Join(sqlite.NewTableJoin(sqlite.INNER_JOIN, "category").Alias("c").On(
				sqlite.NewExpression("p.CategoryId = c.Id"),
			)),
	).GROUP_BY(
		sqlite.NewExpression("c.Name"),
	).HAVING(
		sqlite.NewExpression("COUNT(p.Id) > 2"),
	)

	fmt.Println("SELECT with GROUP BY and HAVING:")
	fmt.Println(select6.Build())
	fmt.Println()

	// SELECT with ORDER BY and LIMIT
	select7 := sqlite.SELECT(sqlite.ALL,
		sqlite.NewWildcardColumn(),
	).FROM(
		sqlite.NewTableFrom("product"),
	).ORDER_BY(
		sqlite.NewOrderBy(sqlite.NewExpression("Price"), sqlite.DESC),
		sqlite.NewOrderBy(sqlite.NewExpression("Name"), sqlite.ASC),
	).Limit(10).Offset(20)

	fmt.Println("SELECT with ORDER BY, LIMIT and OFFSET:")
	fmt.Println(select7.Build())
	fmt.Println()

	// SELECT with subquery
	subquery := sqlite.SELECT(sqlite.ALL,
		sqlite.NewExpressionColumn(sqlite.NewExpression("AVG(Price)")),
	).FROM(
		sqlite.NewTableFrom("product"),
	)

	select8 := sqlite.SELECT(sqlite.ALL,
		sqlite.NewWildcardColumn(),
	).FROM(
		sqlite.NewTableFrom("product"),
	).WHERE(
		sqlite.NewExpression(fmt.Sprintf("Price > (%s)", subquery.Build())),
	)

	fmt.Println("SELECT with subquery in WHERE clause:")
	fmt.Println(select8.Build())
	fmt.Println()

	// SELECT with subquery in FROM
	select9 := sqlite.SELECT(sqlite.ALL,
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("sub.ProductCount"), "CategoryProductCount"),
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("c.Name"), "CategoryName"),
	).FROM(
		sqlite.NewTableFrom("category").Alias("c").
			Join(sqlite.NewSubqueryJoin(sqlite.LEFT_JOIN,
				sqlite.SELECT(sqlite.ALL,
					sqlite.NewExpressionColumn(sqlite.NewExpression("CategoryId")),
					sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("COUNT(*)"), "ProductCount"),
				).FROM(
					sqlite.NewTableFrom("product"),
				).GROUP_BY(
					sqlite.NewExpression("CategoryId"),
				),
			).Alias("sub").On(
				sqlite.NewExpression("c.Id = sub.CategoryId"),
			)),
	)

	fmt.Println("SELECT with subquery in JOIN:")
	fmt.Println(select9.Build())
	fmt.Println()

	// SELECT DISTINCT
	select10 := sqlite.SELECT(sqlite.DISTINCT,
		sqlite.NewExpressionColumn(sqlite.NewExpression("CategoryId")),
	).FROM(
		sqlite.NewTableFrom("product"),
	).WHERE(
		sqlite.NewExpression("Price > 50"),
	)

	fmt.Println("SELECT DISTINCT:")
	fmt.Println(select10.Build())
	fmt.Println()
}

// Function to run queries against a test database
func executeQueries(dbPath string) {
	fmt.Println("===== Executing Queries Against Database =====")

	// Open database connection
	db, err := sql.Open("libsql", fmt.Sprintf("file:%s", dbPath))
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create tables
	createTables(db)

	// Insert test data
	insertTestData(db)

	// Execute select queries and print results
	runSelectQueries(db)
}

// Create database tables
func createTables(db *sql.DB) {
	fmt.Println("Creating tables...")

	// Create schema
	schema := sqlite.NewSchema("test.db").
		Pragma(
			sqlite.ForeignKeys().ValueType("ON"),
			sqlite.JournalModeWAL(""),
		).
		Table(
			sqlite.CREATE_TABLE(Category{}).IfNotExists(),
			sqlite.CREATE_TABLE(Product{}, Category{}).IfNotExists().
				ForeignKey(
					"fk_product_category",
					sqlite.FOREIGN_KEY(&Category{}, "CategoryId").
						ForeighColumns("Id").
						OnDelete(sqlite.SET_NULL).
						OnUpdate(sqlite.CASCADE),
				),
			sqlite.CREATE_TABLE(User{}).IfNotExists(),
			sqlite.CREATE_TABLE(Order{}, User{}).IfNotExists().
				ForeignKey(
					"fk_order_user",
					sqlite.FOREIGN_KEY(&User{}, "UserId").
						ForeighColumns("Id").
						OnDelete(sqlite.CASCADE).
						OnUpdate(sqlite.CASCADE),
				),
			sqlite.CREATE_TABLE(OrderItem{}, Order{}, Product{}).IfNotExists().
				ForeignKey(
					"fk_orderitem_order",
					sqlite.FOREIGN_KEY(&Order{}, "OrderId").
						ForeighColumns("Id").
						OnDelete(sqlite.CASCADE).
						OnUpdate(sqlite.CASCADE),
				).
				ForeignKey(
					"fk_orderitem_product",
					sqlite.FOREIGN_KEY(&Product{}, "ProductId").
						ForeighColumns("Id").
						OnDelete(sqlite.SET_NULL).
						OnUpdate(sqlite.CASCADE),
				),
		)

	// Execute schema creation
	_, err := db.Exec(schema.Build())
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
}

// Insert test data
func insertTestData(db *sql.DB) {
	fmt.Println("Inserting test data...")

	// Insert categories
	_, err := db.Exec(`
		INSERT INTO category (Name) VALUES
		('Electronics'),
		('Clothing'),
		('Books'),
		('Home & Garden')
	`)
	if err != nil {
		log.Fatalf("Failed to insert categories: %v", err)
	}

	// Insert products
	_, err = db.Exec(`
		INSERT INTO product (Name, Description, Price, Quantity, CategoryId) VALUES
		('Smartphone', 'Latest model smartphone', 999.99, 50, 1),
		('Laptop', 'High performance laptop', 1299.99, 30, 1),
		('T-shirt', 'Cotton t-shirt', 19.99, 100, 2),
		('Jeans', 'Blue denim jeans', 49.99, 75, 2),
		('Novel', 'Bestselling fiction book', 14.99, 200, 3),
		('Cookbook', 'Recipes from around the world', 24.99, 150, 3),
		('Plant pot', 'Ceramic pot for plants', 9.99, 80, 4),
		('Garden hose', 'Expandable garden hose', 34.99, 40, 4),
		('Headphones', 'Noise cancelling headphones', 199.99, 60, 1),
		('Tablet', 'Lightweight tablet', 399.99, 45, 1)
	`)
	if err != nil {
		log.Fatalf("Failed to insert products: %v", err)
	}

	// Insert users
	_, err = db.Exec(`
		INSERT INTO user (Username, Email) VALUES
		('john_doe', 'john@example.com'),
		('jane_smith', 'jane@example.com'),
		('mike_jones', 'mike@example.com')
	`)
	if err != nil {
		log.Fatalf("Failed to insert users: %v", err)
	}

	// Insert orders
	_, err = db.Exec(`
		INSERT INTO "order" (UserId, Status, Total) VALUES
		(1, 'completed', 1049.98),
		(1, 'pending', 199.99),
		(2, 'completed', 74.98),
		(3, 'processing', 1699.97)
	`)
	if err != nil {
		log.Fatalf("Failed to insert orders: %v", err)
	}

	// Insert order items
	_, err = db.Exec(`
		INSERT INTO orderitem (OrderId, ProductId, Quantity, Price, Subtotal) VALUES
		(1, 1, 1, 999.99, 999.99),
		(1, 3, 2, 19.99, 39.98),
		(2, 9, 1, 199.99, 199.99),
		(3, 5, 2, 14.99, 29.98),
		(3, 6, 1, 24.99, 24.99),
		(3, 7, 2, 9.99, 19.98),
		(4, 2, 1, 1299.99, 1299.99),
		(4, 9, 2, 199.99, 399.98)
	`)
	if err != nil {
		log.Fatalf("Failed to insert order items: %v", err)
	}
}

// Run SELECT queries and display results
func runSelectQueries(db *sql.DB) {
	fmt.Println("\nExecuting SELECT queries...")

	// Query 1: Products with price > 100
	query1 := sqlite.SELECT(sqlite.ALL,
		sqlite.NewExpressionColumn(sqlite.NewExpression("Id")),
		sqlite.NewExpressionColumn(sqlite.NewExpression("Name")),
		sqlite.NewExpressionColumn(sqlite.NewExpression("Price")),
	).FROM(
		sqlite.NewTableFrom("product"),
	).WHERE(
		sqlite.NewExpression("Price > 100"),
	).ORDER_BY(
		sqlite.NewOrderBy(sqlite.NewExpression("Price"), sqlite.DESC),
	)

	fmt.Println("\nQuery 1: Products with price > 100")
	fmt.Println("SQL:", query1.Build())

	rows, err := db.Query(query1.Build())
	if err != nil {
		log.Printf("Error executing query 1: %v", err)
	} else {
		defer rows.Close()
		fmt.Println("Results:")
		fmt.Println("ID | Name | Price")
		fmt.Println("---------------")
		var id int64
		var name string
		var price float64
		for rows.Next() {
			if err := rows.Scan(&id, &name, &price); err != nil {
				log.Printf("Error scanning row: %v", err)
				continue
			}
			fmt.Printf("%d | %s | %.2f\n", id, name, price)
		}
	}

	// Query 2: Products with their categories
	query2 := sqlite.SELECT(sqlite.ALL,
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("p.Id"), "ProductID"),
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("p.Name"), "ProductName"),
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("c.Name"), "CategoryName"),
	).FROM(
		sqlite.NewTableFrom("product").Alias("p").Join(
			sqlite.NewTableJoin(sqlite.LEFT_JOIN, "category").Alias("c").On(
				sqlite.NewExpression("p.CategoryId = c.Id"),
			),
		),
	)

	fmt.Println("\nQuery 2: Products with their categories")
	fmt.Println("SQL:", query2.Build())

	rows, err = db.Query(query2.Build())
	if err != nil {
		log.Printf("Error executing query 2: %v", err)
	} else {
		defer rows.Close()
		fmt.Println("Results:")
		fmt.Println("ProductID | ProductName | CategoryName")
		fmt.Println("------------------------------------")
		var prodID int64
		var prodName, catName string
		for rows.Next() {
			if err := rows.Scan(&prodID, &prodName, &catName); err != nil {
				log.Printf("Error scanning row: %v", err)
				continue
			}
			fmt.Printf("%d | %s | %s\n", prodID, prodName, catName)
		}
	}

	// Query 3: Order details with customer and product info
	query3 := sqlite.SELECT(sqlite.ALL,
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("o.Id"), "OrderID"),
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("u.Username"), "Customer"),
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("p.Name"), "ProductName"),
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("oi.Quantity"), "Quantity"),
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("oi.Price"), "Price"),
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("oi.Subtotal"), "Subtotal"),
	).FROM(
		sqlite.NewTableFrom("order").Alias("o").
			Join(sqlite.NewTableJoin(sqlite.INNER_JOIN, "user").Alias("u").On(
				sqlite.NewExpression("o.UserId = u.Id"),
			)).
			Join(sqlite.NewTableJoin(sqlite.INNER_JOIN, "orderitem").Alias("oi").On(
				sqlite.NewExpression("oi.OrderId = o.Id"),
			)).
			Join(sqlite.NewTableJoin(sqlite.INNER_JOIN, "product").Alias("p").On(
				sqlite.NewExpression("oi.ProductId = p.Id"),
			)),
	).ORDER_BY(
		sqlite.NewOrderBy(sqlite.NewExpression("o.Id"), sqlite.ASC),
		sqlite.NewOrderBy(sqlite.NewExpression("p.Name"), sqlite.ASC),
	)

	fmt.Println("\nQuery 3: Order details with customer and product info")
	fmt.Println("SQL:", query3.Build())

	rows, err = db.Query(query3.Build())
	if err != nil {
		log.Printf("Error executing query 3: %v", err)
	} else {
		defer rows.Close()
		fmt.Println("Results:")
		fmt.Println("OrderID | Customer | ProductName | Quantity | Price | Subtotal")
		fmt.Println("-------------------------------------------------------------")
		var orderID, qty int64
		var customer, prodName string
		var price, subtotal float64
		for rows.Next() {
			if err := rows.Scan(&orderID, &customer, &prodName, &qty, &price, &subtotal); err != nil {
				log.Printf("Error scanning row: %v", err)
				continue
			}
			fmt.Printf("%d | %s | %s | %d | %.2f | %.2f\n", orderID, customer, prodName, qty, price, subtotal)
		}
	}

	// Query 4: Product count and average price by category
	query4 := sqlite.SELECT(sqlite.ALL,
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("c.Name"), "CategoryName"),
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("COUNT(p.Id)"), "ProductCount"),
		sqlite.NewExpressionColumnWithAlias(sqlite.NewExpression("AVG(p.Price)"), "AveragePrice"),
	).FROM(
		sqlite.NewTableFrom("category").Alias("c").
			Join(sqlite.NewTableJoin(sqlite.LEFT_JOIN, "product").Alias("p").On(
				sqlite.NewExpression("p.CategoryId = c.Id"),
			)),
	).GROUP_BY(
		sqlite.NewExpression("c.Name"),
	).ORDER_BY(
		sqlite.NewOrderBy(sqlite.NewExpression("c.Name"), sqlite.ASC),
	)

	fmt.Println("\nQuery 4: Product count and average price by category")
	fmt.Println("SQL:", query4.Build())

	rows, err = db.Query(query4.Build())
	if err != nil {
		log.Printf("Error executing query 4: %v", err)
	} else {
		defer rows.Close()
		fmt.Println("Results:")
		fmt.Println("CategoryName | ProductCount | AveragePrice")
		fmt.Println("----------------------------------------")
		var catName string
		var prodCount int64
		var avgPrice float64
		for rows.Next() {
			if err := rows.Scan(&catName, &prodCount, &avgPrice); err != nil {
				log.Printf("Error scanning row: %v", err)
				continue
			}
			fmt.Printf("%s | %d | %.2f\n", catName, prodCount, avgPrice)
		}
	}
}

func main() {
	fmt.Println("SQLofi API Test Program")
	fmt.Println("======================")

	// Test expression building
	testExpressions()

	// Test SELECT statement building
	testSelectStatements()

	// Create a test database file
	dbPath := "test.db"

	// Remove previous database if it exists
	os.Remove(dbPath)

	// Execute queries against the database
	executeQueries(dbPath)

	fmt.Println("\nTests completed!")
}
