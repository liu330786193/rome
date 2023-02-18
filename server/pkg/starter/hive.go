package starter

import (
	"fmt"
	"github.com/beltran/gohive"
)

func connHive() {
	// 设置hive客户端
	config := gohive.NewConnectConfiguration()
	config.Username = "hive"
	conn, err := gohive.Connect("172.26.160.120", 10000, "NONE", config)
	if err != nil {
		fmt.Errorf("Connect error %v", err)
	}
	cursor := conn.Cursor()
	context := Contex
	cursor.Exec()
	if err != nil {
		fmt.Errorf("Exec error: %v", err)
	}
	var record string
	for rs.Next() {
		rs.Scan(&record)
		fmt.Println(record)
	}
	conn.Close()
}
