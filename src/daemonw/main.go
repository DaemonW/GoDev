package main

import (
	"log"
	"net/http"
	"daemonw/api"
	"os"
	"os/signal"
	"time"
	"context"
	"strconv"
	"daemonw/conf"
	"path/filepath"
	mylog "daemonw/log"
	"golang.org/x/crypto/acme/autocert"
	"crypto/tls"
	"daemonw/db"
	"syscall"
)

func main() {
	err := conf.ParseConfig("")
	if err != nil {
		log.Fatal(err)
	}
	mylog.InitLog()
	err = db.InitRedis()
	if err!=nil{
		log.Fatal(err)
	}
	err = db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
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
			mylog.Info().Msgf("start http server on %d", cfg.Port)
			err = srv.ListenAndServe()
		} else {
			mylog.Info().Msgf("start https server on %d", cfg.Port)
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
			mylog.Fatal().Err(err).Msg("start server failed")
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
	mylog.Info().Msg("shutdown server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		mylog.Error().Err(err).Msg("shutdown server error")
	}
}

func closeDBConn() {
	if db.GetDB() != nil {
		mylog.Info().Msg("close database connection")
		db.GetDB().Close()
	}
	if db.GetRedis() != nil {
		mylog.Info().Msg("close redis connection")
		db.GetRedis().Close()
	}
}
