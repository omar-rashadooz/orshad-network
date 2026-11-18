package main

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/gdamore/tcell/v2"
)

var (
	screen         tcell.Screen
	currentURL     = "Welcome Page"
	statusMsg      = "STANDBY"
	viewLines      = []string{}
	rawHTMLPayload = ""
	inputBuffer    = ""
	isWelcomePage  = true
)

func main() {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Printf("Failed to create screen: %v\n", err)
		return
	}
	if err = screen.Init(); err != nil {
		fmt.Printf("Failed to init screen: %v\n", err)
		return
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)
	screen.EnableMouse()

	for {
		screen.Clear()
		drawInterface()
		screen.Show()

		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			if ev.Modifiers() == tcell.ModCtrl && ev.Key() == tcell.KeyRune {
				switch ev.Rune() {
				case 'q', 'Q':
					screen.Fini()
					return
				case 's', 'S':
					triggerSearchMode()
					continue
				case 'h', 'H':
					triggerHomeMode()
					continue
				case 'r', 'R':
					triggerRefreshMode()
					continue
				case 'c', 'C':
					inputBuffer = ""
					continue
				}
			}

			switch ev.Key() {
			case tcell.KeyEscape:
				screen.Fini()
				return
			case tcell.KeyF2:
				triggerSearchMode()
				continue
			case tcell.KeyF3:
				triggerHomeMode()
				continue
			case tcell.KeyF4:
				triggerRefreshMode()
				continue
			case tcell.KeyF5:
				inputBuffer = ""
				continue
			}
			
			handleKeyInput(ev)
		}
	}
}

func triggerSearchMode() {
	isWelcomePage = false
	inputBuffer = "search?q="
	statusMsg = "SEARCH MODE"
}

func triggerHomeMode() {
	isWelcomePage = true
	currentURL = "Welcome Page"
	statusMsg = "STANDBY"
	viewLines = []string{}
}

func triggerRefreshMode() {
	if !isWelcomePage && currentURL != "" {
		executeNetworkRequest(currentURL)
	}
}

func drawInterface() {
	width, height := screen.Size()

	// 1. شريط العنوان العلوي
	headerStyle := tcell.StyleDefault.Background(tcell.ColorLightCyan).Foreground(tcell.ColorBlack)
	for x := 0; x < width; x++ {
		screen.SetContent(x, 0, ' ', nil, headerStyle)
	}
	emitString(1, 0, headerStyle, " ORSHAD OS/BROWSER v2.0 | Search Core Engine Active ")
	emitString(width-22, 0, headerStyle, fmt.Sprintf("Status: %-10s", statusMsg))

	addrStyle := tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorYellow)
	for x := 0; x < width; x++ {
		screen.SetContent(x, 1, ' ', nil, addrStyle)
	}
	
	if isWelcomePage {
		emitString(1, 1, addrStyle, "Address -> orsh://welcome.local")
	} else {
		emitString(1, 1, addrStyle, fmt.Sprintf("Address -> orsh://%s", currentURL))
	}

	// 2. مساحة عرض المحتوى الوسطى
	viewStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	for y := 2; y < height-3; y++ {
		for x := 0; x < width; x++ {
			screen.SetContent(x, y, ' ', nil, viewStyle)
		}
	}

	if isWelcomePage {
		drawWelcomePage(width, height)
	} else {
		for i, line := range viewLines {
			if i+2 >= height-3 {
				break
			}
			if strings.HasPrefix(line, "[H1] ") {
				h1Style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorYellow).Bold(true)
				emitString(2, i+2, h1Style, strings.TrimPrefix(line, "[H1] "))
			} else if strings.HasPrefix(line, "[LINK] ") {
				// رندرة رابط الموقع بلؤلؤي ساطع/أخضر كجوجل
				linkStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen).Bold(true)
				emitString(2, i+2, linkStyle, "🔗 "+strings.TrimPrefix(line, "[LINK] "))
			} else if strings.HasPrefix(line, "[SNIPPET] ") {
				// رندرة النبذة بلون رمادي خفيف لتمييزها تماماً كالمحركات الاحترافية
				snippetStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorDarkGray)
				emitString(4, i+2, snippetStyle, strings.TrimPrefix(line, "[SNIPPET] "))
			} else {
				emitString(2, i+2, viewStyle, line)
			}
		}
	}

	// 3. شريط المساعدة السفلي
	ribbonStyle := tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite)
	keyStyle := tcell.StyleDefault.Background(tcell.ColorLightCyan).Foreground(tcell.ColorBlack)
	for x := 0; x < width; x++ {
		screen.SetContent(x, height-3, ' ', nil, ribbonStyle)
	}
	
	emitString(2, height-3, keyStyle, " ESC ")
	emitString(8, height-3, ribbonStyle, "Exit")
	emitString(14, height-3, keyStyle, " F2/Ctrl+S ")
	emitString(25, height-3, ribbonStyle, "Search Node")
	emitString(39, height-3, keyStyle, " F3/Ctrl+H ")
	emitString(50, height-3, ribbonStyle, "Home Engine")
	emitString(64, height-3, keyStyle, " F4 ")
	emitString(69, height-3, ribbonStyle, "Refresh")

	// 4. شريط الإدخال
	inputBarStyle := tcell.StyleDefault.Background(tcell.ColorLightGray).Foreground(tcell.ColorBlack)
	for x := 0; x < width; x++ {
		screen.SetContent(x, height-2, ' ', nil, inputBarStyle)
		screen.SetContent(x, height-1, ' ', nil, inputBarStyle)
	}
	emitString(1, height-2, inputBarStyle, "Query Engine Prompt (e.g., search?q=hardware or search?q=matrix):")
	
	promptStr := "DR-DOS C:\\> "
	emitString(1, height-1, inputBarStyle, promptStr+inputBuffer)

	cursorX := 1 + len(promptStr) + len(inputBuffer)
	cursorY := height - 1
	screen.ShowCursor(cursorX, cursorY)
}

