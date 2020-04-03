package router

import (
	database "github.com/Sterks/Pp.Common.Db/db"
	config2 "github.com/Sterks/fReader/config"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

//WebServer ...
type WebServer struct {
	config *config2.Config
	rr *gin.Engine
	db *database.Database
}

type InfoForm struct {
	TsName string `json:"ts_name"`
	TsDataStart string `json:"ts_data_start"`
	TsRunTimes string `json:"ts_run_times"`
	TsComment string `json:"ts_comment"`
}

func NewWebServer(config *config2.Config, db *database.Database) *WebServer {
	return &WebServer{config: config, rr: &gin.Engine{}, db: db}
}

func (web *WebServer) Start() {
	r := gin.Default()

	r.Use(Cors())

	r.LoadHTMLGlob("views/*")
	r.GET("/ping", func(c *gin.Context) {
		lst := web.db.LastID()
		c.JSON(200, gin.H{
			"message": "pong",
			"lastID": lst,
		})
	})
	r.GET("/GetID", web.GetLastID )
	r.POST("/tasks", web.tasks)
	s := &http.Server{
		Addr:           ":8000",
		Handler:        r,
		//ReadTimeout:    10 * time.Second,
		//WriteTimeout:   10 * time.Second,
		//MaxHeaderBytes: 1 << 20,
	}
	_ = s.ListenAndServe()
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func (web *WebServer) GetLastID(c *gin.Context) {
	lastID := web.db.LastID()
	notif := web.db.QuantityTypeDoc("notifications")
	proto := web.db.QuantityTypeDoc("protocols")
	c.JSON(http.StatusOK, gin.H{
		"LastID": lastID,
		"Notification": notif,
		"Protocols": proto,
	})
}

func (web *WebServer) tasks(c *gin.Context) {
	var infof InfoForm
	_ = c.BindJSON(&infof)
	layout := "2006-01-02"
	tds := infof.TsDataStart
	tsDataStart, err := time.Parse(layout , tds)
	if err != nil {
		log.Printf("Не могу превратить строку в дату - %v", err)
	}
	tsrt := infof.TsRunTimes
	tsRunTimes, err2 := strconv.Atoi(tsrt)
	if err2 != nil {
		log.Printf("Не могу превратить строку в целое число - %v", err)
	}
	web.db.CreateTask(infof.TsName, tsDataStart, tsRunTimes, infof.TsComment)
	c.JSON(200, gin.H{
		"message": "Запись успешно добавлена.",
	})
}

