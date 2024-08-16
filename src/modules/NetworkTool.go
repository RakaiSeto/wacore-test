package modules

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func HTTPConvertHeaderToMap(httpHeader http.Header) map[string]interface{} {
	var respMap = make(map[string]interface{})

	for key, val := range httpHeader {
		respMap[key] = val
	}

	return respMap
}

// noinspection GoUnusedExportedFunction
func HTTPGETString(traceCode string, mainUrl string, headerRequest map[string]interface{}, urlParameter map[string]string) map[string]interface{} {
	DoLog("DEBUG", traceCode, "NetworkTool", "HTTPS_GET_String", fmt.Sprintf("HTTP GET - body: %v, header: %v", urlParameter, headerRequest), false, nil)
	var mapResponse = make(map[string]interface{})
	var headerResponse map[string]interface{}
	var bodyResponse string
	var startDateTime = time.Now()
	var endDateTime = time.Now()
	var timeDiff int64 // millisecond
	var httpStatus string
	var hitStatus = "000"

	// Compose request with body
	startDateTime = time.Now()
	req, err := http.NewRequest("GET", mainUrl, nil)

	if err != nil {
		hitStatus = "240"

		DoLog("DEBUG", traceCode, "NetworkTool", "HTTPS_POST_String",
			"Failed to compose HTTP Request. Error occured.", true, err)
	} else {
		// Compose header
		if headerRequest != nil {
			for key, val := range headerRequest {
				req.Header.Set(key, val.(string))
			}
		}

		// Buat encoded url parameter
		query := url.Values{}

		if urlParameter != nil {
			for key, val := range urlParameter {
				query.Add(key, val)
			}
		}
		req.URL.RawQuery = query.Encode()

		// Buat HTTP Client
		client := &http.Client{}

		// Get the response
		resp, errR := client.Do(req)

		if errR != nil {
			hitStatus = "240"

			DoLog("DEBUG", traceCode, "NetworkTool", "HTTP_POST_String",
				"Failed to execute HTTP Request. Error occured.", true, err)
		} else {
			endDateTime = time.Now()

			body, _ := ioutil.ReadAll(resp.Body)

			// fill bodyResponse, headerResponse, dll
			bodyResponse = string(body)
			headerResponse = HTTPConvertHeaderToMap(resp.Header)
			httpStatus = resp.Status

			_ = resp.Body.Close()
		}
	}

	timeDiff = endDateTime.Sub(startDateTime).Milliseconds()

	// Put into mapResponse
	mapResponse["bodyRequest"] = req.URL.RawQuery
	mapResponse["headerRequest"] = headerRequest
	mapResponse["bodyResponse"] = bodyResponse
	mapResponse["headerResponse"] = headerResponse
	mapResponse["httpStatus"] = httpStatus
	mapResponse["startDateTime"] = startDateTime
	mapResponse["endDateTime"] = endDateTime
	mapResponse["timeDifference"] = timeDiff
	mapResponse["status"] = hitStatus

	DoLog("DEBUG", traceCode, "NetworkTool", "HTTPS_POST_String",
		fmt.Sprintf("mapResponse: %v", mapResponse), true, err)

	return mapResponse
}

