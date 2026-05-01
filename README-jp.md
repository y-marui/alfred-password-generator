# Alfred Password Generator

> **これは日本語版（正本）です。**
> 英語版（参照）は [README.md](README.md) を参照してください。

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![CI](https://github.com/y-marui/alfred-password-generator/actions/workflows/ci.yml/badge.svg)](https://github.com/y-marui/alfred-password-generator/actions/workflows/ci.yml)
[![GitHub Sponsors](https://img.shields.io/github/sponsors/y-marui)](https://github.com/sponsors/y-marui)
[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-ffdd00?logo=buy-me-a-coffee&logoColor=black)](https://www.buymeacoffee.com/y.marui)

| 項目 | 内容 |
|---|---|
| 開発対象 | Alfred 5 Script Filter ワークフロー |
| ライセンス | MIT |
| 動作環境 | Python 3.9+, Alfred 5 |

長さと文字セットを自由に指定してパスワードを生成する Alfred ワークフロー。

## Usage

Alfred を開いて `passgen` に続けてスペースを入力します。

### 基本 (default)

```
passgen [length] [pattern]
```

`length` のデフォルトは 18、`pattern` のデフォルトは `A-Za-z0-9`。

### 記号付き (panc)

```
passgen panc [length] [pattern]
```

デフォルトパターンが `A-Za-z0-9!-*`（`!@#^&*` を含む）になります。

### グループ分割 (split)

```
passgen split [length] [by] [pattern]
```

`length` のデフォルトは 18、`by` のデフォルトは 6（ハイフン区切り: `xxxxxx-xxxxxx-xxxxxx`）。
`length` は `by` の倍数である必要があります。

### 記号付きグループ分割 (panc split)

```
passgen panc split [length] [by] [pattern]
```

記号を含みつつグループ分割します。

Enter キーで選択したパスワードをクリップボードにコピーします。

## パターン構文

文字は直接列挙（例: `ABCabc012!@#`）またはレンジ指定（例: `A-Za-z0-9`）で指定できます。

記号のレンジ: `!-*` は `!@#^&*` に展開されます。

| パターン例 | 展開結果 |
|---|---|
| `A-Z` | 大文字アルファベット |
| `a-z` | 小文字アルファベット |
| `0-9` | 数字 |
| `!-*` | `!@#^&*` |
| `A-Za-z0-9` | 英数字 |
| `A-Za-z0-9!-*` | 英数字 + 記号 |

## Configuration

Alfred の設定（`⌘,`）から以下の項目を設定できます。

| 設定 | デフォルト | 説明 |
|---|---|---|
| Use uv | ON | uv がインストールされている場合に `uv run python` で実行 |
| Clipboard History | OFF | パスワードを Alfred のクリップボード履歴に保存する（セキュリティ上 OFF 推奨） |
| Log Level | WARNING | ログの詳細度（開発時は DEBUG、本番は WARNING） |

## Installation

```bash
make install    # 開発用依存関係をインストール
make build      # ワークフローパッケージをビルド
# → dist/*.alfredworkflow
```

`dist/*.alfredworkflow` をダブルクリックして Alfred にインストールします。

## Project Structure

```
alfred-password-generator/
├── src/
│   ├── alfred/         # Alfred SDK (response, router, config, logger, safe_run)
│   └── app/
│       ├── commands/   # passgen_cmd, config_cmd, help_cmd
│       └── services/   # passgen_service (コアロジック)
├── workflow/           # Alfred パッケージ (info.plist, scripts/entry.py)
└── tests/              # pytest テストスイート
```

## License

MIT — [LICENSE](LICENSE) を参照

---

*この文書には英語版（参照版）[README.md](README.md) があります。編集時は同一コミットで更新してください。*
