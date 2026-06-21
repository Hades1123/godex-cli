# Dev Workflow — Code → Build → Global Install

## 1. Cài đặt Go

```bash
# Kiểm tra Go đã cài chưa
go version   # cần >= 1.26

# Nếu chưa có, tải từ https://go.dev/dl/
wget https://go.dev/dl/go1.26.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.26.4.linux-amd64.tar.gz
```

## 2. Thiết lập biến môi trường (zsh)

Thêm vào `~/.zshrc`:

```bash
# Go
export GOPATH=$HOME/go
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

# DevVM binary (nếu để riêng trong project)
export PATH=$PATH:$HOME/.local/bin
```

Reload:

```bash
source ~/.zshrc
```

> **Giải thích:**
> - `/usr/local/go/bin` — nơi chứa `go` compiler
> - `$GOPATH/bin` (`~/go/bin`) — nơi `go install` đặt binary
> - `~/.local/bin` — nơi ta chọn để cài `godex` (đã có sẵn trong PATH của Ubuntu/zsh)

## 3. Clone & cài dependencies

```bash
cd ~/Code/Core/Golang/CLI
go mod tidy        # tải tất cả dependencies về cache
```

## 4. Quy trình code → build → test

```bash
# ── CODING ──
vim cmd/ports.go             # thêm command mới
vim internal/runtime/xxx.go  # thêm business logic

# ── BUILD (một lệnh duy nhất) ──
./build.sh                   # build + vet + cài vào PATH

# ── TEST ──
godex ports                  # test từ bất kỳ đâu
godex config list
```

## 5. Workflow tóm tắt

```
┌──────────┐     ┌──────────────────────┐
│  viết    │────►│  ./build.sh          │
│  code    │     │  (build + vet + cài) │
└──────────┘     └──────┬───────────────┘
                        │
                   ┌────▼────┐
                   │ godex   │  ← dùng ngay
                   │ <lệnh>  │
                   └─────────┘
```

## 6. Cấu trúc thư mục liên quan

```
~/Code/Core/Golang/CLI/    ← source code (project)
    ├── build.sh           ← build script (chạy sau mỗi lần sửa code)
    ├── main.go
    ├── cmd/                ← cobra commands
    ├── internal/
    │   └── runtime/        ← business logic
    └── docs/               ← tài liệu

~/.local/bin/godex          ← binary toàn cục (trong PATH)
```

## 8. Troubleshooting

| Vấn đề | Cách fix |
|--------|---------|
| `godex: command not found` | `hash -r` (zsh cache), hoặc kiểm tra `echo $PATH \| grep local/bin` |
| `go: module not found` | Chạy `go mod tidy` trong thư mục project |
| `panic: strings.Builder copied` | Đã fix — nhớ pass `*strings.Builder` pointer |
| `fuser: no process found` | Port đó không có process nào listen, hoặc cần sudo với system port (<1024) |
