package modules

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strings"
)

func GetVendorBankCode(transactionId string, redisClient *redis.Client, goContext context.Context, systemBankCode string, vendorId string) string {
	mappingId := strings.TrimSpace(vendorId) + "-" + strings.TrimSpace(systemBankCode)

	// Get data from redis
	redisKey := "vendor_bankcode_map_" + mappingId
	redisVal := RedisGetDataRedis(transactionId, redisKey)

	if len(strings.TrimSpace(redisVal)) == 0 {
		return ""
	} else {
		mapVendorBankCode := ConvertJSONStringToMap(transactionId, redisVal)
		DoLog("DEBUG", transactionId, "DisburseDummy", "getVendorBankCode",
			fmt.Sprintf("mapVendorBankCode: %+v", mapVendorBankCode), false, nil)

		vendorBankCode := GetStringFromMapInterface(mapVendorBankCode, "vendorBankCode")
		DoLog("DEBUG", transactionId, "DisburseDummy", "getVendorBankCode",
			"systemBankCode: "+systemBankCode+" -> vendorId: "+vendorId+" -> vendorBankCode: "+vendorBankCode, false, nil)

		return vendorBankCode
	}
}
