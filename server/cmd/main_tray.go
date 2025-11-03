//go:build windows && tray

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/alexflint/go-arg"
	"fyne.io/systray"
	"golang.org/x/sys/windows/registry"

	"server"
	"server/log"
	"server/settings"
)

var (
	trayParams args
	serverRunning int32
	defaultLogPath string
    startItem *systray.MenuItem
    stopItem *systray.MenuItem
    autostartItem *systray.MenuItem
)

// local copy of args used by tray build (original is behind !tray build tag)
type args struct {
	Port        string `arg:"-p" help:"web server port (default 8090)"`
	IP          string `arg:"-i" help:"web server addr (default empty)"`
	Ssl         bool   `help:"enables https"`
	SslPort     string `help:"web server ssl port, If not set, will be set to default 8091 or taken from db(if stored previously). Accepted if --ssl enabled."`
	SslCert     string `help:"path to ssl cert file. If not set, will be taken from db(if stored previously) or default self-signed certificate/key will be generated. Accepted if --ssl enabled."`
	SslKey      string `help:"path to ssl key file. If not set, will be taken from db(if stored previously) or default self-signed certificate/key will be generated. Accepted if --ssl enabled."`
	Path        string `arg:"-d" help:"database and config dir path"`
	LogPath     string `arg:"-l" help:"server log file path"`
	WebLogPath  string `arg:"-w" help:"web access log file path"`
	RDB         bool   `arg:"-r" help:"start in read-only DB mode"`
	HttpAuth    bool   `arg:"-a" help:"enable http auth on all requests"`
	DontKill    bool   `arg:"-k" help:"don't kill server on signal"`
	UI          bool   `arg:"-u" help:"open torrserver page in browser"`
	TorrentsDir string `arg:"-t" help:"autoload torrents from dir"`
	TorrentAddr string `help:"Torrent client address, like 127.0.0.1:1337 (default :PeersListenPort)"`
	PubIPv4     string `arg:"-4" help:"set public IPv4 addr"`
	PubIPv6     string `arg:"-6" help:"set public IPv6 addr"`
	SearchWA       bool   `arg:"-s" help:"search without auth"`
	MaxSize        string `arg:"-m" help:"max allowed stream size (in Bytes)"`
	TGToken        string `arg:"-T" help:"telegram bot token"`
	JlfnAddr       string `help:"Jellyfin .strm files path (e.g., /media/jellyfin/metadata)"`
	JlfnSrv        string `help:"Jellyfin server URL (e.g., http://127.0.0.1:8096)"`
	JlfnApi        string `help:"Jellyfin API key"`
	JlfnAutoCreate bool   `help:"Auto-create .strm files when adding torrents via web"`
	TMDBApiKey     string `help:"TMDB API key for metadata and posters"`
	TorrServerHost string `help:"Public TorrServer URL for .strm files (e.g., http://192.168.1.197:5665)"`
}

func main() {
	arg.MustParse(&trayParams)

	if trayParams.Path == "" {
		trayParams.Path, _ = os.Getwd()
	}

	// Ensure logs go to file so we can tail in separate terminal
	defaultLogPath = trayParams.LogPath
	if defaultLogPath == "" {
		defaultLogPath = filepath.Join(trayParams.Path, "torrserver.log")
		trayParams.LogPath = defaultLogPath
	}

	settings.Path = trayParams.Path
	settings.HttpAuth = trayParams.HttpAuth
	log.Init(trayParams.LogPath, trayParams.WebLogPath)

	// Start tray UI
	systray.Run(onReady, onExit)
}

