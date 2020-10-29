package main

import "fmt"

//新定义的音乐播放接口
type MusicPlayer interface {
	play(fileType string, fileName string)
}

//历史老的接口
type ExistPlayer struct {
}

func (*ExistPlayer) playerMp3(fileName string) {
	fmt.Printf("play mp3: %s \n", fileName)
}

func (*ExistPlayer) playerWma(fileName string) {
	fmt.Printf("play wma: %s \n", fileName)
}

//适配器
type PlayerAdaptor struct {
	//包含老的的接口
	existPlayer ExistPlayer
}

func (player *PlayerAdaptor) play(fileType string, fileName string) {
	switch fileType {
	case "mp3":
		player.existPlayer.playerMp3(fileName)
	case "wma":
		player.existPlayer.playerWma(fileName)
	default:
		fmt.Println("不支持此类型文件播放")
	}
}

func main() {
	player := PlayerAdaptor{}
	player.play("mp3", "北京欢迎你")
	player.play("wma", "听雨")
	player.play("mp4", "小时代")
}