func drawWelcomePage(width, height int) {
	titleStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorYellow).Bold(true)
	descStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	catalogHeaderStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorDarkCyan).Bold(true)
	catalogRowStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen)
	hintStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed).Bold(true)

	midX := width / 2

	logo :=  "╔═╗╦═╗╔═╗╦ ╦╔═╗╔╦╗  ╔╗╔╔═╗╔╦╗"
	logo2 := "║ ║╠╦╝╚═╗╠═╣╠═╣ ║║  ║║║║╣  ║ "
	logo3 := "╚═╝╩╚═╚═╝╩ ╩╩ ╩═╩╝  ║╚╝╚═╝ ╩ "
	
	emitString(midX-(len(logo)/2), 3, titleStyle, logo)
	emitString(midX-(len(logo2)/2), 4, titleStyle, logo2)
	emitString(midX-(len(logo3)/2), 5, titleStyle, logo3)

	emitString(midX-21, 7, descStyle, "Welcome to the next generation of Decentralized Web.")
	startY := 10
	emitString(midX-25, startY, catalogHeaderStyle, "┌────────────────── SYSTEM CATALOG ──────────────────┐")
	emitString(midX-25, startY+1, catalogHeaderStyle, "│  SHORTCUT     │  OPERATION DESCRIPTION             │")
	emitString(midX-25, startY+2, catalogHeaderStyle, "├────────────────────────────────────────────────────┤")
	emitString(midX-25, startY+3, catalogRowStyle, "│  F2 / Ctrl+S  │  Activate Search Engine Core       │")
	emitString(midX-25, startY+4, catalogRowStyle, "│  F3 / Ctrl+H  │  Return to local Welcome Page      │")
	emitString(midX-25, startY+5, catalogRowStyle, "│  F4 / Ctrl+R  │  Refresh & Reparse Current Node    │")
	emitString(midX-25, startY+6, catalogRowStyle, "│  F5 / Ctrl+C  │  Flush Command Line Buffer         │")
	emitString(midX-25, startY+7, catalogRowStyle, "│  ESC / Q      │  Immediate Core Termination        │")
	emitString(midX-25, startY+8, catalogHeaderStyle, "└────────────────────────────────────────────────────┘")

	emitString(midX-23, startY+10, hintStyle, ">> Press [F2] to initiate Search Engine Query <<")
}

