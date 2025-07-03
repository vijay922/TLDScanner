# TLD Scanner

A high-performance domain enumeration tool written in Go that identifies domains registered by the same organization across different Top-Level Domains (TLDs).

## Features

- **Fast Concurrent Scanning**: Multi-threaded WHOIS lookups with configurable concurrency
- **Comprehensive WHOIS Data**: Extracts organization, registrar, dates, nameservers, and status
- **Rate Limiting**: Built-in rate limiting to avoid overwhelming WHOIS servers
- **Multiple Output Formats**: Text and JSON output options
- **Extensive TLD Coverage**: Includes 500+ common TLDs and country codes
- **Error Handling**: Robust error handling and reporting
- **Progress Tracking**: Real-time progress indicators
- **Colorized Output**: Beautiful terminal output with color coding
- **Flexible Configuration**: Extensive command-line options

## Installation

### Prerequisites
- Go 1.21 or higher

### Install Dependencies
```bash
go mod init tldscanner
go mod tidy
```

### Build
```bash
go build -o tldscanner tldscanner.go
```

## Usage

### Basic Usage
```bash
# Basic scan
./tldscanner -d example.com

# Scan with custom wordlist
./tldscanner -d example.com -w custom_wordlist.txt

# Save results to file
./tldscanner -d example.com -o results.txt

# JSON output
./tldscanner -d example.com -json -o results.json
```

### Advanced Usage
```bash
# High-performance scan with 50 threads
./tldscanner -d example.com -t 50 -r 50

# Verbose output with all domain information
./tldscanner -d example.com -v -all

# Save all scanned domains (not just matches)
./tldscanner -d example.com -all -o complete_results.txt

# Custom timeout and rate limiting
./tldscanner -d example.com -timeout 60 -r 200
```

## Command Line Options

| Option | Description | Default |
|--------|-------------|---------|
| `-d` | Target domain to analyze (required) | - |
| `-w` | Path to TLD wordlist file | `wordlist.txt` |
| `-o` | Output file path | stdout |
| `-t` | Number of concurrent threads | `10` |
| `-timeout` | WHOIS timeout in seconds | `30` |
| `-r` | Rate limit in milliseconds between requests | `100` |
| `-v` | Verbose output | `false` |
| `-json` | Output in JSON format | `false` |
| `-all` | Save all domain results (not just matches) | `false` |
| `-h` | Show help message | - |

## Output Formats

### Text Output
```
=== TLD SCANNER RESULTS ===
Target Domain: example.com
Target Organization: Example Corp
Scan Duration: 2m30s
Total Scanned: 500
Total Matches: 5
Total Errors: 12

=== MATCHING DOMAINS ===
[+] example.net
    Organization: Example Corp
    Registrar: GoDaddy.com
    Created: 2020-01-15
    Expires: 2025-01-15
    Name Servers: ns1.example.com, ns2.example.com
```

### JSON Output
```json
{
  "target_domain": "example.com",
  "target_organization": "Example Corp",
  "matching_domains": [
    {
      "domain": "example.net",
      "organization": "Example Corp",
      "registrar": "GoDaddy.com",
      "created_date": "2020-01-15",
      "expiry_date": "2025-01-15",
      "status": "clientTransferProhibited",
      "name_servers": ["ns1.example.com", "ns2.example.com"],
      "timestamp": "2024-01-15T10:30:00Z"
    }
  ],
  "scan_duration": "2m30s",
  "total_scanned": 500,
  "total_matches": 5,
  "total_errors": 12
}
```

## Wordlist Format

The wordlist file should contain one TLD per line:
```
com
net
org
edu
gov
co.uk
com.au
# Comments start with #
```

## Performance Tips

1. **Adjust Thread Count**: Use `-t` to increase concurrent requests
   ```bash
   ./tldscanner -d example.com -t 50
   ```

2. **Reduce Rate Limiting**: Lower `-r` value for faster scanning (be respectful)
   ```bash
   ./tldscanner -d example.com -r 50
   ```

3. **Increase Timeout**: For slow WHOIS servers
   ```bash
   ./tldscanner
