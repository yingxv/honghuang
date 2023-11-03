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

	"github.com/yingxv/honghuang/user-center/src/app"
	"github.com/yingxv/honghuang/user-center/src/creator"
	"github.com/yingxv/honghuang/util/middleware"
	"github.com/yingxv/honghuang/util/service"

	"github.com/julienschmidt/httprouter"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		addr   = flag.String("l", ":8020", "绑定Host地址")
		dbInit = flag.Bool("i", true, "init database flag")
		mongo  = flag.String("m", "mongodb://localhost:27017", "mongod addr flag")
		mdb    = flag.String("db", "uc", "database name")
		k      = flag.String("k", "f3fa39nui89Wi707", "iv key")
		r      = flag.String("r", "localhost:6379", "rdb addr")
		ePwd   = flag.String("d", "", "email pwd")
	)
	flag.Parse()

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	srv := service.New(&service.ServiceAug{
		RedisAddr: r,
		CipherKey: k,
		DialerPwd: ePwd,
	})

	mongoInit := creator.WithoutInit
	if *dbInit {
		mongoInit = creator.Init
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
	// user ctrl
	router.POST("/login", app.Login)
	router.POST("/register", app.Regsiter)
	router.POST("/forget-pwd", app.ForgetPwd)
	router.GET("/profile", app.JWT(app.Profile))

	// user mgt
	router.POST("/user/create", app.JWT(app.CheckPerm("admin")(app.CreateUser)))
	router.DELETE("/user/remove/:uid", app.JWT(app.CheckPerm("admin")(app.RemoveUser)))
	router.PUT("/user/update", app.JWT(app.CheckPerm("admin")(app.UpdateUser)))
	router.GET("/user/list", app.JWT(app.CheckPerm("admin")(app.UserList)))
	router.GET("/user/validate", app.UserValidateEmail)

	// role mgt
	router.POST("/role/create", app.JWT(app.CheckPerm("admin")(app.RoleCreate)))
	router.DELETE("/role/remove/:id", app.JWT(app.CheckPerm("admin")(app.RoleRemove)))
	router.PUT("/role/update", app.JWT(app.CheckPerm("admin")(app.RoleUpdate)))
	router.GET("/role/list", app.JWT(app.RoleList))
	router.GET("/role/validate", app.JWT(app.RoleValidateKey))

	// perm mgt
	router.POST("/perm/create", app.JWT(app.CheckPerm("admin")(app.PermCreate)))
	router.DELETE("/perm/remove/:id", app.JWT(app.CheckPerm("admin")(app.PermRemove)))
	router.PUT("/perm/update", app.JWT(app.CheckPerm("admin")(app.PermUpdate)))
	router.GET("/perm/list", app.JWT(app.PermList))
	router.GET("/perm/validate", app.JWT(app.PermValidateKey))
	router.GET("/menu", app.JWT(app.Menu))
	router.GET("/micro-app", app.MicroApp)

	// rpc
	router.HEAD("/check-perm-rpc/:perm", app.JWT(app.CheckPermRPC))
	router.HEAD("/rpc/is-login-rpc", app.JWT(app.IsLoginRPC))

	// captcha
	router.GET("/captcha/fetch", app.FetchCaptcha)
	router.GET("/captcha/check", app.CheckCaptcha)

	hSrv := &http.Server{Handler: middleware.CORS(router), ErrorLog: nil}
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
