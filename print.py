# print.py
# Gelişmiş kullanıcı girişi + sistem bilgisi + giriş/çıkış kayıt sistemi

import platform # Şifre girişi için platform modülü eklendi
import getpass  # Şifre girişini daha güvenli yapmak için eklenmiştir
import json
import sys
from pathlib import Path 
from datetime import datetime
import os
import uuid

LOG_FILE = "login_log.txt"
# Event type constants (lint ve tekrar kullanım için)
EVENT_LOGIN = "Giriş"
EVENT_LOGOUT = "Çıkış"
EVENT_FAILED = "Başarısız giriş"

# Basit kullanıcı deposu (örnek). Gerçek uygulamada güvenli bir user store kullanın.
USERS_FILE = "users.json"
USERS = {}
CURRENT_USER = None
CURRENT_SESSION_ID = None
SESSIONS_FILE = "sessions.json"
SESSIONS = []


def _hash_password(password: str) -> tuple:
    """Hash a password with a random salt. Returns (salt_hex, hash_hex)."""
    import hashlib
    import os

    salt = os.urandom(16)
    dk = hashlib.pbkdf2_hmac("sha256", password.encode("utf-8"), salt, 100_000)
    return salt.hex(), dk.hex()


def _verify_password(stored_salt_hex: str, stored_hash_hex: str, password: str) -> bool:
    import hashlib

    salt = bytes.fromhex(stored_salt_hex)
    dk = hashlib.pbkdf2_hmac("sha256", password.encode("utf-8"), salt, 100_000)
    return dk.hex() == stored_hash_hex


def load_user_store():
    """Load users from USERS_FILE or create default admin if not present."""
    global USERS
    if os.path.exists(USERS_FILE):
        try:
            with open(USERS_FILE, "r", encoding="utf-8") as f:
                data = json.load(f)
                # normalize legacy plaintext password entries
                for u, info in data.items():
                    if isinstance(info, dict) and "password" in info:
                        # convert to hash
                        salt, h = _hash_password(info["password"])
                        data[u] = {"salt": salt, "hash": h, "full_name": info.get("full_name")}
                USERS = data
        except Exception:
            USERS = {}
    else:
        # create default admin
        salt, h = _hash_password("1234")
        USERS = {"admin": {"salt": salt, "hash": h, "full_name": "Sistem Yöneticisi"}}
        save_user_store()
    return USERS


def save_user_store():
    try:
        with open(USERS_FILE, "w", encoding="utf-8") as f:
            json.dump(USERS, f, ensure_ascii=False, indent=2)
    except Exception:
        pass


def create_user(username: str, password: str, full_name: str = None) -> bool:
    """Create a new user. Returns True on success, False if user exists."""
    if username in USERS:
        return False
    salt, h = _hash_password(password)
    USERS[username] = {"salt": salt, "hash": h, "full_name": full_name}
    save_user_store()
    return True


def delete_user(username: str) -> bool:
    if username not in USERS:
        return False
    USERS.pop(username)
    save_user_store()
    return True


def list_users():
    return {u: {"full_name": info.get("full_name")} for u, info in USERS.items()}


def add_user_interactive():
    uname = input("Yeni kullanıcı adı: ").strip()
    if not uname:
        print("Kullanıcı adı boş olamaz.")
        return
    if uname in USERS:
        print("Kullanıcı zaten mevcut.")
        return
    pw = getpass.getpass("Şifre: ")
    pw2 = getpass.getpass("Şifre (tekrar): ")
    if pw != pw2:
        print("Şifreler uyuşmuyor.")
        return
    full = input("Ad-Soyad (isteğe bağlı): ").strip() or None
    if create_user(uname, pw, full):
        print("Kullanıcı oluşturuldu.")
    else:
        print("Kullanıcı oluşturulamadı.")


def load_sessions():
    global SESSIONS
    if os.path.exists(SESSIONS_FILE):
        try:
            with open(SESSIONS_FILE, "r", encoding="utf-8") as f:
                SESSIONS = json.load(f)
        except Exception:
            SESSIONS = []
    else:
        SESSIONS = []
    return SESSIONS


