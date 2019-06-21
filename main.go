// PROJECT: GoGUI
//
// MAINTAINED BY: hkdb <hkdb@3df.io>
//
// SPONSORED BY: 3DF OSI - https://osi.3df.io
//
// This application is maintained by volunteers and in no way
// do the maintainers make any guarantees. Use at your own risk.
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/skratchdot/open-golang/open"

	"github.com/zserge/webview"

	b64 "encoding/base64"
)

const (
	windowWidth  = 400
	windowHeight = 600
	title        = "3DF GoGUI"
	version      = "v0.01"
)

// Load Logo
var logo = MustAsset("assets/3DFosi.png")

var indexHTML = `
<!doctype html>
<html>	
	<head>
		<title>3DF GoGUI</title>
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
	</head>
	
	<body>
		<br>
		<center><img src="data:image/png;base64, ` + string(b64.StdEncoding.EncodeToString([]byte(logo))) + `"></center>
		<br>
		<center><font face="Roboto" color="white" size=2>An <a href="javascript:external.invoke('openosi')">OSI</a> application sponsored by <a href="javascript:external.invoke('open3df')">3DF</a></font></center>
		<center><font face="Roboto" color="white" size=2>` + version + `</font></center>
		<br>
		<br>
		<button onclick="external.invoke('close')">Close</button>
		<button onclick="external.invoke('fullscreen')">Fullscreen</button>
		<button onclick="external.invoke('unfullscreen')">Unfullscreen</button>
		<button onclick="external.invoke('open')">Open</button>
		<button onclick="external.invoke('opendir')">Open directory</button>
		<button onclick="external.invoke('save')">Save</button>
		<button onclick="external.invoke('message')">Message</button>
		<button onclick="external.invoke('info')">Info</button>
		<button onclick="external.invoke('warning')">Warning</button>
		<button onclick="external.invoke('error')">Error</button>
		<button onclick="external.invoke('changeTitle:'+document.getElementById('new-title').value)">
			Change title
		</button>
		<input id="new-title" type="text" />
		<button onclick="external.invoke('changeColor:'+document.getElementById('new-color').value)">
			Change color
		</button>
		<input id="new-color" value="#38393b" type="color" />
	</body>
</html>
`

func startServer() string {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer ln.Close()

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(indexHTML))
		})

		log.Fatal(http.Serve(ln, nil))
	}()
	return "http://" + ln.Addr().String()
}

func handleRPC(w webview.WebView, data string) {
	switch {
	case data == "openosi":
		openl("https://osi.3df.io")
	case data == "open3df":
		openl("https://3df.io")
	case data == "close":
		w.Terminate()
	case data == "fullscreen":
		w.SetFullscreen(true)
	case data == "unfullscreen":
		w.SetFullscreen(false)
	case data == "open":
		log.Println("open", w.Dialog(webview.DialogTypeOpen, 0, "Open file", ""))
	case data == "opendir":
		log.Println("open", w.Dialog(webview.DialogTypeOpen, webview.DialogFlagDirectory, "Open directory", ""))
	case data == "save":
		log.Println("save", w.Dialog(webview.DialogTypeSave, 0, "Save file", ""))
	case data == "message":
		w.Dialog(webview.DialogTypeAlert, 0, "Hello", "Hello, world!")
	case data == "info":
		w.Dialog(webview.DialogTypeAlert, webview.DialogFlagInfo, "Hello", "Hello, info!")
	case data == "warning":
		w.Dialog(webview.DialogTypeAlert, webview.DialogFlagWarning, "Hello", "Hello, warning!")
	case data == "error":
		w.Dialog(webview.DialogTypeAlert, webview.DialogFlagError, "Hello", "Hello, error!")
	case strings.HasPrefix(data, "changeTitle:"):
		w.SetTitle(strings.TrimPrefix(data, "changeTitle:"))
	case strings.HasPrefix(data, "changeColor:"):
		hex := strings.TrimPrefix(strings.TrimPrefix(data, "changeColor:"), "#")
		num := len(hex) / 2
		if !(num == 3 || num == 4) {
			log.Println("Color must be RRGGBB or RRGGBBAA")
			return
		}
		i, err := strconv.ParseUint(hex, 16, 64)
		if err != nil {
			log.Println(err)
			return
		}
		if num == 3 {
			r := uint8((i >> 16) & 0xFF)
			g := uint8((i >> 8) & 0xFF)
			b := uint8(i & 0xFF)
			w.SetColor(r, g, b, 255)
			return
		}
		if num == 4 {
			r := uint8((i >> 24) & 0xFF)
			g := uint8((i >> 16) & 0xFF)
			b := uint8((i >> 8) & 0xFF)
			a := uint8(i & 0xFF)
			w.SetColor(r, g, b, a)
			return
		}
	}
}

//Helper for Opening URL with Default Browser
func openl(uri string) {
	err := open.Run(uri)
	fmt.Println(err)

}

func main() {

	url := startServer()

	webview.Debug()

	w := webview.New(webview.Settings{
		Width:     windowWidth,
		Height:    windowHeight,
		Title:     title,
		Resizable: true,
		URL:       url,
		ExternalInvokeCallback: handleRPC,
	})

	w.SetColor(77, 77, 77, 255)
	defer w.Exit()
	w.Run()
}