// noinspection GoUnusedExportedFunction
func HTTPPOSTString(traceCode string, url string, headerRequest map[string]interface{}, bodyRequest string) map[string]interface{} {
	DoLog("DEBUG", traceCode, "NetworkTool", "HTTP_POST_String", fmt.Sprintf("HTTP POST - body: %s, header: %v", bodyRequest, headerRequest), false, nil)
	var mapResponse = make(map[string]interface{})
	var headerResponse map[string]interface{}
	var bodyResponse string
	var startDateTime = time.Now()
	var endDateTime = time.Now()
	var timeDiff int64 // millisecond
	var httpStatus string
	var hitStatus = "000"

	var bodyByte = []byte(bodyRequest)

	// Compose request with body
	startDateTime = time.Now()
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyByte))

	if err != nil {
		hitStatus = "240"

		DoLog("DEBUG", traceCode, "NetworkTool", "HTTP_POST_String",
			"Failed to compose HTTP Request. Error occured.", true, err)
	} else {
		// Compose header
		if headerRequest != nil {
			for key, val := range headerRequest {
				req.Header.Set(key, val.(string))
			}
		}

		client := &http.Client{}

		// Get the response
		resp, errR := client.Do(req)

		if errR != nil {
			hitStatus = "240"

			DoLog("DEBUG", traceCode, "NetworkTool", "HTTP_POST_String",
				"Failed to execute HTTP Request. Error occured.", true, err)
		} else {
			endDateTime = time.Now()

			body, _ := ioutil.ReadAll(resp.Body)

			// fill bodyResponse, headerResponse, dll
			bodyResponse = string(body)
			headerResponse = HTTPConvertHeaderToMap(resp.Header)
			httpStatus = resp.Status

			_ = resp.Body.Close()
		}
	}

	timeDiff = endDateTime.Sub(startDateTime).Milliseconds()

	// Put into mapResponse
	mapResponse["bodyRequest"] = bodyRequest
	mapResponse["headerRequest"] = headerRequest
	mapResponse["bodyResponse"] = bodyResponse
	mapResponse["headerResponse"] = headerResponse
	mapResponse["httpStatus"] = httpStatus
	mapResponse["startDateTime"] = startDateTime
	mapResponse["endDateTime"] = endDateTime
	mapResponse["timeDifference"] = timeDiff
	mapResponse["status"] = hitStatus

	DoLog("DEBUG", traceCode, "NetworkTool", "HTTP_POST_String",
		fmt.Sprintf("mapResponse: %v", mapResponse), true, err)

	return mapResponse
}

// noinspection GoUnusedExportedFunction
func HTTPSPOSTString(traceCode string, url string, headerRequest map[string]interface{}, bodyRequest string) map[string]interface{} {
	DoLog("DEBUG", traceCode, "NetworkTool", "HTTPS_POST_String", fmt.Sprintf("HTTP POST - body: %s, header: %v", bodyRequest, headerRequest), false, nil)
	var mapResponse = make(map[string]interface{})
	var headerResponse map[string]interface{}
	var bodyResponse string
	var startDateTime = time.Now()
	var endDateTime = time.Now()
	var timeDiff int64 // millisecond
	var httpStatus string
	var hitStatus = "000"

	var bodyByte = []byte(bodyRequest)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Compose request with body
	startDateTime = time.Now()
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyByte))

	if err != nil {
		hitStatus = "240"

		DoLog("DEBUG", traceCode, "NetworkTool", "HTTPS_POST_String",
			"Failed to compose HTTP Request. Error occured.", true, err)
	} else {
		// Compose header
		if headerRequest != nil {
			for key, val := range headerRequest {
				req.Header.Set(key, val.(string))
			}
		}

		client := &http.Client{Transport: tr}

		// Get the response
		resp, errR := client.Do(req)

		if errR != nil {
			hitStatus = "240"

			DoLog("DEBUG", traceCode, "NetworkTool", "HTTP_POST_String",
				"Failed to execute HTTP Request. Error occured.", true, err)
		} else {
			endDateTime = time.Now()

			body, _ := ioutil.ReadAll(resp.Body)

			// fill bodyResponse, headerResponse, dll
			bodyResponse = string(body)
			headerResponse = HTTPConvertHeaderToMap(resp.Header)
			httpStatus = resp.Status

			_ = resp.Body.Close()
		}
	}

	timeDiff = endDateTime.Sub(startDateTime).Milliseconds()

	// Put into mapResponse
	mapResponse["bodyRequest"] = bodyRequest
	mapResponse["headerRequest"] = headerRequest
	mapResponse["bodyResponse"] = bodyResponse
	mapResponse["headerResponse"] = headerResponse
	mapResponse["httpStatus"] = httpStatus
	mapResponse["startDateTime"] = startDateTime
	mapResponse["endDateTime"] = endDateTime
	mapResponse["timeDifference"] = timeDiff
	mapResponse["status"] = hitStatus

	DoLog("DEBUG", traceCode, "NetworkTool", "HTTPS_POST_String",
		fmt.Sprintf("mapResponse: %v", mapResponse), true, err)

	return mapResponse
}

