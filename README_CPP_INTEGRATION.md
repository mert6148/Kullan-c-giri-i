# C/C++ Entegrasyonu - Admin Controller

Bu dokümantasyon, Python3 Admin Controller için C/C++ entegrasyonunu açıklar.

## Genel Bakış

C/C++ entegrasyonu, performans kritik işlemler için Python extension modülü sağlar:
- Sistem istatistikleri (cached)
- Hızlı şifre hashleme
- Hızlı input validasyonu
- Timing-safe string karşılaştırma

## Gereksinimler

### Linux/macOS
```bash
sudo apt-get install python3-dev g++ build-essential  # Ubuntu/Debian
brew install python3 gcc                               # macOS
```

### Windows
- Visual Studio Build Tools veya MinGW
- Python3 development headers

## Derleme

### Yöntem 1: Setuptools (Önerilen)
```bash
python3 setup_cpp_extension.py build_ext --inplace
```

### Yöntem 2: Makefile
```bash
make build-setuptools
# veya
make all
```

### Yöntem 3: Manuel Derleme
```bash
# Linux/macOS
g++ -Wall -O3 -fPIC -std=c++11 -shared -o cpp_admin_extension.so \
    cpp_admin_extension.cpp \
    $(python3-config --includes) \
    $(python3-config --ldflags)

# Windows (MinGW)
g++ -Wall -O3 -std=c++11 -shared -o cpp_admin_extension.pyd \
    cpp_admin_extension.cpp \
    -I"C:/Python3X/include" \
    -L"C:/Python3X/libs" \
    -lpython3X
```

## Kullanım

### Python'da Import
```python
from admin_controller import AdminController, get_admin_controller

# Admin Controller otomatik olarak C/C++ extension'ı kullanır
admin_ctrl = get_admin_controller()

# Login
success, session_id, error = admin_ctrl.login("admin", "password", "admin")

# Sistem istatistikleri (C/C++ ile optimize edilmiş)
success, stats, error = admin_ctrl.get_system_stats(session_id)
```

### Doğrudan C/C++ Wrapper Kullanımı
```python
from cpp_admin_wrapper import CppAdminWrapper

# Sistem istatistikleri
stats = CppAdminWrapper.get_system_stats()

# Hızlı şifre hashleme
hash_value = CppAdminWrapper.fast_hash_password("password", "salt")

# Input validasyonu
is_valid, error_msg = CppAdminWrapper.validate_input_fast("username", 3, 20)

# Timing-safe karşılaştırma
is_match = CppAdminWrapper.fast_string_compare("str1", "str2")
```

## API Endpoints

### Admin Login
```bash
POST /api/v1/admin/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password",
  "role": "admin"
}
```

Response:
```json
{
  "success": true,
  "session_id": "uuid-here",
  "username": "admin",
  "role": "admin"
}
```

### Get System Stats
```bash
GET /api/v1/admin/stats
X-Admin-Session: <session_id>
```

### Get Users
```bash
GET /api/v1/admin/users?limit=100
X-Admin-Session: <session_id>
```

### Get Audit Logs
```bash
GET /api/v1/admin/audit-logs?limit=50
X-Admin-Session: <session_id>
```

## Performans

C/C++ extension kullanıldığında:
- Sistem istatistikleri: ~5ms (cached)
- Şifre hashleme: ~0.1ms (Python: ~1ms)
- Input validasyonu: ~0.05ms (Python: ~0.5ms)
- String karşılaştırma: Timing-safe, constant-time

## Fallback Mekanizması

C/C++ extension mevcut değilse, sistem otomatik olarak pure Python implementasyonuna geçer:
- `cpp_admin_wrapper.py` fallback metodları sağlar
- Performans düşer ama işlevsellik korunur
- Log'larda uyarı mesajı görünür

## Test

```bash
# Extension import testi
python3 -c "import cpp_admin_extension; print('OK')"

# Admin Controller testi
python3 -c "from admin_controller import get_admin_controller; print('OK')"

# API server testi
python3 api_server.py
```

## Sorun Giderme

### Import Error
```
ImportError: No module named 'cpp_admin_extension'
```
**Çözüm**: Extension'ı derleyin: `python3 setup_cpp_extension.py build_ext --inplace`

### Compilation Error
```
error: 'Python.h' file not found
```
**Çözüm**: Python development headers'ı yükleyin:
- Linux: `sudo apt-get install python3-dev`
- macOS: `brew install python3`
- Windows: Python installer'dan "Development Tools" seçeneğini işaretleyin

### Runtime Error
```
Fatal Python error: PyThreadState_Get: no current thread
```
**Çözüm**: Extension'ı Python interpreter içinden import edin, doğrudan çalıştırmayın.

## Dosya Yapısı

```
.
├── admin_controller.py          # Admin Controller (Python)
├── cpp_admin_extension.cpp      # C/C++ extension kaynak kodu
├── cpp_admin_wrapper.py         # Python wrapper (fallback)
├── setup_cpp_extension.py       # Setuptools setup scripti
├── Makefile                     # Makefile (derleme)
├── api_server.py                # Flask API server (entegre)
└── README_CPP_INTEGRATION.md    # Bu dosya
```

## Lisans

Bu proje MIT lisansı altında lisanslanmıştır.

