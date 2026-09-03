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
| 動作環境 | Go（ビルド時のみ、`go.mod` 参照）, Alfred 5 |

長さと文字セットを自由に指定してパスワードを生成する Alfred ワークフロー。

## Usage

Alfred を開いて `passgen` に続けてスペースを入力します。

### Basic (default)

```
passgen [length] [pattern]
```

`length` のデフォルトは 18、`pattern` のデフォルトは `A-Za-z0-9`。

### With punctuation (panc)

```
passgen panc [length] [pattern]
```

デフォルトパターンが `A-Za-z0-9!-*`（`!@#^&*` を含む）になります。

### Split into groups (split)

```
passgen split [length] [by] [pattern]
```

`length` のデフォルトは 18、`by` のデフォルトは 6（ハイフン区切り: `xxxxxx-xxxxxx-xxxxxx`）。
`length` は `by` の倍数である必要があります。

### Split with punctuation (panc split)

```
passgen panc split [length] [by] [pattern]
```

記号を含みつつグループ分割します。

Enter キーで選択したパスワードをクリップボードにコピーします。

## Pattern syntax

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
| Clipboard History | OFF | パスワードを Alfred のクリップボード履歴に保存する（セキュリティ上 OFF 推奨） |

## Installation

```bash
make build-workflow   # ワークフローパッケージをビルド
# → dist/*.alfredworkflow
```

`dist/*.alfredworkflow` をダブルクリックして Alfred にインストールします。

## Project Structure

```
alfred-password-generator/
├── cmd/
│   └── password-generator-alfred/  # Alfred が起動するバイナリ
├── internal/
│   ├── passgen/         # パスワード生成ロジック（コア）
│   ├── passgencmd/      # コマンドディスパッチ・引数解析・help
│   └── scriptfilter/    # Alfred Script Filter JSON 型
└── workflow/            # Alfred パッケージ (info.plist, icon.png)
```

## License

MIT — [LICENSE](LICENSE) を参照

---

*この文書には英語版（参照版）[README.md](README.md) があります。編集時は同一コミットで更新してください。*
