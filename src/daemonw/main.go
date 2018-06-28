package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"daemonw/api"
	"os"
	"os/signal"
	"time"
	"context"
	"strconv"
	"daemonw/conf"
	"path/filepath"
	"daemonw/log"
	"golang.org/x/crypto/acme/autocert"
	"crypto/tls"
	"daemonw/db"
)

func main() {
	router := api.GetRouter()
	//router.RunTLS(":"+port, cert, key)
	startServer(router, conf.Config.Port)
}

func startServer(router *gin.Engine, port int) {
	certManager := autocert.Manager{
		Cache:      autocert.DirCache(conf.BinDir + string(filepath.Separator) + ".certs"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("daemonw.cn"),
	}
	var err error
	var tlsConf *tls.Config
	if conf.Config.TLS && conf.Config.UseAutoCert {
		tlsConf = &tls.Config{GetCertificate: certManager.GetCertificate}
	}
	srv := &http.Server{
		Addr:      ":" + strconv.Itoa(port),
		Handler:   router,
		TLSConfig: tlsConf,
	}
	defer closeDBConn()
	go func() {
		if !conf.Config.TLS {
			log.Info().Msgf("start http server on %d", port)
			err = srv.ListenAndServe()
		} else {
			log.Info().Msgf("start https server on %d", port)
			if conf.Config.UseAutoCert {
				err = srv.ListenAndServeTLS("", "")
			} else {
				err = srv.ListenAndServeTLS(conf.Config.TLSCert, conf.Config.TLSKey)
			}
		}
		if err != nil {
			if err == http.ErrServerClosed {
				return
			}
			log.Fatal().Err(err).Msg("start server failed")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("error occurred when try to shutdown server")
	}
	log.Info().Msg("shutdown server")
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
