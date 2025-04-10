package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Nevoral/sqlofi/sqlite"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Product model
type Product struct {
	Id          int64          `sqlofi:"PRIMARY KEY AUTOINCREMENT"`
	Name        string         `sqlofi:"NOT NULL UNIQUE"`
	Description sql.NullString `sqlofi:""`
	Price       float64        `sqlofi:"NOT NULL CHECK(Price >= 0)"`
	Quantity    int            `sqlofi:"NOT NULL DEFAULT 0 CHECK(Quantity >= 0)"`
	CategoryId  sql.NullInt64  `sqlofi:"REFERENCES category (Id)"`
	Created     sql.NullString `sqlofi:"DEFAULT (datetime('now'))"`
	Updated     sql.NullString `sqlofi:""`
}

// Category model
type Category struct {
	Id   int64  `sqlofi:"PRIMARY KEY AUTOINCREMENT"`
	Name string `sqlofi:"NOT NULL UNIQUE"`
}

// Order model
type Order struct {
	Id       int64          `sqlofi:"PRIMARY KEY AUTOINCREMENT"`
	UserId   int64          `sqlofi:"NOT NULL REFERENCES user (Id)"`
	Status   string         `sqlofi:"NOT NULL DEFAULT 'pending'"`
	Total    float64        `sqlofi:"NOT NULL DEFAULT 0"`
	Notes    sql.NullString `sqlofi:""`
	Created  sql.NullString `sqlofi:"DEFAULT (datetime('now'))"`
	Shipping sql.NullString `sqlofi:"DEFAULT (datetime('now', '+3 days'))"`
}

// OrderItem model
type OrderItem struct {
	Id        int64   `sqlofi:"PRIMARY KEY AUTOINCREMENT"`
	OrderId   int64   `sqlofi:"NOT NULL REFERENCES order (Id)"`
	ProductId int64   `sqlofi:"NOT NULL REFERENCES product (Id)"`
	Quantity  int     `sqlofi:"NOT NULL DEFAULT 1 CHECK(Quantity > 0)"`
	Price     float64 `sqlofi:"NOT NULL CHECK(Price >= 0)"`
	Subtotal  float64 `sqlofi:"NOT NULL CHECK(Subtotal >= 0)"`
}

// User model for additional examples
type User struct {
	Id       int64          `sqlofi:"PRIMARY KEY AUTOINCREMENT"`
	Username string         `sqlofi:"NOT NULL UNIQUE"`
	Email    string         `sqlofi:"NOT NULL UNIQUE"`
	Password string         `sqlofi:"NOT NULL"`
	Age      sql.NullInt64  `sqlofi:"CHECK(Age IS NULL OR Age >= 18)"`
	Active   sql.NullBool   `sqlofi:"DEFAULT 1"`
	Created  sql.NullString `sqlofi:"DEFAULT (datetime('now'))"`
}

// Create table functions
func CreateProductTable() *sqlite.Table {
	return sqlite.CREATE_TABLE(Product{}, Category{}).
		IfNotExists().
		ForeignKey(
			"fk_product_category",
			sqlite.FOREIGN_KEY(&Category{}, "Id").
				ForeighColumns("CategoryId").
				OnDelete(sqlite.SET_NULL).
				OnUpdate(sqlite.CASCADE),
		).
		Check(
			"check_product_name_length",
			sqlite.NewExpression("length(Name) >= 3"),
		)
}

func CreateCategoryTable() *sqlite.Table {
	return sqlite.CREATE_TABLE(Category{}).
		IfNotExists().
		Check(
			"check_category_name_length",
			sqlite.NewExpression("length(Name) > 0"),
		)
}
func CreateUserTable() *sqlite.Table {
	return sqlite.CREATE_TABLE(User{}).
		IfNotExists().
		Check(
			"check_username_length",
			sqlite.NewExpression("length(Username) >= 3"),
		).
		Check(
			"check_email_format",
			sqlite.NewExpression("Email LIKE '%@%.%'"),
		)
}
func CreateOrderTable() *sqlite.Table {
	return sqlite.CREATE_TABLE(Order{}, User{}).
		IfNotExists().
		ForeignKey(
			"fk_order_user",
			sqlite.FOREIGN_KEY(&User{}, "Id").
				ForeighColumns("UserId").
				OnDelete(sqlite.CASCADE).
				OnUpdate(sqlite.CASCADE),
		).
		Check(
			"check_order_status",
			sqlite.NewExpression("Status IN ('pending', 'processing', 'shipped', 'delivered', 'cancelled')"),
		)
}

