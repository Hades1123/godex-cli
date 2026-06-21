# Chạy godex trên Windows

Có 2 cách: **(A)** chỉ chép file `.exe` đã build sẵn (nhanh, không cần cài Go), hoặc **(B)** chuyển cả source code sang Windows để build + phát triển tiếp.

---

## Cách A — Nhanh: chỉ chép `.exe` (không cần Go)

Từ Linux, đã build sẵn:

```bash
GOOS=windows GOARCH=amd64 go build -o godex.exe .
```

Chép `godex.exe` sang Windows (USB, Google Drive, `scp`, ...). Trên Windows:

1. Tạo thư mục chứa, ví dụ `C:\Users\<bạn>\bin\`
2. Bỏ `godex.exe` vào đó
3. Thêm thư mục đó vào PATH:
   - Win + R → gõ `sysdm.cpl` → tab **Advanced** → **Environment Variables**
   - Phần *User variables* → chọn `Path` → **Edit** → **New** → dán `C:\Users\<bạn>\bin`
   - OK hết, mở terminal mới
4. Test: `godex --help`

**Hạn chế:** mỗi lần sửa code phải build lại trên Linux + chép lại. Không xem được source trên Windows.

---

## Cách B — Đầy đủ: chuyển source + build trên Windows

### B1. Đưa source code sang Windows

**Cách tốt nhất — qua GitHub** (đồng bộ 2 chiều dễ dàng):

```bash
# Trên Linux: tạo repo trên github.com trước, rồi:
git remote add origin git@github.com:<user>/godex.git
git push -u origin main
```

Trên Windows:

```powershell
git clone https://github.com/<user>/godex.git
cd godex
```

**Hoặc copy trực tiếp** (nếu không dùng GitHub): nén cả thư mục (`tar -czf godex.tar.gz .`) hoặc zip, chép sang Windows rồi giải nén.

### B2. Cài Go trên Windows

Cách nhanh nhất (PowerShell):

```powershell
winget install GoLang.Go
```

Hoặc tải installer từ https://go.dev/dl/ (file `.msi`) → chạy → next hết.

Verify:

```powershell
go version    # cần >= 1.26
```

> Đóng rồi mở lại terminal sau khi cài để Go nhận PATH.

### B3. Build

```powershell
cd godex
go build -o godex.exe .
```

Test:

```powershell
.\godex.exe --help
.\godex.exe config list
```

### B4. Cài global (chạy ở đâu cũng được)

```powershell
# Tạo thư mục bin trong home
mkdir $env:USERPROFILE\bin

# Copy binary vào đó
copy godex.exe $env:USERPROFILE\bin\

# Thêm vào PATH (cập nhật biến User vĩnh viễn)
[Environment]::SetEnvironmentVariable("Path", $env:USERPROFILE\bin + ";" + $env:Path, "User")
```

Đóng + mở terminal lại, test:

```powershell
godex config current
```

> Mỗi lần sửa code: `go build -o $env:USERPROFILE\bin\godex.exe .` — build thẳng vào PATH luôn.

---

## C5. Thiết lập preset trên Windows

Preset dir trên Windows là `%AppData%\godex\presets\` (tức `C:\Users\<bạn>\AppData\Roaming\godex\presets\`).

Tạo thư mục + 2 file preset:

```powershell
mkdir $env:APPDATA\godex\presets
```

Copy nội dung `glm.json` / `deepseek.json` từ Linux sang 2 file cùng tên trong thư mục đó. Hoặc tạo nhanh từ settings hiện tại:

```powershell
copy $env:USERPROFILE\.claude\settings.json $env:APPDATA\godex\presets\deepseek.json
# rồi sửa model/token/base URL trong file đó cho ra preset khác
```

Sau đó:

```powershell
godex config list
godex config use glm
godex config current
```

Claude Code trên Windows lưu settings ở `C:\Users\<bạn>\.claude\settings.json` — đúng chỗ `godex config` ghi đè.

---

## Lưu ý các tính năng trên Windows

| Tính năng | Windows | Ghi chú |
|-----------|---------|---------|
| `config` | ✅ Đầy đủ | File I/O + JSON, chạy tốt |
| `java`/`node` | ⚠️ Hạn chế | Scan path Linux (`/usr/lib/jvm`, `~/.nvm`); cần thêm path Windows (`C:\Program Files\Java`) sau |
| `ports list` | ⚠️ Lỗi runtime | Lệnh `ss` không có trên Windows |
| `ports kill` | ⚠️ Báo "not supported" | Cần implement bằng `netstat` + `taskkill` sau |

→ Trên Windows hiện nên tập trung dùng `config`. Các tính năng khác sẽ cần port thêm nếu muốn dùng đầy đủ.
