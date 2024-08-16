package modules

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gomodule/redigo/redis"
	guuid "github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Global Variable for MapConfig
var MapConfig = make(map[string]string)

// Global Logger for all
var zapLogger *zap.Logger

// Global RedisPooler for All
var RedisPooler *redis.Pool

func InitiateGlobalVariables(isProduction bool) {
	//if isProduction {
	//	//MapConfig = LoadConfigProduction()
	//	zapLogger = initiateZapLogger()
	//	RedisPooler = RedisInitiateRedisPool()
	//} else {
	// MapConfig = LoadConfig()
	//	zapLogger = initiateZapLogger()
	//	RedisPooler = RedisInitiateRedisPool()
	//}

	zapLogger = initiateZapLogger()
	RedisPooler = RedisInitiateRedisPool()
}

func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("Jan  2 15:04:05"))
}

func initiateZapLogger() *zap.Logger {
	// THE APPLICATION NEED TO RUN USING SUPERVISOR BECAUSE IT JUST LOGGING TO STDOUT ONLY.
	// SUPERVISOR WILL HANDLE THE LOGGING ---

	cfg := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: SyslogTimeEncoder,

			//CallerKey:    "caller",
			//EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, _ := cfg.Build()

	return logger
}

var (
	InfoColor  = Teal
	WarnColor  = Yellow
	ErrorColor = Red
)

var (
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}
func DoLog(logLevel string, messageId string, moduleName string, functionName string, message string, isError bool, theError error) {

	if logLevel == "DEBUG" {
		if isError == true {
			message = messageId + "   " + TrimStringLength(strings.ToUpper(moduleName), 10) + "   " + TrimStringLength(strings.ToUpper(functionName), 15) + "   " + message + "." + fmt.Sprintf("Error: %v", theError)
		} else {
			message = messageId + "   " + TrimStringLength(strings.ToUpper(moduleName), 10) + "   " + TrimStringLength(strings.ToUpper(functionName), 15) + "   " + message
		}

		zapLogger.Debug(message)
	} else if logLevel == "ERROR" {
		if isError == true {
			message = messageId + "   " + TrimStringLength(strings.ToUpper(moduleName), 10) + "   " + TrimStringLength(strings.ToUpper(functionName), 15) + "   " + message + "." + fmt.Sprintf("Error: %v", theError)
		} else {
			message = messageId + "   " + TrimStringLength(strings.ToUpper(moduleName), 10) + "   " + TrimStringLength(strings.ToUpper(functionName), 15) + "   " + message
		}

		//zapLogger.Error(message)
		zapLogger.Error(ErrorColor(message))
		//zapLogger.Debug(message)
	} else if logLevel == "WARNING" {
		if isError == true {
			message = messageId + "   " + TrimStringLength(strings.ToUpper(moduleName), 10) + "   " + TrimStringLength(strings.ToUpper(functionName), 15) + "   " + message + "." + fmt.Sprintf("Error: %v", theError)
		} else {
			message = messageId + "   " + TrimStringLength(strings.ToUpper(moduleName), 10) + "   " + TrimStringLength(strings.ToUpper(functionName), 15) + "   " + message
		}

		zapLogger.Warn(message)
		//zapLogger.Debug(message)
	} else if logLevel == "INFO" {
		if isError == true {
			message = messageId + "   " + TrimStringLength(strings.ToUpper(moduleName), 10) + "   " + TrimStringLength(strings.ToUpper(functionName), 15) + "   " + message + "." + fmt.Sprintf("Error: %v", theError)
		} else {
			message = messageId + "   " + TrimStringLength(strings.ToUpper(moduleName), 10) + "   " + TrimStringLength(strings.ToUpper(functionName), 15) + "   " + message
		}

		//zapLogger.Info(message)
		zapLogger.Info(InfoColor(message))
		//zapLogger.Debug(message)
	} else {
		if isError == true {
			message = messageId + "   " + TrimStringLength(strings.ToUpper(moduleName), 10) + "   " + TrimStringLength(strings.ToUpper(functionName), 15) + "   " + message + "." + fmt.Sprintf("Error: %v", theError)
		} else {
			message = messageId + "   " + TrimStringLength(strings.ToUpper(moduleName), 10) + "   " + TrimStringLength(strings.ToUpper(functionName), 15) + "   " + message
		}

		zapLogger.Info(message)
		//zapLogger.Debug(message)
	}
}

