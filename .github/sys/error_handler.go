// Package sys - Hata yönetim sistemi
// Derleyici hataları ve uyarılarını yönetir
package sys

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// ErrorSeverity - Hata şiddeti
type ErrorSeverity int

const (
	SeverityInfo    ErrorSeverity = iota // 0
	SeverityWarning                      // 1
	SeverityError                        // 2
	SeverityCritical                     // 3
)

// String - Severity'nin string temsilini döndür
func (s ErrorSeverity) String() string {
	switch s {
	case SeverityInfo:
		return "INFO"
	case SeverityWarning:
		return "WARNING"
	case SeverityError:
		return "ERROR"
	case SeverityCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

// CompileError - Derleme hatası yapısı
type CompileError struct {
	ID        string
	Severity  ErrorSeverity
	Message   string
	Location  ErrorLocation
	Context   string
	Timestamp time.Time
	Fixable   bool
	Suggestion string
}

// ErrorLocation - Hata konumu
type ErrorLocation struct {
	File   string
	Line   int
	Column int
	Length int
}

// ErrorManager - Hata yönetim sistemi
type ErrorManager struct {
	errors    []*CompileError
	warnings  []*CompileError
	maxErrors int
	mu        sync.RWMutex
}

// NewErrorManager - Yeni hata yöneticisi oluştur
func NewErrorManager(maxErrors int) *ErrorManager {
	return &ErrorManager{
		errors:    make([]*CompileError, 0),
		warnings:  make([]*CompileError, 0),
		maxErrors: maxErrors,
	}
}

// AddError - Hata ekle
func (em *ErrorManager) AddError(severity ErrorSeverity, message string, location ErrorLocation) error {
	em.mu.Lock()
	defer em.mu.Unlock()

	if len(em.errors) >= em.maxErrors {
		return fmt.Errorf("maksimum hata sayısına ulaşıldı (%d)", em.maxErrors)
	}

	err := &CompileError{
		ID:        em.generateErrorID(),
		Severity:  severity,
		Message:   message,
		Location:  location,
		Timestamp: time.Now(),
	}

	if severity == SeverityWarning {
		em.warnings = append(em.warnings, err)
	} else {
		em.errors = append(em.errors, err)
	}

	return nil
}

// AddErrorWithSuggestion - Önerili hata ekle
func (em *ErrorManager) AddErrorWithSuggestion(severity ErrorSeverity, message string, location ErrorLocation, suggestion string) {
	em.mu.Lock()
	defer em.mu.Unlock()

	err := &CompileError{
		ID:         em.generateErrorID(),
		Severity:   severity,
		Message:    message,
		Location:   location,
		Timestamp:  time.Now(),
		Fixable:    true,
		Suggestion: suggestion,
	}

	em.errors = append(em.errors, err)
}

// GetErrors - Tüm hataları döndür
func (em *ErrorManager) GetErrors() []*CompileError {
	em.mu.RLock()
	defer em.mu.RUnlock()

	all := make([]*CompileError, 0, len(em.errors)+len(em.warnings))
	all = append(all, em.errors...)
	all = append(all, em.warnings...)
	return all
}

// GetErrorsByFile - Dosyaya göre hataları döndür
func (em *ErrorManager) GetErrorsByFile(file string) []*CompileError {
	em.mu.RLock()
	defer em.mu.RUnlock()

	var filtered []*CompileError
	for _, err := range em.errors {
		if err.Location.File == file {
			filtered = append(filtered, err)
		}
	}
	for _, warn := range em.warnings {
		if warn.Location.File == file {
			filtered = append(filtered, warn)
		}
	}
	return filtered
}

// GetErrorCount - Hata sayısını döndür
func (em *ErrorManager) GetErrorCount() int {
	em.mu.RLock()
	defer em.mu.RUnlock()
	return len(em.errors)
}

// GetWarningCount - Uyarı sayısını döndür
func (em *ErrorManager) GetWarningCount() int {
	em.mu.RLock()
	defer em.mu.RUnlock()
	return len(em.warnings)
}

// HasErrors - Hata var mı kontrol et
func (em *ErrorManager) HasErrors() bool {
	em.mu.RLock()
	defer em.mu.RUnlock()
	return len(em.errors) > 0
}

// Clear - Tüm hataları temizle
func (em *ErrorManager) Clear() {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.errors = make([]*CompileError, 0)
	em.warnings = make([]*CompileError, 0)
}

// Print - Hataları yazdır
func (em *ErrorManager) Print() string {
	em.mu.RLock()
	defer em.mu.RUnlock()

	if len(em.errors) == 0 && len(em.warnings) == 0 {
		return "Hata yok"
	}

	var output strings.Builder

	// Hataları yazdır
	if len(em.errors) > 0 {
		output.WriteString("=== HATALAR ===\n")
		for _, err := range em.errors {
			output.WriteString(em.formatError(err))
		}
		output.WriteString("\n")
	}

	// Uyarıları yazdır
	if len(em.warnings) > 0 {
		output.WriteString("=== UYARILAR ===\n")
		for _, warn := range em.warnings {
			output.WriteString(em.formatError(warn))
		}
		output.WriteString("\n")
	}

	// Özet
	output.WriteString(fmt.Sprintf("Toplam: %d hata, %d uyarı\n", len(em.errors), len(em.warnings)))

	return output.String()
}

// formatError - Hatayı format et
func (em *ErrorManager) formatError(err *CompileError) string {
	var output strings.Builder

	// Header
	output.WriteString(fmt.Sprintf("[%s] %s\n", err.Severity.String(), err.ID))

	// Location
	output.WriteString(fmt.Sprintf("  Dosya: %s:%d:%d\n", err.Location.File, err.Location.Line, err.Location.Column))

	// Message
	output.WriteString(fmt.Sprintf("  Mesaj: %s\n", err.Message))

	// Context
	if err.Context != "" {
		output.WriteString(fmt.Sprintf("  Bağlam: %s\n", err.Context))
	}

	// Suggestion
	if err.Fixable && err.Suggestion != "" {
		output.WriteString(fmt.Sprintf("  Çözüm: %s\n", err.Suggestion))
	}

	output.WriteString("\n")
	return output.String()
}

// generateErrorID - Benzersiz hata ID'si oluştur
func (em *ErrorManager) generateErrorID() string {
	timestamp := time.Now().Unix()
	count := len(em.errors) + len(em.warnings)
	return fmt.Sprintf("ERR-%d-%d", timestamp%10000, count)
}

// ErrorRecovery - Hata kurtarma sistemi
type ErrorRecovery struct {
	strategies map[string]RecoveryStrategy
	mu         sync.RWMutex
}

// RecoveryStrategy - Kurtarma stratejisi
type RecoveryStrategy struct {
	Name       string
	Condition  func(*CompileError) bool
	Action     func(*CompileError) string
	Automatic  bool
}

// NewErrorRecovery - Yeni kurtarma sistemi oluştur
func NewErrorRecovery() *ErrorRecovery {
	er := &ErrorRecovery{
		strategies: make(map[string]RecoveryStrategy),
	}
	er.registerDefaultStrategies()
	return er
}

// registerDefaultStrategies - Varsayılan kurtarma stratejilerini kaydet
func (er *ErrorRecovery) registerDefaultStrategies() {
	// Missing semicolon recovery
	er.RegisterStrategy("missing_semicolon", RecoveryStrategy{
		Name: "Missing Semicolon",
		Condition: func(err *CompileError) bool {
			return strings.Contains(err.Message, "semicolon")
		},
		Action: func(err *CompileError) string {
			return "Satırın sonuna ';' ekleyin"
		},
		Automatic: true,
	})

	// Undefined variable recovery
	er.RegisterStrategy("undefined_var", RecoveryStrategy{
		Name: "Undefined Variable",
		Condition: func(err *CompileError) bool {
			return strings.Contains(err.Message, "undefined")
		},
		Action: func(err *CompileError) string {
			return "Değişkeni tanımlayın veya bağlantıyı kontrol edin"
		},
		Automatic: false,
	})

	// Type mismatch recovery
	er.RegisterStrategy("type_mismatch", RecoveryStrategy{
		Name: "Type Mismatch",
		Condition: func(err *CompileError) bool {
			return strings.Contains(err.Message, "type")
		},
		Action: func(err *CompileError) string {
			return "Tipleri uyumlu hale getirin veya cast işlemi yapın"
		},
		Automatic: false,
	})
}

// RegisterStrategy - Strateji kaydet
func (er *ErrorRecovery) RegisterStrategy(key string, strategy RecoveryStrategy) {
	er.mu.Lock()
	defer er.mu.Unlock()
	er.strategies[key] = strategy
}

// RecoverError - Hatayı kurtarma stratejisi ile düzeltmeyi dene
func (er *ErrorRecovery) RecoverError(err *CompileError) (bool, string) {
	er.mu.RLock()
	defer er.mu.RUnlock()

	for _, strategy := range er.strategies {
		if strategy.Condition(err) {
			recovery := strategy.Action(err)
			return strategy.Automatic, recovery
		}
	}

	return false, "Otomatik kurtarma stratejisi bulunamadı"
}

// DiagnosticReport - Tanı raporu
type DiagnosticReport struct {
	TotalErrors   int
	TotalWarnings int
	CriticalCount int
	FixableCount  int
	Summary       string
	CreatedAt     time.Time
}

// GenerateDiagnosticReport - Tanı raporu oluştur
func (em *ErrorManager) GenerateDiagnosticReport() *DiagnosticReport {
	em.mu.RLock()
	defer em.mu.RUnlock()

	critical := 0
	fixable := 0

	for _, err := range em.errors {
		if err.Severity == SeverityCritical {
			critical++
		}
		if err.Fixable {
			fixable++
		}
	}

	health := "İYİ"
	if len(em.errors) > 0 {
		health = "UYARI"
	}
	if critical > 0 {
		health = "KRİTİK"
	}

	summary := fmt.Sprintf("Sistem Sağlığı: %s | Hatalar: %d | Uyarılar: %d | Düzeltilebilir: %d",
		health, len(em.errors), len(em.warnings), fixable)

	return &DiagnosticReport{
		TotalErrors:   len(em.errors),
		TotalWarnings: len(em.warnings),
		CriticalCount: critical,
		FixableCount:  fixable,
		Summary:       summary,
		CreatedAt:     time.Now(),
	}
}
