# Programming-Project — Giriş/Çıkış Kayıt ve CLI

Bu proje, yerel bir Python betiği (`print.py`) üzerinden kullanıcı giriş/çıkışlarını kaydetmeyi, oturum yönetimini ve sistem + kod dizini bilgisini JSON-lines formatında loglamayı sağlar.

## Hızlı Özet
- Loglar: `login_log.txt` (JSON-lines; her satır bir JSON objesidir)
- Kullanıcı deposu: `users.json` (salt + hash şeklinde parola saklama)
- Oturumlar: `sessions.json`
- CLI: `print.py` içinde argparse tabanlı komutlar sağlanır (aşağıda kullanım örnekleri).

---

## CLI Kullanımı (kısa)

`print.py` şu komutları sağlar:

- `add-user <username> [-p PASSWORD] [-f FULL_NAME]` : Yeni kullanıcı oluşturur. `-p` verilmezse prompt açılır.
- `del-user <username>` : Kullanıcıyı siler.
- `list-users` : Kayıtlı kullanıcıları listeler.
- `login <username> [-p PASSWORD]` : Non-interactive giriş yapar; parola verilmezse prompt açılır.
- `logout [--username USER]` : Aktif oturumu sonlandırır (veya verilen kullanıcı için çıkış kaydı ekler).
- `show-sessions` : Mevcut ve geçmiş oturumları gösterir.
- `show-log` : `login_log.txt` içindeki JSON-lines kayıtlarını okunur formatta gösterir.
- `seed` : Örnek log kayıtları ekler (test amaçlı).
- `migrate` : Legacy insan okunur logları JSON-lines formatına dönüştürür; orijinali `login_log.txt.bak` olarak yedekleyebilir.
- `normalize` : JSON-lines içindeki satır içi yeni satırları ve fazladan önekleri temizler.

### Hızlı PowerShell örnekleri

Kullanıcı oluşturma (parola argümanlı):

```powershell
python .\print.py add-user alice -p s3cr3t -f "Alice Example"
```

Kullanıcı oluşturma (parola prompt ile):

```powershell
python .\print.py add-user bob
# Parola sorulacak, tekrar onaylanacak
```

Giriş (non-interactive):

```powershell
python .\print.py login alice -p s3cr3t
```

Giriş (prompt ile):

```powershell
python .\print.py login bob
# Parola prompt ile alınır
```

Kayıtları ve oturumları görüntüleme:

```powershell
python .\print.py show-log
python .\print.py show-sessions
```

Legacy logları dönüştürme ve normalize etme (yedekleme yapar):

```powershell
python .\print.py migrate
python .\print.py normalize
```

## Güvenlik Notları

- Komut satırında parola (`-p`) kullanmak rahat olsa da güvenlik riski (shell history) getirir. Mümkünse parola argümanını kullanmayın; prompt kullanın.
- Parolalar `users.json` içinde `salt` + `hash` formatında saklanır. Eğer eski `password` alanları varsa migration önerilir.

## Hızlı Hata Ayıklama

- Eğer `show-log` içinde ham `RAW` satırlar görüyorsanız, `python .\print.py migrate` çalıştırarak legacy kayıtları dönüştürebilirsiniz. Orijinal dosya `login_log.txt.bak` olarak yedeklenir.

---

Bu README, proje ile hızlı çalışmaya başlamanız için kısa bir rehberdir. Daha ayrıntılı açıklama veya başka entegrasyonlar isterseniz söyleyin, ben ekleyeyim.
