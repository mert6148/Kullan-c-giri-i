# 🚀 Şema Validatörü - If Döngüsü Referans Kartı

**Python3 İş Akışı** | **Hızlı Referans** | **Kod Örnekleri**

---

## 📌 Hızlı Başlangıç (3 Satır)

```python
from assests.assert import UserAssetManager, validate_schema_with_conditions
manager = UserAssetManager()
is_valid, err = validate_schema_with_conditions(manager, "profile", "email", "test@example.com")
```

---

## 🔧 Temel Fonksiyonlar

### 1️⃣ Tekil Doğrulama
```python
validate_schema_with_conditions(
    asset_manager: UserAssetManager,
    category: str,           # "profile", "preferences", vb.
    asset_name: str,         # "email", "theme", vb.
    asset_value: Any,        # "user@example.com", vb.
    asset_type: str = "string"
) → (bool, Optional[str])
```

**Örnek**:
```python
is_valid, error = validate_schema_with_conditions(
    manager, "profile", "email", "user@example.com"
)
```

---

### 2️⃣ Toplu Doğrulama
```python
validate_user_assets_batch(
    asset_manager: UserAssetManager,
    assets_dict: Dict[str, Dict[str, Any]]
) → (bool, List[str])
```

**Örnek**:
```python
assets = {
    "profile": {"email": "user@example.com"},
    "preferences": {"theme": "dark"}
}
valid, errors = validate_user_assets_batch(manager, assets)
```

---

### 3️⃣ Kaydet
```python
process_validated_assets(
    asset_manager: UserAssetManager,
    user_id: int,
    assets_dict: Dict[str, Dict[str, Any]]
) → (bool, List[str])
```

**Örnek**:
```python
success, failed = process_validated_assets(manager, 1, assets)
```

---

## 🔄 If Döngüsü Şablonları

### Şablon 1: Kategori Döngüsü
```python
from assests.assest import VALID_ASSET_CATEGORIES

for category in VALID_ASSET_CATEGORIES:
    if category not in data:
        continue
    # İşlem...
```

### Şablon 2: İç İçe Döngü
```python
for category in VALID_ASSET_CATEGORIES:
    if category in data:
        for asset_name, value in data[category].items():
            # İşlem...
```

### Şablon 3: Koşullu Filtre
```python
valid_list = []
for category in VALID_ASSET_CATEGORIES:
    for name, value in data.get(category, {}).items():
        is_valid, _ = validate_schema_with_conditions(
            manager, category, name, value
        )
        if is_valid:
            valid_list.append((category, name, value))
```

### Şablon 4: Hata Yakalama
```python
errors = []
for category in VALID_ASSET_CATEGORIES:
    for asset_name, value in data.get(category, {}).items():
        is_valid, error = validate_schema_with_conditions(
            manager, category, asset_name, value
        )
        if not is_valid:
            errors.append(f"{category}.{asset_name}: {error}")

if errors:
    print(f"✗ {len(errors)} hata")
    for err in errors:
        print(f"  {err}")
```

### Şablon 5: Erken Çıkış
```python
found = False
for category in VALID_ASSET_CATEGORIES:
    if found:
        break
    for asset_name, value in data.get(category, {}).items():
        is_valid, _ = validate_schema_with_conditions(
            manager, category, asset_name, value
        )
        if not is_valid:
            print(f"✗ Hata: {asset_name}")
            found = True
            break
```

---

## 🎯 Kategoriler & Varlıklar

### Kategoriler
```python
ASSET_CATEGORY_PROFILE      # "profile"      - Profil Bilgisi
ASSET_CATEGORY_PREFERENCES  # "preferences"  - Tercihler
ASSET_CATEGORY_SECURITY     # "security"     - Güvenlik
ASSET_CATEGORY_SYSTEM       # "system"       - Sistem
ASSET_CATEGORY_CUSTOM       # "custom"       - Özel
```

### Tipler
```python
ASSET_TYPE_STRING    # "string"   - Metin
ASSET_TYPE_INTEGER   # "integer"  - Sayı
ASSET_TYPE_BOOLEAN   # "boolean"  - True/False
ASSET_TYPE_JSON      # "json"     - JSON
ASSET_TYPE_BINARY    # "binary"   - İkili
ASSET_TYPE_FILE      # "file"     - Dosya
```

---

## 💾 Veri Yapısı

```python
# Doğrulanacak varlıklar
assets_dict = {
    "profile": {
        "first_name": "Ahmet",
        "email": "ahmet@example.com",
    },
    "preferences": {
        "theme": "dark",
        "language": "tr_TR"
    }
}

# Toplu doğrulama
valid, errors = validate_user_assets_batch(manager, assets_dict)

# Kaydetme
success, failed = process_validated_assets(manager, user_id, assets_dict)
```

---

## ⚡ Sık Kullanılan Kombinasyonlar

### Başarı Kontrolü
```python
if is_valid and error is None:
    print("✓ Geçerli")
else:
    print(f"✗ {error}")
```

### Toplu İşlem
```python
valid, errors = validate_user_assets_batch(manager, assets)
if valid:
    success, failed = process_validated_assets(manager, user_id, assets)
    if success:
        print("✓ Başarılı")
    else:
        print(f"✗ {failed}")
else:
    print(f"✗ {errors}")
```

