/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2023-03-19 03:04:11
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-04-16 23:11:32
 * @FilePath: /honghuang/app/flashcard/main.go
 * @Description:
 *
 * Copyright (c) 2023 by ${git_name_email}, All Rights Reserved.
 */
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

	"github.com/NgeKaworu/util/middleware"

	"github.com/julienschmidt/httprouter"
	"github.com/yingxv/flashcard-go/src/controller"
	"github.com/yingxv/flashcard-go/src/creator"

	"github.com/NgeKaworu/util/service"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		addr   = flag.String("l", ":8030", "绑定Host地址")
		dbInit = flag.Bool("i", false, "init database flag")
		mongo  = flag.String("m", "mongodb://localhost:27017", "mongod addr flag")
		mdb    = flag.String("db", "flashcard", "database name")
		ucHost = flag.String("uc", "http://user-center-go-dev", "user center host")
		r      = flag.String("r", "localhost:6379", "rdb addr")
	)
	flag.Parse()

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	srv := service.New(&service.ServiceAug{
		UCHost:    ucHost,
		RedisAddr: r,
	})

	mongoInit := creator.Init
	if *dbInit {
		mongoInit = creator.WithoutInit
	}
	err := srv.Mongo.Open(*mongo, *mdb, mongoInit)

	if err != nil {
		panic(err)
	}
	controller := controller.NewController(srv)

	router := httprouter.New()
	//task ctrl
	router.POST("/record/create", controller.RecordCreate)
	router.DELETE("/record/remove/:id", controller.RecordRemove)
	router.PATCH("/record/update", controller.RecordUpdate)
	router.GET("/record/list", controller.RecordList)
	router.PATCH("/record/review", controller.RecordReview)
	router.GET("/record/review-all", controller.RecordReviewAll)
	router.PATCH("/record/random-review", controller.RecordRandomReview)
	router.PATCH("/record/set-review-result", controller.RecordSetReviewResult)

	hSrv := &http.Server{Handler: srv.IsLogin(middleware.CORS(router)), ErrorLog: nil}
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
			fmt.Println("safe exit")
			cleanupDone <- true
		}
	}()
	<-cleanupDone

}
