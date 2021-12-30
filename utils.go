package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	ejs "encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
	progressbar "github.com/schollz/progressbar/v3"
)

func fileExists(filename string) bool {
  info, err := os.Stat(filename)
  if os.IsNotExist(err) {
    return false
  }
  return !info.IsDir()
}

func getHexTime(t time.Time) string {
	year := t.Format("2006")
	mounth := t.Format("01")
	day := t.Format("02")
	hour := t.Format("15")
	minute := t.Format("04")
	dateTime := []string{year, mounth, day, hour, minute}
	resultDate := ""
	for _, s := range dateTime {
		resultDate += stringToHex(s)
	}
	resultDate = strings.ToTitle(resultDate)
	return resultDate
}

func hexToTime(s string) (time.Time, error) {
	if len(s) < 12 {
		return time.Time{}, errors.New("Hex string is too short - can't decode time. " + s)
	}
	yearb, err := hex.DecodeString(s[0:4])
	mounthb, err := hex.DecodeString(s[4:6])
	dayb, err := hex.DecodeString(s[6:8])
	hourb, err := hex.DecodeString(s[8:10])
	minuteb, err := hex.DecodeString(s[10:12])
	if err != nil {
		log.Println(err.Error())
		return time.Time{}, err
	}
	year := int(yearb[0])*256 + int(yearb[1])
	mounth := time.Month(mounthb[0])
	day := int(dayb[0])
	hour := int(hourb[0])
	minute := int(minuteb[0])
	return time.Date(year, mounth, day, hour, minute, 0, 0, time.UTC), err
}

func stringToHex(s string) string {
	number, _ := strconv.Atoi(s)
	hexString := strconv.FormatInt(int64(number), 16)
	if len(hexString)%2 != 0 {
		hexString = "0" + hexString
	}
	return hexString
}

func getInt(value string) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		log.Panicln(err)
	}
	return i
}

func getFloat(value string) float64 {
	i, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Panicln(err)
	}
	return i
}

func getByte(value string) byte {
	i, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		log.Panicln(err)
	}
	return byte(i)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func size(value interface{}) uint8 {
	switch value.(type) {
	case int8:
	case uint8:
		return 1
	case int16:
	case uint16:
		return 2
	case int32:
	case uint32:
		return 4
	case int64:
	case uint64:
		return 8
	case []byte:
		original, ok := value.([]byte)
		if ok {
			return uint8(len(original))
		}
	}
	panic("unknown type in size estimation")
}

func loadPNG(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return png.Decode(file)
}

func loadJPEG(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return jpeg.Decode(file)
}

func now() time.Time {
	loc, _ := time.LoadLocation("UTC")
	return time.Now().In(loc)
}

func loadJPEGfromB64(b64text string) (image.Image, error) {
	data, err := base64.StdEncoding.DecodeString(b64text)
	if err != nil {
		return nil, err
	}
	return jpeg.Decode(bytes.NewReader(data))
}

func saveJPEGinB64(picture image.Image) (string, error) {
	opt := jpeg.Options{
		Quality: 90,
	}
	var data bytes.Buffer
	if err := jpeg.Encode(&data, picture, &opt); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data.Bytes()), nil
}

func loadFile2B64(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func savePNG(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, img)
}

type callable func(...interface{}) error

func run(procedure callable, timeout time.Duration, args ...interface{}) error {
	channel := make(chan error, 1)

	go func() {
		err := procedure(args...)
		channel <- err
	}()

	// Listen on our channel AND a timeout channel - which ever happens first.
	select {
	case res := <-channel:
		return res
	case <-time.After(timeout):
		return fmt.Errorf("PROCEDURE CALL TIMED OUT")
	}
}

func runWithBar(name string, procedure callable, timeout time.Duration, args ...interface{}) error {
	channel := make(chan error, 1)
	var bar *progressbar.ProgressBar
	if name != "" {
		bar = progressbar.DefaultBytes(
			timeout.Milliseconds(),
			name,
		)
		go func() {
			time.Sleep(50 * time.Millisecond)
			bar.Add(50)
		}()
	}

	go func() {
		err := procedure(args...)
		channel <- err
	}()

	// Listen on our channel AND a timeout channel - which ever happens first.
	select {
	case res := <-channel:
		if name != "" {
			bar.Finish()
			log.Println("")
		}
		return res
	case <-time.After(timeout):
		return fmt.Errorf("PROCEDURE CALL TIMED OUT")
	}
}


