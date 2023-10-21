package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

	wail_runtime "github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/mgo.v2/bson"
)

var BASE_URL = "https://copilot.quan2go.com/api"
var token string
var ua string
var is_window_show = true

func toggleWindow() {
	if is_window_show {
		hideWindow()
	} else {
		showWindow()
	}
}

func hideWindow() {
	is_window_show = false
	wail_runtime.WindowHide(__application.ctx)
}

func showWindow() {
	is_window_show = true
	wail_runtime.WindowShow(__application.ctx)
}

func minWindow() {
	wail_runtime.WindowMinimise(__application.ctx)
}

func isNormal() bool {
	return wail_runtime.WindowIsNormal(__application.ctx)
}

func checkStarted() {
	started := false
	t := time.Now()
	fmt.Println("checkStarted", t)
	wail_runtime.EventsOnce(__application.ctx, "started", func(...interface{}) {
		fmt.Println("started", time.Since(t))
		started = true
	})
	wail_runtime.WindowExecJS(__application.ctx, "runtime.EventsEmit('started')")
	time.Sleep(time.Second * 1)
	if !started {
		fmt.Println("1秒后启动失败，尝试兼容模式启动")
		if runtime.GOOS == "windows" {
			fullpath, err := os.Executable()
			if err != nil {
				fmt.Println(err)
				return
			}
			// 检查是否兼容模式
			const key = "HKEY_CURRENT_USER\\Software\\Microsoft\\Windows NT\\CurrentVersion\\AppCompatFlags\\Layers"
			cmd := spawn("cmd", "/C", "reg", "query", key, "/v", fullpath)
			out, err := cmd.Output()
			if err == nil && bytes.Contains(out, []byte("REG_SZ")) {
				fmt.Println("已经是兼容模式")
				if __application.Confirm("似乎程序启动遇到一些问题, 去反馈一下?") {
					wail_runtime.BrowserOpenURL(__application.ctx, "https://support.qq.com/product/611481")
				}
				return
			}
			if __application.Confirm("程序启动失败, 尝试兼容模式启动?") {
				cmd := spawn("cmd", "/C", "reg", "add", key, "/v", fullpath, "/t", "REG_SZ", "/d", "~ WIN8RTM")
				err := cmd.Run()
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("兼容模式启动成功")
				_ = RestartSelf()
			}
		} else {
			if __application.Confirm("似乎程序启动遇到一些问题, 去反馈一下?") {
				wail_runtime.BrowserOpenURL(__application.ctx, "https://support.qq.com/product/611481")
			}
		}
	}
}

func api_post(url string, params map[string]interface{}) (map[string]interface{}, error) {
	url = BASE_URL + url
	client := &http.Client{}

	buf, err := bson.Marshal(params)

	for i := 0; i < len(buf); i++ {
		buf[i] = buf[i] ^ 0x37
	}

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader((buf)))
	if token != "" {
		req.Header.Set("x-auth-token", token)
	}
	if ua != "" {
		req.Header.Set("User-Agent", ua)
	}
	req.Header.Set("Content-Type", "application/secret")

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	for i := 0; i < len(body); i++ {
		body[i] = body[i] ^ 0x37
	}

	data := make(map[string]interface{})
	err = bson.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	if data["code"] != 0 {
		return nil, errors.New(data["msg"].(string))
	}

	return data["data"].(map[string]interface{}), nil
}
