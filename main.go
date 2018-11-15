package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"lingwei/letsgo"
	"strconv"
	"time"
)

const (
	aeskey = "HIgtcdRUxqT72582"
)

func main() {
	fmt.Println("=======================================")
	fmt.Println("==                                   ==")
	fmt.Println("==               Welcome             ==")
	fmt.Println("==                                   ==")
	fmt.Println("=======================================")
	anykey := ""
	hn := letsgo.GetHardwareNo()
	if hn == "" {
		fmt.Println("无法获取机器码，按回车键退出...")
		fmt.Scanln(&anykey)
		return

	}
	config, err := parseConfig()
	if err != nil {
		fmt.Println("无法获取配置信息，按回车键退出...")
		fmt.Scanln(&anykey)
		return
	}

	baes, err := letsgo.EncryptAES([]byte(hn), []byte(aeskey))
	mySeria := hex.EncodeToString(baes)[7:17]
	if err != nil {
		fmt.Println("序列号错误，按回车键退出...")
		fmt.Scanln(&anykey)
		return
	}
	if mySeria != config.Key {
		fmt.Println("机器码：" + hn)
		fmt.Println("序列号错误，请联系作者获取。")
		fmt.Print("请输入新的序列号：")
		var newSeria string
		fmt.Scanln(&newSeria)
		if newSeria == mySeria {
			config.Key = newSeria
			err := saveConfig(config)
			if err != nil {
				fmt.Println("更新序列号失败...")
			} else {
				fmt.Println("更新序列号成功...")
			}
		} else {
			fmt.Println("序列号错误，按回车键退出...")
			fmt.Scanln(&anykey)
			return
		}
	}
	wt := 5
	fmt.Println("使用说明：")
	fmt.Println("")
	fmt.Println("1. 设置屏幕分辨率为1920 x 1080")
	fmt.Println("2. 游戏显示设置为全屏，标准及以上画质，界面缩放1.0")
	fmt.Println("3. 确保包裹有足够空间，鱼饵充足")
	fmt.Println("4. 找一处有水的地方，进入钓鱼模式，装上鱼饵")
	fmt.Println(fmt.Sprintf("5. 点击开始，在%d秒内切换回游戏", wt))
	fmt.Println("6. 挂机别动鼠标键盘，别切换窗口，可关显示器")
	fmt.Println("")
	fmt.Print("设置钓鱼次数（直接回车表示无限次）：")
	timesstr := ""
	fmt.Scanln(&timesstr)
	times, _ := strconv.Atoi(timesstr)
	if times <= 0 {
		fmt.Println("无限次数，按回车键开始...")
	} else {
		fmt.Println("执行" + timesstr + "次，按回车键开始...")
	}
	fmt.Scanln(&anykey)
	fmt.Println("Enjoy!")
	for i := wt; i > 0; i-- {
		fmt.Println(strconv.Itoa(i))
		time.Sleep(1 * time.Second)
	}
	fish := letsgo.NewFish(*config, times)
	fish.Launch()
	fmt.Println("按回车键退出...")
	fmt.Scanln(&anykey)
}

func parseConfig() (*letsgo.Configeration, error) {
	b, err := ioutil.ReadFile("./config.ini")
	if err != nil {
		return nil, err
	}
	var con letsgo.Configeration
	err = json.Unmarshal(b, &con)
	if err != nil {
		return nil, err
	}
	return &con, nil
}

func saveConfig(config *letsgo.Configeration) error {
	jb, err := json.MarshalIndent(config, "", "	")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("./config.ini", jb, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
