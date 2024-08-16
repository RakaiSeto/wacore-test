package config

import (
	waProto "go.mau.fi/whatsmeow/proto/waCompanionReg"
)

var (
	AppVersion             = "v1.0.0"
	AppPort                = "3000"
	AppDebug               = false
	AppOs                  = "NTE"
	AppPlatform            = waProto.DeviceProps_CHROME
	AppBasicAuthCredential string

	PathQrCode    = "statics/qrcode"
	PathSendItems = "/filestorage"
	PathMedia     = "statics/media"
	PathStorages  = "storages"

	DBName = "whatsapp.db"

	WhatsappAutoReplyMessage    string
	WhatsappWebhook             string
	WhatsappLogLevel                  = "ERROR"
	WhatsappSettingMaxFileSize  int64 = 5000000  // 50MB
	WhatsappSettingMaxVideoSize int64 = 10000000 // 100MB
	WhatsappTypeUser                  = "@s.whatsapp.net"
	WhatsappTypeGroup                 = "@g.us"
)