func handleKeyInput(ev *tcell.EventKey) {
	if ev.Key() == tcell.KeyEnter {
		cmd := strings.TrimSpace(inputBuffer)
		if cmd == "exit" {
			screen.Fini()
			return
		}
		if cmd != "" {
			isWelcomePage = false
			currentURL = cmd
			executeNetworkRequest(cmd)
		}
		inputBuffer = ""
	} else if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
		if len(inputBuffer) > 0 {
			inputBuffer = inputBuffer[:len(inputBuffer)-1]
		}
	} else if ev.Key() == tcell.KeyRune {
		inputBuffer += string(ev.Rune())
	}
}

func executeNetworkRequest(input string) {
	statusMsg = "LOOKUP..."
	viewLines = []string{} 
	rawHTMLPayload = ""

	var domain, path string
	if strings.HasPrefix(input, "search?q=") {
		domain = "127.0.0.1:8585" 
		path = "/" + input
	} else {
		parts := strings.Split(input, "/")
		domain = parts[0]
		if len(parts) > 1 {
			path = "/" + input
		} else {
			path = "/" + domain
		}
	}

	var resolvedAddress string
	if strings.HasPrefix(input, "search?q=") {
		resolvedAddress = "127.0.0.1:8585"
	} else {
		dnsConn, err := net.Dial("tcp", "127.0.0.1:5353")
		if err != nil {
			statusMsg = "DNS_ERR"
			viewLines = append(viewLines, "[CRITICAL] Local DNS Server is offline.")
			return
		}
		fmt.Fprintf(dnsConn, "%s\n", domain)
		dnsReader := bufio.NewReader(dnsConn)
		resolvedAddress, _ = dnsReader.ReadString('\n')
		resolvedAddress = strings.TrimSpace(resolvedAddress)
		dnsConn.Close()
	}

	if resolvedAddress == "NXDOMAIN" {
		statusMsg = "404_ERR"
		viewLines = append(viewLines, "Error: Domain does not exist in DNS registry.")
		return
	}

	statusMsg = "SENDING..."
	conn, err := net.Dial("tcp", resolvedAddress)
	if err != nil {
		statusMsg = "CONN_ERR"
		viewLines = append(viewLines, "[NETWORK ERROR] Connection refused.")
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, "GET %s ORSHP/1.0\n", path)

	statusMsg = "RECEIVING"
	reader := bufio.NewReader(conn)
	
	fullResponse := ""
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fullResponse += line
	}

	parts := strings.Split(fullResponse, "\n\n")
	if len(parts) > 1 {
		rawHTMLPayload = parts[1]
	} else {
		rawHTMLPayload = fullResponse
	}

	parseAndRenderHTML(rawHTMLPayload)
	statusMsg = "200_OK"
}

func parseAndRenderHTML(html string) {
	// تنظيف الرؤوس والوسوم العامة
	html = regexp.MustCompile(`(?i)<html>|</html>|<body>|</body>|<head>.*?</head>`).ReplaceAllString(html, "")

	h1Reg := regexp.MustCompile(`(?i)<h1>(.*?)</h1>`)
	html = h1Reg.ReplaceAllString(html, "[H1] $1")

	// معالجة السلسلة النصية للنتائج وفصل الرابط عن النبذة
	liReg := regexp.MustCompile(`(?i)<li>Link:\s*(.*?)\s*\|\s*(.*?)<br>\s*Snippet:\s*(.*?)</li>`)
	html = liReg.ReplaceAllString(html, "[LINK] $1 -> $2\n[SNIPPET] $3")

	finalCleanReg := regexp.MustCompile("<[^>]*>")
	cleanedText := finalCleanReg.ReplaceAllString(html, "")

	rawLines := strings.Split(cleanedText, "\n")
	for _, line := range rawLines {
		trimmed := strings.TrimSpace(line) // تم التصحيح هنا لاستخدام TrimSpace مباشرة
		if trimmed != "" {
			viewLines = append(viewLines, trimmed)
		}
	}
}
// emitString هي الدالة المساعدة لطباعة النصوص وتلوينها داخل tcell
func emitString(x, y int, style tcell.Style, str string) {
	runes := []rune(str)
	for i := 0; i < len(runes); i++ {
		screen.SetContent(x+i, y, runes[i], nil, style)
	}
}

