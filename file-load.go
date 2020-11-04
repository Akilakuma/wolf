package main

import (
	"io/ioutil"
	"log"
	"regexp"
)

// InitParser 拆解規則初始化
func InitParser() map[string][]byte {
	// 開.lua檔案
	files, err := ioutil.ReadDir("rule/")
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	// 驗證.lua副檔名
	data := make(map[string][]byte)
	re := regexp.MustCompile(`(.+)\.lua`)
	for _, f := range files {

		if !f.IsDir() {

			fs := re.FindStringSubmatch(f.Name())
			if len(fs) != 2 {
				continue
			}

			b, err := ioutil.ReadFile("rule/" + f.Name())
			if err != nil {
				continue
			}
			// key: 檔名, value: lua script
			data[fs[1]] = b
		}
	}

	return data
}
