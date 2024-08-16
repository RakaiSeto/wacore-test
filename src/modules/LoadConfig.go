package modules

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func LoadConfig() map[string]string {
	mapConfig := make(map[string]string)

	// Pasti NEO Staging
	//viper.SetConfigFile("D:\\GoogleDrive\\SyncedProjects\\PastiSwitchingNeoGO\\configuration\\config.json")
	//viper.SetConfigFile("/mnt/hgfs/SyncedProjects/PastiSwitchingNeoGO/configuration/config.json")

	// Pasti NEO Product - excel detail generator
	//viper.SetConfigFile("/pasti/agent/excelDetailTrxGenerator/resources/config.json")

	// PASTI NEO PRODUCTION 2.0
	// viper.SetConfigFile("../configuration/config.json")
	//viper.SetConfigFile("/pasti/config/config.json")

	// PINTAR NEO PRODUCT 2.0
	//viper.SetConfigFile("/pintar/router/config/config.json")
	//D:\Artamaya\Dev\src\Switching\OkeBayar\Go\ArdoraBillerNeo

	//==========Pasti Config================
	//Config Localhost
	//===============
	//viper.SetConfigFile("D:\\Artamaya\\Dev\\src\\Switching\\Pasti\\new pasti\\ArdoraBillerNeo\\configuration\\config.json")
	viper.SetConfigFile("/biller/config/config.json")
	//================
	//Config Production
	//=================
	//viper.SetConfigFile("/pasti/config/pasti_config.json")
	// PINTAR NEO PRODUCT 2.0 - REPORT GENERATOR
	//viper.SetConfigFile("/app/config.json")

	// ORCHESTRA BPJS BRIDGE
	// viper.SetConfigFile("/opt/PCU_BPJS_BRIDGE/config/config.json")

	if err := viper.ReadInConfig(); err != nil {
		log.Panic(err)
	}

	// Confirm which config file is used
	fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())

	// Database
	dbHost := viper.Get("database.host").(string)
	mapConfig["databaseHost"] = dbHost

	dbPort := viper.Get("database.port").(string)
	mapConfig["databasePort"] = dbPort

	dbName := viper.Get("database.name").(string)
	mapConfig["databaseName"] = dbName

	dbUser := viper.Get("database.user").(string)
	mapConfig["databaseUser"] = dbUser

	dbPass := viper.Get("database.pass").(string)
	mapConfig["databasePass"] = dbPass

	redisHost := viper.Get("redis.host").(string)
	mapConfig["redisHost"] = redisHost

	redisPort := viper.Get("redis.port").(string)
	mapConfig["redisPort"] = redisPort

	redisPass := viper.Get("redis.pass").(string)
	mapConfig["redisPass"] = redisPass

	//==============Load Config XmppBridge=============//
	//xmppAddress := viper.Get("xmppbridge.address").(string)
	//mapConfig["xmppAddress"] = xmppAddress
	//
	//xmppUser := viper.Get("xmppbridge.user").(string)
	//mapConfig["xmppUser"] = xmppUser
	//
	//xmppPass := viper.Get("xmppbridge.password").(string)
	//mapConfig["xmppPass"] = xmppPass
	//
	//xmppAddressAgent01 := viper.Get("xmppbridgeAgent01.address").(string)
	//mapConfig["xmppAddressAgent01"] = xmppAddressAgent01
	//
	//xmppUserAgent01 := viper.Get("xmppbridgeAgent01.user").(string)
	//mapConfig["xmppUserAgent01"] = xmppUserAgent01
	//
	//xmppPassAgent01 := viper.Get("xmppbridgeAgent01.password").(string)
	//mapConfig["xmppPassAgent01"] = xmppPassAgent01
	//
	//xmppAddressAgent02 := viper.Get("xmppbridgeAgent02.address").(string)
	//mapConfig["xmppAddressAgent02"] = xmppAddressAgent02
	//
	//xmppUserAgent02 := viper.Get("xmppbridgeAgent02.user").(string)
	//mapConfig["xmppUserAgent02"] = xmppUserAgent02
	//
	//xmppPassAgent02 := viper.Get("xmppbridgeAgent02.password").(string)
	//mapConfig["xmppPassAgent02"] = xmppPassAgent02
	//==============End Config XmppBridge=============//

	//logPathTrace := viper.Get("logging.pathTrace").(string)
	//mapConfig["loggingPathTrace"] = logPathTrace

	//logPathDebug := viper.Get("logging.pathDebug").(string)
	//mapConfig["loggingPathDebug"] = logPathDebug
	//
	//logPathError := viper.Get("logging.pathError").(string)
	//mapConfig["loggingPathError"] = logPathError
	//
	//logPathWarn := viper.Get("logging.pathWarn").(string)
	//mapConfig["loggingPathWarn"] = logPathWarn
	//
	//logPathInfo := viper.Get("logging.pathInfo").(string)
	//mapConfig["loggingPathInfo"] = logPathInfo

	appPath := viper.Get("applicationserver.path").(string)
	mapConfig["applicationPath"] = appPath

	appPort := viper.Get("applicationserver.port").(string)
	mapConfig["applicationPort"] = appPort

	rabbitHost := viper.Get("rabbitmq.host").(string)
	mapConfig["rabbitHost"] = rabbitHost

	rabbitPort := viper.Get("rabbitmq.port").(string)
	mapConfig["rabbitPort"] = rabbitPort
	//fmt.Printf("Using rabbitHost: %s\n", rabbitHost, ":"+rabbitPort)

	rabbitUser := viper.Get("rabbitmq.user").(string)
	mapConfig["rabbitUser"] = rabbitUser

	rabbitPass := viper.Get("rabbitmq.pass").(string)
	mapConfig["rabbitPass"] = rabbitPass

	rabbitVHost := viper.Get("rabbitmq.vhost").(string)
	mapConfig["rabbitVHost"] = rabbitVHost

	callbackReceiverPath := viper.Get("callbackserver.path").(string)
	mapConfig["callbackReceiverPath"] = callbackReceiverPath

	callbackReceiverPort := viper.Get("callbackserver.port").(string)
	mapConfig["callbackReceiverPort"] = callbackReceiverPort

	monitorClientSuccessRateThreshold := viper.Get("monitor.clientPercentSuccessRateThreshold").(string)
	mapConfig["monitorClientSuccessRateThreshold"] = monitorClientSuccessRateThreshold

	monitorProductSuccessRateThreshold := viper.Get("monitor.productPercentSuccessRateThreshold").(string)
	mapConfig["monitorProductSuccessRateThreshold"] = monitorProductSuccessRateThreshold

	monitorVendorSuccessRateThreshold := viper.Get("monitor.vendorPercentSuccessRateThreshold").(string)
	mapConfig["monitorVendorSuccessRateThreshold"] = monitorVendorSuccessRateThreshold

	//detailTrxReportDirectoryPathAdmin := viper.Get("detailtrxreport.directorypathadmin").(string)
	//mapConfig["detailTrxReportDirectoryPathAdmin"] = detailTrxReportDirectoryPathAdmin
	//
	//detailTrxReportDirectoryPathClient := viper.Get("detailtrxreport.directorypathclient").(string)
	//mapConfig["detailTrxReportDirectoryPathClient"] = detailTrxReportDirectoryPathClient

	detailAppPropertyClientGroup := viper.Get("appProperty.clientGroup").(string)
	mapConfig["appPropertyClientGroup"] = detailAppPropertyClientGroup

	systemNoProcessor := viper.Get("system.noProcessor").(string)
	mapConfig["systemNoProcessor"] = systemNoProcessor

	systemSSLCertPath := viper.Get("system.sslCertPath").(string)
	mapConfig["systemSSLCertPath"] = systemSSLCertPath

	systemSSLKeyPath := viper.Get("system.sslKeyPath").(string)
	mapConfig["systemSSLKeyPath"] = systemSSLKeyPath

	vendorDummyPath := viper.Get("vendordummy.path").(string)
	mapConfig["vendorDummyPath"] = vendorDummyPath

	vendorDummyPort := viper.Get("vendordummy.port").(string)
	mapConfig["vendorDummyPort"] = vendorDummyPort

	// bpjsBridgePath := viper.Get("bpjsbridge.path").(string)
	// mapConfig["bpjsBridgePath"] = bpjsBridgePath

	// bpjsBridgePort := viper.Get("bpjsbridge.port").(string)
	// mapConfig["bpjsBridgePort"] = bpjsBridgePort

	return mapConfig
}

