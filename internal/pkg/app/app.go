package app

import (
	"log"
	"net/http"
	"time"
	"fmt"

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
	r.GET("/api/tenders/:tenderId/status", app.GetTenderStatus)
	r.PUT("/api/tenders/:tenderId/status", app.ChangeTenderStatus)
	r.PATCH("/api/tenders/:tenderId/edit", app.ChangeTender)
	r.PUT("/api/tenders/:tenderId/rollback/:version", app.ChangeTenderVersion)
	//Предложения
	r.POST("/api/bids/new", app.AddBid)
	r.GET("/api/bids/my", app.GetBid)
	r.GET("/api/bids/:bidId/list", app.GetTenderBids)
	r.GET("/api/bids/:bidId/status", app.GetBidStatus)
	r.PUT("/api/bids/:bidId/status", app.ChangeBidStatus)
	r.PATCH("/api/bids/:bidId/edit", app.ChangeBid)
	r.PUT("/api/bids/:bidId/submit_decision", app.SubmitBid)
	r.PUT("/api/bids/:bidId/rollback/:version", app.ChangeBidVersion)
	//Отзывы
	r.PUT("/api/bids/:bidId/feedback",app.AddBidFeedback)
	// r.GET("/api/bids/:bidId/reviews", app.GetReviews)

	r.Run(fmt.Sprintf("%s", app.config.ServerAddress)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
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