def save_sessions():
    try:
        with open(SESSIONS_FILE, "w", encoding="utf-8") as f:
            json.dump(SESSIONS, f, ensure_ascii=False, indent=2)
    except Exception:
        pass


def start_session(username, system=None, code_dirs=None):
    """Create and persist a session record, return session_id."""
    sid = str(uuid.uuid4())
    rec = {
        "id": sid,
        "username": username,
        "login_ts": datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
        "logout_ts": None,
        "system": system,
        "code_dirs": code_dirs,
    }
    SESSIONS.append(rec)
    save_sessions()
    return sid


def end_session(session_id):
    for rec in SESSIONS:
        if rec.get("id") == session_id:
            rec["logout_ts"] = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
            save_sessions()
            return True
    return False


def show_sessions(all_sessions=True):
    if not SESSIONS:
        print("Henüz oturum kaydı yok.")
        return
    for i, s in enumerate(SESSIONS, start=1):
        status = "aktif" if s.get("logout_ts") is None else f"kapandı: {s.get('logout_ts')}"
        print(f"[{i}] {s.get('login_ts')} - {s.get('username')} ({status}) id={s.get('id')}")
        if s.get("system"):
            print("  Sistem:", json.dumps(s.get("system"), ensure_ascii=False))
        if s.get("code_dirs"):
            print("  KodDizinleri:", json.dumps(s.get("code_dirs"), ensure_ascii=False))

def _acquire_file_lock(f):
    """Cross-platform file lock (best-effort).

    Uses msvcrt on Windows and fcntl on POSIX. This is a best-effort lock to
    reduce concurrent write races; failures are ignored to avoid crashing the app.
    """
    try:
        if os.name == "nt":
            import msvcrt

            # lock 1 byte at current position
            msvcrt.locking(f.fileno(), msvcrt.LK_LOCK, 1)
        else:
            import fcntl

            fcntl.flock(f.fileno(), fcntl.LOCK_EX)
    except Exception:
        # best-effort: ignore locking errors
        pass


def _release_file_lock(f):
    try:
        if os.name == "nt":
            import msvcrt

            f.seek(0)
            msvcrt.locking(f.fileno(), msvcrt.LK_UNLCK, 1)
        else:
            import fcntl

            fcntl.flock(f.fileno(), fcntl.LOCK_UN)
    except Exception:
        pass


def log_event(event_type, username="Bilinmiyor", full_name="Bilinmiyor", system=None, code_dirs=None):
    """Kimin ne zaman giriş/çıkış yaptığını kaydeder.

    Opsiyonel olarak `system` ve `code_dirs` bilgileri verilebilir; bunlar log içine
    JSON formatında eklenir.
    """
    def _sanitize(v):
        """Basit sanitize: stringlerde yeni satırları tek satıra indirir, dict/list üzerinde rekürsif çalışır."""
        if isinstance(v, str):
            return v.replace("\n", " ").replace("\r", " ")
        if isinstance(v, dict):
            return {k: _sanitize(val) for k, val in v.items()}
        if isinstance(v, list):
            return [_sanitize(x) for x in v]
        return v

    # JSON-only (JSON-lines) format: her kayıt tek satır JSON olarak kaydedilir.
    time_str = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    event_obj = {
        "timestamp": time_str,
        "event": event_type,
        "username": username,
        "full_name": full_name,
        "system": _sanitize(system) if system is not None else None,
        "code_dirs": _sanitize(code_dirs) if code_dirs is not None else None,
    }

    line = json.dumps(event_obj, ensure_ascii=False) + "\n"
    try:
        # Atomic-ish append with file locking: lock, write, flush, fsync, unlock
        with open(LOG_FILE, "a", encoding="utf-8") as f:
            _acquire_file_lock(f)
            try:
                f.write(line)
                f.flush()
                try:
                    os.fsync(f.fileno())
                except Exception:
                    # fsync may not be available on some platforms; ignore if fails
                    pass
            finally:
                _release_file_lock(f)
    except Exception:
        # Eğer dosyaya yazma başarısız olursa basit bir yedek hatırla
        try:
            with open(LOG_FILE, "a", encoding="utf-8", errors="ignore") as f:
                _acquire_file_lock(f)
                try:
                    f.write(json.dumps({"timestamp": time_str, "event": "error_writing_log"}, ensure_ascii=False) + "\n")
                    f.flush()
                finally:
                    _release_file_lock(f)
        except Exception:
            # tam bir başarısızlık: yazmayı atla (uyarı logu konsola)
            print("Hata: log dosyasına yazılamadı.")


