package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"daemonw/conf"
	"daemonw/controller"
	"daemonw/dao"
	"daemonw/router"
	"daemonw/xlog"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	conf.InitConfig()
	xlog.InitLog()
	dao.InitDao()
	defer dao.CloseDao()

	cfg := conf.Config
	r := router.GetRouter()
	var tlsConf *tls.Config
	//tls config
	if cfg.TwoWayAuth {
		pool := x509.NewCertPool()
		caCertPath := cfg.ClientCA

		caCrt, err := ioutil.ReadFile(caCertPath)
		if err != nil {
			xlog.Fatal().Err(err).Msg("parse client ca failed")
			return
		}
		pool.AppendCertsFromPEM(caCrt)
		tlsConf = &tls.Config{
			ClientCAs:  pool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		}
	}
	srv := &http.Server{
		Addr:      ":" + strconv.Itoa(cfg.Port),
		Handler:   r,
		TLSConfig: tlsConf,
	}
	go func() {
		var err error
		if !cfg.TLS {
			xlog.Info().Msgf("start http server on %d", cfg.Port)
			err = srv.ListenAndServe()
		} else {
			xlog.Info().Msgf("start https server on %d", cfg.Port)
			err = srv.ListenAndServeTLS(cfg.TLSCert, cfg.TLSKey)
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

func main0() {
	spider := &controller.GoogleStoreSpider{}
	apkInfo, err := spider.FetchApkInfo("com.wire")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(apkInfo)
}
