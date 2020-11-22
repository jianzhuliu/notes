package command

import (
	"fmt"
	"testing"
	"time"
	"unsafe"
)

//测试 kind sizeof
func TestSizeOf(t *testing.T) {
	var (
		v_int   int
		v_int8  int8
		v_int16 int16
		v_int32 int32

		v_uint   uint
		v_uint8  uint8
		v_uint16 uint16
		v_uint32 uint32

		v_float32 float32
		v_float64 float64

		v_time time.Time

		v_string string

		v_byte_slice  []byte
		v_uint8_slice []uint8
	)

	fmt.Printf("%-20s = %v \n", "sizeof(int)", unsafe.Sizeof(v_int))
	fmt.Printf("%-20s = %v \n", "sizeof(int8)", unsafe.Sizeof(v_int8))
	fmt.Printf("%-20s = %v \n", "sizeof(int16)", unsafe.Sizeof(v_int16))
	fmt.Printf("%-20s = %v \n", "sizeof(int32)", unsafe.Sizeof(v_int32))

	fmt.Printf("%-20s = %v \n", "sizeof(uint)", unsafe.Sizeof(v_uint))
	fmt.Printf("%-20s = %v \n", "sizeof(uint8)", unsafe.Sizeof(v_uint8))
	fmt.Printf("%-20s = %v \n", "sizeof(uint16)", unsafe.Sizeof(v_uint16))
	fmt.Printf("%-20s = %v \n", "sizeof(uint32)", unsafe.Sizeof(v_uint32))

	fmt.Printf("%-20s = %v \n", "sizeof(float32)", unsafe.Sizeof(v_float32))
	fmt.Printf("%-20s = %v \n", "sizeof(float64)", unsafe.Sizeof(v_float64))

	fmt.Printf("%-20s = %v \n", "sizeof(time.Time)", unsafe.Sizeof(v_time))

	fmt.Printf("%-20s = %v \n", "sizeof(string)", unsafe.Sizeof(v_string))

	fmt.Printf("%-20s = %v \n", "sizeof([]byte)", unsafe.Sizeof(v_byte_slice))
	fmt.Printf("%-20s = %v \n", "sizeof([]uint8)", unsafe.Sizeof(v_uint8_slice))
}
