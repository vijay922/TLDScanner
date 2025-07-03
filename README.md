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
   ./tldscanner -d example.com -timeout 60
   ```

4. **Use Smaller Wordlists**: Focus on specific TLDs for faster results
   ```bash
   ./tldscanner -d example.com -w common_tlds.txt
   ```

## Use Cases

### Cybersecurity & Penetration Testing
- **Domain Reconnaissance**: Discover all domains owned by a target organization
- **Attack Surface Mapping**: Identify potential entry points and infrastructure
- **Brand Protection**: Find domains that might be used for phishing or fraud
- **Threat Intelligence**: Track domain registrations by malicious actors

### Brand Monitoring
- **Trademark Protection**: Monitor for unauthorized use of brand names
- **Competitor Analysis**: Discover competitor domain portfolios
- **Domain Squatting Detection**: Find potentially infringing domains

### Business Intelligence
- **Market Research**: Analyze competitor digital presence
- **Due Diligence**: Verify organization domain ownership
- **Asset Discovery**: Complete inventory of digital assets

## Security Considerations

1. **Rate Limiting**: The tool includes built-in rate limiting to avoid overwhelming WHOIS servers
2. **Respectful Usage**: Use reasonable thread counts and delays
3. **Legal Compliance**: Ensure your usage complies with applicable laws and terms of service
4. **Data Privacy**: Be mindful of how you store and share discovered information

## Troubleshooting

### Common Issues

1. **"No organization found"**
   - Some domains don't have organization information in WHOIS
   - Try with a different domain or check WHOIS manually

2. **High error rates**
   - Increase timeout: `-timeout 60`
   - Reduce thread count: `-t 5`
   - Increase rate limit: `-r 500`

3. **Slow performance**
   - Increase thread count: `-t 20`
   - Reduce rate limit: `-r 50`
   - Use smaller wordlist

4. **Memory issues with large scans**
   - Avoid `-all` flag for large wordlists
   - Use streaming output instead of storing all results

### Debug Mode
```bash
# Enable verbose logging
./tldscanner -d example.com -v

# Save all results for analysis
./tldscanner -d example.com -all -o debug_results.json -json
```

## Examples

### Basic Domain Enumeration
```bash
# Scan example.com for related domains
./tldscanner -d example.com -w wordlist.txt -o results.txt

# Expected output:
# [+] MATCH: example.net -> Example Corp
# [+] MATCH: example.org -> Example Corp
# [+] MATCH: example.co.uk -> Example Corp
```

### High-Performance Scan
```bash
# Fast scan with 30 threads and minimal delays
./tldscanner -d target.com -t 30 -r 25 -timeout 45 -o fast_results.txt
```

### Comprehensive Analysis
```bash
# Full scan with detailed output
./tldscanner -d company.com -v -all -json -o comprehensive_analysis.json
```

### Custom Wordlist
```bash
# Create a custom wordlist for specific regions
echo -e "com\nnet\norg\nco.uk\ncom.au\nca\nde\nfr" > custom.txt
./tldscanner -d example.com -w custom.txt
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

This tool is for educational and authorized security testing purposes only. Users are responsible for ensuring their use complies with applicable laws and regulations. The authors are not responsible for any misuse of this tool.

## Changelog

### v2.0.0
- Complete rewrite in Go for better performance
- Added JSON output format
- Implemented concurrent scanning with rate limiting
- Added comprehensive error handling
- Improved WHOIS data extraction
- Added colorized terminal output
- Enhanced command-line interface

### v1.0.0
- Initial Python implementation
- Basic domain enumeration
- Simple text output

## Support

For issues and questions:
- Open an issue on GitHub
- Check the troubleshooting section
- Review existing issues for solutions

## Related Tools

- [Sublist3r](https://github.com/aboul3la/Sublist3r) - Subdomain enumeration
- [Amass](https://github.com/OWASP/Amass) - Network mapping and asset discovery
- [DNSRecon](https://github.com/darkoperator/dnsrecon) - DNS enumeration
- [Fierce](https://github.com/mschwager/fierce) - Domain scanner