func CreateOrderItemTable() *sqlite.Table {
	return sqlite.CREATE_TABLE(OrderItem{}, Order{}, Product{}).
		IfNotExists().
		ForeignKey(
			"fk_orderitem_order",
			sqlite.FOREIGN_KEY(&Order{}, "Id").
				ForeighColumns("OrderId").
				OnDelete(sqlite.CASCADE).
				OnUpdate(sqlite.CASCADE),
		).
		ForeignKey(
			"fk_orderitem_product",
			sqlite.FOREIGN_KEY(&Product{}, "Id").
				ForeighColumns("ProductId").
				OnDelete(sqlite.RESTRICT).
				OnUpdate(sqlite.CASCADE),
		)
}

func main() {
	// Create schema
	schema := sqlite.NewSchema("shop.db").
		Pragma(
			sqlite.ForeignKeys().ValueType("ON"),
			sqlite.JournalModeWAL(""),
			sqlite.SynchronousFull(""),
		).
		Table(
			CreateUserTable(),
			CreateCategoryTable(),
			CreateProductTable(),
			CreateOrderTable(),
			CreateOrderItemTable(),
		).
		Index(
			sqlite.CREATE_INDEX(&Product{}, "idx_product_price").
				IfNotExists(),
			sqlite.CREATE_INDEX(&Product{}, "idx_product_category").
				IfNotExists().
				Where(sqlite.NewExpression("CategoryId IS NOT NULL")),
			sqlite.CREATE_INDEX(&Order{}, "idx_order_user").
				IfNotExists(),
			sqlite.CREATE_INDEX(&Order{}, "idx_order_status").
				IfNotExists(),
			sqlite.CREATE_INDEX(&OrderItem{}, "idx_orderitem_order").
				IfNotExists(),
			sqlite.CREATE_INDEX(&OrderItem{}, "idx_orderitem_product").
				IfNotExists(),
		)

	// Build and execute schema
	schemaSQL := schema.Build()
	fmt.Println("Generated Shop Schema SQL:")
	fmt.Println(schemaSQL)

	db, err := sql.Open("sqlite3", "shop.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(schemaSQL)
	if err != nil {
		log.Fatalf("Failed to execute schema: %v", err)
	}
	fmt.Println("Shop schema executed successfully")

	// Insert sample data using transactions
	insertSampleData(db)

	// Demonstrate advanced queries
	performAdvancedQueries(db)
}

func insertSampleData(db *sql.DB) {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	// Function to handle rollback in case of error
	rollbackOnError := func(err error) {
		if err != nil {
			tx.Rollback()
			log.Fatalf("Transaction failed: %v", err)
		}
	}

	// 1. Insert a user
	userResult, err := tx.Exec(`
		INSERT INTO user (Username, Email, Password, Age, Active)
		VALUES ('shopuser', 'shop@example.com', 'hashed_shop_password', 30, 1)
	`)
	rollbackOnError(err)

	userId, err := userResult.LastInsertId()
	rollbackOnError(err)

	// 2. Insert categories
	categoryStmt, err := tx.Prepare(`
		INSERT INTO category (Name) VALUES (?)
	`)
	rollbackOnError(err)
	defer categoryStmt.Close()

	categories := []string{"Electronics", "Books", "Clothing", "Home & Garden"}
	categoryIds := make([]int64, len(categories))

	for i, category := range categories {
		result, err := categoryStmt.Exec(category)
		rollbackOnError(err)

		categoryIds[i], err = result.LastInsertId()
		rollbackOnError(err)
	}

	// 3. Insert products
	productStmt, err := tx.Prepare(`
		INSERT INTO product (Name, Description, Price, Quantity, CategoryId, Created, Updated)
		VALUES (?, ?, ?, ?, ?, datetime('now'), datetime('now'))
	`)
	rollbackOnError(err)
	defer productStmt.Close()

	products := []struct {
		name        string
		description string
		price       float64
		quantity    int
		categoryIdx int
	}{
		{"Smartphone", "Latest model smartphone", 599.99, 10, 0},
		{"Laptop", "Powerful laptop for professionals", 1299.99, 5, 0},
		{"Novel", "Bestselling fiction novel", 14.99, 20, 1},
		{"T-Shirt", "Cotton t-shirt", 19.99, 50, 2},
		{"Coffee Maker", "Automatic coffee maker", 89.99, 8, 3},
	}

	productIds := make([]int64, len(products))

	for i, product := range products {
		result, err := productStmt.Exec(
			product.name,
			product.description,
			product.price,
			product.quantity,
			categoryIds[product.categoryIdx],
		)
		rollbackOnError(err)

		productIds[i], err = result.LastInsertId()
		rollbackOnError(err)
	}

	// 4. Create an order
	orderResult, err := tx.Exec(`
		INSERT INTO "order" (UserId, Status, Total, Notes, Created)
		VALUES (?, 'pending', 0, 'Test order', datetime('now'))
	`, userId)
	rollbackOnError(err)

	orderId, err := orderResult.LastInsertId()
	rollbackOnError(err)

	// 5. Add items to the order
	orderItemStmt, err := tx.Prepare(`
		INSERT INTO orderitem (OrderId, ProductId, Quantity, Price, Subtotal)
		VALUES (?, ?, ?, ?, ?)
	`)
	rollbackOnError(err)
	defer orderItemStmt.Close()

	// Add two products to the order
	orderItems := []struct {
		productIdx int
		quantity   int
	}{
		{0, 1}, // One smartphone
		{1, 2}, // Two laptops
	}

	var orderTotal float64

	for _, item := range orderItems {
		product := products[item.productIdx]
		price := product.price
		subtotal := price * float64(item.quantity)
		orderTotal += subtotal

		_, err := orderItemStmt.Exec(
			orderId,
			productIds[item.productIdx],
			item.quantity,
			price,
			subtotal,
		)
		rollbackOnError(err)
	}

	// Update order total
	_, err = tx.Exec(`
		UPDATE "order" SET Total = ? WHERE Id = ?
	`, orderTotal, orderId)
	rollbackOnError(err)

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	fmt.Println("Sample shop data inserted successfully")
}

func performAdvancedQueries(db *sql.DB) {
	fmt.Println("\nPerforming advanced queries:")

	// 1. Get order details with all items
	fmt.Println("\n1. Order details with items:")
	rows, err := db.Query(`
		SELECT
			o.Id as OrderId,
			o.Status,
			o.Total,
			u.Username,
			p.Name as ProductName,
			oi.Quantity,
			oi.Price,
			oi.Subtotal
		FROM "order" o
		JOIN user u ON o.UserId = u.Id
		JOIN orderitem oi ON oi.OrderId = o.Id
		JOIN product p ON oi.ProductId = p.Id
		ORDER BY o.Id, p.Name
	`)
	if err != nil {
		fmt.Printf("Error querying order details: %v\n", err)
		return
	}
	defer rows.Close()

	var currentOrderId int64
	var totalAmount float64

	for rows.Next() {
		var orderId int64
		var status, username, productName string
		var total, price, subtotal float64
		var quantity int

		if err := rows.Scan(&orderId, &status, &total, &username, &productName, &quantity, &price, &subtotal); err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			continue
		}

		if currentOrderId != orderId {
			if currentOrderId != 0 {
				fmt.Printf("Total: $%.2f\n\n", totalAmount)
			}
			currentOrderId = orderId
			totalAmount = total
			fmt.Printf("Order #%d (%s) for %s\n", orderId, status, username)
			fmt.Println("--------------------------------------------------")
		}

		fmt.Printf("  - %s: %d x $%.2f = $%.2f\n", productName, quantity, price, subtotal)
	}
	if currentOrderId != 0 {
		fmt.Printf("Total: $%.2f\n\n", totalAmount)
	}

	// 2. Products by category with count
	fmt.Println("\n2. Products by category:")
	rows, err = db.Query(`
		SELECT
			c.Name as CategoryName,
			COUNT(p.Id) as ProductCount,
			SUM(p.Quantity) as TotalStock,
			AVG(p.Price) as AveragePrice
		FROM category c
		LEFT JOIN product p ON p.CategoryId = c.Id
		GROUP BY c.Id
		ORDER BY c.Name
	`)
	if err != nil {
		fmt.Printf("Error querying products by category: %v\n", err)
		return
	}
	defer rows.Close()

	fmt.Printf("%-20s %-15s %-15s %-15s\n", "Category", "Product Count", "Total Stock", "Avg Price")
	fmt.Println("------------------------------------------------------------------")

	for rows.Next() {
		var categoryName string
		var productCount, totalStock int
		var averagePrice float64

		if err := rows.Scan(&categoryName, &productCount, &totalStock, &averagePrice); err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			continue
		}

		fmt.Printf("%-20s %-15d %-15d $%-14.2f\n", categoryName, productCount, totalStock, averagePrice)
	}
	fmt.Println()

	// 3. Perform a stock update in a transaction
	fmt.Println("\n3. Updating product stock with transaction:")

	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("Error starting transaction: %v\n", err)
		return
	}

	// Get a product to update
	var productId int64
	var productName string
	var currentQuantity int

	err = tx.QueryRow(`
		SELECT Id, Name, Quantity FROM product LIMIT 1
	`).Scan(&productId, &productName, &currentQuantity)

	if err != nil {
		tx.Rollback()
		fmt.Printf("Error getting product: %v\n", err)
		return
	}

	// Reduce the quantity (simulate a purchase)
	quantityToReduce := 2
	newQuantity := currentQuantity - quantityToReduce

	fmt.Printf("Updating %s: %d -> %d units\n", productName, currentQuantity, newQuantity)

	_, err = tx.Exec(`
		UPDATE product
		SET Quantity = ?,
		    Updated = datetime('now')
		WHERE Id = ?
	`, newQuantity, productId)

	if err != nil {
		tx.Rollback()
		fmt.Printf("Error updating product: %v\n", err)
		return
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		fmt.Printf("Error committing transaction: %v\n", err)
		return
	}

	fmt.Printf("Stock update completed successfully\n\n")

	// 4. Demonstrate using a prepared statement with parameters
	fmt.Println("\n4. Products above certain price (using prepared statement):")

	minPrice := 50.0
	stmt, err := db.Prepare(`
		SELECT
			p.Name,
			p.Price,
			p.Quantity,
			c.Name as CategoryName
		FROM product p
		JOIN category c ON p.CategoryId = c.Id
		WHERE p.Price > ?
		ORDER BY p.Price DESC
	`)

	if err != nil {
		fmt.Printf("Error preparing statement: %v\n", err)
		return
	}
	defer stmt.Close()

	rows, err = stmt.Query(minPrice)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		return
	}
	defer rows.Close()

	fmt.Printf("Products with price > $%.2f:\n", minPrice)
	fmt.Printf("%-20s %-10s %-10s %-15s\n", "Product", "Price", "Quantity", "Category")
	fmt.Println("----------------------------------------------------------")

	for rows.Next() {
		var name, categoryName string
		var price float64
		var quantity int

		if err := rows.Scan(&name, &price, &quantity, &categoryName); err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			continue
		}

		fmt.Printf("%-20s $%-9.2f %-10d %-15s\n", name, price, quantity, categoryName)
	}
	fmt.Println()

	// 5. Demonstrate a CTE (Common Table Expression) query
	fmt.Println("\n5. Using CTE to find products with below average price:")

	rows, err = db.Query(`
		WITH AvgPrices AS (
			SELECT
				c.Id as CategoryId,
				c.Name as CategoryName,
				AVG(p.Price) as AvgPrice
			FROM category c
			JOIN product p ON p.CategoryId = c.Id
			GROUP BY c.Id
		)
		SELECT
			p.Name as ProductName,
			p.Price,
			ap.CategoryName,
			ap.AvgPrice,
			(p.Price - ap.AvgPrice) as PriceDifference
		FROM product p
		JOIN AvgPrices ap ON p.CategoryId = ap.CategoryId
		WHERE p.Price < ap.AvgPrice
		ORDER BY PriceDifference
	`)

	if err != nil {
		fmt.Printf("Error executing CTE query: %v\n", err)
		return
	}
	defer rows.Close()

	fmt.Printf("%-20s %-10s %-15s %-15s %-15s\n", "Product", "Price", "Category", "Avg Price", "Difference")
	fmt.Println("--------------------------------------------------------------------------------")

	for rows.Next() {
		var productName, categoryName string
		var price, avgPrice, priceDifference float64

		if err := rows.Scan(&productName, &price, &categoryName, &avgPrice, &priceDifference); err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			continue
		}

		fmt.Printf("%-20s $%-9.2f %-15s $%-14.2f $%-14.2f\n",
			productName, price, categoryName, avgPrice, priceDifference)
	}
	fmt.Println()

	// 6. Demonstrate using SQLite's window functions
	fmt.Println("\n6. Using window functions for price ranking within category:")

	rows, err = db.Query(`
		SELECT
			p.Name as ProductName,
			c.Name as CategoryName,
			p.Price,
			RANK() OVER (PARTITION BY p.CategoryId ORDER BY p.Price DESC) as PriceRank,
			AVG(p.Price) OVER (PARTITION BY p.CategoryId) as CategoryAvg,
			p.Price - AVG(p.Price) OVER (PARTITION BY p.CategoryId) as PriceDiff
		FROM product p
		JOIN category c ON p.CategoryId = c.Id
		ORDER BY c.Name, PriceRank
	`)

	if err != nil {
		fmt.Printf("Error executing window function query: %v\n", err)
		return
	}
	defer rows.Close()

	fmt.Printf("%-20s %-15s %-10s %-10s %-15s %-15s\n",
		"Product", "Category", "Price", "Rank", "Category Avg", "Difference")
	fmt.Println("----------------------------------------------------------------------------------")

	var currentCategory string

	for rows.Next() {
		var productName, categoryName string
		var price, categoryAvg, priceDiff float64
		var priceRank int

		if err := rows.Scan(&productName, &categoryName, &price, &priceRank, &categoryAvg, &priceDiff); err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			continue
		}

		if currentCategory != categoryName {
			if currentCategory != "" {
				fmt.Println()
			}
			currentCategory = categoryName
		}

		fmt.Printf("%-20s %-15s $%-9.2f %-10d $%-14.2f $%-14.2f\n",
			productName, categoryName, price, priceRank, categoryAvg, priceDiff)
	}
}
