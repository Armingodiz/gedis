package app

import (
	"fmt"

	"gedis/services/ttlManager"
	"gedis/store"
	"gedis/utils"

	"gedis/config"

	kvcontroller "gedis/controllers/kvController"
	usercontroller "gedis/controllers/userController"
	"gedis/middlewares"

	"gedis/db"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.ForceConsoleColor()
}

type App struct {
	route *gin.Engine
}

func NewApp() *App {
	// Initialize databae
	db, err := initializeDB()
	utils.FailOnError(err, "Database initialization failed, exiting the app with error!")
	r := routing(db)
	return &App{
		route: r,
	}
}

func (a *App) Start(restPort string) error {
	return a.route.Run(restPort)
}

func routing(db *db.DB) *gin.Engine {
	r := gin.Default()
	postgresStore := store.NewStore(db)
	UserController := usercontroller.UserController{Store: postgresStore}
	ttlManager, err := ttlManager.GetManager()
	if err != nil {
		utils.FailOnError(err, "ttlManager initialization failed, exiting the app with error!")
	}
	KvController := kvcontroller.KvController{Store: postgresStore, TtlManager: ttlManager}
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong!",
		})
	})
	r.POST("/users/signup", UserController.Signup())
	r.POST("/users/login", UserController.Login())

	//Protected routes
	r.Use(middlewares.JwtAuthorizationMiddleware())
	r.POST("/kvs", KvController.CreateKv())
	r.GET("/kvs", KvController.GetKvs())
	r.GET("/kvs/:key", KvController.GetKv())
	r.POST("/kvs/ttl", KvController.SetTtl())
	return r
}

func initializeDB() (*db.DB, error) {
	host := config.Configs.Database.Host
	port := config.Configs.Database.Port
	user := config.Configs.Database.User
	password := config.Configs.Database.Password
	dbName := config.Configs.Database.DbName
	extras := config.Configs.Database.Extras
	driver := config.Configs.Database.Driver

	connStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s %s",
		host, port, user, password, dbName, extras)
	db, err := db.Connect(driver, connStr)
	return db, err
}
