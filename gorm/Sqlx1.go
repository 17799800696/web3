package main

import (
    "database/sql"
    "fmt"
    "log"

	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

// Employee 结构体映射数据库表
type Employee struct {
    ID        int    `db:"id"`
    Name      string `db:"name"`
    Department string `db:"department"`
    Salary    float64 `db:"salary"`
}

func main() {
    // 连接数据库
    db, err := sqlx.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    // create employees table
    schema := `
    CREATE TABLE IF NOT EXISTS employees (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        department VARCHAR(100) NOT NULL,
        salary DOUBLE NOT NULL
    );`
    if _, err := db.Exec(schema); err != nil {
        log.Fatalf("create employees table failed: %v", err)
    }

    // insert data
    db.MustExec(`INSERT INTO employees (name, department, salary) VALUES (?, ?, ?)`, "John Doe", "技术部", 100000)
    db.MustExec(`INSERT INTO employees (name, department, salary) VALUES (?, ?, ?)`, "Jane Smith", "技术部", 90000)
    db.MustExec(`INSERT INTO employees (name, department, salary) VALUES (?, ?, ?)`, "Jim Beam", "市场部", 80000)
    db.MustExec(`INSERT INTO employees (name, department, salary) VALUES (?, ?, ?)`, "Jill Johnson", "市场部", 70000)

    // query tech department employees
    techEmployees, err := getTechDepartmentEmployees(db)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("技术部员工:")
    for _, emp := range techEmployees {
        fmt.Printf("ID: %d, 姓名: %s, 工资: %.2f\n", emp.ID, emp.Name, emp.Salary)
    }

    // 2. 查询工资最高的员工
    topEarner, err := getHighestSalaryEmployee(db)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("\n工资最高的员工: %s (%.2f)\n", topEarner.Name, topEarner.Salary)
}

// 查询技术部所有员工
func getTechDepartmentEmployees(db *sqlx.DB) ([]Employee, error) {
    var employees []Employee
    query := "SELECT id, name, department, salary FROM employees WHERE department = ?"
    err := db.Select(&employees, query, "技术部")
    return employees, err
}

// 查询工资最高的员工
func getHighestSalaryEmployee(db *sqlx.DB) (Employee, error) {
    var employee Employee
    query := "SELECT id, name, department, salary FROM employees ORDER BY salary DESC LIMIT 1"
    err := db.Get(&employee, query)
    if err == sql.ErrNoRows {
        return Employee{}, nil // 没有记录时返回空结构体
    }
    return employee, err
}