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
	"os/exec"
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
	
	<body bgcolor="#4D4D4D">
		<br>
		<center><img src="data:image/png;base64, ` + string(b64.StdEncoding.EncodeToString([]byte(logo))) + `"></center>
		<br>
		<center><font face="Roboto" color="white" size=2>An <a href="javascript:external.invoke('openosi')">OSI</a> application sponsored by <a href="javascript:external.invoke('open3df')">3DF</a></font></center>
		<center><font face="Roboto" color="white" size=2>` + version + `</font></center>
		<br>
		<hr size=1>

		<br>
		<center>
		<div style="border: 1px"><input type="file" id="infile" name="infile"></div>
		<p><button onclick="external.invoke('open')">Open</button>
		<button onclick="external.invoke('opendir')">Open Directory</button>
		<button onclick="external.invoke('save')">Save</button>
		<p><button onclick="external.invoke('message')">Message</button>
		<button onclick="external.invoke('info')">Info</button>
		<button onclick="external.invoke('warning')">Warning</button>
		<button onclick="external.invoke('error')">Error</button>
		<p><button onclick="external.invoke('changeTitle:'+document.getElementById('new-title').value)">
    		Change title
		</button>
		<input id="new-title" type="text" />
		</p>
		<p><button onclick="external.invoke('submit')">Submit</button>
		</center>
		
		
		<div id="footer-line" style="position: relative">
			<hr size=1 valign="bottom">
		</div>
		<div id="footer" style="position: relative">
			<p style="position: fixed; bottom: 0; width:100%; text-align: center">
				<button onclick="external.invoke('fullscreen')">Fullscreen</button>
				<button onclick="external.invoke('unfullscreen')">Unfullscreen</button>
				<button onclick="external.invoke('about')">About</button>
				<button onclick="external.invoke('close')">Close</button>
			</p>
		</div>
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
	case data == "about":
		w.Dialog(webview.DialogTypeAlert, webview.DialogFlagInfo, "About", "\nAn OSI application sponsored by 3DF. For more information\n\nVisit:\n\nhttps://osi.3df.io")
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
	case data == "submit":
		submit()

	}
}

//Helper for Opening URL with Default Browser
func openl(uri string) {
	err := open.Run(uri)
	fmt.Println(err)

}

func submit() {
	cmd := exec.Command("echo", "Submitted")
	fmt.Println(cmd)
}

func main() {

	url := startServer()

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
