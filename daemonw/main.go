package main

import (
	"context"
	"crypto/tls"
	"daemonw/api"
	"daemonw/conf"
	"daemonw/dao"
	"daemonw/util"
	"daemonw/xlog"
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
	dao.InitRedis()
	dao.InitDB()
	dao.InitDaoManager()
	defer closeDataSourceConn()

	cfg := conf.Config
	router := api.GetRouter()
	var tlsConf *tls.Config
	//tls config
	srv := &http.Server{
		Addr:      ":" + strconv.Itoa(cfg.Port),
		Handler:   router,
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

func init() {

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

func closeDataSourceConn() {
	if dao.DB() != nil {
		xlog.Info().Msg("close database connection")
		util.FatalIfErr(dao.DB().Close(), false)
	}
	if dao.Redis() != nil {
		xlog.Info().Msg("close redis connection")
		util.FatalIfErr(dao.Redis().Close(), false)
	}
}
