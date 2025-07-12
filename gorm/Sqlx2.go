package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
)

// Book 结构体映射 books 表
type Book struct {
    ID     int     `db:"id"`
    Title  string  `db:"title"`
    Author string  `db:"author"`
    Price  float64 `db:"price"`
}

func main() {
    // 连接数据库
    db, err := sqlx.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // 创建表（如果不存在）
    err = createTable(db)
    if err != nil {
        log.Fatal(err)
    }

    // 插入测试数据
    err = insertSampleData(db)
    if err != nil {
        log.Fatal(err)
    }

    // 查询价格大于50的书籍
    books, err := getBooksOverPrice(db, 50.0)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("价格大于50元的书籍:")
    for _, book := range books {
        fmt.Printf("ID: %d, 书名: %s, 作者: %s, 价格: %.2f\n",
            book.ID, book.Title, book.Author, book.Price)
    }
}

// 创建 books 表
func createTable(db *sqlx.DB) error {
    query := `
    CREATE TABLE IF NOT EXISTS books (
        id INT PRIMARY KEY AUTO_INCREMENT,
        title VARCHAR(255) NOT NULL,
        author VARCHAR(255) NOT NULL,
        price DECIMAL(10, 2) NOT NULL
    )
    `
    _, err := db.Exec(query)
    return err
}

// 插入测试数据
func insertSampleData(db *sqlx.DB) error {
    // 检查是否已有数据
    var count int
    err := db.Get(&count, "SELECT COUNT(*) FROM books")
    if err != nil && err != sql.ErrNoRows {
        return err
    }

    // 如果表为空，则插入测试数据
    if count == 0 {
        books := []Book{
            {Title: "Go语言实战", Author: "Alan A. Donovan", Price: 99.00},
            {Title: "Effective Go", Author: "Robert Griesemer", Price: 79.00},
            {Title: "SQL必知必会", Author: "Ben Forta", Price: 49.00},
            {Title: "Clean Code", Author: "Robert C. Martin", Price: 129.00},
            {Title: "The Art of Computer Programming", Author: "Donald E. Knuth", Price: 299.00},
        }

        // 批量插入
        for _, book := range books {
            _, err := db.NamedExec(`
                INSERT INTO books (title, author, price) 
                VALUES (:title, :author, :price)
            `, book)
            if err != nil {
                return err
            }
        }
        fmt.Println("插入了5条测试数据")
    }
    return nil
}

// 查询价格大于指定值的书籍
func getBooksOverPrice(db *sqlx.DB, price float64) ([]Book, error) {
    var books []Book
    query := `
        SELECT 
            id, 
            title, 
            author, 
            price 
        FROM 
            books 
        WHERE 
            price > ? 
        ORDER BY 
            price DESC
    `
    err := db.Select(&books, query, price)
    return books, err
}