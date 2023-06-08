package config

// 設定ファイル格納用構造体
type Config struct {
	Target_path string `json:"target_path"`
	Logs []string `json:"logs"`
	File_name string `json:"file_name"`
}