package main

import (
	"fmt"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/test/init_project/config"
)

// Account account model
type Account struct {
	ID      uint    `gorm:"primaryKey"`
	Balance float64 `gorm:"type:decimal(10,2)"`
}

// Transaction transaction model
type Transaction struct {
	ID            uint `gorm:"primaryKey"`
	FromAccountID uint
	ToAccountID   uint
	Amount        float64 `gorm:"type:decimal(10,2)"`
}

func transferMoney(db *gorm.DB, fromAccountID, toAccountID uint, amount float64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// query the balance of the from account (pessimistic lock)
		var fromAccount Account
		if err := tx.Model(&Account{}).Where("id = ?", fromAccountID).
			Clauses(clause.Locking{Strength: "UPDATE"}).Find(&fromAccount).Error; err != nil {
			return err
		}

		// check the balance
		if fromAccount.Balance < amount {
			return fmt.Errorf("balance not enough, current balance: %.2f", fromAccount.Balance)
		}

		// deduct the balance of the from account
		if err := tx.Model(&Account{}).
			Where("id = ? AND balance >= ?", fromAccountID, amount).
			Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}

		// increase the balance of the to account
		if err := tx.Model(&Account{}).
			Where("id = ?", toAccountID).
			Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}

		// record transaction
		transaction := Transaction{
			FromAccountID: fromAccountID,
			ToAccountID:   toAccountID,
			Amount:        amount,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		return nil
	})
}

// init test data (optional)
func initTestData(db *gorm.DB) {
	// check if there is data
	var count int64
	db.Model(&Account{}).Count(&count)
	if count > 0 {
		return
	}

	// create test accounts
	accounts := []Account{
		{ID: 1, Balance: 1000}, // account A
		{ID: 2, Balance: 500},  // account B
	}
	db.Create(&accounts)
	log.Println("init test accounts success")
}

func main() {
	// 使用全局配置获取数据库连接
	db := config.GetDB()

	// 自动迁移表结构
	if err := db.AutoMigrate(&Account{}, &Transaction{}); err != nil {
		log.Fatal("table structure migration failed:", err)
	}

	// init test data (optional)
	initTestData(db)

	// transfer money
	err := transferMoney(db, 1, 2, 100.0)
	if err != nil {
		log.Fatal("transfer money failed:", err)
	}
	log.Println("transfer money success")
}