// noinspection GoUnusedExportedFunction
func HTTPSGETString(traceCode string, mainUrl string, headerRequest map[string]interface{}, urlParameter map[string]string) map[string]interface{} {
	DoLog("DEBUG", traceCode, "NetworkTool", "HTTPS_GET_String", fmt.Sprintf("HTTPS GET - body: %v, header: %v", urlParameter, headerRequest), false, nil)
	var mapResponse = make(map[string]interface{})
	var headerResponse map[string]interface{}
	var bodyResponse string
	var startDateTime = time.Now()
	var endDateTime = time.Now()
	var timeDiff int64 // millisecond
	var httpStatus string
	var hitStatus = "000"

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Compose request with body
	startDateTime = time.Now()
	req, err := http.NewRequest("GET", mainUrl, nil)

	if req != nil {
		req.Close = true // Close once done
	}

	if err != nil {
		hitStatus = "240"

		DoLog("DEBUG", traceCode, "NetworkTool", "HTTPS_GET_String",
			"Failed to compose HTTP Request. Error occured.", true, err)
	} else {
		// Compose header
		if headerRequest != nil {
			for key, val := range headerRequest {
				req.Header.Set(key, val.(string))
			}
		}

		// Buat encoded url parameter
		query := url.Values{}

		if urlParameter != nil {
			for key, val := range urlParameter {
				query.Add(key, val)
			}
		}
		req.URL.RawQuery = query.Encode()

		DoLog("DEBUG", traceCode, "NetworkTool", "HTTPS_GET_String",
			fmt.Sprintf("Request URL: %v", req.URL.String()), false, nil)

		// Buat HTTP Client
		client := &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		}

		// Get the response
		resp, errR := client.Do(req)

		if errR != nil {
			hitStatus = "240"

			DoLog("DEBUG", traceCode, "NetworkTool", "HTTPS_GET_String",
				"Failed to execute HTTP Request. Error occured.", true, err)
		} else {
			endDateTime = time.Now()

			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)

			// fill bodyResponse, headerResponse, dll
			bodyResponse = string(body)
			headerResponse = HTTPConvertHeaderToMap(resp.Header)
			httpStatus = resp.Status
		}
	}

	timeDiff = endDateTime.Sub(startDateTime).Milliseconds()

	// Put into mapResponse
	mapResponse["bodyRequest"] = req.URL.RawQuery
	mapResponse["headerRequest"] = headerRequest
	mapResponse["bodyResponse"] = bodyResponse
	mapResponse["headerResponse"] = headerResponse
	mapResponse["httpStatus"] = httpStatus
	mapResponse["startDateTime"] = startDateTime
	mapResponse["endDateTime"] = endDateTime
	mapResponse["timeDifference"] = timeDiff
	mapResponse["status"] = hitStatus

	DoLog("DEBUG", traceCode, "NetworkTool", "HTTPS_GET_String",
		fmt.Sprintf("mapResponse: %v", mapResponse), true, err)

	return mapResponse
}