func TrimStringLength(theString string, length int) string {
	hasil := theString

	if len(theString) > length {
		hasil = fmt.Sprintf("%."+strconv.Itoa(length)+"s", theString)
	} else {
		hasil = theString

		for x := 0; x < length-len(theString); x++ {
			hasil = hasil + " "
		}
	}

	return hasil
}

func GenerateUUID() string {
	theGuuid := guuid.New().String()
	theGuuid = strings.Replace(theGuuid, "-", "", -1)

	return theGuuid
}

func ConvertMapStringToJSON(theMap map[string]string) string {
	jsonString, err := json.Marshal(theMap)

	if err != nil {
		return "{}"
	}

	return string(jsonString)
}

// noinspection GoUnusedExportedFunction
func ConvertMapInterfaceToJSON(theMap map[string]interface{}) string {
	jsonString, err := json.Marshal(theMap)

	if err != nil {
		return "{}"
	}

	return string(jsonString)
}

func ConvertInt64ToStringFixLength(emptySpaceFiller string, fillerPosition string, theExpectedLength int, theInteger int64) string {
	// fillerPosition = "RIGHT" or "LEFT"
	response := ""

	// Convert theInteger to str
	theStrInteger := strconv.FormatInt(theInteger, 10)

	// Check if original lenth of theInteger > expectedLength
	if len(theStrInteger) >= theExpectedLength {
		response = theStrInteger
	} else {
		theFiller := emptySpaceFiller
		for i := 0; i < theExpectedLength-len(theStrInteger); i++ {
			theFiller = theFiller + emptySpaceFiller
		}

		if fillerPosition == "RIGHT" {
			response = theStrInteger + theFiller
		} else {
			response = theFiller + theStrInteger
		}
	}

	return response
}

//goland:noinspection GoUnusedExportedFunction
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

//goland:noinspection GoUnusedExportedFunction
func AppendStringToSlice(slice []string, data ...string) []string {
	m := len(slice)
	n := m + len(data)
	if n > cap(slice) { // if necessary, reallocate
		// allocate double what's needed, for future growth.
		newSlice := make([]string, (n+1)*2)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0:n]
	copy(slice[m:n], data)
	return slice
}

func DoFormatDateTime(dateTimeFormat string, theTime time.Time) string {
	const (
		stdLongMonth   = "January"
		stdMonth       = "Jan"
		stdNumMonth    = "1"
		stdZeroMonth   = "01"
		stdLongWeekDay = "Monday"
		stdWeekDay     = "Mon"
		stdDay         = "2"
		//stdUnderDay       = "_2"
		stdZeroDay    = "02"
		stdHour       = "15"
		stdHour12     = "3"
		stdZeroHour12 = "03"
		stdMinute     = "4"
		stdZeroMinute = "04"
		stdSecond     = "5"
		stdZeroSecond = "05"
		stdLongYear   = "2006"
		stdYear       = "06"
		//stdPM             = "PM"
		//stdpm             = "pm"
		//stdTZ             = "MST"
		//stdISO8601TZ      = "Z0700"  // prints Z for UTC
		//stdISO8601ColonTZ = "Z07:00" // prints Z for UTC
		//stdNumTZ          = "-0700"  // always numeric
		//stdNumShortTZ     = "-07"    // always numeric
		//stdNumColonTZ     = "-07:00" // always numeric
		//stdMilliSecond    = "000"
	)

	theFormat := dateTimeFormat

	// replace YYYY with stdLongYear
	theFormat = strings.Replace(theFormat, "YYYY", stdLongYear, -1)

	// replace YY with stdYear
	theFormat = strings.Replace(theFormat, "YY", stdYear, -1)

	// replace MMMM with stdLongMonth - January
	theFormat = strings.Replace(theFormat, "MMMM", stdLongMonth, -1)

	// replace MM with stdMonth - Jan
	theFormat = strings.Replace(theFormat, "MM", stdMonth, -1)

	// replace 0M with zeroMonth - 01
	theFormat = strings.Replace(theFormat, "0M", stdZeroMonth, -1)

	// replace M with oneMonth - 1
	theFormat = strings.Replace(theFormat, "M", stdNumMonth, -1)

	// replace DDDD with stdLongWeekDay - Monday
	theFormat = strings.Replace(theFormat, "DDDD", stdLongWeekDay, -1)

	// replace DD with stdWeekDay - Mon
	theFormat = strings.Replace(theFormat, "DD", stdWeekDay, -1)

	// replace 0D with stdZeroDay - 01
	theFormat = strings.Replace(theFormat, "0D", stdZeroDay, -1)

	// replace D with stdNumDate - 1
	theFormat = strings.Replace(theFormat, "D", stdDay, -1)

	// replace HH with 24 hours hour
	theFormat = strings.Replace(theFormat, "HH", stdHour, -1)

	// replace 0H with 2 digits 12 hour
	theFormat = strings.Replace(theFormat, "0H", stdZeroHour12, -1)

	// replace H with num hour 12
	theFormat = strings.Replace(theFormat, "H", stdHour12, -1)

	// repalce mm with 2 digits minute start with 0
	theFormat = strings.Replace(theFormat, "mm", stdZeroMinute, -1)

	// replace m with number minute
	theFormat = strings.Replace(theFormat, "m", stdMinute, -1)

	// replace ss with 2 digits seconds
	theFormat = strings.Replace(theFormat, "ss", stdZeroSecond, -1)

	// replace s with number second
	theFormat = strings.Replace(theFormat, "s", stdSecond, -1)

	// replace S with millisecond after theFormat is implemented
	theReturn := theTime.Format(theFormat)
	//fmt.Println("theReturn before milisecond: " + theReturn)

	// Create the milist of current theTime
	milis := (theTime.UnixNano()) / 1000000

	nowMilis := milis - (theTime.Unix() * 1000)

	strMilis := ConvertInt64ToStringFixLength("0", "LEFT", 3, nowMilis)

	if len(strMilis) > 3 {
		theStrMilisRune := []rune(strMilis)
		strMilis = string(theStrMilisRune[:3])
	}

	// Replace S with millisecond
	theReturn = strings.Replace(theReturn, "S", strMilis, -1)

	return theReturn
}

