package main

import (
	"market-service/custom_config"
	"market-service/registry"

	"github.com/New-Tatthep/microservice"
	"github.com/New-Tatthep/microservice/log"
)

func main() {
	// Read specific config
	appSettingLocation := "./conf/service_server.yaml"
	appConf, err := microservice.NewConfig(appSettingLocation)
	if err != nil {
		log.Fatalf("%v", err)
	}

	apiConfig, apiConfigFound := appConf.GetAPIConfig()
	if !apiConfigFound {
		log.Fatal("api config not found")
	}

	logConfig, logCfgFound := appConf.GetLogConfig()
	if !logCfgFound {
		log.Fatal("log config not found")
	}

	dbConfigs, dbCfgFound := appConf.GetDBConfigs()
	if !dbCfgFound {
		log.Fatal("db config not found")
	}

	redisConfigs, redisCfgFound := appConf.GetRedisConfigs()
	if !redisCfgFound {
		log.Fatal("redis config not found")
	}

	kafkaConfigs, kafkaConfigsFound := appConf.GetKafkaConfigs()
	if !kafkaConfigsFound {
		log.Fatal("kafka config not found")
	}

	// asyncTaskConfig, asyncTaskConfigFound := appConf.GetAsyncTaskConfig()
	// if !asyncTaskConfigFound {
	// 	log.Fatalf("%v", err)
	// }

	// Read custom config
	customConfigFileLocation := "./conf/service_conf.yaml"
	_, err = custom_config.New(customConfigFileLocation)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// New server instance
	ms, err := microservice.New(
		microservice.WithAPIConfig(apiConfig),
		microservice.WithLogConfig(logConfig),
		microservice.WithDBConfigs(dbConfigs...),
		microservice.WithRedisConfigs(redisConfigs...),
		microservice.WithKafkaConfigs(kafkaConfigs...),
		// microservice.WithAsyncTaskConfig(asyncTaskConfig),
	)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Register route
	registry.APIRegister(ms)
	registry.APIProductRegister(ms)
	registry.ApiLoginRegister(ms)

	// Start server
	err = ms.Start()
	if err != nil {
		ms.Log(microservice.ErrorLevel, err.Error())
	}
}
