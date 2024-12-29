# OVERVIEW 
[![Go Reference](https://pkg.go.dev/badge/paepcke.de/uniex.svg)](https://pkg.go.dev/paepcke.de/uniex) 
[![Go Report Card](https://goreportcard.com/badge/paepcke.de/uniex)](https://goreportcard.com/report/paepcke.de/uniex) 
[![Go Build](https://github.com/paepckehh/uniex/actions/workflows/golang.yml/badge.svg)](https://github.com/paepckehh/uniex/actions/workflows/golang.yml)
[![License](https://img.shields.io/github/license/paepckehh/uniex)](https://github.com/paepckehh/uniex/blob/master/LICENSE)
[![SemVer](https://img.shields.io/github/v/release/paepckehh/uniex)](https://github.com/paepckehh/uniex/releases/latest)
<br>[![built with nix](https://builtwithnix.org/badge.svg)](https://search.nixos.org/packages?channel=unstable&from=0&size=50&sort=relevance&type=packages&query=uniex)

[paepcke.de/uniex](https://paepcke.de/uniex)

# UNIEX

- Exports UNIFI Controller Inventory Database (mongoDB) to [csv|json]
- Adds missing attributes (latest used hostname, latest used ip, latest seen, ...) via parsing all stats snippets
- Converts Timestamps (eg. unix nano time) into RFC3339 (parseable by excel, human readable) format
- Fast, even for large corp inventory (in-memory, parallel processing)
 
# SCREENSHOT CLI

![UNIEX SAMPLE SCREENSHOT](https://github.com/paepckehh/uniex/blob/main/resources/screenshot01.png "SCREEN")

# ⚡️QUICK START

```
go run paepcke.de/uniex/cmd/uniex@main
```

# ⚡️PRETTY PRINT OUTPUT VIA [BAT](https://github.com/sharkdp/bat) / [JQ](https://github.com/jqlang/jq) 

```
go run paepcke.de/uniex/cmd/uniex@main | bat -l csv
UNIEX_FORMAT=json go run paepcke.de/uniex/cmd/uniex@main | jq
```

# ⚡️HOW TO INSTALL

```
go install paepcke.de/uniex/cmd/uniex@main
```

# ⚡️PRE-BUILD BINARIES (DOWNLOAD)
[https://github.com/paepckehh/uniex/releases](https://github.com/paepckehh/uniex/releases)


# SUPPORTED OPTIONS 

```
# Optional (via env variables)
- UNIEX_MONGODB   - mongodb uri, default: mongodb://127.0.0.1:27117
- UNIEX_FORMAT    - export format, default: csv [csv|json]
```

# DOCS

[pkg.go.dev/paepcke.de/uniex](https://pkg.go.dev/paepcke.de/uniex)

# 🛡 License

[![License](https://img.shields.io/github/license/paepckehh/uniex)](https://github.com/paepckehh/uniex/blob/master/LICENSE)

This project is licensed under the terms of the `BSD 3-Clause License` license. See [LICENSE](https://github.com/paepckehh/uniex/blob/master/LICENSE) for more details.

# 📃 Citation

```bibtex
@misc{uniex,
  author = {Michael Paepcke},
  title = {Export UNIFI MongoDB Inventory Database},
  year = {2024},
  publisher = {GitHub},
  journal = {GitHub repository},
  howpublished = {\url{https://paepcke.de/uniex}}
}
```

# CONTRIBUTION

Yes, Please! PRs Welcome! 
