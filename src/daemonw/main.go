package main

import (
	dlog "log"
	"net/http"
	"daemonw/api"
	"os"
	"os/signal"
	"time"
	"context"
	"strconv"
	"daemonw/conf"
	"path/filepath"
	log "daemonw/log"
	"golang.org/x/crypto/acme/autocert"
	"crypto/tls"
	"daemonw/db"
	"syscall"
	"daemonw/dao"
	"fmt"
)

func main(){
	err := conf.ParseConfig("")
	if err != nil {
		dlog.Fatal(err)
	}
	log.InitLog()
	t1:=time.Now()
	for i:=0;i<1000;i++{
		log.Info().Msg("test")
	}
	d:=time.Now().Sub(t1)
	fmt.Printf("cost = %d",d.Nanoseconds()/1000)
}

func tset() {
	err := conf.ParseConfig("")
	if err != nil {
		dlog.Fatal(err)
	}
	log.InitLog()
	err = db.InitRedis()
	if err!=nil{
		dlog.Fatal(err)
	}
	err = db.InitDB()
	if err != nil {
		dlog.Fatal(err)
	}
	dao.InitDaoManager()
	defer closeDBConn()

	cfg := conf.Config
	router := api.GetRouter()
	var tlsConf *tls.Config
	//let's encrypt auto cert
	if cfg.TLS && cfg.UseAutoCert {
		certManager := autocert.Manager{
			Cache:      autocert.DirCache(conf.BinDir + string(filepath.Separator) + ".certs"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(cfg.Domain),
		}
		tlsConf = &tls.Config{GetCertificate: certManager.GetCertificate}
	}
	//tls config
	srv := &http.Server{
		Addr:      ":" + strconv.Itoa(cfg.Port),
		Handler:   router,
		TLSConfig: tlsConf,
	}
	go func() {
		if !cfg.TLS {
			log.Info().Msgf("start http server on %d", cfg.Port)
			err = srv.ListenAndServe()
		} else {
			log.Info().Msgf("start https server on %d", cfg.Port)
			if cfg.UseAutoCert {
				err = srv.ListenAndServeTLS("", "")
			} else {
				err = srv.ListenAndServeTLS(cfg.TLSCert, cfg.TLSKey)
			}
		}
		if err != nil {
			if err == http.ErrServerClosed {
				return
			}
			log.Fatal().Err(err).Msg("start server failed")
		}
	}()

	listenShutdownSignal(srv)
}

func listenShutdownSignal(srv *http.Server) {
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	log.Info().Msg("shutdown server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("shutdown server error")
	}
}

func closeDBConn() {
	if db.GetDB() != nil {
		log.Info().Msg("close database connection")
		db.GetDB().Close()
	}
	if db.GetRedis() != nil {
		log.Info().Msg("close redis connection")
		db.GetRedis().Close()
	}
}
