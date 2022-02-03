package middlewares

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type transactionKey struct{}

var TransactionKey = transactionKey{}

func TransactionMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		log.Print("beginning database transaction")
		transaction := db.Begin()
		if transaction.Error != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		SetDBToContext(c, transaction)

		defer func() {
			if r := recover(); r != nil {
				c.Writer.WriteHeader(http.StatusInternalServerError)
				transaction.Rollback()
				log.Print("transaction has been rolled back")
				panic(r)
			}
		}()

		c.Next()

		err := c.Errors.Last()
		if err != nil {
			//c.Writer.WriteHeader(http.StatusInternalServerError)
			transaction.Rollback()
			log.Print("transaction has been rolled back")
			return
		}
		log.Print("committing transactions")
		trErr := transaction.Commit().Error

		if trErr != nil && trErr != sql.ErrTxDone {
			log.Print("transaction commit error: ", trErr)
			transaction.Rollback()
			log.Print("transaction has been rolled back")
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Print("transaction has been committed")
	}
}

func SetDBToContext(c *gin.Context, transaction *gorm.DB) {
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), TransactionKey, transaction))
}

func GetDBFromContext(c *gin.Context) *gorm.DB {
	db := c.Request.Context().Value(TransactionKey).(*gorm.DB)
	return db
}
