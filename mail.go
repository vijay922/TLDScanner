package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

// Config holds the application configuration
type Config struct {
	Domain      string
	Wordlist    string
	Output      string
	Threads     int
	Timeout     int
	Verbose     bool
	JSONOutput  bool
	SaveAll     bool
	RateLimit   int
}

// DomainInfo represents domain information
type DomainInfo struct {
	Domain       string    `json:"domain"`
	Organization string    `json:"organization"`
	Registrar    string    `json:"registrar"`
	CreatedDate  string    `json:"created_date"`
	ExpiryDate   string    `json:"expiry_date"`
	Status       string    `json:"status"`
	NameServers  []string  `json:"name_servers"`
	Error        string    `json:"error,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}

// Result holds the scan results
type Result struct {
	TargetDomain     string       `json:"target_domain"`
	TargetOrg        string       `json:"target_organization"`
	MatchingDomains  []DomainInfo `json:"matching_domains"`
	AllDomains       []DomainInfo `json:"all_domains,omitempty"`
	ScanDuration     string       `json:"scan_duration"`
	TotalScanned     int          `json:"total_scanned"`
	TotalMatches     int          `json:"total_matches"`
	TotalErrors      int          `json:"total_errors"`
}

// Colors for terminal output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

func main() {
	config := parseFlags()
	
	if config.Domain == "" {
		fmt.Fprintf(os.Stderr, "%s[ERROR]%s Domain is required. Use -h for help.\n", ColorRed, ColorReset)
		os.Exit(1)
	}

	// Print banner
	printBanner()

	// Get target domain organization
	fmt.Printf("%s[INFO]%s Analyzing target domain: %s\n", ColorBlue, ColorReset, config.Domain)
	targetInfo, err := getWhoisInfo(config.Domain, config.Timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s[ERROR]%s Failed to get WHOIS info for %s: %v\n", ColorRed, ColorReset, config.Domain, err)
		os.Exit(1)
	}

	if targetInfo.Organization == "" {
		fmt.Fprintf(os.Stderr, "%s[WARNING]%s No organization found for %s\n", ColorYellow, ColorReset, config.Domain)
		os.Exit(1)
	}

	fmt.Printf("%s[INFO]%s Target organization: %s%s%s\n", ColorBlue, ColorReset, ColorGreen, targetInfo.Organization, ColorReset)

	// Load TLD wordlist
	tlds, err := loadWordlist(config.Wordlist)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s[ERROR]%s Failed to load wordlist: %v\n", ColorRed, ColorReset, err)
		os.Exit(1)
	}

	fmt.Printf("%s[INFO]%s Loaded %d TLDs from wordlist\n", ColorBlue, ColorReset, len(tlds))

	// Generate domain list
	baseDomain := extractBaseDomain(config.Domain)
	domains := generateDomains(baseDomain, tlds)

	fmt.Printf("%s[INFO]%s Starting scan of %d domains with %d threads...\n", ColorBlue, ColorReset, len(domains), config.Threads)

	// Perform scan
	startTime := time.Now()
	allResults, matchingResults := scanDomains(domains, targetInfo.Organization, config)
	scanDuration := time.Since(startTime)

	// Prepare results
	result := Result{
		TargetDomain:    config.Domain,
		TargetOrg:       targetInfo.Organization,
		MatchingDomains: matchingResults,
		ScanDuration:    scanDuration.String(),
		TotalScanned:    len(domains),
		TotalMatches:    len(matchingResults),
		TotalErrors:     countErrors(allResults),
	}

	if config.SaveAll {
		result.AllDomains = allResults
	}

	// Output results
	if config.JSONOutput {
		outputJSON(result, config.Output)
	} else {
		outputText(result, config.Output, config.Verbose)
	}

	// Print summary
	printSummary(result)
}

func parseFlags() Config {
	var config Config

	flag.StringVar(&config.Domain, "d", "", "Target domain to analyze (required)")
	flag.StringVar(&config.Wordlist, "w", "wordlist.txt", "Path to TLD wordlist file")
	flag.StringVar(&config.Output, "o", "", "Output file path (optional)")
	flag.IntVar(&config.Threads, "t", 10, "Number of concurrent threads")
	flag.IntVar(&config.Timeout, "timeout", 30, "WHOIS timeout in seconds")
	flag.BoolVar(&config.Verbose, "v", false, "Verbose output")
	flag.BoolVar(&config.JSONOutput, "json", false, "Output in JSON format")
	flag.BoolVar(&config.SaveAll, "all", false, "Save all domain results (not just matches)")
	flag.IntVar(&config.RateLimit, "r", 100, "Rate limit in milliseconds between requests")

	flag.Usage = func() {
		fmt.Printf("%sTLD Scanner - Domain Enumeration Tool%s\n\n", ColorCyan, ColorReset)
		fmt.Printf("Usage: %s [OPTIONS]\n\n", os.Args[0])
		fmt.Printf("Options:\n")
		flag.PrintDefaults()
		fmt.Printf("\nExample:\n")
		fmt.Printf("  %s -d example.com -w wordlist.txt -o results.txt -t 20 -v\n", os.Args[0])
		fmt.Printf("  %s -d example.com -json -o results.json -all\n", os.Args[0])
	}

	flag.Parse()
	return config
}

func printBanner() {
	banner := `
████████╗██╗     ██████╗     ███████╗ ██████╗ █████╗ ███╗   ██╗███╗   ██╗███████╗██████╗ 
╚══██╔══╝██║     ██╔══██╗    ██╔════╝██╔════╝██╔══██╗████╗  ██║████╗  ██║██╔════╝██╔══██╗
   ██║   ██║     ██║  ██║    ███████╗██║     ███████║██╔██╗ ██║██╔██╗ ██║█████╗  ██████╔╝
   ██║   ██║     ██║  ██║    ╚════██║██║     ██╔══██║██║╚██╗██║██║╚██╗██║██╔══╝  ██╔══██╗
   ██║   ███████╗██████╔╝    ███████║╚██████╗██║  ██║██║ ╚████║██║ ╚████║███████╗██║  ██║
   ╚═╝   ╚══════╝╚═════╝     ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝
`
	fmt.Printf("%s%s%s\n", ColorCyan, banner, ColorReset)
	fmt.Printf("%s                    Domain Enumeration Tool v2.0%s\n", ColorYellow, ColorReset)
	fmt.Printf("%s                    github.com/yourusername/tldscanner%s\n\n", ColorPurple, ColorReset)
}

func getWhoisInfo(domain string, timeout int) (*DomainInfo, error) {
	whoisRaw, err := whois.Whois(domain)
	if err != nil {
		return nil, fmt.Errorf("whois query failed: %w", err)
	}

	result, err := whoisparser.Parse(whoisRaw)
	if err != nil {
		return nil, fmt.Errorf("whois parsing failed: %w", err)
	}

	var nameServers []string
	for _, ns := range result.Domain.NameServers {
		nameServers = append(nameServers, ns)
	}

	return &DomainInfo{
		Domain:       domain,
		Organization: result.Registrant.Organization,
		Registrar:    result.Registrar.Name,
		CreatedDate:  result.Domain.CreatedDate,
		ExpiryDate:   result.Domain.ExpirationDate,
		Status:       strings.Join(result.Domain.Status, ", "),
		NameServers:  nameServers,
		Timestamp:    time.Now(),
	}, nil
}

func loadWordlist(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open wordlist file: %w", err)
	}
	defer file.Close()

	var tlds []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tld := strings.TrimSpace(scanner.Text())
		if tld != "" && !strings.HasPrefix(tld, "#") {
			if !strings.HasPrefix(tld, ".") {
				tld = "." + tld
			}
			tlds = append(tlds, tld)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading wordlist: %w", err)
	}

	return tlds, nil
}

func extractBaseDomain(domain string) string {
	parts := strings.Split(domain, ".")
	if len(parts) >= 2 {
		return parts[0]
	}
	return domain
}

func generateDomains(baseDomain string, tlds []string) []string {
	var domains []string
	for _, tld := range tlds {
		domains = append(domains, baseDomain+tld)
	}
	return domains
}

func scanDomains(domains []string, targetOrg string, config Config) ([]DomainInfo, []DomainInfo) {
	var allResults []DomainInfo
	var matchingResults []DomainInfo
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Create a channel to limit concurrency
	semaphore := make(chan struct{}, config.Threads)
	
	// Rate limiting
	rateLimiter := time.NewTicker(time.Duration(config.RateLimit) * time.Millisecond)
	defer rateLimiter.Stop()

	processed := 0
	total := len(domains)

	for _, domain := range domains {
		wg.Add(1)
		
		go func(d string) {
			defer wg.Done()
			
			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			
			// Rate limiting
			<-rateLimiter.C
			
			info, err := getWhoisInfo(d, config.Timeout)
			if err != nil {
				info = &DomainInfo{
					Domain:    d,
					Error:     err.Error(),
					Timestamp: time.Now(),
				}
			}
			
			mu.Lock()
			allResults = append(allResults, *info)
			processed++
			
			// Check if organization matches
			if info.Organization != "" && strings.EqualFold(info.Organization, targetOrg) {
				matchingResults = append(matchingResults, *info)
				if !config.JSONOutput {
					fmt.Printf("%s[+] MATCH:%s %s -> %s%s%s\n", 
						ColorGreen, ColorReset, info.Domain, ColorYellow, info.Organization, ColorReset)
				}
			}
			
			if config.Verbose && !config.JSONOutput {
				if info.Error != "" {
					fmt.Printf("%s[!] ERROR:%s %s -> %s\n", ColorRed, ColorReset, info.Domain, info.Error)
				} else if info.Organization != "" {
					fmt.Printf("%s[-] CHECKED:%s %s -> %s\n", ColorWhite, ColorReset, info.Domain, info.Organization)
				}
			}
			
			// Progress indicator
			if !config.JSONOutput && !config.Verbose {
				fmt.Printf("\r%s[INFO]%s Progress: %d/%d domains scanned (%d matches)", 
					ColorBlue, ColorReset, processed, total, len(matchingResults))
			}
			mu.Unlock()
		}(domain)
	}

	wg.Wait()
	
	if !config.JSONOutput && !config.Verbose {
		fmt.Println() // New line after progress
	}

	// Sort results by domain name
	sort.Slice(allResults, func(i, j int) bool {
		return allResults[i].Domain < allResults[j].Domain
	})
	sort.Slice(matchingResults, func(i, j int) bool {
		return matchingResults[i].Domain < matchingResults[j].Domain
	})

	return allResults, matchingResults
}

func countErrors(results []DomainInfo) int {
	count := 0
	for _, result := range results {
		if result.Error != "" {
			count++
		}
	}
	return count
}

func outputJSON(result Result, outputFile string) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return
	}

	if outputFile != "" {
		err := os.WriteFile(outputFile, data, 0644)
		if err != nil {
			log.Printf("Error writing to file: %v", err)
			return
		}
		fmt.Printf("%s[INFO]%s Results saved to %s\n", ColorBlue, ColorReset, outputFile)
	} else {
		fmt.Println(string(data))
	}
}

func outputText(result Result, outputFile string, verbose bool) {
	var output strings.Builder
	
	output.WriteString(fmt.Sprintf("\n%s=== TLD SCANNER RESULTS ===%s\n", ColorCyan, ColorReset))
	output.WriteString(fmt.Sprintf("Target Domain: %s\n", result.TargetDomain))
	output.WriteString(fmt.Sprintf("Target Organization: %s\n", result.TargetOrg))
	output.WriteString(fmt.Sprintf("Scan Duration: %s\n", result.ScanDuration))
	output.WriteString(fmt.Sprintf("Total Scanned: %d\n", result.TotalScanned))
	output.WriteString(fmt.Sprintf("Total Matches: %d\n", result.TotalMatches))
	output.WriteString(fmt.Sprintf("Total Errors: %d\n\n", result.TotalErrors))

	if len(result.MatchingDomains) > 0 {
		output.WriteString(fmt.Sprintf("%s=== MATCHING DOMAINS ===%s\n", ColorGreen, ColorReset))
		for _, domain := range result.MatchingDomains {
			output.WriteString(fmt.Sprintf("[+] %s\n", domain.Domain))
			output.WriteString(fmt.Sprintf("    Organization: %s\n", domain.Organization))
			output.WriteString(fmt.Sprintf("    Registrar: %s\n", domain.Registrar))
			output.WriteString(fmt.Sprintf("    Created: %s\n", domain.CreatedDate))
			output.WriteString(fmt.Sprintf("    Expires: %s\n", domain.ExpiryDate))
			if len(domain.NameServers) > 0 {
				output.WriteString(fmt.Sprintf("    Name Servers: %s\n", strings.Join(domain.NameServers, ", ")))
			}
			output.WriteString("\n")
		}
	}

	if verbose && len(result.AllDomains) > 0 {
		output.WriteString(fmt.Sprintf("%s=== ALL SCANNED DOMAINS ===%s\n", ColorYellow, ColorReset))
		for _, domain := range result.AllDomains {
			if domain.Error != "" {
				output.WriteString(fmt.Sprintf("[!] %s -> ERROR: %s\n", domain.Domain, domain.Error))
			} else {
				output.WriteString(fmt.Sprintf("[-] %s -> %s\n", domain.Domain, domain.Organization))
			}
		}
	}

	if outputFile != "" {
		err := os.WriteFile(outputFile, []byte(output.String()), 0644)
		if err != nil {
			log.Printf("Error writing to file: %v", err)
			return
		}
		fmt.Printf("%s[INFO]%s Results saved to %s\n", ColorBlue, ColorReset, outputFile)
	} else {
		fmt.Print(output.String())
	}
}

func printSummary(result Result) {
	fmt.Printf("\n%s=== SCAN SUMMARY ===%s\n", ColorCyan, ColorReset)
	fmt.Printf("Domains Scanned: %s%d%s\n", ColorWhite, result.TotalScanned, ColorReset)
	fmt.Printf("Matches Found: %s%d%s\n", ColorGreen, result.TotalMatches, ColorReset)
	fmt.Printf("Errors: %s%d%s\n", ColorRed, result.TotalErrors, ColorReset)
	fmt.Printf("Duration: %s%s%s\n", ColorYellow, result.ScanDuration, ColorReset)
	fmt.Printf("Rate: %s%.2f domains/second%s\n", ColorPurple, 
		float64(result.TotalScanned)/time.Since(time.Now().Add(-parseDuration(result.ScanDuration))).Seconds(), ColorReset)
}

func parseDuration(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}
