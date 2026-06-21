# godex config — Chuyển đổi Claude Code Settings

Chuyển đổi nhanh giữa các model/API preset cho Claude Code mà không cần mở file `~/.claude/settings.json` sửa tay.

## Preset files

Presets nằm ở `~/.config/godex/presets/*.json`. Mỗi file `.json` là một bản sao hoàn chỉnh của `settings.json` mà Claude Code dùng.

```
~/.config/godex/presets/
├── deepseek.json    # DeepSeek API
└── glm.json         # GLM (z.ai) API
```

Muốn thêm preset mới (VD: Anthropic chính chủ, OpenAI...), chỉ cần thả file `.json` vào thư mục này.

## Lệnh

```
godex config list       # Liệt kê tất cả presets, dấu * là preset đang active
godex config current    # Xem thông tin settings hiện tại (model, API URL, preset)
godex config use <tên>  # Chuyển sang preset khác, tự backup settings cũ
```

## Ví dụ

```bash
# Xem danh sách
$ godex config list
* deepseek
  glm

# Xem đang dùng gì
$ godex config current
Model:  deepseek-v4-pro[1m]
API:    https://api.deepseek.com/anthropic
Preset: deepseek

# Chuyển sang GLM
$ godex config use glm
Backed up previous settings to /home/hades/.claude/settings.json.bak
Switched to preset "glm" → /home/hades/.claude/settings.json

# Kiểm tra lại
$ godex config current
Model:  glm-5.2[1m]
API:    https://api.z.ai/api/anthropic
Preset: glm
```

## Cơ chế

| Việc | Chi tiết |
|------|----------|
| `use` | Copy file preset → `~/.claude/settings.json` |
| Backup | Trước khi ghi đè, settings cũ được copy ra `settings.json.bak` |
| Detect preset | So khớp `model` + `ANTHROPIC_BASE_URL` giữa settings hiện tại và các preset |

Sau khi chạy `config use`, đóng Claude Code và mở lại để model mới có hiệu lực.

## Thêm preset mới

Tạo file JSON với nội dung giống `settings.json`:

```bash
cp ~/.claude/settings.json ~/.config/godex/presets/anthropic.json
# Sửa model, token, base URL... trong file vừa copy
```

Sau đó `godex config list` sẽ thấy preset mới, và dùng `godex config use anthropic` để chuyển.
