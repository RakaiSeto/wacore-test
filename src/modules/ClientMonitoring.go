package modules

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

var durationInMinutes = 10

func handleNullInt64(parameter sql.NullInt64) int64 {
	if parameter.Valid {
		return parameter.Int64
	} else {
		return 0
	}
}

func MonitoringClientMinute(mapConfig map[string]string, db *sql.DB) {
	DoLog("INFO", "", "pasti_ClientMonitoring", "MonitoringClient",
		"Processing Monitoring Client agent.", false, nil)

	// Query every 5 minutes, for last 15 minutes data
	nowTime := time.Now()
	minutesBefore := nowTime.Add(-15 * time.Minute)

	query := "select cc.client_id as client_id_one, trx.* from client as cc " +
		"left join (" +
		"select client_id, " +
		"count(*) as total_traffic, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NOT NULL and query_code = 'inquiry') then 1 else 0 end) as traffic_prepaid_inquiry, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NOT NULL and query_code = 'inquiry') and final_trx_status = '000' then 1 else 0 end) as traffic_prepaid_inquiry_success, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NOT NULL and query_code = 'inquiry') and final_trx_status = '002' then 1 else 0 end) as traffic_prepaid_inquiry_pending, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NOT NULL and query_code = 'inquiry') and final_trx_status not like '00%' then 1 else 0 end) as traffic_prepaid_inquiry_failed, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NOT NULL and query_code = 'inquiry') and final_trx_status = '201' then 1 else 0 end) as traffic_prepaid_inquiry_failed_201, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NOT NULL and query_code = 'inquiry') and final_trx_status = '203' then 1 else 0 end) as traffic_prepaid_inquiry_failed_203, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NOT NULL and query_code = 'inquiry') and final_trx_status not like '00%' and final_trx_status != '201' and final_trx_status != '203' then 1 else 0 end) as traffic_prepaid_inquiry_failed_other, " +

		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NULL or (query_code IS NOT NULL and query_code != 'inquiry')) then 1 else 0 end) as traffic_prepaid_purchase, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NULL or (query_code IS NOT NULL and query_code != 'inquiry')) and final_trx_status = '000' then 1 else 0 end) as traffic_prepaid_purchase_success, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NULL or (query_code IS NOT NULL and query_code != 'inquiry')) and final_trx_status = '002' then 1 else 0 end) as traffic_prepaid_purchase_pending, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NULL or (query_code IS NOT NULL and query_code != 'inquiry')) and final_trx_status not like '00%' then 1 else 0 end) as traffic_prepaid_purchase_failed, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NULL or (query_code IS NOT NULL and query_code != 'inquiry')) and final_trx_status = '201' then 1 else 0 end) as traffic_prepaid_purchase_failed_201, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NULL or (query_code IS NOT NULL and query_code != 'inquiry')) and final_trx_status = '203' then 1 else 0 end) as traffic_prepaid_purchase_failed_203, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NULL or (query_code IS NOT NULL and query_code != 'inquiry')) and final_trx_status not like '00%' and final_trx_status != '201' and final_trx_status != '203' then 1 else 0 end) as traffic_prepaid_purchase_failed_other," +

		"sum(case when transaction_code = 'pasti-inquiry' then 1 else 0 end) as traffic_postpaid_inquiry, " +
		"sum(case when transaction_code = 'pasti-inquiry' and final_trx_status = '000' then 1 else 0 end) as traffic_postpaid_inquiry_success, " +
		"sum(case when transaction_code = 'pasti-inquiry' and final_trx_status = '002' then 1 else 0 end) as traffic_postpaid_inquiry_pending, " +
		"sum(case when transaction_code = 'pasti-inquiry' and final_trx_status not like '00%' then 1 else 0 end) as traffic_postpaid_inquiry_failed, " +
		"sum(case when transaction_code = 'pasti-inquiry' and final_trx_status = '201' then 1 else 0 end) as traffic_postpaid_inquiry_failed_201, " +
		"sum(case when transaction_code = 'pasti-inquiry' and final_trx_status = '203' then 1 else 0 end) as traffic_postpaid_inquiry_failed_203, " +
		"sum(case when transaction_code = 'pasti-inquiry' and final_trx_status not like '00%' and final_trx_status != '201' and final_trx_status != '203' then 1 else 0 end) as traffic_postpaid_inquiry_failed_other, " +

		"sum(case when transaction_code = 'pasti-pembayaran' then 1 else 0 end) as traffic_postpaid_payment, " +
		"sum(case when transaction_code = 'pasti-pembayaran' and final_trx_status = '000' then 1 else 0 end) as traffic_postpaid_payment_success, " +
		"sum(case when transaction_code = 'pasti-pembayaran' and final_trx_status = '002' then 1 else 0 end) as traffic_postpaid_payment_pending, " +
		"sum(case when transaction_code = 'pasti-pembayaran' and final_trx_status not like '00%' then 1 else 0 end) as traffic_postpaid_payment_failed, " +
		"sum(case when transaction_code = 'pasti-pembayaran' and final_trx_status = '201' then 1 else 0 end) as traffic_postpaid_payment_failed_201, " +
		"sum(case when transaction_code = 'pasti-pembayaran' and final_trx_status = '203' then 1 else 0 end) as traffic_postpaid_payment_failed_203, " +
		"sum(case when transaction_code = 'pasti-pembayaran' and final_trx_status not like '00%' and final_trx_status != '201' and final_trx_status != '203' then 1 else 0 end) as traffic_postpaid_payment_failed_other " +

		"from transaction where " +
		"(transaction_code = 'pasti-prabayar' or transaction_code = 'pasti-inquiry' or transaction_code = 'pasti-pembayaran') and " +
		"trx_datetime >= '" + DoFormatDateTime("YYYY-0M-0D HH:mm:ss", minutesBefore) + "' and trx_datetime <= '" + DoFormatDateTime("YYYY-0M-0D HH:mm:ss", nowTime) + "' " +
		"group by client_id) as trx " +
		"on cc.client_id = trx.client_id " +
		"order by cc.client_id"

	rows, err := db.Query(query)
	defer rows.Close()

	if err != nil {
		DoLog("INFO", "", "pasti_ClientMonitoring", "MonitoringClient",
			"Failed to read database transaction. Error occured.", true, err)

	} else {
		for rows.Next() {
			var clientId string
			var clientIdIgnore sql.NullString
			var trafficTotal sql.NullInt64

			var trafficPrepaidInquiry sql.NullInt64
			var trafficPrepaidInquirySuccess sql.NullInt64
			var trafficPrepaidInquiryPending sql.NullInt64
			var trafficPrepaidInquiryFailed sql.NullInt64
			var trafficPrepaidInquiryFailed201 sql.NullInt64
			var trafficPrepaidInquiryFailed203 sql.NullInt64
			var trafficPrepaidInquiryFailedOther sql.NullInt64

			var trafficPrepaidPurchase sql.NullInt64
			var trafficPrepaidPurchaseSuccess sql.NullInt64
			var trafficPrepaidPurchasePending sql.NullInt64
			var trafficPrepaidPurchaseFailed sql.NullInt64
			var trafficPrepaidPurchaseFailed201 sql.NullInt64
			var trafficPrepaidPurchaseFailed203 sql.NullInt64
			var trafficPrepaidPurchaseFailedOther sql.NullInt64

			var trafficPostpaidInquiry sql.NullInt64
			var trafficPostpaidInquirySuccess sql.NullInt64
			var trafficPostpaidInquiryPending sql.NullInt64
			var trafficPostpaidInquiryFailed sql.NullInt64
			var trafficPostpaidInquiryFailed201 sql.NullInt64
			var trafficPostpaidInquiryFailed203 sql.NullInt64
			var trafficPostpaidInquiryFailedOther sql.NullInt64

			var trafficPostpaidPayment sql.NullInt64
			var trafficPostpaidPaymentSuccess sql.NullInt64
			var trafficPostpaidPaymentPending sql.NullInt64
			var trafficPostpaidPaymentFailed sql.NullInt64
			var trafficPostpaidPaymentFailed201 sql.NullInt64
			var trafficPostpaidPaymentFailed203 sql.NullInt64
			var trafficPostpaidPaymentFailedOther sql.NullInt64

			err := rows.Scan(&clientId, &clientIdIgnore, &trafficTotal, &trafficPrepaidInquiry, &trafficPrepaidInquirySuccess, &trafficPrepaidInquiryPending, &trafficPrepaidInquiryFailed,
				&trafficPrepaidInquiryFailed201, &trafficPrepaidInquiryFailed203, &trafficPrepaidInquiryFailedOther, &trafficPrepaidPurchase,
				&trafficPrepaidPurchaseSuccess, &trafficPrepaidPurchasePending, &trafficPrepaidPurchaseFailed, &trafficPrepaidPurchaseFailed201,
				&trafficPrepaidPurchaseFailed203, &trafficPrepaidPurchaseFailedOther, &trafficPostpaidInquiry, &trafficPostpaidInquirySuccess,
				&trafficPostpaidInquiryPending, &trafficPostpaidInquiryFailed, &trafficPostpaidInquiryFailed201, &trafficPostpaidInquiryFailed203,
				&trafficPostpaidInquiryFailedOther, &trafficPostpaidPayment, &trafficPostpaidPaymentSuccess, &trafficPostpaidPaymentPending,
				&trafficPostpaidPaymentFailed, &trafficPostpaidPaymentFailed201, &trafficPostpaidPaymentFailed203, &trafficPostpaidPaymentFailedOther)

			// Convert to golang native data type
			xClientId := clientId
			xTrafficTotal := handleNullInt64(trafficTotal)

			xTrafficPrepaidInquiry := handleNullInt64(trafficPrepaidInquiry)
			xTrafficPrepaidInquirySuccess := handleNullInt64(trafficPrepaidInquirySuccess)
			xTrafficPrepaidInquiryPending := handleNullInt64(trafficPrepaidInquiryPending)
			xTrafficPrepaidInquiryFailed := handleNullInt64(trafficPrepaidInquiryFailed)
			xTrafficPrepaidInquiryFailed201 := handleNullInt64(trafficPrepaidInquiryFailed201)
			xTrafficPrepaidInquiryFailed203 := handleNullInt64(trafficPrepaidInquiryFailed203)
			xTrafficPrepaidInquiryFailedOther := handleNullInt64(trafficPrepaidInquiryFailedOther)

			// Calculate prepaid inquiry success rate
			var successRatePrepaidInquiry float64 = 0.00
			if xTrafficPrepaidInquiry > 0 {
				successRatePrepaidInquiry = ((float64(xTrafficPrepaidInquirySuccess) + float64(xTrafficPrepaidInquiryPending)) / float64(xTrafficPrepaidInquiry)) * 100

				if successRatePrepaidInquiry < 95 {
					contentPrepaidInquiry := "ALERT! Success rate PREPAID INQUIRY dari CLIENT ID " + clientId + ": <b>" + fmt.Sprintf("%.2f", successRatePrepaidInquiry) +
						" %</b> dari traffic inquiry prepaid: <b>" + strconv.Itoa(int(xTrafficPrepaidInquiry)) + "</b>. Please follow up!"

					SendAlert(3, time.Now(), "successRatePrepaidInquiry", "", clientId, "SYSADMIN", contentPrepaidInquiry, db)
				}
			}

			xTrafficPrepaidPurchase := handleNullInt64(trafficPrepaidPurchase)
			xTrafficPrepaidPurchaseSuccess := handleNullInt64(trafficPrepaidPurchaseSuccess)
			xTrafficPrepaidPurchasePending := handleNullInt64(trafficPrepaidPurchasePending)
			xTrafficPrepaidPurchaseFailed := handleNullInt64(trafficPrepaidPurchaseFailed)
			xTrafficPrepaidPurchaseFailed201 := handleNullInt64(trafficPrepaidPurchaseFailed201)
			xTrafficPrepaidPurchaseFailed203 := handleNullInt64(trafficPrepaidPurchaseFailed203)
			xTrafficPrepaidPurchaseFailedOther := handleNullInt64(trafficPrepaidPurchaseFailedOther)

			// Calculate prepaid purchase succses rate
			var successRatePrepaidPurchase float64 = 0.00
			if xTrafficPrepaidPurchase > 0 {
				successRatePrepaidPurchase = ((float64(xTrafficPrepaidPurchaseSuccess) + float64(xTrafficPrepaidPurchasePending)) / float64(xTrafficPrepaidPurchase)) * 100

				if successRatePrepaidPurchase < 95 {
					contentPrepaidPurchase := "ALERT! Success rate PREPAID PURCHASE dari CLIENT ID " + clientId + ": <b>" + fmt.Sprintf("%.2f", successRatePrepaidPurchase) +
						" %</b> dari traffic PREPAID PURCHASE: <b>" + strconv.Itoa(int(xTrafficPrepaidPurchase)) + "</b>. Please follow up!"

					SendAlert(3, time.Now(), "successRatePrepaidPurchase", "", clientId, "SYSADMIN", contentPrepaidPurchase, db)
				}
			}

			xTrafficPostpaidInquiry := handleNullInt64(trafficPostpaidInquiry)
			xTrafficPostpaidInquirySuccess := handleNullInt64(trafficPostpaidInquirySuccess)
			xTrafficPostpaidInquiryPending := handleNullInt64(trafficPostpaidInquiryPending)
			xTrafficPostpaidInquiryFailed := handleNullInt64(trafficPostpaidInquiryFailed)
			xTrafficPostpaidInquiryFailed201 := handleNullInt64(trafficPostpaidInquiryFailed201)
			xTrafficPostpaidInquiryFailed203 := handleNullInt64(trafficPostpaidInquiryFailed203)
			xTrafficPostpaidInquiryFailedOther := handleNullInt64(trafficPostpaidInquiryFailedOther)

			// Calculate postpaid inquiry success rate
			var successRatePostpaidInquiry float64 = 0.00
			if xTrafficPostpaidInquiry > 0 {
				successRatePostpaidInquiry = ((float64(xTrafficPostpaidInquirySuccess) + float64(xTrafficPostpaidInquiryPending)) / float64(xTrafficPostpaidInquiry)) * 100

				if successRatePostpaidInquiry < 95 {
					contentPostpaidInquiry := "ALERT! Success rate POSTPAID INQUIRY dari CLIENT ID " + clientId + ": <b>" + fmt.Sprintf("%.2f", successRatePostpaidInquiry) +
						" %</b> dari traffic POSTPAID INQUIRY: <b>" + strconv.Itoa(int(xTrafficPostpaidInquiry)) + "</b>. Please follow up!"

					SendAlert(3, time.Now(), "successRatePostpaidInquiry", "", clientId, "SYSADMIN", contentPostpaidInquiry, db)
				}
			}

			xTrafficPostpaidPayment := handleNullInt64(trafficPostpaidPayment)
			xTrafficPostpaidPaymentSuccess := handleNullInt64(trafficPostpaidPaymentSuccess)
			xTrafficPostpaidPaymentPending := handleNullInt64(trafficPostpaidPaymentPending)
			xTrafficPostpaidPaymentFailed := handleNullInt64(trafficPostpaidPaymentFailed)
			xTrafficPostpaidPaymentFailed201 := handleNullInt64(trafficPostpaidPaymentFailed201)
			xTrafficPostpaidPaymentFailed203 := handleNullInt64(trafficPostpaidPaymentFailed203)
			xTrafficPostpaidPaymentFailedOther := handleNullInt64(trafficPostpaidPaymentFailedOther)

			// Calculate postpaid payment success rate
			var successRatePostpaidPayment float64 = 0.00
			if xTrafficPostpaidPayment > 0 {
				successRatePostpaidPayment = ((float64(xTrafficPostpaidPaymentSuccess) + float64(xTrafficPostpaidPaymentPending)) / float64(xTrafficPostpaidPayment)) * 100

				if successRatePostpaidPayment < 95 {
					contentPostpaidPayment := "ALERT! Success rate POSTPAID PAYMENT dari CLIENT ID " + clientId + ": <b>" + fmt.Sprintf("%.2f", successRatePostpaidPayment) +
						" %</b> dari traffic POSTPAID PAYMENT: <b>" + strconv.Itoa(int(xTrafficPostpaidPayment)) + "</b>. Please follow up!"

					SendAlert(3, time.Now(), "successRatePostpaidPayment", "", clientId, "SYSADMIN", contentPostpaidPayment, db)
				}
			}

			fmt.Printf("clientId: %s, totalTraffic: %d \n", xClientId, xTrafficTotal)
			if err != nil {
				DoLog("INFO", "", "pasti_ClientMonitoring", "MonitoringClient",
					"Failed to read row of table transaction. Error occured.", true, err)
			} else {
				// Insert into table monitor_client
				queryInsert := "INSERT INTO monitor_client_minute(monitor_id, client_id, added_datetime, traffics_prepaid_inquiry, " +
					"prepaid_inquiry_success_traffic, prepaid_inquiry_pending_traffic, prepaid_inquiry_failed_traffic, prepaid_inquiry_failed_201, " +
					"prepaid_inquiry_failed_203, prepaid_inquiry_failed_other, traffic_prepaid_purchase, prepaid_purchase_success_traffic, " +
					"prepaid_purchase_pending_traffic, prepaid_purchase_failed_traffic, prepaid_purchase_failed_201, prepaid_purchase_failed_203, " +
					"prepaid_purchase_failed_other, traffic_postpaid_inquiry, postpaid_inquiry_success_traffic, postpaid_inquiry_pending_traffic, " +
					"postpaid_inquiry_failed_traffic, postpaid_inquiry_failed_201, postpaid_inquiry_failed_203, postpaid_inquiry_failed_other, " +
					"traffic_postpaid_payment, postpaid_payment_success_traffic, postpaid_payment_pending_traffic, postpaid_payment_failed_traffic, " +
					"postpaid_payment_failed_201, postpaid_payment_failed_203, postpaid_payment_failed_other, duration, monitor_group) VALUES (" +
					"$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, " +
					"$28, $29, $30, $31, $32, $33);"

				// Set monitor group
				monitorGroup := DoFormatDateTime("YYYY-0M-0D HH:mm", time.Now())

				// Set monitorId
				monitorId := clientId + "-" + monitorGroup

				_, errInsert := db.Exec(queryInsert, monitorId, xClientId, time.Now(), xTrafficPrepaidInquiry,
					xTrafficPrepaidInquirySuccess, xTrafficPrepaidInquiryPending, xTrafficPrepaidInquiryFailed, xTrafficPrepaidInquiryFailed201,
					xTrafficPrepaidInquiryFailed203, xTrafficPrepaidInquiryFailedOther, xTrafficPrepaidPurchase, xTrafficPrepaidPurchaseSuccess,
					xTrafficPrepaidPurchasePending, xTrafficPrepaidPurchaseFailed, xTrafficPrepaidPurchaseFailed201, xTrafficPrepaidPurchaseFailed203,
					xTrafficPrepaidPurchaseFailedOther, xTrafficPostpaidInquiry, xTrafficPostpaidInquirySuccess, xTrafficPostpaidInquiryPending,
					xTrafficPostpaidInquiryFailed, xTrafficPostpaidInquiryFailed201, xTrafficPostpaidInquiryFailed203, xTrafficPostpaidInquiryFailedOther,
					xTrafficPostpaidPayment, xTrafficPostpaidPaymentSuccess, xTrafficPostpaidPaymentPending, xTrafficPostpaidPaymentFailed,
					xTrafficPostpaidPaymentFailed201, xTrafficPostpaidPaymentFailed203, xTrafficPostpaidPaymentFailedOther, "5 MINUTES", monitorGroup)

				if errInsert != nil {
					DoLog("INFO", "", "pasti_ClientMonitoring", "MonitoringClient",
						"Failed to save data monitoring. Error occured.", true, err)
				} else {
					DoLog("INFO", "", "pasti_ClientMonitoring", "MonitoringClient",
						"Save to save data monitoring. ", true, err)
				}

			}
		}
	}
}

