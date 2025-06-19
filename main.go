package main

import (
	"embed"
	"fmt"
	"github.com/energye/energy/v2/cef"
	"github.com/energye/energy/v2/cef/ipc"
	"github.com/energye/energy/v2/pkgs/assetserve"
	"github.com/energye/golcl/lcl"
	"github.com/energye/golcl/lcl/rtl/version"
	"github.com/yankeguo/rg"
	"helpme/split_excel"
	"io"
)

//go:embed resources
var resources embed.FS

func main() {
	//Global initialization must be called
	cef.GlobalInit(nil, &resources)
	//Create an application
	app := cef.NewApplication()

	//http's url
	cef.BrowserWindow.Config.Url = "http://localhost:22022/index.html"
	//Security key and value settings for built-in static resource services
	assetserve.AssetsServerHeaderKeyName = "energy"
	assetserve.AssetsServerHeaderKeyValue = "energy"
	cef.SetBrowserProcessStartAfterCallback(func(b bool) {
		server := assetserve.NewAssetsHttpServer() //Built in HTTP service
		server.PORT = 22022                        //Service Port Number
		server.AssetsFSName = "resources"          //Resource folder with the same name
		server.Assets = &resources                 //Assets resources
		go server.StartHttpServer()
	})
	//config := cef.BrowserWindow.Config.ChromiumConfig()
	//config.SetEnableDevTools(true) //启用开发者工具
	// run main process and main thread
	cef.BrowserWindow.SetBrowserInit(browserInit)
	//run app
	cef.Run(app)
}

// run main process and main thread
func browserInit(event *cef.BrowserEvent, window cef.IBrowserWindow) {
	// index.html ipc.emit("count", [count++])
	window.SetTitle("helpme")

	ipc.On("select_file_return", handler)

	// page load end
	event.SetOnLoadEnd(func(sender lcl.IObject, browser *cef.ICefBrowser, frame *cef.ICefFrame, httpStatusCode int32, window cef.IBrowserWindow) {
		// index.html, ipc.on("osInfo", function(){...})
		println("osInfo", version.OSVersion.ToString())
		ipc.Emit("osInfo", version.OSVersion.ToString())
		var windowType string
		if window.IsLCL() {
			windowType = "LCL"
		} else {
			windowType = "VF"
		}
		// index.html, ipc.on("windowType", function(){...});
		ipc.Emit("windowType", windowType)
	})

}

func handler(value []byte, fileName string) (_bytes []byte) {
	var err error
	defer func() {
		if err != nil {
			fmt.Printf("select_file_return failed: %v\n", err)
			ipc.Emit("error", err.Error())
		}
	}()
	defer rg.Guard(&err)
	//fmt.Printf("select_file  %+v\n", value)
	file, err := split_excel.Split(fileName, value)
	if err != nil {
		fmt.Printf("Split failed: %v\n", err)
		return
	}
	//fmt.Println("split done")
	// file is closed
	_bytes, err = io.ReadAll(file)
	if err != nil {
		fmt.Printf("ReadAll failed: %v\n", err)
		return
	}
	//fmt.Println("zipFile", _bytes)
	//ipc.Emit("zipFile", _bytes)
	return
}
