package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	//时间格式化配置
	C_time_layout string = "2006-01-02 15:04:05.000"

	C_machine_id_bits uint8 = 10 //10 位，机器编码，支持 2 ^ 10 = 1024 台
	C_sequence_bits   uint8 = 14 //14 位，存储序列编码，支持  2 ^ 14 = 16384
	C_time_bits       uint8 = 39 //39 位，存储时间戳，精确到毫秒，支持 2 ^ 39 /1000 * 365 * 24 * 3600 ~= 17 年

	C_machine_id_max = uint16(1<<C_machine_id_bits - 1)
	C_sequence_max   = uint16(1<<C_sequence_bits - 1)

	//时间精确单位毫秒，相对于纳秒，
	C_time_unit = 1e6
)

//默认开始时间
var default_start_time = time.Date(2020, 12, 11, 0, 0, 0, 0, time.UTC)

//1毫秒内可以支持16384个序列号，qps 可达到 1 千万
//分布式id雪花算法结构体
type SnowFlake struct {
	mu           sync.Mutex
	machineID    uint16 //机器id
	sequence     uint16 //序列号id
	startTime    int64  //开始时间，保留毫秒级别
	intervalTime int64  //间隔时间
}

func NewSnowFlake(machineID uint16) *SnowFlake {
	if machineID < 1 || C_machine_id_max < machineID {
		panic(fmt.Sprintf("machineID is out of range 1 ~ %d", C_machine_id_max))
	}

	return &SnowFlake{
		machineID: machineID,
		sequence:  0,
		startTime: toSnowFlakeTime(default_start_time),
	}
}

//获取序列号id
func (s *SnowFlake) NextId() (uint64, error) {
	//加锁，并发安全控制
	s.mu.Lock()
	defer s.mu.Unlock()

	curIntervalTime := curIntervalTime(s.startTime)
	if s.intervalTime < curIntervalTime {
		s.sequence = 0
		s.intervalTime = curIntervalTime
	} else {
		s.sequence = (s.sequence + 1) & C_sequence_max
		//毫秒内，数量已经达到上限，或者系统时间缩小了
		if s.sequence == 0 || curIntervalTime < s.intervalTime {
			s.intervalTime++
			overtime := s.intervalTime - curIntervalTime
			time.Sleep(sleepTime(overtime))
		}
	}

	return s.genId()
}

//生成 id
func (s *SnowFlake) genId() (uint64, error) {
	if s.intervalTime >= 1<<C_time_bits {
		return 0, fmt.Errorf("time out of range")
	}

	return uint64(s.intervalTime)<<(C_machine_id_bits+C_sequence_bits) | uint64(s.machineID)<<C_sequence_bits | uint64(s.sequence), nil
}

//解析id
func Parse(id uint64) map[string]interface{} {
	//1位符号 + 39位时间 + 10位机器 + 14位序列号
	intervalTime := id >> (C_machine_id_bits + C_sequence_bits)
	machineID := id & (uint64(C_machine_id_max) << C_sequence_bits) >> C_sequence_bits
	sequence := id & uint64(C_sequence_max)

	return map[string]interface{}{
		"id":           id,
		"intervalTime": intervalTime,
		"machineID":    machineID,
		"sequence":     sequence,
		"created":      default_start_time.Add(time.Duration(intervalTime * C_time_unit)).Local().Format(C_time_layout),
	}
}

//当前时间间隔
func curIntervalTime(startTime int64) int64 {
	return toSnowFlakeTime(time.Now()) - startTime
}

//转换为指定单位时间
func toSnowFlakeTime(t time.Time) int64 {
	return t.UTC().UnixNano() / C_time_unit
}

//获取间隔时间，对应的 sleep 时间
func sleepTime(overtime int64) time.Duration {
	return time.Duration(overtime)*time.Millisecond -
		time.Duration(time.Now().UTC().UnixNano()%C_time_unit)*time.Nanosecond
}

func main() {
	//fmt.Println(time.Now().Format(C_time_layout))

	var wg sync.WaitGroup
	s := NewSnowFlake(3)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		if i == 8 {
			//模拟下一毫秒
			time.Sleep(1 * time.Millisecond)
		}
		go func(i int) {
			defer wg.Done()
			id, err := s.NextId()
			if err != nil {
				fmt.Println(i, err)
				return
			}

			fmt.Println("\r\n========", i, Parse(id))
		}(i)
	}

	wg.Wait()
}