// return the source filename after the last slash
func chopPath(original string) string {
	i := strings.LastIndex(original, "/")
	if i == -1 {
		return original
	} else {
		return original[i+1:]
	}
}

// return a string containing the file name, function name
// and the line number of a specified entry on the call stack
func FileLine(depthList ...int) string {
	var depth int
	if depthList == nil {
		depth = 1
	} else {
		depth = depthList[0]
	}
	function, file, line, _ := runtime.Caller(depth)
	return fmt.Sprintf("%s:%d => %s()", chopPath(file), line, runtime.FuncForPC(function).Name())
}

func makeError(err error, fileline string) error {
	return fmt.Errorf("ERROR: %s => %s", fileline, err.Error())
}

func make_timestamp(year, month, day, hour, minute, second, hundredth uint8) time.Time {
	timestamp := time.Date(int(year), time.Month(month), int(day), int(hour), int(minute),
		int(second), int(hundredth)*10000, time.UTC)
	return timestamp
}

func parse(buffer *[]byte, from, size int64, data interface{}) string {
	if from < 0 || size <= 0 {
		panic("from and size must be positive")
	}
	var chunk = make([]byte, size)
	copy(chunk, (*buffer)[from:from+size])
	reader := bytes.NewReader(chunk)
	err := binary.Read(reader, binary.BigEndian, data)
	check(err)
	return fmt.Sprintf("parse [%d, %d]: %s", from, from+size, hex.Dump(chunk))
}

func binary_insert(buffer *[]byte, from int, data interface{}) []byte {
	if from < 0 {
		panic("from must be positive")
	}

	bsize := int(size(data))

	var chunk = make([]byte, from)
	if from > 0 {
		copy(chunk, (*buffer)[0:from])
	}
	buf := bytes.NewBuffer(chunk)

	var bs []byte
	switch data.(type) {
	case int8:
	case uint8:
		bs = make([]byte, 1)
		bs[0] = data.(uint8)
	case int16:
	case uint16:
		bs = make([]byte, 2)
		binary.BigEndian.PutUint16(bs, data.(uint16))
	case int32:
	case uint32:
		bs = make([]byte, 4)
		binary.BigEndian.PutUint32(bs, data.(uint32))
	case int64:
	case uint64:
		bs = make([]byte, 8)
		binary.BigEndian.PutUint64(bs, data.(uint64))
	case []byte:
		bs, _ = data.([]byte)
	}

	buf.Write(bs)

	bufsize := (len(*buffer) - from)
	if bufsize > 0 {
		var chunk2 = make([]byte, bufsize)
		copy(chunk2, (*buffer)[from:])
		buf.Write(chunk2)
	}

	fmt.Printf("BINARY INSERT [%d, %d]: %s", from, from+bsize, hex.Dump(buf.Bytes()))
	return buf.Bytes()
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func asInteger(result interface{}) int {
	value64, ok := result.(uint64)
	if ok != true {
		log.Printf("ERROR: asInteger() failed! => %v", result)
	}
	return int(value64)
}

func asByteArray(result interface{}) []byte {
	array, ok := result.([]byte)
	if ok != true {
		log.Printf("ERROR: asByteArray() failed! => %v", result)
	}
	return array
}

func asIntegerArray(data string) []int {
	var array []int
	err := ejs.Unmarshal([]byte(data), &array)
	if err != nil {
		log.Printf("ERROR: asIntegerArray() failed! => %s", err.Error())
		return array
	}
	return array
}

func asStringsMap(data string) map[string]string {
	sm := map[string]string{}
	var base interface{}
	err := ejs.Unmarshal([]byte(data), &base)
	if err != nil {
		log.Printf("ERROR: asStringsMap() failed! => %s", err.Error())
		return sm
	}
	values := base.(map[string]interface{})
	for k, v := range values {
		switch v := v.(type) {
		case string:
			sm[k] = v
		}
	}
	return sm
}

func asIntegersMap(data string) map[string]int {
	sm := map[string]int{}
	var base interface{}
	err := ejs.Unmarshal([]byte(data), &base)
	if err != nil {
		log.Printf("ERROR: asIntegersMap() failed! => %s", err.Error())
		return sm
	}
	values := base.(map[string]interface{})
	for k, v := range values {
		switch v := v.(type) {
		case float64:
			sm[k] = int(v)
		}
	}
	return sm
}

func asString(result interface{}) string {
	return fmt.Sprintf("%q", result)
}

func asString2(result interface{}) string {
	return fmt.Sprintf("%v", result)
}

func IntInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func StrInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
