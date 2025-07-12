package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/test/init_project/config"
)

type Student struct {
	ID     uint   `gorm:"primarykey"`
	Name   string `gorm:"size:255;not null"`
	Age    int    `gorm:"not null"`
	Grade  string `gorm:"size:10;not null"`
	Class  string `gorm:"size:50;not null"`
	Gender string `gorm:"size:20;not null"`
}

func main() {
	db := config.GetDB()

	// Auto migrate the schema
	if err := db.AutoMigrate(&Student{}); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	// Create a new student
	student := &Student{Name: "John", Age: 20, Grade: "A", Class: "101", Gender: "Male"}
	if err := db.Create(student).Error; err != nil {
		log.Printf("Failed to create student: %v", err)
		return
	}
	log.Printf("Created student with ID: %d", student.ID)

	// Update student age
	if err := db.Model(&Student{}).Where("name = ?", "John").Update("age", 21).Error; err != nil {
		log.Printf("Failed to update student: %v", err)
		return
	}
	log.Println("Updated student age to 21")

	// Query student
	var result Student
	if err := db.Where("name = ?", "John").First(&result).Error; err != nil {
		log.Printf("Failed to query student: %v", err)
		return
	}

	// Marshal to JSON with proper error handling
	jsonData, err := json.Marshal(&result)
	if err != nil {
		log.Printf("Failed to marshal JSON: %v", err)
		return
	}
	fmt.Printf("Student data: %s\n", string(jsonData))

	// Delete student
	if err := db.Where("name = ?", "John").Delete(&Student{}).Error; err != nil {
		log.Printf("Failed to delete student: %v", err)
		return
	}
	log.Println("Deleted student successfully")
}
