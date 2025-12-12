// Package assets - Derleme sistemine ait kaynaklar ve varlıkları yönet
// Geliştirilmiş assets yönetim sistemi derleyici için gerekli dosyaları işler
package assets

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// AssetMetadata - Varlık meta bilgilerini tutar
type AssetMetadata struct {
	Path      string
	Size      int64
	Hash      string
	Modified  time.Time
	IsDir     bool
	Cached    bool
	CacheTime time.Time
}

// AssetManager - Varlık yönetim sistemi
type AssetManager struct {
	basePath  string
	cache     map[string]*AssetMetadata
	cacheLock sync.RWMutex
	cacheDir  string
}

// NewAssetManager - Yeni varlık yöneticisi oluştur
func NewAssetManager(basePath string) *AssetManager {
	return &AssetManager{
		basePath: basePath,
		cache:    make(map[string]*AssetMetadata),
		cacheDir: filepath.Join(basePath, ".asset_cache"),
	}
}

// ProcessAssets - Range döngüsü ile varlıkları işle
func (am *AssetManager) ProcessAssets(assetPaths []string) error {
	if len(assetPaths) == 0 {
		return fmt.Errorf("varlık yolları boş")
	}

	// Foreach döngüsü - Go'da range kullanılır
	for idx, path := range assetPaths {
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("dosya kontrolü başarısız [%d/%d]: %w", idx+1, len(assetPaths), err)
		}

		// Dizin kontrolü
		if info.IsDir() {
			if err := am.processDirectory(path); err != nil {
				return err
			}
			continue
		}

		// Dosya boyutu kontrolü
		if info.Size() == 0 {
			return fmt.Errorf("boş dosya bulundu: %s", path)
		}

		// Meta veri oluştur
		metadata := &AssetMetadata{
			Path:     path,
			Size:     info.Size(),
			Modified: info.ModTime(),
			IsDir:    false,
		}

		// Hash hesapla
		if err := am.calculateHash(metadata); err != nil {
			return fmt.Errorf("hash hesaplama hatası: %w", err)
		}

		am.cacheLock.Lock()
		am.cache[path] = metadata
		am.cacheLock.Unlock()
	}

	return nil
}

// processDirectory - Dizin içeriğini işle (While döngüsü örneği)
func (am *AssetManager) processDirectory(dirPath string) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	// While döngüsü - Go'da for kullanılır (index kullanarak)
	i := 0
	for i < len(entries) {
		entry := entries[i]
		fullPath := filepath.Join(dirPath, entry.Name())

		if entry.IsDir() {
			// Dizini işle
			if err := am.processDirectory(fullPath); err != nil {
				return err
			}
		} else {
			// Dosyayı işle
			info, _ := entry.Info()
			metadata := &AssetMetadata{
				Path:     fullPath,
				Size:     info.Size(),
				Modified: info.ModTime(),
				IsDir:    false,
			}

			if err := am.calculateHash(metadata); err != nil {
				return err
			}

			am.cacheLock.Lock()
			am.cache[fullPath] = metadata
			am.cacheLock.Unlock()
		}
		i++
	}

	return nil
}

// calculateHash - Dosya hash'ini hesapla
func (am *AssetManager) calculateHash(metadata *AssetMetadata) error {
	file, err := os.Open(metadata.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}

	metadata.Hash = fmt.Sprintf("%x", hash.Sum(nil))
	return nil
}

// GetAssetMetadata - Varlık meta bilgisini al
func (am *AssetManager) GetAssetMetadata(path string) *AssetMetadata {
	am.cacheLock.RLock()
	defer am.cacheLock.RUnlock()
	return am.cache[path]
}

// CacheStats - Cache istatistiklerini döndür
func (am *AssetManager) CacheStats() map[string]interface{} {
	am.cacheLock.RLock()
	defer am.cacheLock.RUnlock()

	totalSize := int64(0)
	for _, metadata := range am.cache {
		totalSize += metadata.Size
	}

	return map[string]interface{}{
		"cached_files": len(am.cache),
		"total_size":   totalSize,
		"cache_path":   am.cacheDir,
	}
}

// ValidateAssets - Varlıkları doğrula (async işlem)
func (am *AssetManager) ValidateAssets() error {
	am.cacheLock.RLock()
	paths := make([]string, 0, len(am.cache))
	for path := range am.cache {
		paths = append(paths, path)
	}
	am.cacheLock.RUnlock()

	// Parallel doğrulama
	errChan := make(chan error, len(paths))
	var wg sync.WaitGroup

	for _, path := range paths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			info, err := os.Stat(p)
			if err != nil {
				errChan <- fmt.Errorf("varlık doğrulama hatası (%s): %w", p, err)
				return
			}

			// Meta veri güncelle
			am.cacheLock.Lock()
			metadata := am.cache[p]
			if metadata != nil {
				metadata.Modified = info.ModTime()
				metadata.Size = info.Size()
				metadata.CacheTime = time.Now()
				metadata.Cached = true
			}
			am.cacheLock.Unlock()
		}(path)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// ClearCache - Cache temizle
func (am *AssetManager) ClearCache() {
	am.cacheLock.Lock()
	defer am.cacheLock.Unlock()
	am.cache = make(map[string]*AssetMetadata)
}