### Geri Alma
```python
all_assets = manager.get_all_assets(user_id)
for category in VALID_ASSET_CATEGORIES:
    category_assets = all_assets.get(category, {})
    for asset_name, asset in category_assets.items():
        print(f"{asset_name}: {asset.asset_value}")
```

---

## 🛡️ Hata Türleri

| Hata | Nedenи | Çözüm |
|------|--------|-------|
| Validatör aktif değil | Manager None | Manager oluştur |
| Geçersiz kategori | Yanlış kategori | VALID_ASSET_CATEGORIES kullan |
| Değer boş | None/empty value | Değer sağla |
| Geçersiz tip | Tip eşleşmemiş | Tipi düzelt |
| Pattern hatası | Regex eşleşmemiş | Format'ı düzelt |
| Max length | String çok uzun | Kısalt |
| Max value | Sayı çok büyük | Azalt |

---

## 🧪 Test Komutları

### Tek Varlık Testi
```bash
python3 -c "
from assests.assest import *
mgr = UserAssetManager()
is_valid, err = validate_schema_with_conditions(mgr, 'profile', 'email', 'test@example.com')
print('✓' if is_valid else f'✗ {err}')
"
```

### Toplu Test
```bash
python3 -c "
from assests.assest import *
mgr = UserAssetManager()
assets = {'profile': {'email': 'test@example.com'}}
valid, errs = validate_user_assets_batch(mgr, assets)
print('✓' if valid else f'✗ {errs}')
"
```

### Full Test Suite
```bash
python test_validation_workflow.py
```

---

## 📊 Döngü Karmaşıklığı

```python
# Basit döngü
for cat in categories:           # O(n)
    pass

# İç içe döngü
for cat in categories:           # O(n*m)
    for asset in assets:
        pass

# Koşullu döngü
for cat in categories:           # O(n*m) + doğrulama
    for asset in assets:
        if validate():
            pass
```

---

## 🔍 Debugging İpuçları

### Kategoriyi Kontrol Et
```python
if category not in VALID_ASSET_CATEGORIES:
    print(f"✗ Geçersiz kategori: {category}")
```

### Varlığı Kontrol Et
```python
validator = manager.validator
if category in validator.schema:
    if asset_name in validator.schema[category]:
        print(f"✓ Varlık tanımlı")
```

### Tüm Varlıkları Listele
```python
for cat in VALID_ASSET_CATEGORIES:
    assets = manager.get_assets_by_category(user_id, cat)
    print(f"{cat}: {len(assets)}")
```

---

## 📈 Performance İpuçları

1. **Batch işlem yap**: Tek tek yerine toplu doğrula
2. **If koşullarını minimize et**: Erken döngüden çık
3. **Kategorileri önceden filtreле**: Gereksiz döngüleri atla
4. **Hata sayını sınırla**: İlk N hatayı göster

---

## 🎓 Tam Örnek: Kullanıcı Kaydı

```python
def register(user_id, data):
    manager = UserAssetManager()
    
    assets = {
        "profile": {
            "first_name": data["first_name"],
            "email": data["email"]
        },
        "security": {
            "two_factor_enabled": "false"
        }
    }
    
    # Doğrula
    valid, errors = validate_user_assets_batch(manager, assets)
    if not valid:
        return False, errors
    
    # Kaydet
    success, failed = process_validated_assets(manager, user_id, assets)
    if not success:
        return False, failed
    
    # Kontrol
    saved = manager.get_all_assets(user_id)
    return bool(saved), None
```

---

## 📚 Dosya Referansları

| Dosya | İçerik |
|-------|--------|
| `assests/assest.py` | Kaynak kod |
| `test_validation_workflow.py` | 5 tam test |
| `VALIDATION_WORKFLOW.md` | Detaylı dokümantasyon |
| `PYTHON3_QUICKSTART.md` | Hızlı rehber |
| `VALIDATION_WORKFLOW_SUMMARY.md` | Özet rapor |

---

## ⌨️ Kısayollar

```python
# İmportlar
from assests.assest import (
    UserAssetManager,
    validate_schema_with_conditions,
    validate_user_assets_batch,
    process_validated_assets,
    VALID_ASSET_CATEGORIES,
)

# Manager
m = UserAssetManager()

# Doğrula
validate_schema_with_conditions(m, cat, name, val)

# Toplu Doğrula
validate_user_assets_batch(m, assets)

# Kaydet
process_validated_assets(m, uid, assets)

# Geri Al
m.get_all_assets(uid)
```

---

## ✨ Son İpuçları

✅ **Her zaman doğrula**: Kaydeden önce  
✅ **If döngüsü kullan**: Kategoriler üzerinde  
✅ **Hataları yakalama**: try-except ile  
✅ **Test et**: Her değişiklikten sonra  
✅ **Loglama ekle**: Debugging için  

---

**Sürüm**: 1.0.0  
**Son Güncelleme**: 10 Aralık 2025  
**Durum**: ✅ Hazır  

---

## 🔗 Hızlı Bağlantılar

- `VALIDATION_WORKFLOW.md` → Detaylı Rehber
- `PYTHON3_QUICKSTART.md` → Başlangıç
- `test_validation_workflow.py` → Testler
- `assests/assest.py` → Kod

**Sorular?** Dokümantasyona bakın!