def login(max_attempts: int = 3) -> bool:
    """Etkileşimli giriş.

    - Boş kullanıcı adı girişini engeller.
    - `USERS` içinde tanımlı parolalarla eşleştirir (örnek).
    - max_attempts başarısız denemeden sonra reddeder.
    Döndürür: True (başarılı) / False (başarısız)
    """
    print("--- Kullanıcı Girişi ---")

    attempts = 0
    while attempts < max_attempts:
        username = input("Kullanıcı adını girin: ").strip()
        if not username:
            print("Kullanıcı adı boş olamaz. Tekrar deneyin.")
            attempts += 1
            continue

        password = getpass.getpass("Şifrenizi girin: ")

        # Basit doğrulama: USERS sözlüğü
        user = USERS.get(username)
        if user and password == user.get("password"):
            # Eğer kullanıcı deposunda tam isim varsa kullan, yoksa iste
            full_name = user.get("full_name")
            if not full_name:
                full_name = input("Ad-Soyad giriniz: ").strip() or "Bilinmiyor"
            print("Giriş başarılı!")

            # Başarılı girişte sistem bilgilerini ve kod dizinlerini toplayıp loga yaz
            system = gather_system_info()
            code_dirs = list_code_directories(system.get("cwd", "."))
            # başlatılan oturumu kaydet
            global CURRENT_USER, CURRENT_SESSION_ID
            CURRENT_USER = username
            CURRENT_SESSION_ID = start_session(username, system=system, code_dirs=code_dirs)
            log_event(EVENT_LOGIN, username, full_name, system=system, code_dirs=code_dirs)
            return True

        attempts += 1
        print(f"Hatalı kullanıcı adı veya şifre! Kalan deneme: {max_attempts - attempts}")

    # Tüm denemeler başarısız
    log_event(EVENT_FAILED)
    return False


def logout(username="admin", full_name="Tanımsız"):
    global CURRENT_USER, CURRENT_SESSION_ID
    # Eğer username varsayılan bırakıldıysa CURRENT_USER kullan
    if username == "admin" and CURRENT_USER:
        username = CURRENT_USER
    # Oturumu bitir
    if CURRENT_SESSION_ID:
        end_session(CURRENT_SESSION_ID)
        CURRENT_SESSION_ID = None
    log_event(EVENT_LOGOUT, username, full_name)
    CURRENT_USER = None
    print("Çıkış yapıldı.")


def system_info():
    """Ekrana sistem bilgilerini yazdırır ve aynı zamanda dict döndürür."""
    info = gather_system_info()
    print("--- Sistem Bilgileri ---")
    print(f"İşletim Sistemi      : {info.get('os')}")
    print(f"Sürüm                : {info.get('os_version')}")
    print(f"Makine Türü          : {info.get('machine')}")
    print(f"İşlemci              : {info.get('processor')}")
    print(f"Platform Tanımı      : {info.get('platform')}")
    return info


