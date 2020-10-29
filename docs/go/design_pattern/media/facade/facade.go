package main

import "fmt"

//统一接口,对外提供一个公用方法
type Api struct {
	cpu    CpuApi
	memory MemoryApi
	disk   DiskApi
}

//对外使用这一个方法创建对象
func NewApi() Api {
	return Api{
		cpu:    NewCpuApi(),
		memory: NewMemoryApi(),
		disk:   NewDiskApi(),
	}
}

//统一调用子系统各个方法
func (api Api) Start() {
	api.cpu.Start()
	api.memory.Start()
	api.disk.Start()
}

//cpu 子系统
type CpuApi struct {
}

func NewCpuApi() CpuApi {
	return CpuApi{}
}

func (cpu CpuApi) Start() {
	fmt.Println("cpu starting")
}

//memory 子系统
type MemoryApi struct {
}

func NewMemoryApi() MemoryApi {
	return MemoryApi{}
}

func (memory MemoryApi) Start() {
	fmt.Println("memory starting")
}

//disk 子系统
type DiskApi struct {
}

func NewDiskApi() DiskApi {
	return DiskApi{}
}

func (disk DiskApi) Start() {
	fmt.Println("disk starting")
}

func main() {
	api := NewApi()
	api.Start()
}