func GetBeginOfDayDateTime(t time.Time) time.Time {
	year, month, day := t.Date()

	loc, _ := time.LoadLocation("Asia/Jakarta")

	return time.Date(year, month, day, 0, 0, 0, 0, loc)
}

func ConvertSQLNullStringToString(theNullStr sql.NullString) string {
	var response string

	if theNullStr.Valid == true {
		response = theNullStr.String
	} else {
		response = ""
	}

	return response
}

func ConvertSQLNullFloat64ToFloat64(theNullFloat64 sql.NullFloat64) float64 {
	var response float64

	if theNullFloat64.Valid == true {
		response = theNullFloat64.Float64
	} else {
		response = 0.00000
	}

	return response
}

func ConvertSQLNullBoolToBool(theNullBool sql.NullBool) bool {
	response := false

	if theNullBool.Valid == true {
		response = theNullBool.Bool
	} else {
		response = false
	}

	return response
}

func ConvertSQLNullTimeToTime(theNullTime sql.NullTime) time.Time {
	response := time.Now()

	if theNullTime.Valid == true {
		response = theNullTime.Time
	} else {
		response = time.Now()
	}

	return response
}

func ConvertJSONStringToMap(messageId string, theJSON string) map[string]interface{} {
	resultMap := make(map[string]interface{})

	err := json.Unmarshal([]byte(theJSON), &resultMap)

	if err != nil {
		log.Debugln(messageId + " ." + fmt.Sprintf("Failed to convert json to map for json content: %s", theJSON))
		resultMap = nil
	} else {
		log.Debugln(fmt.Sprintf("Success converting json %s to hashmap: %+v", theJSON, resultMap))
	}

	return resultMap
}

func GetStringFromMapInterface(theMap map[string]interface{}, theKey string) string {
	theValue := ""

	_, exist := theMap[theKey]

	if exist == true {
		theValue = strings.TrimSpace(theMap[theKey].(string))
	}

	return theValue
}

