package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/tomohiro-hata/log-analysis/pkg/config"
	"github.com/tomohiro-hata/log-analysis/pkg/totalling"
)

func main() {
	// 設定ファイルの読み込み
	file, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// 設定の読み込み
	var config config.Config
	json.Unmarshal(file, &config)
	fmt.Println(config)
	fmt.Println("start")
	totalling.Log_to_csv(config)
	fmt.Println("end")
}