func MonitoringClientDay(db *sql.DB) {
	DoLog("INFO", "", "pasti_ClientMonitoring", "MonitoringClientDay",
		"Processing Daily Monitoring Client agent.", false, nil)

	// rows, err := db.Query(query)
	nowTime := time.Now()
	startOfDayTime := GetBeginOfDayDateTime(nowTime)

	query := "select cc.client_id as client_id_one, trx.* from client as cc " +
		"left join (" +
		"select client_id, " +
		"count(*) as total_traffic, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NOT NULL and query_code = 'inquiry') then 1 else 0 end) as traffic_prepaid_inquiry, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NOT NULL and query_code = 'inquiry') and final_trx_status = '000' then 1 else 0 end) as traffic_prepaid_inquiry_success, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NOT NULL and query_code = 'inquiry') and final_trx_status = '002' then 1 else 0 end) as traffic_prepaid_inquiry_pending, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NOT NULL and query_code = 'inquiry') and final_trx_status not like '00%' then 1 else 0 end) as traffic_prepaid_inquiry_failed, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NOT NULL and query_code = 'inquiry') and final_trx_status = '201' then 1 else 0 end) as traffic_prepaid_inquiry_failed_201, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NOT NULL and query_code = 'inquiry') and final_trx_status = '203' then 1 else 0 end) as traffic_prepaid_inquiry_failed_203, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NOT NULL and query_code = 'inquiry') and final_trx_status not like '00%' and final_trx_status != '201' and final_trx_status != '203' then 1 else 0 end) as traffic_prepaid_inquiry_failed_other, " +

		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NULL or (query_code IS NOT NULL and query_code != 'inquiry')) then 1 else 0 end) as traffic_prepaid_purchase, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NULL or (query_code IS NOT NULL and query_code != 'inquiry')) and final_trx_status = '000' then 1 else 0 end) as traffic_prepaid_purchase_success, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NULL or (query_code IS NOT NULL and query_code != 'inquiry')) and final_trx_status = '002' then 1 else 0 end) as traffic_prepaid_purchase_pending, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NULL or (query_code IS NOT NULL and query_code != 'inquiry')) and final_trx_status not like '00%' then 1 else 0 end) as traffic_prepaid_purchase_failed, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NULL or (query_code IS NOT NULL and query_code != 'inquiry')) and final_trx_status = '201' then 1 else 0 end) as traffic_prepaid_purchase_failed_201, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NULL or (query_code IS NOT NULL and query_code != 'inquiry')) and final_trx_status = '203' then 1 else 0 end) as traffic_prepaid_purchase_failed_203, " +
		"sum(case when transaction_code = 'pasti-prabayar' and (query_code IS NULL or (query_code IS NOT NULL and query_code != 'inquiry')) and final_trx_status not like '00%' and final_trx_status != '201' and final_trx_status != '203' then 1 else 0 end) as traffic_prepaid_purchase_failed_other," +

		"sum(case when transaction_code = 'pasti-inquiry' then 1 else 0 end) as traffic_postpaid_inquiry, " +
		"sum(case when transaction_code = 'pasti-inquiry' and final_trx_status = '000' then 1 else 0 end) as traffic_postpaid_inquiry_success, " +
		"sum(case when transaction_code = 'pasti-inquiry' and final_trx_status = '002' then 1 else 0 end) as traffic_postpaid_inquiry_pending, " +
		"sum(case when transaction_code = 'pasti-inquiry' and final_trx_status not like '00%' then 1 else 0 end) as traffic_postpaid_inquiry_failed, " +
		"sum(case when transaction_code = 'pasti-inquiry' and final_trx_status = '201' then 1 else 0 end) as traffic_postpaid_inquiry_failed_201, " +
		"sum(case when transaction_code = 'pasti-inquiry' and final_trx_status = '203' then 1 else 0 end) as traffic_postpaid_inquiry_failed_203, " +
		"sum(case when transaction_code = 'pasti-inquiry' and final_trx_status not like '00%' and final_trx_status != '201' and final_trx_status != '203' then 1 else 0 end) as traffic_postpaid_inquiry_failed_other, " +

		"sum(case when transaction_code = 'pasti-pembayaran' then 1 else 0 end) as traffic_postpaid_payment, " +
		"sum(case when transaction_code = 'pasti-pembayaran' and final_trx_status = '000' then 1 else 0 end) as traffic_postpaid_payment_success, " +
		"sum(case when transaction_code = 'pasti-pembayaran' and final_trx_status = '002' then 1 else 0 end) as traffic_postpaid_payment_pending, " +
		"sum(case when transaction_code = 'pasti-pembayaran' and final_trx_status not like '00%' then 1 else 0 end) as traffic_postpaid_payment_failed, " +
		"sum(case when transaction_code = 'pasti-pembayaran' and final_trx_status = '201' then 1 else 0 end) as traffic_postpaid_payment_failed_201, " +
		"sum(case when transaction_code = 'pasti-pembayaran' and final_trx_status = '203' then 1 else 0 end) as traffic_postpaid_payment_failed_203, " +
		"sum(case when transaction_code = 'pasti-pembayaran' and final_trx_status not like '00%' and final_trx_status != '201' and final_trx_status != '203' then 1 else 0 end) as traffic_postpaid_payment_failed_other " +

		"from transaction where " +
		"(transaction_code = 'pasti-prabayar' or transaction_code = 'pasti-inquiry' or transaction_code = 'pasti-pembayaran') and " +
		"trx_datetime >= '" + DoFormatDateTime("YYYY-0M-0D HH:mm:ss", startOfDayTime) + "' and trx_datetime <= '" + DoFormatDateTime("YYYY-0M-0D HH:mm:ss", nowTime) + "' " +
		"group by client_id) as trx " +
		"on cc.client_id = trx.client_id " +
		"order by cc.client_id"

	rows, err := db.Query(query)
	defer rows.Close()

	if err != nil {
		DoLog("INFO", "", "pasti_ClientMonitoring", "MonitoringClient",
			"Failed to read database transaction. Error occured.", true, err)

	} else {
		for rows.Next() {
			var clientId string
			var clientIdIgnore sql.NullString
			var trafficTotal sql.NullInt64

			var trafficPrepaidInquiry sql.NullInt64
			var trafficPrepaidInquirySuccess sql.NullInt64
			var trafficPrepaidInquiryPending sql.NullInt64
			var trafficPrepaidInquiryFailed sql.NullInt64
			var trafficPrepaidInquiryFailed201 sql.NullInt64
			var trafficPrepaidInquiryFailed203 sql.NullInt64
			var trafficPrepaidInquiryFailedOther sql.NullInt64

			var trafficPrepaidPurchase sql.NullInt64
			var trafficPrepaidPurchaseSuccess sql.NullInt64
			var trafficPrepaidPurchasePending sql.NullInt64
			var trafficPrepaidPurchaseFailed sql.NullInt64
			var trafficPrepaidPurchaseFailed201 sql.NullInt64
			var trafficPrepaidPurchaseFailed203 sql.NullInt64
			var trafficPrepaidPurchaseFailedOther sql.NullInt64

			var trafficPostpaidInquiry sql.NullInt64
			var trafficPostpaidInquirySuccess sql.NullInt64
			var trafficPostpaidInquiryPending sql.NullInt64
			var trafficPostpaidInquiryFailed sql.NullInt64
			var trafficPostpaidInquiryFailed201 sql.NullInt64
			var trafficPostpaidInquiryFailed203 sql.NullInt64
			var trafficPostpaidInquiryFailedOther sql.NullInt64

			var trafficPostpaidPayment sql.NullInt64
			var trafficPostpaidPaymentSuccess sql.NullInt64
			var trafficPostpaidPaymentPending sql.NullInt64
			var trafficPostpaidPaymentFailed sql.NullInt64
			var trafficPostpaidPaymentFailed201 sql.NullInt64
			var trafficPostpaidPaymentFailed203 sql.NullInt64
			var trafficPostpaidPaymentFailedOther sql.NullInt64

			err := rows.Scan(&clientId, &clientIdIgnore, &trafficTotal, &trafficPrepaidInquiry, &trafficPrepaidInquirySuccess, &trafficPrepaidInquiryPending, &trafficPrepaidInquiryFailed,
				&trafficPrepaidInquiryFailed201, &trafficPrepaidInquiryFailed203, &trafficPrepaidInquiryFailedOther, &trafficPrepaidPurchase,
				&trafficPrepaidPurchaseSuccess, &trafficPrepaidPurchasePending, &trafficPrepaidPurchaseFailed, &trafficPrepaidPurchaseFailed201,
				&trafficPrepaidPurchaseFailed203, &trafficPrepaidPurchaseFailedOther, &trafficPostpaidInquiry, &trafficPostpaidInquirySuccess,
				&trafficPostpaidInquiryPending, &trafficPostpaidInquiryFailed, &trafficPostpaidInquiryFailed201, &trafficPostpaidInquiryFailed203,
				&trafficPostpaidInquiryFailedOther, &trafficPostpaidPayment, &trafficPostpaidPaymentSuccess, &trafficPostpaidPaymentPending,
				&trafficPostpaidPaymentFailed, &trafficPostpaidPaymentFailed201, &trafficPostpaidPaymentFailed203, &trafficPostpaidPaymentFailedOther)

			// Convert to golang native data type
			xClientId := clientId
			xTrafficTotal := handleNullInt64(trafficTotal)

			xTrafficPrepaidInquiry := handleNullInt64(trafficPrepaidInquiry)
			xTrafficPrepaidInquirySuccess := handleNullInt64(trafficPrepaidInquirySuccess)
			xTrafficPrepaidInquiryPending := handleNullInt64(trafficPrepaidInquiryPending)
			xTrafficPrepaidInquiryFailed := handleNullInt64(trafficPrepaidInquiryFailed)
			xTrafficPrepaidInquiryFailed201 := handleNullInt64(trafficPrepaidInquiryFailed201)
			xTrafficPrepaidInquiryFailed203 := handleNullInt64(trafficPrepaidInquiryFailed203)
			xTrafficPrepaidInquiryFailedOther := handleNullInt64(trafficPrepaidInquiryFailedOther)

			xTrafficPrepaidPurchase := handleNullInt64(trafficPrepaidPurchase)
			xTrafficPrepaidPurchaseSuccess := handleNullInt64(trafficPrepaidPurchaseSuccess)
			xTrafficPrepaidPurchasePending := handleNullInt64(trafficPrepaidPurchasePending)
			xTrafficPrepaidPurchaseFailed := handleNullInt64(trafficPrepaidPurchaseFailed)
			xTrafficPrepaidPurchaseFailed201 := handleNullInt64(trafficPrepaidPurchaseFailed201)
			xTrafficPrepaidPurchaseFailed203 := handleNullInt64(trafficPrepaidPurchaseFailed203)
			xTrafficPrepaidPurchaseFailedOther := handleNullInt64(trafficPrepaidPurchaseFailedOther)

			xTrafficPostpaidInquiry := handleNullInt64(trafficPostpaidInquiry)
			xTrafficPostpaidInquirySuccess := handleNullInt64(trafficPostpaidInquirySuccess)
			xTrafficPostpaidInquiryPending := handleNullInt64(trafficPostpaidInquiryPending)
			xTrafficPostpaidInquiryFailed := handleNullInt64(trafficPostpaidInquiryFailed)
			xTrafficPostpaidInquiryFailed201 := handleNullInt64(trafficPostpaidInquiryFailed201)
			xTrafficPostpaidInquiryFailed203 := handleNullInt64(trafficPostpaidInquiryFailed203)
			xTrafficPostpaidInquiryFailedOther := handleNullInt64(trafficPostpaidInquiryFailedOther)

			xTrafficPostpaidPayment := handleNullInt64(trafficPostpaidPayment)
			xTrafficPostpaidPaymentSuccess := handleNullInt64(trafficPostpaidPaymentSuccess)
			xTrafficPostpaidPaymentPending := handleNullInt64(trafficPostpaidPaymentPending)
			xTrafficPostpaidPaymentFailed := handleNullInt64(trafficPostpaidPaymentFailed)
			xTrafficPostpaidPaymentFailed201 := handleNullInt64(trafficPostpaidPaymentFailed201)
			xTrafficPostpaidPaymentFailed203 := handleNullInt64(trafficPostpaidPaymentFailed203)
			xTrafficPostpaidPaymentFailedOther := handleNullInt64(trafficPostpaidPaymentFailedOther)

			fmt.Printf("clientId: %s, totalTraffic: %d \n", xClientId, xTrafficTotal)
			if err != nil {
				DoLog("INFO", "", "pasti_ClientMonitoring", "MonitoringClient",
					"Failed to read row of table transaction. Error occured.", true, err)
			} else {
				// Insert into table monitor_client
				queryInsert := "INSERT INTO monitor_client_minute(monitor_id, client_id, added_datetime, traffics_prepaid_inquiry, " +
					"prepaid_inquiry_success_traffic, prepaid_inquiry_pending_traffic, prepaid_inquiry_failed_traffic, prepaid_inquiry_failed_201, " +
					"prepaid_inquiry_failed_203, prepaid_inquiry_failed_other, traffic_prepaid_purchase, prepaid_purchase_success_traffic, " +
					"prepaid_purchase_pending_traffic, prepaid_purchase_failed_traffic, prepaid_purchase_failed_201, prepaid_purchase_failed_203, " +
					"prepaid_purchase_failed_other, traffic_postpaid_inquiry, postpaid_inquiry_success_traffic, postpaid_inquiry_pending_traffic, " +
					"postpaid_inquiry_failed_traffic, postpaid_inquiry_failed_201, postpaid_inquiry_failed_203, postpaid_inquiry_failed_other, " +
					"traffic_postpaid_payment, postpaid_payment_success_traffic, postpaid_payment_pending_traffic, postpaid_payment_failed_traffic, " +
					"postpaid_payment_failed_201, postpaid_payment_failed_203, postpaid_payment_failed_other, duration) VALUES (" +
					"$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32);"

				// Set monitorId
				monitorId := clientId + "-" + DoFormatDateTime("YYYY-0M-0D HH:mm", time.Now())

				_, errInsert := db.Exec(queryInsert, monitorId, xClientId, time.Now(), xTrafficPrepaidInquiry,
					xTrafficPrepaidInquirySuccess, xTrafficPrepaidInquiryPending, xTrafficPrepaidInquiryFailed, xTrafficPrepaidInquiryFailed201,
					xTrafficPrepaidInquiryFailed203, xTrafficPrepaidInquiryFailedOther, xTrafficPrepaidPurchase, xTrafficPrepaidPurchaseSuccess,
					xTrafficPrepaidPurchasePending, xTrafficPrepaidPurchaseFailed, xTrafficPrepaidPurchaseFailed201, xTrafficPrepaidPurchaseFailed203,
					xTrafficPrepaidPurchaseFailedOther, xTrafficPostpaidInquiry, xTrafficPostpaidInquirySuccess, xTrafficPostpaidInquiryPending,
					xTrafficPostpaidInquiryFailed, xTrafficPostpaidInquiryFailed201, xTrafficPostpaidInquiryFailed203, xTrafficPostpaidInquiryFailedOther,
					xTrafficPostpaidPayment, xTrafficPostpaidPaymentSuccess, xTrafficPostpaidPaymentPending, xTrafficPostpaidPaymentFailed,
					xTrafficPostpaidPaymentFailed201, xTrafficPostpaidPaymentFailed203, xTrafficPostpaidPaymentFailedOther, "5 MINUTES")

				if errInsert != nil {
					DoLog("INFO", "", "pasti_ClientMonitoring", "MonitoringClient",
						"Failed to save data monitoring. Error occured.", true, err)

					fmt.Printf("query: " + queryInsert)
				} else {
					DoLog("INFO", "", "pasti_ClientMonitoring", "MonitoringClient",
						"Save to save data monitoring. ", true, err)
				}

			}
		}
	}
}
