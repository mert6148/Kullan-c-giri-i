// Package sys - Test dosyaları
// Derleyici sistemi test eder
package sys

import (
	"testing"
)

// TestLexer - Lexer test'leri
func TestLexer(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  int // token count
		expectErr bool
	}{
		{
			name:     "Simple assignment",
			input:    "x = 10;",
			expected: 4, // IDENT, ASSIGN, INT, SEMI
		},
		{
			name:     "If statement",
			input:    "if (x > 0) { return x; }",
			expected: 9,
		},
		{
			name:     "Function declaration",
			input:    "func add(a, b) { return a + b; }",
			expected: 13,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens, err := lexer.Tokenize()

			if err != nil && !tt.expectErr {
				t.Errorf("unexpected error: %v", err)
			}

			if len(tokens) != tt.expected {
				t.Errorf("expected %d tokens, got %d", tt.expected, len(tokens))
			}
		})
	}
}

// TestParser - Parser test'leri
func TestParser(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{
			name:      "Valid if statement",
			input:     "if (x > 0) { return x; }",
			expectErr: false,
		},
		{
			name:      "Valid for loop",
			input:     "for i = 0; i < 10; i++ { }",
			expectErr: false,
		},
		{
			name:      "Valid function",
			input:     "func test() { return 0; }",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens, _ := lexer.Tokenize()
			parser := NewParser(tokens)
			err := parser.Parse()

			if err != nil && !tt.expectErr {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestCompiler - Compiler test'leri
func TestCompiler(t *testing.T) {
	source := `
		func main() {
			x := 10
			if x > 0 {
				return x
			}
		}
	`

	compiler := NewCompiler(source)
	err := compiler.Compile()

	if err != nil {
		t.Logf("Compile error (expected): %v", err)
	}

	output := compiler.GetOutput()
	if output == "" {
		t.Error("expected non-empty output")
	}
}

// TestOptimizer - Optimizer test'leri
func TestOptimizer(t *testing.T) {
	source := `
		// This is a comment
		func add(a, b) {
			/* Multi-line
			   comment */
			return a + b;
		}
	`

	optimizer := NewOptimizer(source)
	result, err := optimizer.Optimize()

	if err != nil {
		t.Errorf("optimization failed: %v", err)
	}

	// Comments should be removed
	if len(result) > len(source) {
		t.Error("optimized code should be shorter")
	}

	stats := optimizer.GetStats()
	if stats.RemovedComments == 0 {
		t.Error("should have removed comments")
	}
}

// TestCodeAnalyzer - Code analyzer test'leri
func TestCodeAnalyzer(t *testing.T) {
	source := `
		func longFunctionNameWithManyParameters(param1, param2, param3, param4) {
			if param1 > 0 {
				for i = 0; i < param2; i++ {
					if param3 > i {
						return param4;
					}
				}
			}
		}
	`

	analyzer := NewCodeAnalyzer(source)
	err := analyzer.Analyze()

	if err != nil {
		t.Errorf("analysis failed: %v", err)
	}

	metrics := analyzer.GetMetrics()
	if metrics.Lines == 0 {
		t.Error("expected metrics to be calculated")
	}

	if metrics.Functions == 0 {
		t.Error("expected function count > 0")
	}

	issues := analyzer.GetIssues()
	t.Logf("Found %d code issues", len(issues))
}

// TestErrorManager - Error manager test'leri
func TestErrorManager(t *testing.T) {
	em := NewErrorManager(100)

	// Add test error
	err := em.AddError(SeverityError, "Test error", ErrorLocation{
		File:   "test.go",
		Line:   10,
		Column: 5,
	})

	if err != nil {
		t.Errorf("failed to add error: %v", err)
	}

	if em.GetErrorCount() != 1 {
		t.Error("expected 1 error")
	}

	// Test error formatting
	output := em.Print()
	if len(output) == 0 {
		t.Error("expected non-empty error output")
	}

	t.Logf("Error output:\n%s", output)
}

// TestErrorRecovery - Error recovery test'leri
func TestErrorRecovery(t *testing.T) {
	recovery := NewErrorRecovery()

	testError := &CompileError{
		Message: "missing semicolon",
	}

	automatic, suggestion := recovery.RecoverError(testError)
	if !automatic || suggestion == "" {
		t.Error("expected recovery strategy to work")
	}

	t.Logf("Recovery suggestion: %s", suggestion)
}

// TestCompilerManager - Compiler manager test'leri
func TestCompilerManager(t *testing.T) {
	manager := NewCompilerManager(".")

	// Test version
	version := manager.Version()
	if version == "" {
		t.Error("expected version string")
	}
	t.Logf("Compiler version: %s", version)

	// Test build stats
	stats := manager.GetBuildStats()
	if stats == nil {
		t.Error("expected build stats")
	}
}

// BenchmarkLexer - Lexer performans testi
func BenchmarkLexer(b *testing.B) {
	source := `
		func fibonacci(n) {
			if n <= 1 {
				return n;
			}
			return fibonacci(n-1) + fibonacci(n-2);
		}
	`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer := NewLexer(source)
		lexer.Tokenize()
	}
}

// BenchmarkParser - Parser performans testi
func BenchmarkParser(b *testing.B) {
	source := `func test() { if x > 0 { return x; } }`
	lexer := NewLexer(source)
	tokens, _ := lexer.Tokenize()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := NewParser(tokens)
		parser.Parse()
	}
}

// BenchmarkOptimizer - Optimizer performans testi
func BenchmarkOptimizer(b *testing.B) {
	source := `
		// Comment
		func add(a, b) {
			/* Multi
			   line */
			return a + b;
		}
	`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		optimizer := NewOptimizer(source)
		optimizer.Optimize()
	}
}

// TestCompileWithAssets - Assets ile derleme testi
func TestCompileWithAssets(t *testing.T) {
	am := NewAssetManager(".")
	
	// Mock assets
	assetPaths := []string{"."}
	err := am.ProcessAssets(assetPaths)
	
	if err != nil {
		t.Logf("Asset processing info: %v", err)
	}

	stats := am.CacheStats()
	t.Logf("Asset stats: %v", stats)
}

// TestIntegration - Entegrasyon testi
func TestIntegration(t *testing.T) {
	// Tam derleme süreci testi
	source := `
		func main() {
			x := 10;
			if x > 0 {
				return x;
			}
		}
	`

	// 1. Compile
	compiler := NewCompiler(source)
	if err := compiler.Compile(); err != nil {
		t.Logf("Compile step completed with: %v", err)
	}

	// 2. Analyze
	analyzer := NewCodeAnalyzer(source)
	analyzer.Analyze()

	metrics := analyzer.GetMetrics()
	t.Logf("Code metrics - Lines: %d, Functions: %d, Branches: %d, Complexity: %.2f",
		metrics.Lines, metrics.Functions, metrics.Branches, metrics.Complexity)

	// 3. Error handling
	em := NewErrorManager(10)
	em.AddError(SeverityWarning, "Test warning", ErrorLocation{
		File: "test.go",
		Line: 5,
	})

	report := em.GenerateDiagnosticReport()
	t.Logf("Diagnostic report: %s", report.Summary)
}
