package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"syscall"

	"github.com/inu1255/go-selfupdate/selfupdate"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

type ProxyLoader struct {
	http.Handler
	target string
}

func NewProxyLoader(target string) *ProxyLoader {
	return &ProxyLoader{
		target: target,
	}
}

func (h *ProxyLoader) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if runtime.GOOS == "windows" {
		w.Header().Set("Location", h.target)
		w.WriteHeader(302)
	} else {
		_, _ = w.Write([]byte("<script>window.location.href='" + h.target + "'</script>"))
	}
}

func main() {
	self_update()
	// Create an instance of the app structure
	__application = &App{}
	ua = __application.GetUA()

	// Create application with optionsgo
	err := wails.Run(&options.App{
		Title:  title,
		Width:  500,
		Height: 630,
		AssetServer: &assetserver.Options{
			Handler: NewProxyLoader("https://copilot.quan2go.com/"),
		},
		BackgroundColour: &options.RGBA{R: 250, G: 252, B: 255, A: 1},
		OnStartup:        __application.onstartup,
		OnBeforeClose:    __application.onbeforeclose,
		Debug: options.Debug{
			OpenInspectorOnStartup: true,
		},
		EnableFraudulentWebsiteDetection: true,
		Linux: &linux.Options{
			WebviewGpuPolicy: linux.WebviewGpuPolicyNever,
		},
		Windows:           &windows.Options{},
		HideWindowOnClose: true,
		OnDomReady:        __application.onDomReady,
		Bind: []interface{}{
			__application,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func self_update() {
	url := "https://2go.inu1255.cn/pc/"
	var updater = &selfupdate.Updater{
		CurrentVersion: version,
		ApiURL:         url,
		BinURL:         url,
		DiffURL:        url,
		Dir:            os.TempDir() + "/copilot-activator-updater/",
		CmdName:        exename,
		OnSuccessfulUpdate: func() {
			_ = RestartSelf()
		},
	}

	go func() {
		err := updater.BackgroundRun()
		fmt.Println(err)
	}()

}

func RestartSelf() error {
	self, err := os.Executable()
	if err != nil {
		return err
	}
	args := os.Args
	env := os.Environ()
	// Windows does not support exec syscall.
	if runtime.GOOS == "windows" {
		cmd := spawn("cmd", "/C", "start", "/b", "", self)
		err := cmd.Run()
		if err == nil {
			quitTray()
		}
		return err
	}
	return syscall.Exec(self, args, env)
}
