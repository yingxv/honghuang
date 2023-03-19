/*
 * @Author: fuRan NgeKaworu@gmail.com
 * @Date: 2023-03-19 02:57:10
 * @LastEditors: fuRan NgeKaworu@gmail.com
 * @LastEditTime: 2023-03-19 23:16:17
 * @FilePath: /honghuang/app/todolist/main.go
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

	"github.com/NgeKaworu/to-do-list-go/src/app"
	"github.com/NgeKaworu/to-do-list-go/src/creator"
	"github.com/NgeKaworu/util/middleware"
	"github.com/NgeKaworu/util/service"
	"github.com/julienschmidt/httprouter"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		addr   = flag.String("l", ":8040", "绑定Host地址")
		dbinit = flag.Bool("i", false, "init database flag")
		mongo  = flag.String("m", "mongodb://localhost:27017", "mongod addr flag")
		mdb    = flag.String("db", "to-do-list", "database name")
		ucHost = flag.String("uc", "http://user-center-go", "user center host")
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

	router := httprouter.New()
	//task ctrl
	router.POST("/v1/task/create", app.AddTask)
	router.PATCH("/v1/task/update", app.SetTask)
	router.GET("/v1/task/list", app.ListTask)
	router.DELETE("/v1/task/:id", app.RemoveTask)

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