func ConvertMapInterfaceToXML(theMap map[string]interface{}) string {
	// theXML := `<?xml version="1.0" encoding="UTF-8"?>`
	theXML := ""

	for key, val := range theMap {
		fmt.Printf("key: %s - val %+v", key, val)

		theXML = theXML + `<` + key + `>`

		if reflect.ValueOf(val).Kind() == reflect.Map {
			fmt.Println("val is MAP")

			theSubXML := ""
			for key01, val01 := range val.(map[string]interface{}) {
				// Sub xml
				theSubXML = theSubXML + `<` + key01 + `>`

				if reflect.ValueOf(val01).Kind() == reflect.Map {
					fmt.Println("val01 is MAP")

					theBelowSubSML := ""

					for key02, val02 := range val01.(map[string]interface{}) {
						theBelowSubSML = theBelowSubSML + `<` + key02 + `>`

						if reflect.ValueOf(val02).Kind() == reflect.Int64 || reflect.ValueOf(val02).Kind() == reflect.Int {
							fmt.Println("val02 is INT")

							theBelowSubSML = theBelowSubSML + fmt.Sprintf("%d", val02)
						} else if reflect.ValueOf(val02).Kind() == reflect.Float64 {
							fmt.Println("val02 is FLOAT")

							theBelowSubSML = theBelowSubSML + fmt.Sprintf("%.2f", val02)
						} else if reflect.ValueOf(val02).Kind() == reflect.Bool {
							fmt.Println("val02 is BOOL")

							theBelowSubSML = theBelowSubSML + fmt.Sprintf("%t", val02)
						} else {
							fmt.Println("val02 is STRING (SHOULD BE)")

							theBelowSubSML = theBelowSubSML + val02.(string)
						}

						theBelowSubSML = theBelowSubSML + `</` + key02 + `>`
					}

					theSubXML = theSubXML + theBelowSubSML
				} else if reflect.ValueOf(val01).Kind() == reflect.Int64 || reflect.ValueOf(val01).Kind() == reflect.Int {
					fmt.Println("val01 is INT")

					theSubXML = theSubXML + fmt.Sprintf("%d", val01)
				} else if reflect.ValueOf(val).Kind() == reflect.Float64 {
					fmt.Println("val01 is FLOAT")

					theSubXML = theSubXML + fmt.Sprintf("%.2f", val01)
				} else if reflect.ValueOf(val).Kind() == reflect.Bool {
					fmt.Println("val01 is BOOL")

					theSubXML = theSubXML + fmt.Sprintf("%t", val01)
				} else {
					fmt.Println("val01 is STRING (SHOULD BE)")

					theSubXML = theSubXML + val01.(string)
				}

				theSubXML = theSubXML + `</` + key01 + `>`
			}

			theXML = theXML + theSubXML

		} else if reflect.ValueOf(val).Kind() == reflect.Int64 || reflect.ValueOf(val).Kind() == reflect.Int {
			fmt.Println("val is INT")

			theXML = theXML + fmt.Sprintf("%d", val)
		} else if reflect.ValueOf(val).Kind() == reflect.Float64 {
			fmt.Println("val is FLOAT")

			theXML = theXML + fmt.Sprintf("%.2f", val)
		} else {
			fmt.Println("val is STRING")

			theXML = theXML + val.(string)
		}

		theXML = theXML + `</` + key + `>`

		fmt.Print("theXML: " + theXML)
	}

	return theXML
}

func GetFloatFromMapInterface(theMap map[string]interface{}, theKey string) float64 {
	theValue := 0.00000000

	floatVal, ok := theMap[theKey].(float64)

	if !ok {
		// Bukan float64 Check apakah string
		stringVal, strOk := theMap[theKey].(string)
		if !strOk {
			// Bukan string, apakah int64
			intVal, intOk := theMap[theKey].(int64)

			if !intOk {
				// Bukan integer, apa dong ya?
				theValue = 0.00000000
			} else {
				// Integer64, convert it to float64
				theValue = float64(intVal)
			}
		} else {
			// String, convert it to float64
			theValueF, errStr := strconv.ParseFloat(stringVal, 64)

			if errStr == nil {
				theValue = theValueF
			}
		}
	} else {
		theValue = floatVal
	}

	return theValue
}

// noinspection GoUnusedExportedFunction
func GetStringFromMapString(theMap map[string]string, theKey string) string {
	theValue := ""

	_, exist := theMap[theKey]

	if exist == true {
		theValue = strings.TrimSpace(theMap[theKey])
	}

	return theValue
}

// noinspection GoUnusedExportedFunction
func ConvertMapInterfaceToMapString(mapInterface map[string]interface{}) map[string]string {
	mapHasil := make(map[string]string)

	for k := range mapInterface {
		if fmt.Sprintf("%T", reflect.TypeOf(mapInterface[k])) == "bool" {
			// variable is bool
			mapHasil[k] = strconv.FormatBool(mapInterface[k].(bool))
		} else if fmt.Sprintf("%T", reflect.TypeOf(mapInterface[k])) == "string" {
			// variable is string
			mapHasil[k] = mapInterface[k].(string)
		} else if fmt.Sprintf("%T", reflect.TypeOf(mapInterface[k])) == "int" {
			// variable is int
			mapHasil[k] = strconv.Itoa(mapInterface[k].(int))
		} else if fmt.Sprintf("%T", reflect.TypeOf(mapInterface[k])) == "float64" {
			// variable is float64
			mapHasil[k] = fmt.Sprintf("%.2f", mapInterface[k].(float64))
		} else if fmt.Sprintf("%T", reflect.TypeOf(mapInterface[k])) == "float32" {
			// variable is float64
			mapHasil[k] = fmt.Sprintf("%.2f", mapInterface[k].(float32))
		}
	}

	return mapHasil
}

