package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/NgeKaworu/stock/src/app"
	"github.com/NgeKaworu/stock/src/creator"

	"github.com/NgeKaworu/util/middleware"
	"github.com/NgeKaworu/util/service"

	"github.com/julienschmidt/httprouter"
	"github.com/robfig/cron/v3"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		addr   = flag.String("l", ":8060", "绑定Host地址")
		dbinit = flag.Bool("i", false, "init database flag")
		mongo  = flag.String("m", "mongodb://localhost:27017", "mongod addr flag")
		mdb    = flag.String("db", "stock", "database name")
		// ucHost = flag.String("uc", "http://localhost:8020", "user center host")
		ucHost = flag.String("uc", "https://api.furan.xyz/user-center", "user center host")
		r      = flag.String("r", "localhost:6379", "rdb addr")
	)
	flag.Parse()

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	srv := service.New(ucHost, r)

	mongoInit := creator.Init
	if *dbinit {
		mongoInit = creator.WithoutInit
	}
	err := srv.Mongo.Open(*mongo, *mdb, mongoInit)

	if err != nil {
		panic(err)
	}
	app := app.New(srv)

	if err != nil {
		log.Println(err.Error())
	}

	c := cron.New(cron.WithParser(cron.NewParser(cron.Hour | cron.Dow)))
	c.AddFunc("19 MON-FRI", func() { app.StockCrawlManyService() })

	go c.Start()
	router := httprouter.New()

	// 爬+计算所有年报
	router.GET("/stockCrawlMany", srv.CheckPerm("admin")(app.StockCrawlMany))
	router.GET("/stock-list", app.StockList)

	router.GET("/exchange/:code", app.ExchangeList)
	router.POST("/exchange", app.ExchangeUpsert)
	router.PATCH("/exchange/:id", app.ExchangeUpsert)
	router.DELETE("/exchange/:id", app.ExchangeDelete)

	router.GET("/position", app.Position)
	router.GET("/position/:code", app.Position)
	router.PATCH("/position/:code", app.PositionUpsert)

	hSrv := &http.Server{Handler: middleware.CORS(srv.IsLogin(router)), ErrorLog: nil}
	hSrv.Addr = *addr

	go func() {
		if err := hSrv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("server on http port", hSrv.Addr)

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	cleanup := make(chan bool)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range signalChan {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			go func() {
				_ = hSrv.Shutdown(ctx)
				cleanup <- true
			}()
			<-cleanup
			srv.Destroy()
			c.Stop()
			fmt.Println("safe exit")
			cleanupDone <- true
		}
	}()
	<-cleanupDone

}
