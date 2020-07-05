package router

import (
	"log"
	"net/http"
	"strconv"
	"time"

	database "github.com/Sterks/Pp.Common.Db/db"
	config2 "github.com/Sterks/fReader/config"
	"github.com/gin-gonic/gin"
)

//WebServer ...
type WebServer struct {
	config *config2.Config
	rr     *gin.Engine
	db     *database.Database
}

type InfoForm struct {
	TsName      string `form:"ts_name"`
	TsDataStart string `form:"ts_data_start"`
	TsRunTimes  string `form:"ts_run_times"`
	TsComment   string `form:"ts_comment"`
}

//DateCheck - Даты для передачи с формы
type DateCheck struct {
	From time.Time `form:"from"`
	To   time.Time `form:"to"`
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
			"lastID":  lst,
		})
	})
	r.GET("/GetID", web.GetLastID)
	r.GET("/tasks", web.tasksGET)
	r.POST("/tasks", web.tasks)
	r.GET("/Info", web.Info)
	r.POST("/Info", web.Info)
	s := &http.Server{
		Addr:    ":8000",
		Handler: r,
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
		"LastID":       lastID,
		"Notification": notif,
		"Protocols":    proto,
	})
}

// Info метод для отображения сводной информации
func (web *WebServer) Info(c *gin.Context) {
	str := "2020-07-04"
	from, _ := time.Parse(time.RFC3339, str)
	to := time.Now()
	notification44zip := web.db.CountDocument(from, to, 1, 1)
	notification223zip := web.db.CountDocument(from, to, 3, 1)
	protocol44zip := web.db.CountDocument(from, to, 2, 1)
	protocol223zip := web.db.CountDocument(from, to, 4, 1)

	notification44 := web.db.CountDocument(from, to, 1, 2)
	notification223 := web.db.CountDocument(from, to, 3, 2)
	protocol44 := web.db.CountDocument(from, to, 2, 2)
	protocol223 := web.db.CountDocument(from, to, 4, 2)
	c.JSON(http.StatusOK, gin.H{
		"Notification44zip":  notification44zip,
		"protocol44zip":      protocol44zip,
		"Notification223zip": notification223zip,
		"protocol223zip":     protocol223zip,
		"Notification44":     notification44,
		"protocol44":         protocol44,
		"Notification223":    notification223,
		"protocol223":        protocol223,
	})
}

func (web *WebServer) tasks(c *gin.Context) {
	var infof InfoForm
	_ = c.ShouldBind(&infof)
	layout := "2006-01-02"
	tds := infof.TsDataStart
	tsDataStart, err := time.Parse(layout, tds)
	if err != nil {
		log.Printf("Не могу превратить строку в дату - %v", err)
	}
	tsrt := infof.TsRunTimes
	tsRunTimes, err2 := strconv.Atoi(tsrt)
	if err2 != nil {
		log.Printf("Не могу превратить строку в целое число - %v", err)
	}
	web.db.CreateTask(infof.TsName, tsDataStart, tsRunTimes, infof.TsComment)
	c.HTML(200, "form.html", nil)
}

func (web *WebServer) tasksGET(c *gin.Context) {
	c.HTML(200, "form.html", nil)
}
