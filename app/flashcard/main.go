/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2023-03-19 03:04:11
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-03-19 17:35:26
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

	mid "github.com/NgeKaworu/util/middleware"
	"github.com/NgeKaworu/util/tool"
	"github.com/julienschmidt/httprouter"
	"github.com/yingxv/flashcard-go/src/controller"
	"github.com/yingxv/flashcard-go/src/db"
	"github.com/yingxv/flashcard-go/src/middleware"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		addr   = flag.String("l", ":8041", "绑定Host地址")
		dbinit = flag.Bool("i", false, "init database flag")
		mongo  = flag.String("m", "mongodb://localhost:27017", "mongod addr flag")
		mdb    = flag.String("db", "to-do-list", "database name")
		ucHost = flag.String("uc", "http://user-center-go-dev", "user center host")
	)
	flag.Parse()

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	auth := middleware.NewAuth(ucHost)
	mongoClient := db.NewMongoClient()
	err := mongoClient.Open(*mongo, *mdb, *dbinit)
	validate := tool.NewValidator()
	trans := tool.NewValidatorTranslator(validate)

	controller := controller.NewController(validate, trans, auth, mongoClient)
	if err != nil {
		panic(err)
	}

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

	srv := &http.Server{Handler: auth.IsLogin(mid.CORS(router)), ErrorLog: nil}
	srv.Addr = *addr

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("server on http port", srv.Addr)

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	cleanup := make(chan bool)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range signalChan {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			go func() {
				_ = srv.Shutdown(ctx)
				cleanup <- true
			}()
			<-cleanup
			mongoClient.Close()
			fmt.Println("safe exit")
			cleanupDone <- true
		}
	}()
	<-cleanupDone

}