// Custom : Get AUth Header Param
func HTTPSGETStringAuth(traceCode string, mainUrl string, headerRequest map[string]interface{}, urlParameter map[string]string) map[string]interface{} {
	DoLog("DEBUG", traceCode, "NetworkTool", "HTTPS_GET_String", fmt.Sprintf("HTTP GET - body: %v, header: %v", urlParameter, headerRequest), false, nil)
	var mapResponse = make(map[string]interface{})
	var headerResponse map[string]interface{}
	var bodyResponse string
	var startDateTime = time.Now()
	var endDateTime = time.Now()
	var timeDiff int64 // millisecond
	var httpStatus string
	var hitStatus = "000"
	var realm, qop, nonce, opaque string

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Compose request with body
	startDateTime = time.Now()
	req, err := http.NewRequest("GET", mainUrl, nil)

	if err != nil {
		hitStatus = "240"

		DoLog("DEBUG", traceCode, "NetworkTool", "HTTPS_POST_String",
			"Failed to compose HTTP Request. Error occured.", true, err)
	} else {
		// Compose header
		if headerRequest != nil {
			for key, val := range headerRequest {
				req.Header.Set(key, val.(string))
			}
		}

		// Buat encoded url parameter
		query := url.Values{}

		if urlParameter != nil {
			for key, val := range urlParameter {
				query.Add(key, val)
			}
		}
		req.URL.RawQuery = query.Encode()

		// Buat HTTP Client
		client := &http.Client{Transport: tr}

		// Get the response
		resp, errR := client.Do(req)

		if errR != nil {
			hitStatus = "240"

			DoLog("DEBUG", traceCode, "NetworkTool", "HTTP_POST_String",
				"Failed to execute HTTP Request. Error occured.", true, err)
		} else {
			endDateTime = time.Now()

			body, _ := ioutil.ReadAll(resp.Body)

			// fill bodyResponse, headerResponse, dll
			bodyResponse = string(body)
			headerResponse = HTTPConvertHeaderToMap(resp.Header)
			httpStatus = resp.Status
			// Get String Auth
			header := resp.Header.Get("www-authenticate")
			parts := strings.SplitN(header, " ", 2)
			parts = strings.Split(parts[1], ",")
			fmt.Println("Parts: ", parts)
			opts := make(map[string]string)
			//newParts = strings.ReplaceAll(parts,`"`,"")
			i := 0
			//res := make(map[string]interface{})
			//temp := make([]string, 2)
			for _, part := range parts {
				fmt.Println("Part: ", part)
				vals := strings.SplitN(part, "=", 2)
				key := vals[0]
				val := strings.Trim(vals[1], "\",")
				opts[key] = val
				//trying split
				//temp = strings.Split(part,",")
				//res[temp[0]] = temp[1]
				i++
			}
			realm = GetStringFromMapString(opts, "realm")
			qop = GetStringFromMapString(opts, "qop")
			nonce = GetStringFromMapString(opts, "nonce")
			opaque = GetStringFromMapString(opts, "opaque")

			//fmt.Println("realm: ", realm)

			_ = resp.Body.Close()
		}
	}

	timeDiff = endDateTime.Sub(startDateTime).Milliseconds()

	// Put into mapResponse
	mapResponse["bodyRequest"] = req.URL.RawQuery
	mapResponse["headerRequest"] = headerRequest
	mapResponse["bodyResponse"] = bodyResponse
	mapResponse["headerResponse"] = headerResponse
	mapResponse["httpStatus"] = httpStatus
	mapResponse["startDateTime"] = startDateTime
	mapResponse["endDateTime"] = endDateTime
	mapResponse["timeDifference"] = timeDiff
	mapResponse["status"] = hitStatus
	mapResponse["realm"] = realm
	mapResponse["qop"] = qop
	mapResponse["nonce"] = nonce
	mapResponse["opaque"] = opaque

	DoLog("DEBUG", traceCode, "NetworkTool", "HTTPS_POST_String",
		fmt.Sprintf("mapResponse: %v", mapResponse), true, err)

	return mapResponse
}

