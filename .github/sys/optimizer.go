// Package sys - Kod optimizasyon modülü
// Performance ve boyut optimizasyonu sağlar
package sys

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

// Optimizer - Kod optimizasyon motoru
type Optimizer struct {
	source        string
	optimizations []OptimizationRule
	stats         *OptimizationStats
	mu            sync.RWMutex
}

// OptimizationStats - Optimizasyon istatistikleri
type OptimizationStats struct {
	RemovedLines       int
	RemovedComments    int
	RemovedWhitespace  int
	InlinedFunctions   int
	DeadCodeRemoved    int
	ConstantsOptimized int
}

// OptimizationRule - Optimizasyon kuralı
type OptimizationRule struct {
	Name        string
	Pattern     *regexp.Regexp
	Replacement string
	Active      bool
}

// NewOptimizer - Yeni optimizer oluştur
func NewOptimizer(source string) *Optimizer {
	opt := &Optimizer{
		source: source,
		stats:  &OptimizationStats{},
		optimizations: []OptimizationRule{
			{
				Name:    "RemoveComments",
				Pattern: regexp.MustCompile(`//.*?(?:\n|$)`),
				Active:  true,
			},
			{
				Name:    "RemoveMultiLineComments",
				Pattern: regexp.MustCompile(`/\*.*?\*/`),
				Active:  true,
			},
			{
				Name:    "RemoveExtraWhitespace",
				Pattern: regexp.MustCompile(`\s+`),
				Active:  true,
			},
			{
				Name:    "SimplifyExpressions",
				Pattern: regexp.MustCompile(`\(([^()]+)\)`),
				Active:  true,
			},
		},
	}
	return opt
}

// Optimize - Kodu optimize et
func (o *Optimizer) Optimize() (string, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	result := o.source

	// 1. Boş satırları kaldır
	result = o.removeEmptyLines(result)

	// 2. Yorumları kaldır
	result = o.removeComments(result)

	// 3. Whitespace'i optimize et
	result = o.optimizeWhitespace(result)

	// 4. Dead code'u kaldır
	result = o.removeDeadCode(result)

	// 5. Constantları optimize et
	result = o.optimizeConstants(result)

	return result, nil
}

// removeEmptyLines - Boş satırları kaldır
func (o *Optimizer) removeEmptyLines(source string) string {
	lines := strings.Split(source, "\n")
	var filtered []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			filtered = append(filtered, line)
			o.stats.RemovedLines++
		}
	}

	return strings.Join(filtered, "\n")
}

// removeComments - Yorumları kaldır
func (o *Optimizer) removeComments(source string) string {
	// Single-line comments
	result := regexp.MustCompile(`//[^\n]*`).ReplaceAllString(source, "")
	o.stats.RemovedComments += strings.Count(source, "//")

	// Multi-line comments
	result = regexp.MustCompile(`(?s)/\*.*?\*/`).ReplaceAllString(result, "")
	o.stats.RemovedComments += strings.Count(source, "/*")

	return result
}

// optimizeWhitespace - Whitespace'i optimize et
func (o *Optimizer) optimizeWhitespace(source string) string {
	// Multiple spaces to single space
	result := regexp.MustCompile(`[ \t]+`).ReplaceAllString(source, " ")

	// Trailing whitespace remove
	lines := strings.Split(result, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
		o.stats.RemovedWhitespace += len(line) - len(lines[i])
	}

	return strings.Join(lines, "\n")
}

// removeDeadCode - Dead code'u kaldır
func (o *Optimizer) removeDeadCode(source string) string {
	// Unreachable code after return
	result := regexp.MustCompile(`(?m)return\s+.*?\n\s+.*?\n`).ReplaceAllStringFunc(result, func(match string) string {
		lines := strings.Split(match, "\n")
		if len(lines) > 1 {
			o.stats.DeadCodeRemoved++
			return lines[0] + "\n"
		}
		return match
	})

	return result
}

// optimizeConstants - Constants'ı optimize et
func (o *Optimizer) optimizeConstants(source string) string {
	// const x = 5; y = x + 1; => y = 6;
	pattern := regexp.MustCompile(`const\s+(\w+)\s*=\s*(\d+)`)
	matches := pattern.FindAllStringSubmatch(source, -1)

	result := source
	for _, match := range matches {
		constName := match[1]
		constValue := match[2]

		// Replace const usage with value
		usage := regexp.MustCompile(`\b` + constName + `\b`)
		if usage.MatchString(result) {
			result = usage.ReplaceAllString(result, constValue)
			o.stats.ConstantsOptimized++
		}
	}

	return result
}

// GetStats - Optimizasyon istatistiklerini döndür
func (o *Optimizer) GetStats() *OptimizationStats {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.stats
}

