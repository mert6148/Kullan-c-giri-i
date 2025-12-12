// Package sys - Geliştirilmiş GO derleyici ve test sistemi
// Derleme, çapraz derleme, test ve optimizasyon işlevlerini sağlar
package sys

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// BuildConfig - Derleme yapılandırması
type BuildConfig struct {
	GOOS       string
	GOARCH     string
	CGOEnabled bool
	Flags      []string
	Tags       []string
	Optimize   bool
}

// CompileResult - Derleme sonuçları
type CompileResult struct {
	Success      bool
	Duration     time.Duration
	Output       string
	Errors       string
	WarningCount int
	ErrorCount   int
}

// CompilerManager - Derleyici yönetim sistemi
type CompilerManager struct {
	basePath   string
	buildCache map[string]*CompileResult
	cacheLock  sync.RWMutex
	verbose    bool
	buildLock  sync.Mutex
}

// NewCompilerManager - Yeni derleyici yöneticisi oluştur
func NewCompilerManager(basePath string) *CompilerManager {
	return &CompilerManager{
		basePath:   basePath,
		buildCache: make(map[string]*CompileResult),
		verbose:    os.Getenv("VERBOSE_BUILD") == "1",
	}
}

// CompileSys - Tüm modülü derle
func (cm *CompilerManager) CompileSys() error {
	cm.buildLock.Lock()
	defer cm.buildLock.Unlock()

	start := time.Now()
	cmd := exec.Command("go", "build", "./...")

	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()
	duration := time.Since(start)

	result := &CompileResult{
		Success:  err == nil,
		Duration: duration,
		Output:   out.String(),
		Errors:   errOut.String(),
	}

	if cm.verbose {
		fmt.Printf("[BUILD] Derleme tamamlandı (%.2fs)\n", duration.Seconds())
	}

	cm.cacheLock.Lock()
	cm.buildCache["main"] = result
	cm.cacheLock.Unlock()

	if err != nil {
		return fmt.Errorf("sys derleme başarısız: %v\nHata: %s", err, errOut.String())
	}
	return nil
}

// CompilePath - Belirtilen yolu derle
func (cm *CompilerManager) CompilePath(path string) (*CompileResult, error) {
	cm.buildLock.Lock()
	defer cm.buildLock.Unlock()

	start := time.Now()
	cmd := exec.Command("go", "build", path)

	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()
	duration := time.Since(start)

	// Hata sayısını say
	errorCount := strings.Count(errOut.String(), "error")
	warningCount := strings.Count(errOut.String(), "warning")

	result := &CompileResult{
		Success:      err == nil,
		Duration:     duration,
		Output:       out.String(),
		Errors:       errOut.String(),
		ErrorCount:   errorCount,
		WarningCount: warningCount,
	}

	cm.cacheLock.Lock()
	cm.buildCache[path] = result
	cm.cacheLock.Unlock()

	if err != nil {
		return result, fmt.Errorf("%s derleme başarısız: %v", path, err)
	}
	return result, nil
}

// Version - Derleme sürümünü döndür
func (cm *CompilerManager) Version() string {
	// Git commit hash'ini al (varsa)
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()

	commit := strings.TrimSpace(out.String())
	if commit == "" {
		commit = "unknown"
	}

	return fmt.Sprintf("sys-compiler-v1.0.0 (commit: %s)", commit)
}

// CrossCompile - Çapraz derleme (multiple platforms)
func (cm *CompilerManager) CrossCompile(platforms []string) map[string]*CompileResult {
	results := make(map[string]*CompileResult)
	var wg sync.WaitGroup

	for _, platform := range platforms {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			result := cm.crossCompilePlatform(p)
			results[p] = result
		}(platform)
	}

	wg.Wait()
	return results
}