func MapInterfaceHasKey(theMap map[string]interface{}, theKey string) bool {
	_, exist := theMap[theKey]

	return exist
}

// noinspection GoUnusedExportedFunction
func MapStringHasKey(theMap map[string]string, theKey string) bool {
	_, exist := theMap[theKey]

	return exist
}

func DoMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))

	return hex.EncodeToString(hasher.Sum(nil))
}

// noinspection GoUnusedExportedFunction
func GetStringInBetweenInsideBoundary(str string, start string, end string) (result string) {
	fmt.Println("GetStringInBetweenInsideBoundary - start: " + start + ", end: " + end + ", from string: " + str)
	s := strings.Index(str, start)
	if s == -1 {
		return ""
	}
	s += len(start)

	withoutStartStr := str[s:]
	fmt.Println("GetStringInBetweenInsideBoundary - withoutStartStr: " + withoutStartStr)
	e := strings.Index(withoutStartStr, end) // ambil batas akhir paling dalam
	fmt.Printf("GetStringInBetweenInsideBoundary = e: " + fmt.Sprintf("%d", e))
	if e == -1 {
		return ""
	}
	return withoutStartStr[0:e]
}

// noinspection GoUnusedExportedFunction
func GetStringInBetweenOutsideBoundary(str string, start string, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return ""
	}
	s += len(start)
	e := strings.LastIndex(str, end) // ambil batas akhir paling luar
	return str[s:e]
}

// noinspection GoUnusedExportedFunction
func GenerateRandomNumericString(strLength int) string {
	var letters = []rune("1234567890")

	rand.Seed(time.Now().UnixNano())
	b := make([]rune, strLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

// noinspection GoUnusedExportedFunction
func GenerateRandomAlphaNumericString(strLength int) string {
	var letters = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXZ")

	rand.Seed(time.Now().UnixNano())
	b := make([]rune, strLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

// noinspection GoUnusedExportedFunction
func GenerateRandomAlphabeticalString(strLength int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXZ")

	rand.Seed(time.Now().UnixNano())
	b := make([]rune, strLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func SHA1(traceCode string, theStr string) string {
	DoLog("DEBUG", traceCode, "PastiTool", "SHA1",
		"String to SHA1: "+theStr, false, nil)

	h := sha1.New()
	h.Write([]byte(theStr))
	bs := h.Sum(nil)
	sha1Result := hex.EncodeToString(bs)

	DoLog("DEBUG", traceCode, "PastiTool", "SHA1",
		"String to SHA1: "+theStr+" -> SHA1 result: "+sha1Result, false, nil)

	return sha1Result
}

func GetMD5HashWithSum(text string) string {
	DoLog("DEBUG", "", "PastiTool", "MD5",
		"String to MD5: "+text, false, nil)
	return fmt.Sprintf("%x", md5.Sum([]byte(text)))
}

func StringPaddding(input string, padLength int, padString string, padType string) string {
	var output string

	inputLength := len(input)
	padStringLength := len(padString)

	if inputLength >= padLength {
		return input
	}

	repeat := math.Ceil(float64(1) + (float64(padLength-padStringLength))/float64(padStringLength))

	switch padType {
	case "RIGHT":
		output = input + strings.Repeat(padString, int(repeat))
		output = output[:padLength]
	case "LEFT":
		output = strings.Repeat(padString, int(repeat)) + input
		output = output[len(output)-padLength:]
	case "BOTH":
		length := (float64(padLength - inputLength)) / float64(2)
		repeat = math.Ceil(length / float64(padStringLength))
		output = strings.Repeat(padString, int(repeat))[:int(math.Floor(float64(length)))] + input + strings.Repeat(padString, int(repeat))[:int(math.Ceil(float64(length)))]
	}

	return output
}

func HandlingGracefulShutdown(e *echo.Echo) {
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	fmt.Println("received shutdown command")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	fmt.Println("service down")
}
