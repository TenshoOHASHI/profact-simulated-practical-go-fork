package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"github.com/yamu-studio/profact-simulated-practical-go/internal/handler/api"
	"github.com/yamu-studio/profact-simulated-practical-go/internal/handler/validator"
	"github.com/yamu-studio/profact-simulated-practical-go/internal/repository"
	"github.com/yamu-studio/profact-simulated-practical-go/internal/usecase"
)

func main() {
	// DB接続設定
	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = "user:password@tcp(127.0.0.1:3306)/profact_simulated_practical_db?parseTime=true"
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// DI: リポジトリ
	customerRepo := repository.NewCustomerRepository(db)
	propertyRepo := repository.NewPropertyRepository(db)
	dealRepo := repository.NewDealRepository(db)

	// DI: ユースケース
	customerUC := usecase.NewCustomerUsecase(customerRepo)
	propertyUC := usecase.NewPropertyUsecase(propertyRepo)
	dealUC := usecase.NewDealUsecase(dealRepo)

	// Validator: バリデーション
	v := validator.NewValidator()
	// DI: ハンドラ
	customerHandler := api.NewCustomerHandler(customerUC, v)
	propertyHandler := api.NewPropertyHandler(propertyUC)
	dealHandler := api.NewDealHandler(dealUC)

	r := gin.Default()

	// ヘルスチェックAPI
	r.GET("/connect", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Now System is running"})
	})

	// API ルーティング (Phase 1)
	apiRoutes := r.Group("/api")
	{
		customers := apiRoutes.Group("/customers")
		{
			customers.GET("", customerHandler.ListCustomers)
			customers.POST("", customerHandler.CreateCustomer)
			customers.GET("/:id", customerHandler.GetCustomer)
			customers.PUT("/:id", customerHandler.UpdateCustomer)
			customers.DELETE("/:id", customerHandler.DeleteCustomer)
		}

		properties := apiRoutes.Group("/properties")
		{
			properties.GET("", propertyHandler.ListProperties)
			properties.POST("", propertyHandler.CreateProperty)
			properties.GET("/:id", propertyHandler.GetProperty)
			properties.PUT("/:id", propertyHandler.UpdateProperty)
			properties.DELETE("/:id", propertyHandler.DeleteProperty)
		}

		deals := apiRoutes.Group("/deals")
		{
			deals.GET("", dealHandler.ListDeals)
			deals.POST("", dealHandler.CreateDeal)
			deals.GET("/:id", dealHandler.GetDeal)
			deals.PUT("/:id", dealHandler.UpdateDeal)
			deals.DELETE("/:id", dealHandler.DeleteDeal)
			deals.PATCH("/:id/status", dealHandler.UpdateDealStatus)
		}
	}

	r.Run(":8080") // port 8080 で起動
}
