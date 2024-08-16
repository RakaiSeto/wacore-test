package modules

import (
	"database/sql"
	"strconv"
	"time"
)

func SendAlert(priority int, datetime time.Time, module string, groupId string, clientId string, to string, content string, db *sql.DB) {
	// priority:
	// 1: Information
	// 2: Warning
	// 3: Alert, Margin Lost
	// 4: Company Impact, business stop

	// module:
	// clientBalance
	// clientSuccessRate
	// productSuccessRate
	// vendorSuccessRate
	// vendorBalance
	// application

	// Save to database
	alertId := clientId + "-" + module + "-" + strconv.Itoa(priority) + "-" + DoFormatDateTime("YYYY0M0DHHmmss", datetime)
	queryInsert := "INSERT INTO alert(alert_id, priority, datetime, module, group_id, alert_to, content) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	_, errInsert := db.Exec(queryInsert, alertId, priority, datetime, module, groupId, to, content)

	if errInsert != nil {
		DoLog("INFO", "", "pasti_Alert", "SendAlert",
			"Failed to save data alert. Error occured.", true, errInsert)

		panic(errInsert)
	} else {
		DoLog("INFO", "", "pasti_Alert", "SendAlert",
			"Success to save data alert "+alertId, true, errInsert)
	}

}
