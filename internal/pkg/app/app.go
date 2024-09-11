package app

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"avito_project/internal/app/config"
	"avito_project/internal/app/dsn"
	"avito_project/internal/app/repository"
)

type Application struct {
	repo        *repository.Repository
	config      *config.Config
	// dsn string
}

func (app *Application) Run() {
	log.Println("Server start up")

	r := gin.Default()

	r.Use(ErrorHandler())

	//Пинг
	r.GET("/api/ping", app.Ping)    
	//Тендеры
	r.GET("/api/tenders", app.GetAllTenders)
	r.POST("/api/tenders/new", app.AddTender)
	r.GET("/api/tenders/my", app.GetTender)
	r.GET("/api/tenders/:tender_id/status", app.GetTenderStatus)
	r.PUT("/api/tenders/:tender_id/status", app.ChangeTenderStatus)
	r.PUT("/api/tenders/:tender_id/edit", app.ChangeTender)
	// r.PUT("/api/tenders/:tender_id/rollback/:version", app.RollbackTender)
	//Предложения
	r.POST("/api/bids/new", app.AddBid)
	r.GET("/api/bids/my", app.GetBid)
	// r.GET("/api/bids/:tender_id/list", app.GetBidTenders)
	// r.PUT("/api/bids/:bid_id/edit", app.ChangeBid)
	// r.PUT("/api/bids/:bid_id/rollback/:version", app.RollbackBid)
	//Отзывы
	// r.GET("/api/bids/:tender_id/reviews", app.GetBidsReviews)

	r.Run("localhost:80") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	log.Println("Server down")
}

func New() (*Application, error) {
	var err error
	loc, _ := time.LoadLocation("UTC")
	time.Local = loc
	app := Application{}
	app.config, err = config.NewConfig()
	if err != nil {
		return nil, err
	}
	app.repo, err = repository.New(dsn.FromEnv())
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &app, nil
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			log.Println(err.Err)
		}
		lastError := c.Errors.Last()
		if lastError != nil {
			switch c.Writer.Status() {
			case http.StatusBadRequest:
				c.JSON(-1, gin.H{"error": "wrong request"})
			case http.StatusNotFound:
				c.JSON(-1, gin.H{"error": lastError.Error()})
			case http.StatusMethodNotAllowed:
				c.JSON(-1, gin.H{"error": lastError.Error()})
			default:
				c.Status(-1)
			}
		}
	}
}
