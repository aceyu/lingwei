package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"keybd_event"
	"lingwei/letsgo"
	"strconv"
	"time"
)

const (
	aeskey = "HIgtcdRUxqT72582"
)

type configeration struct {
	Round     int     `json:"round"`
	StartX    int     `json:"startX"`
	StartY    int     `json:"startY"`
	Rmin      int     `json:"rmin"`
	Rmax      int     `json:"rmax"`
	Gmin      int     `json:"gmin"`
	Gmax      int     `json:"gmax"`
	Bmin      int     `json:"bmin"`
	Bmax      int     `json:"bmax"`
	RoundTime float32 `json:"roundTime"`
	Waiting   int     `json:"waiting"`
	Interval  int     `json:"interval"`
	Key       string  `json:"key"`
}

func main() {
	fmt.Println("=======================================")
	fmt.Println("==                                   ==")
	fmt.Println("==               Welcome             ==")
	fmt.Println("==                                   ==")
	fmt.Println("=======================================")
	anykey := ""
	hn := letsgo.GetHardwareNo()
	if hn == "" {
		fmt.Println("无法获取机器码，按任意键退出...")
		fmt.Scanln(&anykey)
		return

	}
	config, err := parseConfig()
	if err != nil {
		fmt.Println("无法获取配置信息，按任意键退出...")
		fmt.Scanln(&anykey)
		return
	}

	baes, err := letsgo.EncryptAES([]byte(hn), []byte(aeskey))
	mySeria := hex.EncodeToString(baes)[7:17]
	if err != nil {
		fmt.Println("序列号错误，按任意键退出...")
		fmt.Scanln(&anykey)
		return
	}
	if mySeria != config.Key {
		fmt.Println("机器码：" + hn)
		fmt.Println("序列号错误，请联系作者获取。")
		fmt.Println("获取后，请用记事本打开config.ini，将倒数第二行的thisismykey替换为真正的序列号")
		fmt.Println("按任意键退出...")
		fmt.Scanln(&anykey)
		return
	}

	fmt.Println("使用说明：")
	fmt.Println("")
	fmt.Println("1. 设置屏幕分辨率为1920 x 1080")
	fmt.Println("2. 游戏显示设置为全屏，中等画质，界面缩放1.0")
	fmt.Println("3. 确保包裹有足够空间，鱼饵充足")
	fmt.Println("4. 找一处有水的地方，进入钓鱼模式，装上鱼饵")
	fmt.Println("5. 点击开始，在3秒内切换回游戏")
	fmt.Println("6. 挂机别动鼠标键盘，别切换窗口，可关显示器")
	fmt.Println("")
	fmt.Print("设置钓鱼次数（直接回车表示无限次）：")
	timesstr := ""
	fmt.Scanln(&timesstr)
	times, _ := strconv.Atoi(timesstr)
	if times == 0 {
		fmt.Println("无限次数，点击任意键开始...")
	} else {
		fmt.Println("执行" + timesstr + "次，点击任意键开始...")
	}

	fmt.Scanln(&anykey)
	fmt.Println("Enjoy!")
	time.Sleep(10 * time.Second)
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}
	kb.SetKeys(keybd_event.VK_A)
	err = kb.Launching()
	kb.Launching()
	kb.Launching()
	kb.Launching()
	if err != nil {
		panic(err)
	}
	fmt.Println("按任意键退出...")
	fmt.Scanln(&anykey)
}

func parseConfig() (*configeration, error) {
	b, err := ioutil.ReadFile("./config.ini")
	if err != nil {
		return nil, err
	}
	var con configeration
	err = json.Unmarshal(b, &con)
	if err != nil {
		return nil, err
	}
	return &con, nil
}
