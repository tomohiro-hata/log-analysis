package totalling

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tomohiro-hata/log-analysis/pkg/config"
)

func Log_to_csv(conf config.Config) {
	// 結果格納変数
	var result []Totalling
	// 月格納変数
	var dates Log_date
	// ログの個数分読み取り
	for _, name := range conf.Logs {
		// CSVファイルを開く
		file, err := os.Open(name)
		// エラーハンドリング
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		// reader作成
		reader := bufio.NewScanner(file)
		var rows [][]string
		var flag bool = true
		// 行読み出し(初期処理)
		for reader.Scan() {
			if flag {
				flag = false
				continue
			}
			row := strings.Split(reader.Text(), ",")
			for index, tmp := range row {
				row[index] = strings.Replace(tmp, "\"", "", -1)
			}
			rows = append(rows, row)
			// 結果変数が空の場合(追加)
			if len(result) == 0 {
				var tmp_result Totalling
				tmp_result.Username = row[2]
				result = append(result, tmp_result)
			} else {
				// ユーザ名が存在しない場合追加
				if username_find(row[2], result) {
					var tmp_result Totalling
					tmp_result.Username = row[2]
					result = append(result, tmp_result)
				}
			}
			// 月の収集
			tmp_sec, _ := strconv.ParseInt(row[0], 10, 64)
			tmp_date := time.Unix(tmp_sec/1000, 0)
			if len(dates.Dates) == 0 {
				dates.Dates = append(dates.Dates, tmp_date)
			} else {
				if date_find(tmp_date, dates) {
					dates.Dates = append(dates.Dates, tmp_date)
				}
			}
		}
		// カウント配列初期化
		for i := 0; i < len(result); i++ {
			for j := 0; j < len(dates.Dates); j++ {
				result[i].Result = append(result[i].Result, 0)
			}
		}
		// ログ集計(月別)
		for _, tmp_row := range rows {
			// Unix時間をtime型へ変換
			tmp_sec, _ := strconv.ParseInt(tmp_row[0], 10, 64)
			target_date := time.Unix(tmp_sec/1000, 0)
			// 加算Indexの捜索
			user_index := username_find_index(tmp_row[2], result)
			date_index := date_find_index(target_date, dates)
			// 配列内に存在するもののみ加算
			if user_index != -1 {
				if date_index != -1 {
					tmp_num, _ := strconv.Atoi(tmp_row[3])
					result[user_index].Result[date_index] += int32(tmp_num)
				}
			}
		}
	}
	// CSV出力用配列作成
	var records [][]string
	// ヘッダ用配列
	var tmp_record_date []string
	tmp_record_date = append(tmp_record_date, "")
	// ヘッダ作成
	for _, record_date := range dates.Dates {
		tmp_record_date = append(tmp_record_date, strconv.Itoa(record_date.Year())+"/"+record_date.Month().String())
	}
	records = append(records, tmp_record_date)
	// データ配列
	var tmp_record_data []string
	// データ作成
	for _, record_data := range result {
		tmp_record_data = nil
		tmp_record_data = append(tmp_record_data, record_data.Username)
		for _, record_num := range record_data.Result {
			tmp_record_data = append(tmp_record_data, strconv.Itoa(int(record_num)))
		}
		records = append(records, tmp_record_data)
	}

	// 書き込み先作成
	output, err := os.Create(conf.Target_path + "/" + conf.File_name)
	// エラーハンドリング
	if err != nil {
		fmt.Println(err)
	}
	defer output.Close()

	// ライター作成
	writer := csv.NewWriter(output)
	// 全書き込み
	writer.WriteAll(records)
}

// 存在確認関数
func username_find(username string, targets []Totalling) bool {
	for _, target := range targets {
		if target.Username == username {
			return false
		}
	}
	return true
}

// 存在確認関数
func date_find(tmp_time time.Time, targets Log_date) bool {
	for _, target := range targets.Dates {
		if target.Year() == tmp_time.Year() && target.Month() == tmp_time.Month() {
			return false
		}
	}
	return true
}

// 位置返却関数
func username_find_index(username string, targets []Totalling) int {
	for index, target := range targets {
		if target.Username == username {
			return index
		}
	}
	return -1
}

// 存在確認関数
func date_find_index(tmp_time time.Time, targets Log_date) int {
	for index, target := range targets.Dates {
		if target.Year() == tmp_time.Year() && target.Month() == tmp_time.Month() {
			return index
		}
	}
	return -1
}
