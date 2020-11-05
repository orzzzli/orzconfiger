package orzconfiger

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type ConfigerInfo struct {
	Path        string
	HotLoading  bool
	ScanSec     int
	Debug       bool
	CommentChar string
}

var ConfigerSection map[string]map[string]string
var ConfigerMap map[string]string
var ConfigerInfoObj ConfigerInfo

const DefaultScanSec = 60
const DefaultHotLoading = true
const DefaultDebug = false
const DefaultCommentChar = ";"

func InitConfiger(path string) error {
	ConfigerInfoObj.Debug = DefaultDebug
	ConfigerInfoObj.ScanSec = DefaultScanSec
	ConfigerInfoObj.HotLoading = DefaultHotLoading
	ConfigerInfoObj.CommentChar = DefaultCommentChar

	if path == "" {
		return errors.New("config path is empty")
	}

	exist,err := fileExists(path)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("config file is not exist")
	}
	ConfigerInfoObj.Path = path

	invokeConfigerObj()
	go hotLoadingConfiger()
	return nil
}

/*
	Get function return value and is ok.
*/
func GetString(section string, key string) (string,bool) {
	value, ok := ConfigerSection[section][key]
	return value, ok
}

/*
	Get function return int and is ok.
	If not exist return -1 and false.
	If convert err return 0 and false.
*/
func GetInt(section string, key string) (int,bool) {
	value, ok := ConfigerSection[section][key]
	if ok {
		valueInt,err := strconv.Atoi(value)
		if err == nil {
			return valueInt,true
		}
		return 0,false
	}
	return -1,false
}

/*
	Hot Loading config while file md5 changed.
*/
func hotLoadingConfiger() error {
	lastMD5 := ""
	for {
		if ConfigerInfoObj.HotLoading {
			file, err := os.Open(ConfigerInfoObj.Path)
			if err != nil {
				return errors.New("open config file err : " + err.Error())
			}
			md5Obj := md5.New()
			_, err = io.Copy(md5Obj, file)
			if err != nil {
				return errors.New("io copy file error : " + err.Error())
			}
			md5Str := hex.EncodeToString(md5Obj.Sum(nil))
			//first time
			if lastMD5 == "" {
				lastMD5 = md5Str
			} else if lastMD5 != md5Str { //config file changed
				invokeConfigerObj()
			}
			file.Close()
		}

		time.Sleep(time.Duration(DefaultScanSec) * time.Second)
	}
}

/*
	Check file is dir and check file exist.
*/
func fileExists(path string) (bool,error) {
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return false, errors.New("file path is a folder.")
		}
		return true,nil
	}
	if os.IsNotExist(err) {
		return false,nil
	}
	return false,nil
}

/*
	Read ini config and instantiate configer map.
*/
func invokeConfigerObj() {
	configFile := ConfigerInfoObj.Path

	file, err := os.Open(configFile)
	if err != nil {
		log.Fatal("open config file err : " + err.Error())
	}
	defer file.Close()

	ConfigerSection = make(map[string]map[string]string)
	ConfigerMap = make(map[string]string)

	reader := bufio.NewReader(file)
	lastSection := ""
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		lineStr := string(line)

		exist, first, start := processComment(lineStr)
		//whole line is comment
		if exist && first {
			continue
		}

		section, sectionExist := processSection(lineStr, start)
		//whole line is section
		if sectionExist {
			_, ok := ConfigerSection[section]
			//section already exist
			if ok {
				continue
			}
			ConfigerSection[section] = make(map[string]string)
			lastSection = section
		}

		key, value, ok := processLine(lineStr, start)
		if ok {
			ConfigerMap[key] = value
			//dont have section
			if lastSection == "" {
				continue
			}
			ConfigerSection[lastSection][key] = value
		}
	}
}

func processLine(context string, commentStart int) (key string, value string, ok bool) {
	start := strings.Index(context, "=")
	if start != -1 && start != 0 {
		//comment exist
		if commentStart != -1 {
			if commentStart < start {
				return "", "", false
			}
		}
		key = string([]rune(context)[:start])
		key = strings.Replace(key, " ", "", -1)
		if key == "" {
			return "", "", false
		}
		value = string([]rune(context)[start+1:])
		value = strings.Replace(value, " ", "", -1)
		if value == "" {
			return key, "", true
		}
		//check comment in value
		commentStart2 := strings.Index(value, DefaultCommentChar)
		if commentStart2 != -1 && commentStart2 != 0 { //comment in middle of value
			value = string([]rune(value)[:commentStart2])
			return key, value, true
		} else if commentStart2 == 0 { //comment at first of value
			return key, "", true
		} else { //comment not exist
			return key, value, true
		}
	} else {
		return "", "", false
	}
}

func processSection(context string, commentStart int) (str string, exist bool) {
	section := ""
	start := strings.Index(context, "[")
	end := strings.Index(context, "]")
	if start != -1 && end != -1 {
		//comment exist
		if commentStart != -1 {
			if commentStart < start {
				return section, false
			}
		}
		section = string([]rune(context)[start+1 : end])
		return section, true
	} else {
		return section, false
	}
}

func processComment(context string) (exist bool, isFirst bool, index int) {
	start := strings.Index(context, DefaultCommentChar)
	if start == 0 {
		return true, true, start
	} else if start == -1 {
		return false, false, start
	} else {
		return true, false, start
	}
}
