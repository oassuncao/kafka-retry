package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"github.com/oassuncao/kafka-retry/pkg/config"
	"github.com/oassuncao/kafka-retry/pkg/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var started bool
var manager *service.Manager

func getConfiguration() (*config.Application, error) {
	viper.SetEnvPrefix("RETRY_")
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if pathEnv := os.Getenv("CONFIG_PATH"); pathEnv != "" {
		paths := strings.Split(pathEnv, ",")
		for _, path := range paths {
			viper.AddConfigPath(strings.TrimSpace(path))
		}
	}

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var app = config.Application{}
	err = viper.Unmarshal(&app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func configureLogger(conf config.Application) {
	logConf := zap.NewProductionConfig()
	if conf.Log.Environment == "" {
		logConf = zap.NewDevelopmentConfig()
		logConf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	if conf.Log.Level != "" {
		var level zapcore.Level = 0
		err := level.Set(conf.Log.Level)
		if err != nil {
			panic(err)
		}
		logConf.Level = zap.NewAtomicLevelAt(level)
	}

	logger, err := logConf.Build()
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)
}

func StatusOk(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusOK)
}

func Health(writer http.ResponseWriter, _ *http.Request) {
	defaultStatus := http.StatusOK
	if started {
		conf := config.GetContext().Config.Retry
		if conf.GetDesiredRetries() != manager.GetRetriesSize() || conf.GetDesiredAttempts() != manager.GetAttemptsSize() {
			defaultStatus = http.StatusServiceUnavailable
		}
	}
	writer.WriteHeader(defaultStatus)
}

func startHttpServer() *http.Server {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/ready", StatusOk)
	http.HandleFunc("/health", Health)
	server := &http.Server{Addr: ":" + httpPort}
	go func() {
		zap.S().Info("Starting HTTP Server on port ", httpPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.S().Error("Error on start HTTP Server: ", err)
		}
	}()
	return server
}

func main() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	configuration, err := getConfiguration()
	if err != nil {
		panic(err)
	}

	config.SetContext(*configuration)
	configureLogger(*configuration)

	zap.S().Info("Creating manager")
	manager = service.NewManager()
	defer manager.Close()

	server := startHttpServer()

	zap.S().Info("Starting services")
	manager.Start()
	started = true

	zap.S().Infof("Started %d retries and %d attempts consumers", manager.GetRetriesSize(), manager.GetAttemptsSize())
	if manager.GetRetriesSize() != configuration.Retry.GetDesiredRetries() {
		zap.S().Warnf("The desired retries is %d and available is %d", configuration.Retry.GetDesiredRetries(), manager.GetRetriesSize())
	}

	if manager.GetAttemptsSize() != configuration.Retry.GetDesiredAttempts() {
		zap.S().Warnf("The desired attempts is %d and available is %d", configuration.Retry.GetDesiredAttempts(), manager.GetAttemptsSize())
	}

	//waiting to a sig channel
	sig := <-sigchan
	zap.S().Infof("Caught signal %v: terminating", sig)

	zap.S().Info("Stopping HTTP Server")
	if err = server.Close(); err != nil {
		zap.S().Error("Error on shutdown HTTP Server: ", err)
	}

	manager.Stop()
}