// crossCompilePlatform - Tek platform için çapraz derleme
func (cm *CompilerManager) crossCompilePlatform(platform string) *CompileResult {
	parts := strings.Split(platform, "/")
	if len(parts) != 2 {
		return &CompileResult{
			Success: false,
			Errors:  fmt.Sprintf("Geçersiz platform formatı: %s (beklenen: os/arch)", platform),
		}
	}

	goos, goarch := parts[0], parts[1]
	start := time.Now()

	cmd := exec.Command("go", "build", "-o", fmt.Sprintf("build-%s-%s", goos, goarch), "./...")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%s", goos),
		fmt.Sprintf("GOARCH=%s", goarch),
		"CGO_ENABLED=0",
	)

	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()
	duration := time.Since(start)

	result := &CompileResult{
		Success:    err == nil,
		Duration:   duration,
		Output:     out.String(),
		Errors:     errOut.String(),
		ErrorCount: strings.Count(errOut.String(), "error"),
	}

	if cm.verbose {
		fmt.Printf("[CROSS-COMPILE] %s → %.2fs\n", platform, duration.Seconds())
	}

	return result
}

// CompileWithOptimization - Optimizasyon ile derleme
func (cm *CompilerManager) CompileWithOptimization(optimize string) (*CompileResult, error) {
	cm.buildLock.Lock()
	defer cm.buildLock.Unlock()

	start := time.Now()

	flags := []string{"build"}

	// Optimizasyon türüne göre flag ekle
	switch strings.ToLower(optimize) {
	case "size":
		flags = append(flags, "-ldflags", "-s -w")
	case "speed":
		flags = append(flags, "-gcflags", "-O=3")
	case "full":
		flags = append(flags, "-ldflags", "-s -w", "-gcflags", "-O=3")
	}

	cmd := exec.Command("go", flags...)
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()
	duration := time.Since(start)

	result := &CompileResult{
		Success:    err == nil,
		Duration:   duration,
		Output:     out.String(),
		Errors:     errOut.String(),
		ErrorCount: strings.Count(errOut.String(), "error"),
	}

	if cm.verbose {
		fmt.Printf("[OPTIMIZE] %s → %.2fs\n", optimize, duration.Seconds())
	}

	if err != nil {
		return result, fmt.Errorf("optimize derleme başarısız: %v", err)
	}
	return result, nil
}

// CleanBuildCache - Build cache temizle
func (cm *CompilerManager) CleanBuildCache() error {
	cmd := exec.Command("go", "clean", "-cache")
	var out bytes.Buffer
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build cache temizleme başarısız: %v\nDetails: %s", err, out.String())
	}

	// İç cache'i de temizle
	cm.cacheLock.Lock()
	cm.buildCache = make(map[string]*CompileResult)
	cm.cacheLock.Unlock()

	if cm.verbose {
		fmt.Println("[CACHE] Build cache temizlendi")
	}
	return nil
}

// GetBuildStats - Derleme istatistiklerini al
func (cm *CompilerManager) GetBuildStats() map[string]interface{} {
	cm.cacheLock.RLock()
	defer cm.cacheLock.RUnlock()

	totalTime := time.Duration(0)
	successCount := 0
	failCount := 0

	for _, result := range cm.buildCache {
		totalTime += result.Duration
		if result.Success {
			successCount++
		} else {
			failCount++
		}
	}

	return map[string]interface{}{
		"total_builds":   len(cm.buildCache),
		"success":        successCount,
		"failed":         failCount,
		"total_duration": totalTime.String(),
		"avg_duration":   (totalTime / time.Duration(len(cm.buildCache))).String(),
	}
}

// TestModule - Test çalıştır
func (cm *CompilerManager) TestModule(pattern string) (*CompileResult, error) {
	start := time.Now()

	cmd := exec.Command("go", "test", "-v", "-race", pattern)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	duration := time.Since(start)

	result := &CompileResult{
		Success:  err == nil,
		Duration: duration,
		Output:   out.String(),
	}

	if cm.verbose {
		fmt.Printf("[TEST] %s → %.2fs\n", pattern, duration.Seconds())
	}

	return result, err
}

// GetBuildCache - Cache verilerini al
func (cm *CompilerManager) GetBuildCache(key string) *CompileResult {
	cm.cacheLock.RLock()
	defer cm.cacheLock.RUnlock()
	return cm.buildCache[key]
}
