package main

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"keybd_event"
	"lingwei/letsgo"
	"os"
	"time"

	"github.com/andlabs/ui"

	_ "github.com/andlabs/ui/winmanifest"
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
	hn := letsgo.GetHardwareNo()
	config, err := parseConfig()
	err = ui.Main(func() {
		if hn == "" {
			showError("无法获取机器码")
			return
		}
		if err != nil {
			showError("无法获取配置信息")
			return
		}

		window := ui.NewWindow("IG牛逼", 600, 280, true)
		vbox := ui.NewVerticalBox()

		vbox.SetPadded(true)

		hbox := ui.NewHorizontalBox()
		hbox.SetPadded(true)
		vbox.Append(hbox, false)
		vbox.Append(ui.NewLabel("机器号："+hn), false)
		serialForm := ui.NewForm()
		serialForm.SetPadded(true)
		serialEntry := ui.NewEntry()
		if config.Key != "" {
			serialEntry.SetText(config.Key)
		}
		serialForm.Append("序列号", serialEntry, false)
		hbox.Append(serialForm, false)

		btnSer := ui.NewButton("设置序列号")
		btnSer.OnClicked(func(*ui.Button) {
			config.Key = serialEntry.Text()
			err := saveConfig(config)
			if err != nil {
				ui.MsgBox(window, "错误", "设置序列号失败："+err.Error())
			} else {
				ui.MsgBox(window, "提示", "设置序列号成功")
			}
		})

		hbox.Append(btnSer, false)
		vbox.Append(ui.NewLabel("使用说明："), false)
		vbox.Append(ui.NewLabel("1. 先设置屏幕分辨率为1920 x 1080"), false)
		vbox.Append(ui.NewLabel("2. 游戏显示设置为全屏，中等画质，界面缩放1.0"), false)
		vbox.Append(ui.NewLabel("3. 找一处有水的地方，进入钓鱼模式，装上鱼饵"), false)
		vbox.Append(ui.NewLabel("4. 点击开始，在3秒内切换回游戏"), false)
		vbox.Append(ui.NewLabel("5. 挂机别动鼠标键盘，别切换窗口，可关显示器。Enjoy!"), false)
		btnStart := ui.NewButton("开始")
		btnStart.OnClicked(func(*ui.Button) {
			baes, err := letsgo.EncryptAES([]byte(hn), []byte(aeskey))
			if err != nil {
				ui.MsgBox(window, "错误", "序列号错误："+err.Error())
			}
			mySeria := hex.EncodeToString(baes)[7:17]
			if mySeria == config.Key {
				go func() {
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
				}()
			} else {
				ui.MsgBox(window, "错误", "序列号错误，请联系作者获取。")
			}
		})

		vbox.Append(btnStart, false)

		window.SetChild(vbox)

		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.SetMargined(true)
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}

func saveConfig(config *configeration) error {
	jb, err := json.MarshalIndent(config, "", "	")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("config.ini", jb, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func parseConfig() (*configeration, error) {
	b, err := ioutil.ReadFile("config.ini")
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

func showError(msg string) {
	w := ui.NewWindow("错误", 100, 100, true)
	v := ui.NewVerticalBox()
	v.Append(ui.NewLabel(msg), false)
	b := ui.NewButton("关闭")
	v.Append(b, false)
	b.OnClicked(func(*ui.Button) {
		ui.Quit()
	})
	w.SetChild(v)
	w.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	w.SetMargined(true)
	w.Show()
}
