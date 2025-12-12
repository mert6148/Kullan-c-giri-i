# Python3 Şema Validatörü İş Akışı

## Genel Bakış

Bu dokümant, `AssetSchemaValidator` sınıfını ve if döngüleri ile varlık doğrulama iş akışını açıklar.

---

## 1. Temel Fonksiyonlar

### 1.1 `validate_schema_with_conditions()`

**Amaç**: Tekil bir varlığı şema doğrulamasından geçirmek

**Kullanım**:
```python
# İmport
from assest import UserAssetManager, validate_schema_with_conditions
from assest import ASSET_CATEGORY_PROFILE, ASSET_TYPE_STRING

# UserAssetManager oluştur
manager = UserAssetManager()

# Varlığı doğrula
is_valid, error_msg = validate_schema_with_conditions(
    asset_manager=manager,
    category=ASSET_CATEGORY_PROFILE,
    asset_name="email",
    asset_value="user@example.com",
    asset_type=ASSET_TYPE_STRING
)

if is_valid:
    print("✓ Varlık geçerli")
else:
    print(f"✗ Hata: {error_msg}")
```

**If döngüsü örneği**:
```python
# Birden fazla varlığı kontrol et
assets_to_check = [
    ("email", "test@example.com"),
    ("phone", "05301234567"),
    ("bio", "Biyografim")
]

for asset_name, asset_value in assets_to_check:
    is_valid, error = validate_schema_with_conditions(
        manager,
        ASSET_CATEGORY_PROFILE,
        asset_name,
        asset_value
    )
    
    if is_valid:
        print(f"✓ {asset_name} doğru")
    else:
        print(f"✗ {asset_name} hatalı: {error}")
```

---

### 1.2 `validate_user_assets_batch()`

**Amaç**: Toplu varlık doğrulaması (tüm kategorilerde)

**Kullanım**:
```python
from assest import validate_user_assets_batch
from assest import ASSET_CATEGORY_PROFILE, ASSET_CATEGORY_PREFERENCES

# Doğrulanacak varlıklar
assets_to_validate = {
    ASSET_CATEGORY_PROFILE: {
        "first_name": "Ahmet",
        "email": "ahmet@example.com",
    },
    ASSET_CATEGORY_PREFERENCES: {
        "theme": "dark",
        "language": "tr_TR",
    }
}

# Toplu doğrulama
all_valid, errors = validate_user_assets_batch(
    asset_manager=manager,
    assets_dict=assets_to_validate
)

if all_valid:
    print("✓ Tüm varlıklar geçerli")
else:
    print(f"✗ {len(errors)} hata:")
    for error in errors:
        print(f"  - {error}")
```

**If döngüsü iç yapısı**:
```python
# Kategoriler üzerinde döngü
for category in VALID_ASSET_CATEGORIES:
    if category not in assets_to_validate:
        continue  # Kategori yoksa atla
    
    category_assets = assets_to_validate[category]
    
    # Her kategorideki varlıklar üzerinde döngü
    if isinstance(category_assets, dict):
        for asset_name, asset_value in category_assets.items():
            # Doğrulama yapılır
```

---

### 1.3 `process_validated_assets()`

**Amaç**: Doğrulanmış varlıkları veritabanına kaydet

**Kullanım**:
```python
from assest import process_validated_assets

assets_dict = {
    ASSET_CATEGORY_PROFILE: {
        "first_name": "Ahmet",
        "last_name": "Yılmaz",
    },
    ASSET_CATEGORY_SECURITY: {
        "two_factor_enabled": "true",
    }
}

user_id = 1

success, failed_assets = process_validated_assets(
    asset_manager=manager,
    user_id=user_id,
    assets_dict=assets_dict
)

if success:
    print(f"✓ Varlıklar kaydedildi (User ID: {user_id})")
else:
    print(f"✗ Kaydedilemeyen varlıklar:")
    for failed in failed_assets:
        print(f"  - {failed}")
```

---

## 2. If Döngüsü Kullanım Örnekleri

### 2.1 Kategoriler Üzerinde Döngü

```python
from assest import VALID_ASSET_CATEGORIES

assets = manager.get_all_assets(user_id=1)

# Her kategori için döngü
for category in VALID_ASSET_CATEGORIES:
    if category in assets and assets[category]:
        print(f"\n--- {category.upper()} ---")
        
        # Her varlık için döngü
        for asset_name, asset in assets[category].items():
            print(f"  {asset_name}: {asset.asset_value}")
```

### 2.2 Hata Kontrolü

```python
errors = []

# Birden fazla varlık doğrula
for asset_data in asset_list:
    is_valid, error = validate_schema_with_conditions(
        manager,
        asset_data["category"],
        asset_data["name"],
        asset_data["value"]
    )
    
    if not is_valid:
        errors.append(f"{asset_data['category']}.{asset_data['name']}: {error}")

if errors:
    print("Hata bulundu:")
    for err in errors:
        print(f"  {err}")
else:
    print("Tüm varlıklar geçerli")
```

