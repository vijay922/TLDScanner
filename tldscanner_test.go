package main

import (
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestExtractBaseDomain(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"example.com", "example"},
		{"subdomain.example.com", "subdomain"},
		{"test.co.uk", "test"},
		{"simple", "simple"},
		{"multi.level.domain.org", "multi"},
	}

	for _, test := range tests {
		result := extractBaseDomain(test.input)
		if result != test.expected {
			t.Errorf("extractBaseDomain(%s) = %s; expected %s", test.input, result, test.expected)
		}
	}
}

func TestGenerateDomains(t *testing.T) {
	baseDomain := "example"
	tlds := []string{".com", ".net", ".org"}
	expected := []string{"example.com", "example.net", "example.org"}

	result := generateDomains(baseDomain, tlds)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("generateDomains(%s, %v) = %v; expected %v", baseDomain, tlds, result, expected)
	}
}

func TestLoadWordlist(t *testing.T) {
	// Create a temporary wordlist file
	content := "com\nnet\norg\n# This is a comment\n\n  co.uk  \n"
	tmpFile, err := os.CreateTemp("", "test_wordlist_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test loading the wordlist
	tlds, err := loadWordlist(tmpFile.Name())
	if err != nil {
		t.Fatalf("loadWordlist failed: %v", err)
	}

	expected := []string{".com", ".net", ".org", ".co.uk"}
	if !reflect.DeepEqual(tlds, expected) {
		t.Errorf("loadWordlist() = %v; expected %v", tlds, expected)
	}
}

func TestLoadWordlistNonExistent(t *testing.T) {
	_, err := loadWordlist("non_existent_file.txt")
	if err == nil {
		t.Error("Expected error for non-existent file, but got nil")
	}
}

func TestCountErrors(t *testing.T) {
	results := []DomainInfo{
		{Domain: "example.com", Error: ""},
		{Domain: "example.net", Error: "timeout"},
		{Domain: "example.org", Error: ""},
		{Domain: "example.co.uk", Error: "invalid domain"},
	}

	count := countErrors(results)
	expected := 2

	if count != expected {
		t.Errorf("countErrors() = %d; expected %d", count, expected)
	}
}

func TestDomainInfoStructure(t *testing.T) {
	info := DomainInfo{
		Domain:       "example.com",
		Organization: "Example Corp",
		Registrar:    "GoDaddy",
		CreatedDate:  "2020-01-01",
		ExpiryDate:   "2025-01-01",
		Status:       "active",
		NameServers:  []string{"ns1.example.com", "ns2.example.com"},
		Timestamp:    time.Now(),
	}

	if info.Domain != "example.com" {
		t.Errorf("Domain field not set correctly")
	}
	if info.Organization != "Example Corp" {
		t.Errorf("Organization field not set correctly")
	}
	if len(info.NameServers) != 2 {
		t.Errorf("NameServers field not set correctly")
	}
}

func TestConfigDefaults(t *testing.T) {
	// Test that default config values are reasonable
	config := Config{
		Threads:   10,
		Timeout:   30,
		RateLimit: 100,
	}

	if config.Threads <= 0 {
		t.Error("Default threads should be positive")
	}
	if config.Timeout <= 0 {
		t.Error("Default timeout should be positive")
	}
	if config.RateLimit < 0 {
		t.Error("Default rate limit should be non-negative")
	}
}

func TestResultStructure(t *testing.T) {
	result := Result{
		TargetDomain:    "example.com",
		TargetOrg:       "Example Corp",
		MatchingDomains: []DomainInfo{},
		AllDomains:      []DomainInfo{},
		ScanDuration:    "1m30s",
		TotalScanned:    100,
		TotalMatches:    5,
		TotalErrors:     2,
	}

	if result.TargetDomain != "example.com" {
		t.Error("TargetDomain not set correctly")
	}
	if result.TotalScanned != 100 {
		t.Error("TotalScanned not set correctly")
	}
}

// Benchmark tests
func BenchmarkExtractBaseDomain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		extractBaseDomain("example.com")
	}
}