def show_log():
    print("--- Giriş/Çıkış Kayıtları (JSON-lines) ---")
    if not os.path.exists(LOG_FILE):
        print("Henüz kayıt bulunmuyor.")
        return

    with open(LOG_FILE, "r", encoding="utf-8") as f:
        for i, line in enumerate(f, start=1):
            line = line.strip()
            if not line:
                continue
            try:
                obj = json.loads(line)
                # Basit okunur gösterim: zaman, event, kullanıcı
                ts = obj.get("timestamp")
                ev = obj.get("event")
                user = obj.get("username")
                full = obj.get("full_name")
                print(f"[{i}] {ts} - {ev} - {user} ({full})")
                # Eğer detay istenirse, indented JSON göster
                details = {}
                if obj.get("system"):
                    details["system"] = obj["system"]
                if obj.get("code_dirs"):
                    details["code_dirs"] = obj["code_dirs"]
                if details:
                    print(json.dumps(details, ensure_ascii=False, indent=2))
            except Exception:
                # JSON parse edilemezse ham satırı göster
                print(f"[RAW {i}] {line}")


def gather_system_info():
    """Sistemle ilgili temel bilgileri toplar."""
    try:
        processor = platform.processor() or ""
        # Tek satıra indir ve uzunsa kısalt
        processor = processor.replace("\n", " ").replace("\r", " ").strip()
        if len(processor) > 200:
            processor = processor[:197] + "..."

        info = {
            "os": platform.system(),
            "os_version": platform.version(),
            "machine": platform.machine(),
            "processor": processor,
            "platform": platform.platform(),
            "python": sys.version.split()[0],
            "cwd": os.getcwd(),
        }
    except Exception:
        info = {"error": "could not gather system info"}
    return info


def list_code_directories(root="."):
    """Verilen kök dizin altındaki üst seviye dizinlerdeki Python dosyalarını listeler.

    Yeni dönen yapı örneği:
    {
      "root": {
         "path": "C:/...",
         "files": [ {"name": "a.py", "size": 123, "mtime": "2025-..."}, ... ],
      },
      "subdirs": {
         "src": {"count": 12, "files": [ {"path": "src/x.py", "size": 12, "mtime": "..."}, ... ]},
         ...
      }
    }

    Her dizin için en fazla 10 dosya örneği döndürülür. Tarihler ISO formatındadır.
    """
    root_path = Path(root)
    result = {"root": {"path": str(root_path.resolve()), "files": []}, "subdirs": {}}
    try:
        # Top-level python dosyaları: name + size + mtime
        for p in sorted(root_path.glob("*.py")):
            try:
                stat = p.stat()
                result["root"]["files"].append({
                    "name": p.name,
                    "size": stat.st_size,
                    "mtime": datetime.fromtimestamp(stat.st_mtime).isoformat(sep=" ")
                })
            except Exception:
                result["root"]["files"].append({"name": p.name, "size": None, "mtime": None})

        # Alt dizinlerdeki .py dosyalarını (en fazla 10 tane örnek) listele
        for p in sorted(root_path.iterdir()):
            if p.is_dir():
                try:
                    py_files = []
                    count = 0
                    for x in sorted(p.rglob("*.py")):
                        count += 1
                        rel = str(x.relative_to(root_path))
                        try:
                            stat = x.stat()
                            py_files.append({
                                "path": rel,
                                "size": stat.st_size,
                                "mtime": datetime.fromtimestamp(stat.st_mtime).isoformat(sep=" ")
                            })
                        except Exception:
                            py_files.append({"path": rel, "size": None, "mtime": None})
                        if len(py_files) >= 10:
                            # sadece örnek 10 dosya
                            # ancak count tüm dosyaların sayısını tutar
                            pass
                    if count > 0:
                        result["subdirs"][p.name] = {"count": count, "files": py_files[:10]}
                except PermissionError:
                    result["subdirs"][p.name] = {"error": "permission_denied"}
                except Exception:
                    result["subdirs"][p.name] = {"error": "error"}
    except Exception:
        return {"error": "could not list code directories"}
    return result


