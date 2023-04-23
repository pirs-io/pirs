package config

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"pirs.io/commons"
	"pirs.io/process-storage/storage"
)

var (
	log                                           = commons.GetLoggerFor("config")
	trackerApplicationContext *ApplicationContext = nil
)

type ProcessStorageAppConfig struct {
	commons.BaseConfig
	GrpcPort        int              `mapstructure:"GRPC_PORT"`
	StorageProvider storage.Provider `mapstructure:"STORAGE_PROVIDER"`
	RepoRootPath    string           `mapstructure:"GIT_ROOT"`
	Tenant          string           `mapstructure:"TENANT"`
	ChunkSize       int64            `mapstructure:"CHUNK_SIZE_BYTES"`
	PrometheusPort  int              `mapstructure:"PROMETHEUS_PORT"`
}

type ApplicationContext struct {
	AppConfig ProcessStorageAppConfig
}

func InitApp(configFilePath string) (conf *ProcessStorageAppConfig) {
	// config loading
	log.Info().Msg("Starting application")
	conf, confErr := commons.GetAppConfig(configFilePath, &ProcessStorageAppConfig{})
	if confErr != nil {
		log.Fatal().Msgf("Unable to load application config! %s", confErr)
	}

	// create app context
	appCtx, contextErr := createApplicationContext(*conf)
	if contextErr != nil {
		log.Fatal().Msgf("Error initializing app context")
	}
	trackerApplicationContext = appCtx

	err := initMetrics(conf.PrometheusPort)
	if err != nil {
		log.Warn().Msgf("Failed to initialize metric config! %s", err.Error())
	} else {
		log.Info().Msgf(fmt.Sprintf("Prometheus endpoint available at: %d/metrics", conf.PrometheusPort))
	}

	log.Info().Msg("Application started!")
	trackerApplicationContext.AppConfig = *conf
	return conf
}

func initMetrics(prometheusPort int) error {
	log.Info().Msg("Initializing metric instrumentation")
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewBuildInfoCollector(),
	)

	server := &http.Server{Addr: fmt.Sprintf(":%d", prometheusPort)}
	serveMux := http.NewServeMux()
	promhttp.Handler()
	serveMux.Handle("/metrics", promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{EnableOpenMetrics: true},
	))
	server.Handler = serveMux
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Err(err)
		}
	}()
	return nil
}

func GetContext() *ApplicationContext {
	return trackerApplicationContext
}

func createApplicationContext(conf ProcessStorageAppConfig) (appContext *ApplicationContext, err error) {
	return &ApplicationContext{conf}, nil
}
