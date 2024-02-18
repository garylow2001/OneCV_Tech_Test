package tests

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setUpRouters(apiEndpoint string, handler func(db *gorm.DB) gin.HandlerFunc) (*gin.Engine, sqlmock.Sqlmock) {
	mockDb, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	r := gin.Default()
	r.POST(apiEndpoint, func(c *gin.Context) { handler(db)(c) })
	r.GET(apiEndpoint, func(c *gin.Context) { handler(db)(c) })

	return r, mock
}