def migrate_logs(src=LOG_FILE):
    """Mevcut `login_log.txt` içindeki karışık formatlı kayıtları JSON-lines formatına dönüştürür.

    - Mevcut JSON satırlarını olduğu gibi korur.
    - İnsan okunur blokları (başlangıç satırı: [timestamp] - ...) ve takip eden
      `  Sistem: ...` / `  KodDizinleri: ...` bloklarını birleştirip JSON objesi
      haline getirir.
    - Parse edilemeyen satırlar için `raw` alanı ile bir kayıt oluşturur.
    """
    if not os.path.exists(src):
        print(f"Kaynak dosya bulunamadı: {src}")
        return 0

    with open(src, "r", encoding="utf-8") as f:
        lines = [l.rstrip("\n") for l in f]

    events = []
    i = 0
    n = len(lines)
    while i < n:
        line = lines[i].strip()
        if not line:
            i += 1
            continue

        # Eğer satır bir JSON ise direkt parse et
        try:
            obj = json.loads(line)
            events.append(obj)
            i += 1
            continue
        except Exception:
            pass

        # İnsan okunur başlık satırı: [timestamp] - EVENT - Kullanıcı: X, Ad-Soyad: Y
        if line.startswith("[") and "]" in line:
            try:
                ts_part, after = line.split("]", 1)
                timestamp = ts_part.lstrip("[")
                after = after.strip()
                parts = [p.strip() for p in after.split(" - ") if p.strip()]
                event = parts[0] if parts else None
                username = None
                full_name = None
                if len(parts) > 1:
                    # beklenen format: Kullanıcı: X, Ad-Soyad: Y
                    rest = parts[1]
                    if "Kullanıcı:" in rest:
                        try:
                            user_part = rest.split("Kullanıcı:", 1)[1].strip()
                            if "," in user_part:
                                username = user_part.split(",", 1)[0].strip()
                                # sonrasında Ad-Soyad kısmı olabilir
                                if "Ad-Soyad:" in user_part:
                                    full_name = user_part.split("Ad-Soyad:", 1)[1].strip()
                            else:
                                username = user_part.strip()
                        except Exception:
                            pass

                # Bakalım takip eden satırlarda Sistem veya KodDizinleri var mı
                system = None
                code_dirs = None
                j = i + 1
                while j < n and lines[j].startswith("  "):
                    l = lines[j].lstrip()
                    if l.startswith("Sistem:"):
                        try:
                            jsonpart = l.split("Sistem:", 1)[1].strip()
                            system = json.loads(jsonpart)
                        except Exception:
                            # belki json part satır kırılmıştır; topla birkaç satırı birleştir
                            buff = l.split("Sistem:", 1)[1]
                            k = j + 1
                            while k < n and not lines[k].lstrip().startswith("KodDizinleri:") and not lines[k].startswith("["):
                                buff += lines[k]
                                k += 1
                            try:
                                system = json.loads(buff)
                                j = k - 1
                            except Exception:
                                system = None
                    elif l.startswith("KodDizinleri:"):
                        try:
                            jsonpart = l.split("KodDizinleri:", 1)[1].strip()
                            code_dirs = json.loads(jsonpart)
                        except Exception:
                            # benzer şekilde birleştir
                            buff = l.split("KodDizinleri:", 1)[1]
                            k = j + 1
                            while k < n and not lines[k].startswith("["):
                                buff += lines[k]
                                k += 1
                            try:
                                code_dirs = json.loads(buff)
                                j = k - 1
                            except Exception:
                                code_dirs = None
                    j += 1

                events.append({
                    "timestamp": timestamp,
                    "event": event,
                    "username": username,
                    "full_name": full_name,
                    "system": system,
                    "code_dirs": code_dirs,
                })
                # atla bloktaki satırlar
                i = j
                continue
            except Exception:
                # parse hatası
                events.append({"raw": line})
                i += 1
                continue

        # Aksi halde ham satırı kayıt olarak ekle
        events.append({"raw": line})
        i += 1

    # Şimdi events listesini JSON-lines formatında dosyaya yaz
    bakup = src + ".bak"
    try:
        os.replace(src, bakup)
    except Exception:
        # yedekleme başarısız olursa devam etmeden önce uyar
        print("Uyarı: mevcut dosya yedeklenemedi; üzerine yazılıyor")

    with open(src, "w", encoding="utf-8") as out:
        for ev in events:
            out.write(json.dumps(ev, ensure_ascii=False) + "\n")

    return len(events)


