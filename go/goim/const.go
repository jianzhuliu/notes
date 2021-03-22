package goim

import "time"

//常量配置
const (
	//读取连接数据缓冲大小
	C_READ_BUF_SIZE = 1024

	//连接超时时间
	C_CONN_TIMEOUT = 300 * time.Second
)
