package main

import (
	"context"
	"crypto/tls"
	"daemonw/api"
	"daemonw/conf"
	"daemonw/dao"
	"daemonw/xlog"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

func main() {
	err := conf.ParseConfig("")
	if err != nil {
		log.Fatal(err)
	}
	xlog.InitLog()
	err = dao.InitRedis()
	if err != nil {
		log.Fatal(err)
	}
	err = dao.InitDB()
	if err != nil {
		log.Fatal(err)
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
			xlog.Info().Msgf("start http server on %d", cfg.Port)
			err = srv.ListenAndServe()
		} else {
			xlog.Info().Msgf("start https server on %d", cfg.Port)
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
			xlog.Fatal().Err(err).Msg("start server failed")
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
	xlog.Info().Msg("shutdown server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		xlog.Error().Err(err).Msg("shutdown server error")
	}
}

func closeDBConn() {
	if dao.DB() != nil {
		xlog.Info().Msg("close database connection")
		fatalErr(dao.DB().Close())
	}
	if dao.Redis() != nil {
		xlog.Info().Msg("close redis connection")
		fatalErr(dao.Redis().Close())
	}
}

func fatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
