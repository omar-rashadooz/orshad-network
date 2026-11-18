# Orshad Network

> A self-contained, experimental web network: its own DNS, its own browser, its own search engine, and its own site format (`.osh`).

![Go](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=white)
![Rust](https://img.shields.io/badge/Rust-000000?logo=rust&logoColor=white)
![Status](https://img.shields.io/badge/status-experimental-orange)

## Overview

**Orshad Network** is a small, independent network stack built from scratch. Instead of relying on the regular web infrastructure, it defines its own way of naming, serving, browsing and searching sites:

- Sites live in a simple folder format ending in **`.osh`**.
- A custom **DNS server** resolves Orshad names.
- A custom **server** delivers site content.
- A custom **browser** loads and renders those sites.
- A **search engine** written in Rust indexes the network.

## Architecture

```
┌──────────────┐   resolve name   ┌──────────────┐
│   Browser    │ ───────────────▶ │  DNS Server  │
│ (browser.go) │ ◀─────────────── │(dns_server.go)│
└──────┬───────┘     address      └──────────────┘
       │ request site
       ▼
┌──────────────┐    reads    ┌─────────────────┐
│    Server    │ ──────────▶ │ database/*.osh  │
│  (server.go) │             │  (site folders) │
└──────────────┘             └────────▲────────┘
                                      │ indexes
                             ┌────────┴────────┐
                             │  orshad_search  │
                             │     (Rust)      │
                             └─────────────────┘
```

## Project Structure

```
orshad-network/
├── browser.go          # Orshad browser
├── dns_server.go       # DNS server for Orshad names
├── server.go           # Server that serves .osh sites
├── go.mod / go.sum     # Go module files
├── database/           # Hosted sites (one .osh folder per site)
│   └── orshad.osh/
│       ├── index.html  # Site entry page
│       ├── main.js     # Site scripts
│       ├── style.css   # Site styles
│       └── meta.json   # Site metadata
└── orshad_search/      # Search engine (Rust / Cargo)
    ├── Cargo.toml
    └── src/main.rs
```

## The `.osh` Site Format

Every site is a folder named `<name>.osh` inside `database/`, containing standard web files plus a `meta.json` describing the site.

```
database/
└── mysite.osh/
    ├── index.html
    ├── main.js
    ├── style.css
    └── meta.json
```

To publish a new site, create a new `.osh` folder with at least an `index.html` and a `meta.json`.

## Getting Started

### Requirements

- [Go](https://go.dev/dl/) 1.20+
- [Rust & Cargo](https://rustup.rs/)

### Clone

```bash
git clone https://github.com/omar-rashadooz/orshad-network.git
cd orshad-network
```

### Run the components

Start each component in its own terminal. (Adjust the commands if your entry points differ.)

```bash
# 1. DNS server
go run dns_server.go

# 2. Content server
go run server.go

# 3. Browser
go run browser.go

# 4. Search engine
cd orshad_search
cargo run --release
```

> **Note:** if Go reports multiple `main` functions in one folder, run each file separately as shown above, or move each component into its own subfolder.

## Roadmap

- [ ] Document the DNS record format
- [ ] Define the full `meta.json` specification
- [ ] Add tests for the DNS server and the server
- [ ] Add a crawler so the search engine indexes every `.osh` site automatically
- [ ] Add configuration files / environment variables for ports and paths

## Contributing

Ideas and pull requests are welcome.

1. Fork the repository
2. Create a branch: `git checkout -b feature/my-feature`
3. Commit your changes: `git commit -m "Add my feature"`
4. Push: `git push origin feature/my-feature`
5. Open a Pull Request

## Author

**Omar Rashad** — [@omar-rashadooz](https://github.com/omar-rashadooz)

---

## 🇪🇬 نبذة بالعربي

**شبكة أرشاد (Orshad Network)** مشروع تجريبي بيبني شبكة ويب مستقلة من الصفر: سيرفر DNS خاص، وسيرفر لعرض المواقع، ومتصفح مخصص، ومحرك بحث مكتوب بلغة Rust. المواقع بتتخزن في فولدرات بامتداد `.osh` جوه مجلد `database`، وكل موقع فيه `index.html` و`main.js` و`style.css` وملف `meta.json` للبيانات الوصفية.