func LoadConfigProduction() map[string]string {
	mapConfig := make(map[string]string)

	// Pasti NEO Staging
	//viper.SetConfigFile("D:\\GoogleDrive\\SyncedProjects\\PastiSwitchingNeoGO\\configuration\\config.json")
	//viper.SetConfigFile("/mnt/hgfs/SyncedProjects/PastiSwitchingNeoGO/configuration/config.json")

	// Pasti NEO Product - excel detail generator
	//viper.SetConfigFile("/pasti/agent/excelDetailTrxGenerator/resources/config.json")

	// PASTI NEO PRODUCTION 2.0
	// viper.SetConfigFile("../configuration/config.json")
	//viper.SetConfigFile("/pasti/config/config.json")

	// PINTAR NEO PRODUCT 2.0
	//viper.SetConfigFile("/pintar/router/config/config.json")
	//D:\Artamaya\Dev\src\Switching\OkeBayar\Go\ArdoraBillerNeo

	//==========Pasti Config================
	//Config Localhost
	//===============
	//viper.SetConfigFile("D:\\Artamaya\\Dev\\src\\Switching\\Pasti\\new pasti\\ArdoraBillerNeo\\configuration\\config.json")
	//================
	//Config Production
	//=================
	viper.SetConfigFile("/biller/config/config.json")
	// PINTAR NEO PRODUCT 2.0 - REPORT GENERATOR
	//viper.SetConfigFile("/app/config.json")

	// ORCHESTRA BPJS BRIDGE
	// viper.SetConfigFile("/opt/PCU_BPJS_BRIDGE/config/config.json")

	if err := viper.ReadInConfig(); err != nil {
		log.Panic(err)
	}

	// Confirm which config file is used
	fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())

	// Database
	dbHost := viper.Get("database.host").(string)
	mapConfig["databaseHost"] = dbHost

	dbPort := viper.Get("database.port").(string)
	mapConfig["databasePort"] = dbPort

	dbName := viper.Get("database.name").(string)
	mapConfig["databaseName"] = dbName

	dbUser := viper.Get("database.user").(string)
	mapConfig["databaseUser"] = dbUser

	dbPass := viper.Get("database.pass").(string)
	mapConfig["databasePass"] = dbPass

	redisHost := viper.Get("redis.host").(string)
	mapConfig["redisHost"] = redisHost

	redisPort := viper.Get("redis.port").(string)
	mapConfig["redisPort"] = redisPort

	redisPass := viper.Get("redis.pass").(string)
	mapConfig["redisPass"] = redisPass

	//==============Load Config XmppBridge=============//
	//xmppAddress := viper.Get("xmppbridge.address").(string)
	//mapConfig["xmppAddress"] = xmppAddress
	//
	//xmppUser := viper.Get("xmppbridge.user").(string)
	//mapConfig["xmppUser"] = xmppUser
	//
	//xmppPass := viper.Get("xmppbridge.password").(string)
	//mapConfig["xmppPass"] = xmppPass
	//
	//xmppAddressAgent01 := viper.Get("xmppbridgeAgent01.address").(string)
	//mapConfig["xmppAddressAgent01"] = xmppAddressAgent01
	//
	//xmppUserAgent01 := viper.Get("xmppbridgeAgent01.user").(string)
	//mapConfig["xmppUserAgent01"] = xmppUserAgent01
	//
	//xmppPassAgent01 := viper.Get("xmppbridgeAgent01.password").(string)
	//mapConfig["xmppPassAgent01"] = xmppPassAgent01
	//
	//xmppAddressAgent02 := viper.Get("xmppbridgeAgent02.address").(string)
	//mapConfig["xmppAddressAgent02"] = xmppAddressAgent02
	//
	//xmppUserAgent02 := viper.Get("xmppbridgeAgent02.user").(string)
	//mapConfig["xmppUserAgent02"] = xmppUserAgent02
	//
	//xmppPassAgent02 := viper.Get("xmppbridgeAgent02.password").(string)
	//mapConfig["xmppPassAgent02"] = xmppPassAgent02
	//==============End Config XmppBridge=============//

	//logPathTrace := viper.Get("logging.pathTrace").(string)
	//mapConfig["loggingPathTrace"] = logPathTrace

	//logPathDebug := viper.Get("logging.pathDebug").(string)
	//mapConfig["loggingPathDebug"] = logPathDebug
	//
	//logPathError := viper.Get("logging.pathError").(string)
	//mapConfig["loggingPathError"] = logPathError
	//
	//logPathWarn := viper.Get("logging.pathWarn").(string)
	//mapConfig["loggingPathWarn"] = logPathWarn
	//
	//logPathInfo := viper.Get("logging.pathInfo").(string)
	//mapConfig["loggingPathInfo"] = logPathInfo

	appPath := viper.Get("applicationserver.path").(string)
	mapConfig["applicationPath"] = appPath

	appPort := viper.Get("applicationserver.port").(string)
	mapConfig["applicationPort"] = appPort

	rabbitHost := viper.Get("rabbitmq.host").(string)
	mapConfig["rabbitHost"] = rabbitHost

	rabbitPort := viper.Get("rabbitmq.port").(string)
	mapConfig["rabbitPort"] = rabbitPort
	//fmt.Printf("Using rabbitHost: %s\n", rabbitHost, ":"+rabbitPort)

	rabbitUser := viper.Get("rabbitmq.user").(string)
	mapConfig["rabbitUser"] = rabbitUser

	rabbitPass := viper.Get("rabbitmq.pass").(string)
	mapConfig["rabbitPass"] = rabbitPass

	rabbitVHost := viper.Get("rabbitmq.vhost").(string)
	mapConfig["rabbitVHost"] = rabbitVHost

	callbackReceiverPath := viper.Get("callbackserver.path").(string)
	mapConfig["callbackReceiverPath"] = callbackReceiverPath

	callbackReceiverPort := viper.Get("callbackserver.port").(string)
	mapConfig["callbackReceiverPort"] = callbackReceiverPort

	monitorClientSuccessRateThreshold := viper.Get("monitor.clientPercentSuccessRateThreshold").(string)
	mapConfig["monitorClientSuccessRateThreshold"] = monitorClientSuccessRateThreshold

	monitorProductSuccessRateThreshold := viper.Get("monitor.productPercentSuccessRateThreshold").(string)
	mapConfig["monitorProductSuccessRateThreshold"] = monitorProductSuccessRateThreshold

	monitorVendorSuccessRateThreshold := viper.Get("monitor.vendorPercentSuccessRateThreshold").(string)
	mapConfig["monitorVendorSuccessRateThreshold"] = monitorVendorSuccessRateThreshold

	//detailTrxReportDirectoryPathAdmin := viper.Get("detailtrxreport.directorypathadmin").(string)
	//mapConfig["detailTrxReportDirectoryPathAdmin"] = detailTrxReportDirectoryPathAdmin
	//
	//detailTrxReportDirectoryPathClient := viper.Get("detailtrxreport.directorypathclient").(string)
	//mapConfig["detailTrxReportDirectoryPathClient"] = detailTrxReportDirectoryPathClient

	detailAppPropertyClientGroup := viper.Get("appProperty.clientGroup").(string)
	mapConfig["appPropertyClientGroup"] = detailAppPropertyClientGroup

	systemNoProcessor := viper.Get("system.noProcessor").(string)
	mapConfig["systemNoProcessor"] = systemNoProcessor

	systemSSLCertPath := viper.Get("system.sslCertPath").(string)
	mapConfig["systemSSLCertPath"] = systemSSLCertPath

	systemSSLKeyPath := viper.Get("system.sslKeyPath").(string)
	mapConfig["systemSSLKeyPath"] = systemSSLKeyPath

	vendorDummyPath := viper.Get("vendordummy.path").(string)
	mapConfig["vendorDummyPath"] = vendorDummyPath

	vendorDummyPort := viper.Get("vendordummy.port").(string)
	mapConfig["vendorDummyPort"] = vendorDummyPort

	// bpjsBridgePath := viper.Get("bpjsbridge.path").(string)
	// mapConfig["bpjsBridgePath"] = bpjsBridgePath

	// bpjsBridgePort := viper.Get("bpjsbridge.port").(string)
	// mapConfig["bpjsBridgePort"] = bpjsBridgePort

	return mapConfig
}