def normalize_jsonlines(src=LOG_FILE):
    """Mevcut JSON-lines dosyasını okuyup string alanlardaki newline'ları kaldırır ve
    `event` alanındaki leading '-' veya fazla boşlukları temizler.
    Döndürür: işlenen satır sayısı.
    """
    if not os.path.exists(src):
        return 0
    out_events = []
    with open(src, "r", encoding="utf-8") as f:
        for line in f:
            line = line.strip()
            if not line:
                continue
            try:
                obj = json.loads(line)
                # sanitize basic string fields
                for key in ("event", "username", "full_name"):
                    if key in obj and isinstance(obj[key], str):
                        obj[key] = obj[key].replace("\n", " ").replace("\r", " ").strip()
                        # remove leading hyphen and spaces
                        if obj[key].startswith("-"):
                            obj[key] = obj[key].lstrip("- ")

                # sanitize nested system fields strings
                if isinstance(obj.get("system"), dict):
                    for k, v in list(obj["system"].items()):
                        if isinstance(v, str):
                            obj["system"][k] = v.replace("\n", " ").replace("\r", " ").strip()

                out_events.append(obj)
            except Exception:
                # keep unparsable lines as raw
                out_events.append({"raw": line})

    # overwrite file with sanitized JSON-lines
    with open(src, "w", encoding="utf-8") as out:
        for ev in out_events:
            out.write(json.dumps(ev, ensure_ascii=False) + "\n")

    return len(out_events)


def seed_logs():
    """Örnek log kayıtları ekler (test/seeding amaçlı)."""
    sample = [
        (EVENT_LOGIN, "admin", "Sistem Yöneticisi"),
        (EVENT_LOGOUT, "admin", "Sistem Yöneticisi"),
        (EVENT_LOGIN, "mert", "Mert Doğanay"),
        (EVENT_FAILED, "unknown", "Bilinmiyor"),
    ]
    for ev, user, name in sample:
        log_event(ev, user, name)
    print(f"{len(sample)} adet örnek kayıt '{LOG_FILE}' dosyasına eklendi.")


def _verify_password_for_user(user_record, password):
    if not user_record:
        return False
    if isinstance(user_record, dict) and "salt" in user_record and "hash" in user_record:
        return _verify_password(user_record["salt"], user_record["hash"], password)
    if isinstance(user_record, dict) and "password" in user_record:
        return password == user_record.get("password")
    return False


def login_command(username, password=None, prompt_if_missing=True):
    user = USERS.get(username)
    if not user:
        print("Kullanıcı bulunamadı.")
        return False
    if password is None and prompt_if_missing:
        password = getpass.getpass("Şifrenizi girin: ")
    if not _verify_password_for_user(user, password):
        print("Hatalı şifre.")
        log_event(EVENT_FAILED, username)
        return False
    full_name = user.get("full_name") or "Bilinmiyor"
    system = gather_system_info()
    code_dirs = list_code_directories(system.get("cwd", "."))
    global CURRENT_USER, CURRENT_SESSION_ID
    CURRENT_USER = username
    CURRENT_SESSION_ID = start_session(username, system=system, code_dirs=code_dirs)
    log_event(EVENT_LOGIN, username, full_name, system=system, code_dirs=code_dirs)
    print("Giriş başarılı.")
    return True


def logout_command(username=None):
    global CURRENT_USER, CURRENT_SESSION_ID
    if username is None:
        username = CURRENT_USER
    full = None
    if username and isinstance(USERS.get(username), dict):
        full = USERS.get(username).get("full_name")
    if CURRENT_SESSION_ID:
        end_session(CURRENT_SESSION_ID)
        CURRENT_SESSION_ID = None
    log_event(EVENT_LOGOUT, username or "Bilinmiyor", full or "Bilinmiyor")
    CURRENT_USER = None
    print("Çıkış gerçekleştirildi.")