// noinspection GoUnusedExportedFunction
func HTTPSPOSTForm(traceCode string, httpUrl string, headerRequest map[string]interface{}, mapFormContent map[string]string) map[string]interface{} {
	DoLog("DEBUG", traceCode, "NetworkTool", "HTTPS_POST_Form", fmt.Sprintf("HTTP POST - body: %s, header: %v", mapFormContent, headerRequest), false, nil)
	var mapResponse = make(map[string]interface{})
	var headerResponse map[string]interface{}
	var bodyResponse string
	var startDateTime = time.Now()
	var endDateTime = time.Now()
	var timeDiff int64 // millisecond
	var httpStatus string
	var hitStatus = "000"

	// make form data
	form := url.Values{}

	for key, val := range mapFormContent {
		form.Add(key, val)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Compose request with body
	startDateTime = time.Now()
	req, err := http.NewRequest("POST", httpUrl, strings.NewReader(form.Encode()))

	if err != nil {
		hitStatus = "240"

		DoLog("DEBUG", traceCode, "NetworkTool", "HTTPS_POST_String",
			"Failed to compose HTTP Request. Error occured.", true, err)
	} else {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// Compose header
		if headerRequest != nil {
			for key, val := range headerRequest {
				req.Header.Set(key, val.(string))
			}
		}

		client := &http.Client{Transport: tr}

		// Get the response
		resp, errR := client.Do(req)

		if errR != nil {
			hitStatus = "240"

			DoLog("DEBUG", traceCode, "NetworkTool", "HTTP_POST_String",
				"Failed to execute HTTP Request. Error occured.", true, err)
		} else {
			endDateTime = time.Now()

			body, _ := ioutil.ReadAll(resp.Body)

			// fill bodyResponse, headerResponse, dll
			bodyResponse = string(body)
			headerResponse = HTTPConvertHeaderToMap(resp.Header)
			httpStatus = resp.Status

			_ = resp.Body.Close()
		}
	}

	timeDiff = endDateTime.Sub(startDateTime).Milliseconds()

	// Put into mapResponse
	mapResponse["bodyRequest"] = fmt.Sprintf("%+v", mapFormContent)
	mapResponse["headerRequest"] = headerRequest
	mapResponse["bodyResponse"] = bodyResponse
	mapResponse["headerResponse"] = headerResponse
	mapResponse["httpStatus"] = httpStatus
	mapResponse["startDateTime"] = startDateTime
	mapResponse["endDateTime"] = endDateTime
	mapResponse["timeDifference"] = timeDiff
	mapResponse["status"] = hitStatus

	DoLog("DEBUG", traceCode, "NetworkTool", "HTTPS_POST_String",
		fmt.Sprintf("mapResponse: %v", mapResponse), true, err)

	return mapResponse
}

func HTTPPOSTStringFAST(traceCode string, url string, headerRequest map[string]interface{}, bodyRequest string) map[string]interface{} {
	var mapResponse = make(map[string]interface{})

	// Default response
	bodyResponse := ""
	headerResponse := make(map[string]interface{})
	httpStatus := ""
	opStatus := "000"

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.SetBody([]byte(bodyRequest))
	req.Header.SetMethod("POST") // Method

	// Set Headers
	if headerRequest != nil {
		for key, val := range headerRequest {
			req.Header.Set(key, val.(string))
		}
	}

	// Prepare response
	res := fasthttp.AcquireResponse()

	// Do hit
	startDateTime := time.Now()
	endDateTime := time.Now()
	var timeDiff int64 = 0
	if err := fasthttp.Do(req, res); err != nil {
		// Release request
		fasthttp.ReleaseRequest(req)

		DoLog("DEBUG", traceCode, "NetworkTool", "HTTP_POST_String",
			"Error while hitting HTTP Client. Error occur", true, err)

		opStatus = "201"
	} else {
		// Proses
		endDateTime = time.Now()
		timeDiff = endDateTime.Sub(startDateTime).Milliseconds()

		// Release request
		fasthttp.ReleaseRequest(req)

		bodyResponse = string(res.Body())
		origheaderResponse := res.Header
		httpStatus = strconv.Itoa(res.StatusCode())

		origheaderResponse.VisitAll(func(key, value []byte) {
			headerResponse["key"] = value
		})

		// Release response
		fasthttp.ReleaseResponse(res)

		opStatus = "000"
	}

	// Prepare the response.
	mapResponse["bodyRequest"] = bodyRequest
	mapResponse["headerRequest"] = headerRequest
	mapResponse["bodyResponse"] = bodyResponse
	mapResponse["headerResponse"] = headerResponse
	mapResponse["httpStatus"] = httpStatus
	mapResponse["startDateTime"] = startDateTime
	mapResponse["endDateTime"] = endDateTime
	mapResponse["timeDifference"] = timeDiff
	mapResponse["status"] = opStatus

	DoLog("DEBUG", traceCode, "NetworkTool", "HTTP_POST_String",
		fmt.Sprintf("mapResponse: %v", mapResponse), false, nil)

	return mapResponse
}

func HTTPSPOSTStringFast(traceCode string, url string, headerRequest map[string]interface{}, bodyRequest string) map[string]interface{} {
	DoLog("DEBUG", traceCode, "NetworkTool", "HTTPS_POST_String_FAST",
		fmt.Sprintf("HTTP POST - URL: %s, body: %s, header: %v", url, bodyRequest, headerRequest), false, nil)
	var mapResponse = make(map[string]interface{})
	var headerResponse = make(map[string]interface{})
	var bodyResponse string
	var startDateTime = time.Now()
	var endDateTime = time.Now()
	var timeDiff int64 // millisecond
	var httpStatus string
	var hitStatus = "000"

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.SetBody([]byte(bodyRequest))
	req.Header.SetMethod("POST") // Method

	// Set Headers
	if headerRequest != nil {
		for key, val := range headerRequest {
			req.Header.Set(key, val.(string))
		}
	}

	// Prepare response
	res := fasthttp.AcquireResponse()

	// Do hit
	if err := fasthttp.DoTimeout(req, res, 10*time.Second); err != nil {
		// Release request
		fasthttp.ReleaseRequest(req)

		DoLog("DEBUG", traceCode, "NetworkTool", "HTTPS_POST_String",
			"Error while hitting HTTP Client. Error occur", true, err)

		hitStatus = "901"
	} else {
		// Proses
		endDateTime = time.Now()
		timeDiff = endDateTime.Sub(startDateTime).Milliseconds()

		// Release request
		fasthttp.ReleaseRequest(req)

		bodyResponse = string(res.Body())
		origheaderResponse := res.Header
		httpStatus = strconv.Itoa(res.StatusCode())

		origheaderResponse.VisitAll(func(key, value []byte) {
			headerResponse[string(key)] = string(value)
		})

		// Release response
		fasthttp.ReleaseResponse(res)

		hitStatus = "000"
	}

	// Put into mapResponse
	mapResponse["bodyRequest"] = bodyRequest
	mapResponse["headerRequest"] = headerRequest
	mapResponse["bodyResponse"] = bodyResponse
	mapResponse["headerResponse"] = headerResponse
	mapResponse["httpStatus"] = httpStatus
	mapResponse["startDateTime"] = startDateTime
	mapResponse["endDateTime"] = endDateTime
	mapResponse["timeDifference"] = timeDiff
	mapResponse["status"] = hitStatus

	DoLog("DEBUG", traceCode, "NetworkTool", "HTTPS_POST_String",
		fmt.Sprintf("mapResponse: %v", mapResponse), false, nil)

	return mapResponse
}

func HTTPSGETStringFast(traceCode string, mainUrl string, headerRequest map[string]interface{}, urlParameter map[string]string) map[string]interface{} {
	DoLog("DEBUG", traceCode, "NetworkTool", "HTTPS_POST_String_FAST",
		fmt.Sprintf("HTTP POST - body: %v, header: %v", urlParameter, headerRequest), false, nil)
	var mapResponse = make(map[string]interface{})
	var headerResponse = make(map[string]interface{})
	var bodyResponse string
	var startDateTime = time.Now()
	var endDateTime = time.Now()
	var timeDiff int64 // millisecond
	var httpStatus string
	var hitStatus = "000"

	// Set complete GET URL
	completeURL := mainUrl
	count := 0
	for key := range urlParameter {
		if count == 0 {
			completeURL = completeURL + "?" + key + "=" + url.QueryEscape(urlParameter[key])
		} else {
			completeURL = completeURL + "&" + key + "=" + urlParameter[key]
		}

		count = count + 1
	}

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(mainUrl)
	req.Header.SetMethod("GET") // Method

	// Set Headers
	if headerRequest != nil {
		for key, val := range headerRequest {
			req.Header.Set(key, val.(string))
		}
	}

	// Prepare response
	res := fasthttp.AcquireResponse()

	// Do hit
	if err := fasthttp.Do(req, res); err != nil {
		// Release request
		fasthttp.ReleaseRequest(req)

		hitStatus = "901"
	} else {
		// Proses
		endDateTime = time.Now()
		timeDiff = endDateTime.Sub(startDateTime).Milliseconds()

		// Release request
		fasthttp.ReleaseRequest(req)

		bodyResponse = string(res.Body())
		origheaderResponse := res.Header
		httpStatus = strconv.Itoa(res.StatusCode())

		origheaderResponse.VisitAll(func(key, value []byte) {
			headerResponse[string(key)] = string(value)
		})

		// Release response
		fasthttp.ReleaseResponse(res)

		hitStatus = "000"
	}

	// Put into mapResponse
	mapResponse["bodyRequest"] = completeURL
	mapResponse["headerRequest"] = headerRequest
	mapResponse["bodyResponse"] = bodyResponse
	mapResponse["headerResponse"] = headerResponse
	mapResponse["httpStatus"] = httpStatus
	mapResponse["startDateTime"] = startDateTime
	mapResponse["endDateTime"] = endDateTime
	mapResponse["timeDifference"] = timeDiff
	mapResponse["status"] = hitStatus

	DoLog("DEBUG", traceCode, "NetworkTool", "HTTPS_POST_String",
		fmt.Sprintf("mapResponse: %v", mapResponse), false, nil)

	return mapResponse
}