### 2.3 Koşullu Işleme

```python
# Geçersiz varlıkları ayrı işle
invalid_assets = []
valid_assets = []

for category, assets in assets_dict.items():
    for asset_name, asset_value in assets.items():
        is_valid, error = validate_schema_with_conditions(
            manager, category, asset_name, asset_value
        )
        
        if is_valid:
            valid_assets.append((category, asset_name, asset_value))
        else:
            invalid_assets.append((category, asset_name, error))

# Geçerlileri kaydet
if valid_assets:
    for category, name, value in valid_assets:
        manager.set_asset(1, name, value, category=category)

# Geçersizleri raporla
if invalid_assets:
    print("Geçersiz varlıklar:")
    for category, name, error in invalid_assets:
        print(f"  {category}.{name}: {error}")
```

---

## 3. Tam Örnek: Kullanıcı Kaydı Akışı

```python
def register_new_user_with_validation(manager, user_data):
    """Yeni kullanıcı kaydı - varlık doğrulaması ile"""
    
    user_id = user_data["id"]
    assets_dict = {
        ASSET_CATEGORY_PROFILE: {
            "first_name": user_data.get("first_name"),
            "last_name": user_data.get("last_name"),
            "email": user_data.get("email"),
        },
        ASSET_CATEGORY_SECURITY: {
            "two_factor_enabled": "false",
            "login_attempts": "0",
        }
    }
    
    # Adım 1: Toplu doğrulama
    print(f"[1/3] Varlıklar doğrulanıyor...")
    all_valid, errors = validate_user_assets_batch(manager, assets_dict)
    
    if not all_valid:
        print(f"✗ Doğrulama başarısız:")
        for err in errors:
            print(f"  - {err}")
        return False
    
    print("✓ Varlıklar doğru")
    
    # Adım 2: Varlıkları kaydet
    print(f"[2/3] Varlıklar kaydediliyor...")
    success, failed = process_validated_assets(manager, user_id, assets_dict)
    
    if not success:
        print(f"✗ Kayıt başarısız:")
        for fail in failed:
            print(f"  - {fail}")
        return False
    
    print("✓ Varlıklar kaydedildi")
    
    # Adım 3: Doğrulama
    print(f"[3/3] Varlıklar doğrulanıyor...")
    saved_assets = manager.get_all_assets(user_id)
    
    verification_ok = True
    for category in ASSET_CATEGORY_PROFILE, ASSET_CATEGORY_SECURITY:
        if category not in saved_assets or not saved_assets[category]:
            print(f"✗ {category} varlıkları bulunamadı")
            verification_ok = False
    
    if verification_ok:
        print("✓ Varlıklar doğrulandı")
        return True
    else:
        return False
```

---

## 4. CLI Komutları

### Doğrulama Sonuçlarını Göster

```bash
# Tek bir varlığı doğrula
python -c "
from assests.assest import *
manager = UserAssetManager()
is_valid, err = validate_schema_with_conditions(
    manager, 'profile', 'email', 'test@example.com'
)
print('Valid' if is_valid else f'Error: {err}')
"
```

### Toplu Doğrulama

```bash
# Toplu doğrulama
python -c "
from assests.assest import *

manager = UserAssetManager()
assets = {
    'profile': {'email': 'test@example.com', 'phone': '05301234567'},
    'preferences': {'theme': 'dark'}
}

valid, errors = validate_user_assets_batch(manager, assets)
print('✓ Valid' if valid else f'✗ Errors: {errors}')
"
```

---

## 5. Hata Türleri

| Hata | Açıklama | Çözüm |
|------|----------|-------|
| Geçersiz kategori | Kategoriye `VALID_ASSET_CATEGORIES`'de tanımlanmıştı | Kategoriyi kontrol et |
| Null değer | Varlık değeri boş | Değer sağla |
| Geçersiz tip | Varlık tipi eşleşmedi | Tipi düzelt |
| Max length aşımı | String çok uzun | Stringi kısalt |
| Pattern eşleşmedi | Regex pattern'e uymadı | Format'ı düzelt |

---

## 6. Özet: If Döngüsü Yapısı

```python
# Kategoriler üzerinde döngü
for category in VALID_ASSET_CATEGORIES:
    if category not in data:
        continue  # Kategori yoksa atla
    
    # Varlıklar üzerinde döngü
    for asset_name, asset_value in data[category].items():
        # Doğrulama
        is_valid, error = validate_schema_with_conditions(
            manager, category, asset_name, asset_value
        )
        
        # Sonuca göre işlem yap
        if is_valid:
            # Başarılı işlem
            pass
        else:
            # Hata işleme
            errors.append(error)
```

---

**Son Güncelleme**: 10 Aralık 2025
**Python Versiyonu**: 3.8+