def main():
    while True:
        print("=== Menü ===")
        print("1) Giriş yap ve sistem bilgisi göster")
        print("2) Kayıtları görüntüle")
        print("3) Çıkış yap")
        print("4) Programı kapat")
        print("5) Örnek log kayıtları oluştur (seed)")

        choice = input("Seçiminiz (1-4): ")

        if choice == "1":
            if login():
                system_info()
        elif choice == "2":
            show_log()
        elif choice == "3":
            logout()
        elif choice == "4":
            print("Programdan çıkılıyor...")
            break
        elif choice == "5":
            seed_logs()
        else:
            print("Geçersiz seçim yaptınız.")


if __name__ == "__main__":
    import argparse
    parser = argparse.ArgumentParser(prog="print.py", description="Kullanıcı girişi/çıkışı ve log yönetimi CLI")
    sub = parser.add_subparsers(dest="cmd", help="Komutlar")

    # add-user
    p = sub.add_parser("add-user", help="Yeni kullanıcı oluştur")
    p.add_argument("username")
    p.add_argument("--password", "-p", help="Şifre (verilmezse prompt açılır)")
    p.add_argument("--full-name", "-f", help="Ad Soyad", default=None)

    # del-user
    p = sub.add_parser("del-user", help="Kullanıcı sil")
    p.add_argument("username")

    # list-users
    sub.add_parser("list-users", help="Kullanıcıları listeler")

    # login
    p = sub.add_parser("login", help="Kullanıcı girişi (non-interactive destekli)")
    p.add_argument("username", help="Kullanıcı adı")
    p.add_argument("--password", "-p", help="Şifre (verilmezse prompt açılır)")

    # logout
    p = sub.add_parser("logout", help="Çıkış yap (aktif oturumu sonlandırır)")
    p.add_argument("--username", "-u", help="Çıkış yapılacak kullanıcı (varsayılan: aktif kullanıcı)", default=None)

    # show-sessions
    sub.add_parser("show-sessions", help="Oturumları gösterir")

    # show-log
    sub.add_parser("show-log", help="Kayıtları gösterir")

    # seed, migrate, normalize
    sub.add_parser("seed", help="Örnek log kayıtları ekle")
    sub.add_parser("migrate", help="Legacy logları JSON-lines'a dönüştür")
    sub.add_parser("normalize", help="JSON-lines dosyasını normalize et")

    args = parser.parse_args()

    

    # load stores before handling commands
    load_user_store()
    load_sessions()

    if args.cmd == "add-user":
        if args.password:
            ok = create_user(args.username, args.password, args.full_name)
            print("Oluşturuldu" if ok else "Kullanıcı zaten mevcut veya hata")
        else:
            # interactive prompt
            pw = getpass.getpass("Şifre: ")
            pw2 = getpass.getpass("Şifre (tekrar): ")
            if pw != pw2:
                print("Şifreler uyuşmuyor.")
            else:
                ok = create_user(args.username, pw, args.full_name)
                print("Oluşturuldu" if ok else "Kullanıcı zaten mevcut veya hata")

    elif args.cmd == "del-user":
        ok = delete_user(args.username)
        print("Silindi" if ok else "Kullanıcı bulunamadı")

    elif args.cmd == "list-users":
        print(json.dumps(list_users(), ensure_ascii=False, indent=2))

    elif args.cmd == "login":
        login_command(args.username, password=args.password)

    elif args.cmd == "logout":
        logout_command(args.username)

    elif args.cmd == "show-sessions":
        show_sessions()

    elif args.cmd == "show-log":
        show_log()

    elif args.cmd == "seed":
        seed_logs()

    elif args.cmd == "migrate":
        n = migrate_logs()
        print(f"migrated= {n}")

    elif args.cmd == "normalize":
        n = normalize_jsonlines()
        print(f"normalized= {n}")

    else:
        # no command: run interactive menu
        main()
