package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const DB_PATH = "./database"

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8585")
	if err != nil {
		return
	}
	defer listener.Close()
	fmt.Println("=== ORSHAD GATEWAY OPERATIONAL (ROUTING FIXED) ===")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleNetworkTraffic(conn)
	}
}

func handleNetworkTraffic(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	requestLine = strings.TrimSpace(requestLine)
	parts := strings.Split(requestLine, " ")
	if len(parts) < 3 || parts[0] != "GET" {
		sendResponse(conn, 400, "BAD_REQUEST", "text/plain", "Header Syntax Error.")
		return
	}

	url := parts[1]

	// معالجة محرك البحث
	if strings.HasPrefix(url, "search?q=") {
		query := strings.TrimPrefix(url, "search?q=")
		searchConn, err := net.Dial("tcp", "127.0.0.1:8686")
		if err != nil {
			sendResponse(conn, 502, "BAD_GATEWAY", "text/plain", "Search Engine Unreachable.")
			return
		}

		fmt.Fprintf(searchConn, "%s\n", query)
		searchReader := bufio.NewReader(searchConn)
		rawResults, _ := searchReader.ReadString('\n')
		rawResults = strings.TrimSpace(rawResults)
		searchConn.Close()

		if rawResults == "NO_RESULTS" || rawResults == "" {
			sendResponse(conn, 204, "NO_CONTENT", "text/plain", "No results found.")
			return
		}

		html := "<html><head><title>Search Results</title></head><body><h1>Results</h1><ul>"
		sites := strings.Split(rawResults, ",")
		for _, site := range sites {
			html += fmt.Sprintf("<li><a href=\"/index.html\">%s</a></li>", site)
		}
		html += "</ul></body></html>"

		sendResponse(conn, 200, "ORSHAD_OK", "text/html", html)
		return
	}

	// معالجة الـ Routing الذكي للمجلدات والملفات الفرعية
	cleanedUrl := strings.TrimPrefix(url, "/")
	var targetFile string

	if cleanedUrl == "" || !strings.Contains(cleanedUrl, ".") {
		sendResponse(conn, 400, "BAD_REQUEST", "text/plain", "Invalid Domain/Path Context.")
		return
	}

	// تتبع ما إذا كان الطلب مجرد دومين صافي أم ملف فرعي داخل الدومين
	urlParts := strings.Split(cleanedUrl, "/")
	if len(urlParts) == 1 {
		// طلب دومين صافي مثل orshad.osh -> نوجهه لـ index.html
		targetFile = filepath.Join(DB_PATH, urlParts[0], "index.html")
	} else {
		// طلب ملف فرعي مثل orshad.osh/style.css -> نوجهه للمسار الصحيح مباشرة
		targetFile = filepath.Join(DB_PATH, cleanedUrl)
	}

	body, err := os.ReadFile(targetFile)
	if err != nil {
		sendResponse(conn, 404, "NOT_FOUND", "text/plain", "Resource Not Found: "+targetFile)
		return
	}

	contentType := "text/plain"
	ext := filepath.Ext(targetFile)
	switch ext {
	case ".html":
		contentType = "text/html"
	case ".css":
		contentType = "text/css"
	case ".js":
		contentType = "application/javascript"
	}

	sendResponse(conn, 200, "ORSHAD_OK", contentType, string(body))
}

func sendResponse(conn net.Conn, statusCode int, statusMsg string, contentType string, body string) {
	conn.Write([]byte(fmt.Sprintf("ORSHP/1.0 %d %s\n", statusCode, statusMsg)))
	conn.Write([]byte(fmt.Sprintf("Content-Type: %s\n\n", contentType)))
	conn.Write([]byte(body + "\n"))
}
