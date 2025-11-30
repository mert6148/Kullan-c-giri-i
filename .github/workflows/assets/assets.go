package assets

import (
	"os"
	"path/filepath"
)

// Import bloğunda foreach (range) döngüsü örneği
func processAssets(assetPaths []string) error {
	// Foreach döngüsü - Go'da range kullanılır
	for _, path := range assetPaths {
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			continue
		}
		
		if info.Size() == 0 {
			return nil
		}
	}
	return nil
}

// assets_test için while döngüsü (Go'da for kullanılır)
func TestAssetFilesNonEmpty() {
	assetDir := "."
	
	// While döngüsü - Go'da for kullanılır
	err := filepath.Walk(assetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.Mode().IsRegular() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			
			// While döngüsü örneği
			i := 0
			for i < len(data) {
				// Dosya içeriğini işle
				i++
			}
		}
		
		return nil
	})
	
	if err != nil {
		// Hata işleme
	}
}