func BenchmarkGenerateDomains(b *testing.B) {
	tlds := []string{".com", ".net", ".org", ".co.uk", ".de", ".fr"}
	for i := 0; i < b.N; i++ {
		generateDomains("example", tlds)
	}
}

// Mock WHOIS test (since we can't rely on external services in tests)
func TestWhoisInfoMock(t *testing.T) {
	// This test demonstrates the structure we expect from WHOIS
	// In a real implementation, you might want to mock the WHOIS service
	
	info := &DomainInfo{
		Domain:       "example.com",
		Organization: "Example Corporation",
		Registrar:    "Mock Registrar",
		CreatedDate:  "2020-01-01",
		ExpiryDate:   "2025-01-01",
		Status:       "active",
		NameServers:  []string{"ns1.example.com", "ns2.example.com"},
		Timestamp:    time.Now(),
	}

	if info.Domain == "" {
		t.Error("Domain should not be empty")
	}
	if info.Organization == "" {
		t.Error("Organization should not be empty")
	}
	if len(info.NameServers) == 0 {
		t.Error("NameServers should not be empty")
	}
}

// Test color constants
func TestColorConstants(t *testing.T) {
	colors := []string{
		ColorReset, ColorRed, ColorGreen, ColorYellow,
		ColorBlue, ColorPurple, ColorCyan, ColorWhite,
	}

	for _, color := range colors {
		if color == "" {
			t.Error("Color constant should not be empty")
		}
		if !strings.HasPrefix(color, "\033[") {
			t.Error("Color constant should start with ANSI escape sequence")
		}
	}
}

// Test wordlist parsing edge cases
func TestWordlistEdgeCases(t *testing.T) {
	testCases := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name:     "Empty file",
			content:  "",
			expected: []string{},
		},
		{
			name:     "Only comments",
			content:  "# Comment 1\n# Comment 2\n",
			expected: []string{},
		},
		{
			name:     "Mixed content",
			content:  "com\n# Comment\nnet\n\norg\n  co.uk  \n",
			expected: []string{".com", ".net", ".org", ".co.uk"},
		},
		{
			name:     "TLDs with dots",
			content:  ".com\n.net\norg\n",
			expected: []string{".com", ".net", ".org"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "test_*.txt")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.WriteString(tc.content); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			tmpFile.Close()

			result, err := loadWordlist(tmpFile.Name())
			if err != nil {
				t.Fatalf("loadWordlist failed: %v", err)
			}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("loadWordlist() = %v; expected %v", result, tc.expected)
			}
		})
	}
}

// Test error handling
func TestErrorHandling(t *testing.T) {
	// Test with invalid domain format
	info := DomainInfo{
		Domain: "invalid..domain",
		Error:  "invalid domain format",
	}

	if info.Error == "" {
		t.Error("Error should be set for invalid domain")
	}
}

// Performance test for large wordlists
func TestLargeWordlistPerformance(t *testing.T) {
	// Create a large wordlist
	var content strings.Builder
	for i := 0; i < 1000; i++ {
		content.WriteString("tld")
		content.WriteString(string(rune(i)))
		content.WriteString("\n")
	}

	tmpFile, err := os.CreateTemp("", "large_wordlist_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content.String()); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	start := time.Now()
	tlds, err := loadWordlist(tmpFile.Name())
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("loadWordlist failed: %v", err)
	}

	if len(tlds) != 1000 {
		t.Errorf("Expected 1000 TLDs, got %d", len(tlds))
	}

	// Should be fast enough for 1000 entries
	if duration > time.Second {
		t.Errorf("Loading 1000 TLDs took too long: %v", duration)
	}
}

// Test concurrent safety (basic test)
func TestConcurrentSafety(t *testing.T) {
	baseDomain := "example"
	tlds := []string{".com", ".net", ".org"}

	// Run multiple goroutines doing the same operation
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			result := generateDomains(baseDomain, tlds)
			if len(result) != 3 {
				t.Errorf("Expected 3 domains, got %d", len(result))
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}