// SetOptimization - Specific optimizasyonu aç/kapat
func (o *Optimizer) SetOptimization(name string, active bool) {
	o.mu.Lock()
	defer o.mu.Unlock()

	for i, opt := range o.optimizations {
		if opt.Name == name {
			o.optimizations[i].Active = active
			break
		}
	}
}

// CodeAnalyzer - Kod analiz sistemi
type CodeAnalyzer struct {
	source       string
	metrics      *CodeMetrics
	issues       []CodeIssue
	mu           sync.RWMutex
}

// CodeMetrics - Kod metrikleri
type CodeMetrics struct {
	Lines        int
	Functions    int
	Branches     int
	Complexity   float64
	CoverageRate float64
}

// CodeIssue - Kod problemi
type CodeIssue struct {
	Type     string // warning, error, info
	Message  string
	Line     int
	Severity int // 1-10
}

// NewCodeAnalyzer - Yeni analizci oluştur
func NewCodeAnalyzer(source string) *CodeAnalyzer {
	return &CodeAnalyzer{
		source: source,
		metrics: &CodeMetrics{},
		issues: make([]CodeIssue, 0),
	}
}

// Analyze - Kodu analiz et
func (ca *CodeAnalyzer) Analyze() error {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	// Metrik hesapla
	ca.calculateMetrics()

	// Problemleri tespit et
	ca.detectIssues()

	return nil
}

// calculateMetrics - Kod metriklerini hesapla
func (ca *CodeAnalyzer) calculateMetrics() {
	lines := strings.Split(ca.source, "\n")
	ca.metrics.Lines = len(lines)

	// Function sayısı
	ca.metrics.Functions = strings.Count(ca.source, "func ")

	// Branch sayısı
	ca.metrics.Branches = strings.Count(ca.source, "if ") + 
		strings.Count(ca.source, "for ") + 
		strings.Count(ca.source, "else")

	// Cyclomatic complexity (basit)
	ca.metrics.Complexity = float64(ca.metrics.Branches) / float64(ca.metrics.Functions + 1)

	// Coverage rate (tahmin)
	ca.metrics.CoverageRate = 0.75 // Default
}

// detectIssues - Problemleri tespit et
func (ca *CodeAnalyzer) detectIssues() {
	lines := strings.Split(ca.source, "\n")

	for i, line := range lines {
		// Long lines
		if len(line) > 120 {
			ca.issues = append(ca.issues, CodeIssue{
				Type:     "warning",
				Message:  fmt.Sprintf("Çok uzun satır (%d karakter)", len(line)),
				Line:     i + 1,
				Severity: 3,
			})
		}

		// Multiple statements in one line
		if strings.Count(line, ";") > 2 {
			ca.issues = append(ca.issues, CodeIssue{
				Type:     "warning",
				Message:  "Bir satırda çok fazla statement",
				Line:     i + 1,
				Severity: 2,
			})
		}

		// Unused variables (simple check)
		if strings.Contains(line, "var ") && !strings.Contains(ca.source[strings.Index(ca.source, line):], "=") {
			ca.issues = append(ca.issues, CodeIssue{
				Type:     "warning",
				Message:  "Potential unused variable",
				Line:     i + 1,
				Severity: 4,
			})
		}
	}

	// Complexity check
	if ca.metrics.Complexity > 15 {
		ca.issues = append(ca.issues, CodeIssue{
			Type:     "error",
			Message:  fmt.Sprintf("Çok yüksek cyclomatic complexity: %.2f", ca.metrics.Complexity),
			Line:     1,
			Severity: 8,
		})
	}
}

// GetMetrics - Metrikleri döndür
func (ca *CodeAnalyzer) GetMetrics() *CodeMetrics {
	ca.mu.RLock()
	defer ca.mu.RUnlock()
	return ca.metrics
}

// GetIssues - Problemleri döndür
func (ca *CodeAnalyzer) GetIssues() []CodeIssue {
	ca.mu.RLock()
	defer ca.mu.RUnlock()
	return ca.issues
}

// ReportIssues - Problemleri rapor et
func (ca *CodeAnalyzer) ReportIssues() string {
	ca.mu.RLock()
	defer ca.mu.RUnlock()

	if len(ca.issues) == 0 {
		return "Sorun tespit edilmedi"
	}

	var report strings.Builder
	report.WriteString(fmt.Sprintf("=== Kod Analiz Raporu ===\n"))
	report.WriteString(fmt.Sprintf("Toplam Sorun: %d\n\n", len(ca.issues)))

	for _, issue := range ca.issues {
		report.WriteString(fmt.Sprintf("[%s] (Satır %d, Önem: %d/10)\n", strings.ToUpper(issue.Type), issue.Line, issue.Severity))
		report.WriteString(fmt.Sprintf("  → %s\n\n", issue.Message))
	}

	return report.String()
}
