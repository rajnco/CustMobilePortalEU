package handler

import (
        "net/http"

        "cust-mobile-eu/database"
        "cust-mobile-eu/database/model"

        "gorm.io/gorm"
	"gorm.io/gorm/clause"
        log "github.com/sirupsen/logrus"
)

// ProductsTableName - table name for Product model in database
const ProductsTableName = "products"

// Product - data struct to access Products
type Product struct {
        db *gorm.DB
}

// Build - connect/re-use db connection
func (p *Product) Build(db *gorm.DB) {
        if db == nil {
                p.db = database.GetDB()
        } else {
                p.db = db
        }
}

// TotalCount - returns total count from table
func (p *Product) TotalCount() int64 {
        var count int64 
        if err := p.db.Table(ProductsTableName).Count(&count).Error; err != nil {
                return 0
        }
        return count
}


// GetProduct - get particular product from products list
func (p *Product) GetProduct(id string) (response model.Product, status int) {
        //if err := p.db.Model(&model.Product{}).Where("id = ?", id).Find(&response).Error; err != nil {
        if err := p.db.Where("id = ?", id).Find(&response).Error; err != nil {
                log.Printf("error retrive product from DB : %v", err)
                status = http.StatusInternalServerError
                return
        }
        status = http.StatusOK
        return
}


// GetProducts - list all the products 
func (p *Product) GetProducts() (response []model.Product, status int) {
        if err := p.db.Find(&response).Error; err != nil {
                log.Printf("error getting all products from DB : %v", err)
                status = http.StatusInternalServerError
                return
        }
        status = http.StatusOK
        return
}

// Clean - clean the table - delete all the records
func (p *Product) Clean() error {
        result := p.db.Exec("DELETE FROM " + ProductsTableName)

        if result.Error != nil {
                return result.Error
        }
        return nil
}


// UpsertProduct - create new product in database table products
func(p *Product) UpsertProduct(cproduct *model.Product) (status int) {  
        tx := p.db.Begin()
        
        if err := p.db.Model(&model.Product{}).Clauses(clause.OnConflict{UpdateAll: true, }).Create(&cproduct).Error; err != nil {
                log.Printf(" error upsert into products table  : %v ", err)
                status = http.StatusInternalServerError
                tx.Rollback()
                return
        }
        tx.Commit()
        status = http.StatusCreated
        return
}

