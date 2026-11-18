use std::io::{BufRead, BufReader, Write};
use std::net::{TcpListener, TcpStream};
use std::thread;

const DB_PATH: &str = "../database"; 

// الحل السحري هنا: نمنح الهيكل خصائص الاستنساخ التلقائي في الذاكرة
#[derive(Clone, Copy)]
struct WebSite {
    url: &'static str,
    title: &'static str,
    snippet: &'static str,
    content: &'static str,
}

fn main() {
    let listener = TcpListener::bind("127.0.0.1:8585").expect("Failed to bind search server");
    println!("Sovereign Search Engine Core running on port 8585...");

    let database = vec![
        WebSite {
            url: "orshad.osh",
            title: "ORSHAD Hardware & Tech Hub",
            snippet: "The premier portal for hardware metrics, low-level OS configurations, and CPU/GPU analysis.",
            content: "Welcome to ORSHAD Hardware & Tech Hub\nLoading system metrics...",
        },
        WebSite {
            url: "news.osh",
            title: "Sovereign Network News Channel",
            snippet: "Breaking updates on decentralized protocols, DR-DOS kernel developments, and node sync status.",
            content: "Welcome to Sovereign Network News\nAll nodes are currently healthy.",
        },
        WebSite {
            url: "games.osh",
            title: "Retro Matrix Terminal Games",
            snippet: "Extreme difficulty algorithmic logic games designed strictly for native command-line systems.",
            content: "Welcome to Retro Matrix Games\nSelect game core to initiate setup.",
        },
    ];

    for stream in listener.incoming() {
        match stream {
            Ok(stream) => {
                let db_clone = database.clone(); // الآن ستعمل بسلاسة تامة دون أي معارضة من الـ Compiler
                thread::spawn(move || {
                    handle_client(stream, db_clone);
                });
            }
            Err(e) => eprintln!("Connection failed: {}", e),
        }
    }
}

fn handle_client(mut stream: TcpStream, db: Vec<WebSite>) {
    let reader = BufReader::new(&stream);
    let mut lines = reader.lines();
    
    if let Some(Ok(request_line)) = lines.next() {
        let parts: Vec<&str> = request_line.split_whitespace().collect();
        if parts.len() < 2 { return; }
        
        let path = parts[1];
        let mut html_response = String::new();

        if path.starts_with("/search?q=") {
            let query = path.trim_start_matches("/search?q=").to_lowercase();
            html_response.push_str("<html><body>\n");
            html_response.push_str(&format!("<h1>Sovereign Search Results for: [{}]</h1>\n", query));
            html_response.push_str("<p>Found relevant nodes registered in the sovereign grid:</p>\n");

            let mut hits = 0;
            for site in db {
                if site.title.to_lowercase().contains(&query) || site.snippet.to_lowercase().contains(&query) || site.url.to_lowercase().contains(&query) {
                    hits += 1;
                    html_response.push_str(&format!(
                        "<li>Link: {} | {}<br>  Snippet: {}</li>\n", 
                        site.url, site.title, site.snippet
                    ));
                }
            }

            if hits == 0 {
                html_response.push_str("<p>No micro-nodes matched your routing query parameters.</p>\n");
            }
            html_response.push_str("</body></html>\n");
        } else {
            html_response.push_str("ORSHP/1.0 404 NOT_FOUND\nContent-Type: text/plain\n\nResource Not Found.");
            let _ = stream.write_all(html_response.as_bytes());
            return;
        }

        let full_response = format!(
            "ORSHP/1.0 200 ORSHAD_OK\nContent-Type: text/html\nContent-Length: {}\n\n{}",
            html_response.len(),
            html_response
        );
        let _ = stream.write_all(full_response.as_bytes());
    }
}