func onReady() {
	// Try to load icon from common locations
	exeDir := func() string { p, _ := os.Executable(); d, _ := filepath.Abs(filepath.Dir(p)); return d }()
	iconCandidates := []string{
		filepath.Join(exeDir, "favicon.ico"),
		filepath.Join(exeDir, "..", "web", "public", "favicon.ico"),
		filepath.Join(exeDir, "..", "..", "web", "public", "favicon.ico"),
	}
	for _, p := range iconCandidates {
		if b, err := os.ReadFile(p); err == nil && len(b) > 0 {
			systray.SetIcon(b)
			break
		}
	}
	systray.SetTitle("TorrServer")
	systray.SetTooltip("TorrServer Jellyfin")

    startItem = systray.AddMenuItem("Start server", "Start server")
    stopItem = systray.AddMenuItem("Stop server", "Stop server")
    autostartItem = systray.AddMenuItemCheckbox("Autostart", "Run at user login", isAutostartEnabled())
    logsItem := systray.AddMenuItem("Logs", "Open logs window")
	systray.AddSeparator()
	quitItem := systray.AddMenuItem("Exit", "Exit tray")

	go func() {
		for {
			select {
			case <-startItem.ClickedCh:
				startServer()
			case <-stopItem.ClickedCh:
				stopServer()
            case <-autostartItem.ClickedCh:
                desired := !isAutostartEnabled()
                setAutostart(desired)
                refreshMenuState()
			case <-logsItem.ClickedCh:
				openLogsTerminal()
			case <-quitItem.ClickedCh:
				stopServer()
				systray.Quit()
				return
			}
		}
	}()

	// Sync initial state and auto start server when tray UI is ready
	refreshMenuState()
	startServer()
}

func refreshMenuState() {
    if atomic.LoadInt32(&serverRunning) == 1 {
        if startItem != nil { startItem.Disable() }
        if stopItem != nil { stopItem.Enable() }
    } else {
        if startItem != nil { startItem.Enable() }
        if stopItem != nil { stopItem.Disable() }
    }
    if autostartItem != nil {
        if isAutostartEnabled() { autostartItem.Check() } else { autostartItem.Uncheck() }
    }
}

func onExit() {
	log.Close()
}

func startServer() {
	if !atomic.CompareAndSwapInt32(&serverRunning, 0, 1) {
		return
	}
	// Minimal defaults matching non-tray main
	if trayParams.Port == "" { trayParams.Port = "8090" }
	// Run server in goroutine
	go func() {
		server.Start(trayParams.Port, trayParams.IP, trayParams.SslPort, trayParams.SslCert, trayParams.SslKey, trayParams.Ssl, trayParams.RDB, trayParams.SearchWA, trayParams.TGToken, trayParams.JlfnAddr, trayParams.JlfnSrv, trayParams.JlfnApi, trayParams.JlfnAutoCreate, trayParams.TMDBApiKey, trayParams.TorrServerHost)
		_ = server.WaitServer()
		atomic.StoreInt32(&serverRunning, 0)
        refreshMenuState()
	}()
    refreshMenuState()
}

func stopServer() {
	if !atomic.CompareAndSwapInt32(&serverRunning, 1, 0) {
		return
	}
	server.Stop()
    refreshMenuState()
}

func openLogsTerminal() {
	// Prefer file tail if we have a file path
	logPath := defaultLogPath
	if logPath == "" {
		// fallback to simple echo to keep window open
		_ = exec.Command("cmd", "/c", "start", "", "powershell", "-NoExit", "-Command", "Write-Host 'No log file configured. Set --logpath' ").Start()
		return
	}
	// Ensure file exists so Get-Content -Wait won't fail immediately
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		_ = os.MkdirAll(filepath.Dir(logPath), 0o755)
		if f, err := os.Create(logPath); err == nil { f.Close() }
	}
	psCmd := fmt.Sprintf("Get-Content -Path '%s' -Wait -Encoding UTF8", strings.ReplaceAll(logPath, "'", "''"))
	// Try PowerShell 7 first (pwsh), then Windows PowerShell
	if err := exec.Command("pwsh", "-NoExit", "-Command", psCmd).Start(); err != nil {
		_ = exec.Command("powershell", "-NoExit", "-Command", psCmd).Start()
	}
}

func isAutostartEnabled() bool {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\\Microsoft\\Windows\\CurrentVersion\\Run`, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer key.Close()
	_, _, err = key.GetStringValue("TorrServer")
	return err == nil
}

func setAutostart(enable bool) {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\\Microsoft\\Windows\\CurrentVersion\\Run`, registry.SET_VALUE|registry.QUERY_VALUE)
	if err != nil {
		return
	}
	defer key.Close()
	if enable {
		exe, _ := os.Executable()
		exe, _ = filepath.Abs(exe)
		cmd := fmt.Sprintf("\"%s\" --tray", exe)
		_ = key.SetStringValue("TorrServer", cmd)
	} else {
		_ = key.DeleteValue("TorrServer")
	}
}


