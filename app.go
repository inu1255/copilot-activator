package main

import (
	"context"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"runtime"

	"github.com/denisbrodbeck/machineid"
	wail_runtime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

var __application *App

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) onstartup(ctx context.Context) {
	a.ctx = ctx
	go createTray()
}

func (a *App) onbeforeclose(ctx context.Context) bool {
	if runtime.GOOS == "windows" {
		hideWindow()
	} else {
		minWindow()
	}
	return true
}

// 白屏时也会触发
func (a *App) onDomReady(ctx context.Context) {
	go checkStarted()
}

func (a *App) GetUA() string {
	var ua string
	currentUser, _ := user.Current()
	if currentUser != nil {
		ua = currentUser.Username
	}
	ua = ua + "/" + runtime.GOOS + "/" + runtime.GOARCH
	return ua
}

func (a *App) GetDeviceID() string {
	id, _ := machineid.ID()
	return id
}

func (a *App) SetToken(baseURL, t string) {
	BASE_URL = baseURL
	token = t
}

func (a *App) CheckDevice() (map[string]interface{}, error) {
	return api_post("/copilot/check-device", nil)
}

func (a *App) SwitchDevice(cur map[string]string) (map[string]interface{}, error) {
	id, _ := machineid.ID()
	envs, err := api_post("/copilot/switch", map[string]interface{}{
		"device": id,
	})
	if err != nil {
		return nil, err
	}
	if envs != nil {
		envs1 := make(map[string]string)
		for k, v := range envs {
			s := v.(string)
			if s != cur[k] {
				envs1[k] = s
			}
		}
		if runtime.GOOS == "windows" {
			for k, v := range envs1 {
				cmd := spawn("setx", k, v)
				err := cmd.Run()
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return envs, nil
}

func (a *App) GetEnv(name string) string {
	return os.Getenv(name)
}

func (a *App) SetEnv(name, value string) {
	os.Setenv(name, value)
}

func (a *App) GetUserEnv(name string) (string, error) {
	if runtime.GOOS == "windows" {
		// 通过注册表获取环境变量
		cmd := spawn("reg", "query", "HKCU\\Environment", "/v", name)
		data, err := cmd.CombinedOutput()
		if err != nil {
			return "", err
		}
		regexp1 := regexp.MustCompile(name + `\s+(REG_SZ|REG_EXPAND_SZ)\s+([^\r\n]*)`)
		match := regexp1.FindStringSubmatch(string(data))
		if len(match) > 2 {
			return match[2], nil
		}
		return "", nil
	}
	return "", nil
}

func (a *App) Exec(exe string, args ...string) ([]byte, error) {
	cmd := exec.Command(exe, args...)
	out, err := cmd.CombinedOutput()
	return out, err
}

func (a *App) ReadDir(dir string) ([]map[string]interface{}, error) {
	dirs, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	files := make([]map[string]interface{}, len(dirs))
	for i, d := range dirs {
		files[i] = map[string]interface{}{
			"name":  d.Name(),
			"isDir": d.IsDir(),
		}
	}
	return files, nil
}

func (a *App) ReadFile(file string) ([]byte, error) {
	return os.ReadFile(file)
}

func (a *App) WriteFile(file string, data []byte) error {
	return os.WriteFile(file, data, 0644)
}

func (a *App) GetVersion() string {
	return version
}

func (a *App) GetArgs() []string {
	return os.Args
}

func (a *App) TempDir() string {
	return os.TempDir()
}

func (a *App) HomeDir() (string, error) {
	return os.UserHomeDir()
}

func (a *App) Confirm(msg string) bool {
	s, _ := wail_runtime.MessageDialog(__application.ctx, wail_runtime.MessageDialogOptions{
		Type:          wail_runtime.QuestionDialog,
		Title:         "提示",
		Message:       msg,
		DefaultButton: "确定",
		CancelButton:  "取消",
	})
	return s == "确定" || s == "Ok" || s == "Yes" || s == "是"
}
