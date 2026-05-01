# Alfred Workflow Template

> **これは日本語版（正本）です。**
> 英語版（参照）は [README.md](README.md) を参照してください。

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![CI](https://github.com/y-marui/alfred-workflow-template/actions/workflows/ci.yml/badge.svg)](https://github.com/y-marui/alfred-workflow-template/actions/workflows/ci.yml)
[![Charter Check](https://github.com/y-marui/alfred-workflow-template/actions/workflows/dev-charter-check.yml/badge.svg)](https://github.com/y-marui/alfred-workflow-template/actions/workflows/dev-charter-check.yml)
[![GitHub Sponsors](https://img.shields.io/github/sponsors/y-marui)](https://github.com/sponsors/y-marui)
[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-ffdd00?logo=buy-me-a-coffee&logoColor=black)](https://www.buymeacoffee.com/y.marui)

| 項目 | 内容 |
|---|---|
| 開発対象 | Alfred 5 Script Filter ワークフロー |
| 開発環境 | 個人〜小規模チーム（1〜3人） |
| 主言語 | 英語（OSS） |
| ライセンス | MIT |
| 動作環境 | Python 3.9+, Alfred 5 |
| AI ツール | Claude Code / GitHub Copilot / Gemini CLI |

> Alfred 5 Script Filter ワークフローのプロダクションレディなテンプレート。
> 10分で開発を開始できます。

## Features

- ✅ **レイヤードアーキテクチャ** — Alfred 境界とビジネスロジックを分離
- ✅ **軽量 Alfred SDK** — レスポンスビルダー、ルーター、キャッシュ、設定、ロガー
- ✅ **コマンドベース UX** — `wf search`、`wf open`、`wf config`、`wf help`
- ✅ **フルテストスイート** — pytest で Alfred なしにテスト実行可能
- ✅ **CI/CD** — GitHub Actions でリント・テスト・ビルド・リリースを自動化
- ✅ **ベンダーパッケージング** — サードパーティ依存を `vendor/` にバンドル
- ✅ **AI 対応** — `AI_CONTEXT.md` + `CLAUDE.md` で AI アシスタントのコンテキストを管理

## Requirements

- Alfred 5（Script Filter には Powerpack が必要）
- Python 3.9+
- [pre-commit](https://pre-commit.com/)（セキュリティフック用）

## Quick Start

### Using this template

1. GitHub で **"Use this template"** → **"Create a new repository"** をクリック
2. 新しいリポジトリをクローンして AI ツールで開く
3. AI に「AI_CONTEXT.md の初期セットアップを実行して」と伝える — AI が以下を実行します:
   - GitHub リポジトリ設定を適用
   - `README_TEMPLATE.md` → `README.md`、`README_TEMPLATE-jp.md` → `README-jp.md` にリネーム
   - `{user}`、`{repo}`、`{keyword}` プレースホルダを置換
   - dev-charter の cron スケジュールをランダム化
4. ワークフローをカスタマイズする（`DEVELOPING.md` 参照）

### Development (this template)

```bash
git clone https://github.com/yourname/alfred-workflow-template
cd alfred-workflow-template

# 開発用依存関係をインストール
make install

# Alfred をローカルでシミュレート
make run Q="search foo"
make run Q="help"

# テストを実行
make test

# ワークフローパッケージをビルド
make build
# → dist/workflow-template-0.1.0.alfredworkflow
```

`dist/*.alfredworkflow` をダブルクリックして Alfred にインストールします。

## Usage

Alfred を開いて `wf` に続けてスペースを入力します。

### Search (default)

```
wf <query>
wf search <query>
```

クエリを入力して検索します。Enter を押すと結果を開きます。

| キー | 操作 |
|---|---|
| ↩ Enter | 結果を開く |
| ⌘C | 結果の URL をコピー |

### Open

```
wf open <name>
```

ショートカットを開きます。利用可能なショートカット: `repo`、`docs`、`issues`

### Config

```
wf config
wf config reset
```

現在の設定を確認、またはすべての設定をリセットします。

### Help

```
wf help
```

利用可能なコマンド一覧を表示します。

### Tips

- ワークフローは最もよく使った結果を記憶します（Alfred の学習機能）。
- API 呼び出しを最小化するため、結果は 5 分間キャッシュされます。
- `⌘,` で Alfred のワークフロー設定にアクセスできます。

### Troubleshooting

**結果が表示されない場合**
- Alfred のデバッガーを確認: Alfred を開いて `⌘D`
- ログを確認: `~/Library/Logs/Alfred/Workflow/<bundle-id>.log`

**結果が古い場合**
- キャッシュ TTL は 5 分です。期限切れを待つか、手動でクリア: `wf config reset`

## Project Structure

```
alfred-workflow-template/
├── src/
│   ├── alfred/         # Alfred SDK (response, router, cache, config, logger, safe_run)
│   └── app/            # アプリケーション層 (commands, services, clients)
├── workflow/           # Alfred パッケージ (info.plist, scripts/entry.py, vendor/)
├── tests/              # pytest テストスイート
├── scripts/            # build.sh, dev.sh, release.sh, vendor.sh
└── docs/               # アーキテクチャ・リファレンスドキュメント
```

## Documentation

| ドキュメント | 内容 |
|---|---|
| [docs/architecture.md](docs/architecture.md) | モジュール・レイヤー構造 |
| [docs/file-map.md](docs/file-map.md) | ファイルレベルの依存関係マップ |
| [docs/specification.md](docs/specification.md) | 機能仕様・データフロー |
| [docs/ui-design.md](docs/ui-design.md) | Alfred 結果アイテムの UI 設計指針 |
| [docs/configuration-builder.md](docs/configuration-builder.md) | Alfred Configuration Builder リファレンス |

## AI-Assisted Development

このテンプレートは AI 支援開発に対応しています。

| ツール | 役割 |
|---|---|
| Claude Code | アーキテクチャ設計・大規模変更・リファクタリング |
| GitHub Copilot | バグ修正・細かな実装・単体テスト作成 |
| Gemini CLI | ドキュメント管理 |

セッションコンテキスト: [`AI_CONTEXT.md`](AI_CONTEXT.md)、[`CLAUDE.md`](CLAUDE.md)

## Customizing This Template

上記 Quick Start の初期セットアップが完了したら、ワークフローをカスタマイズします:

1. `workflow/info.plist` を編集:
   - `bundleid` を自分のバンドル ID に変更（例: `com.yourname.workflowname`）
   - キーワード（`wf`）を自分のトリガーキーワードに変更
   - `uuidgen` で生成した UUID に置き換え
2. `src/app/clients/api_client.py` を実際の API クライアントに置き換え
3. `pyproject.toml` のワークフロー名を更新
4. `src/app/commands/open_cmd.py` のショートカットを更新
5. `workflow/icon.png` を追加

## Release

```bash
# 1. pyproject.toml のバージョンを更新
# 2. タグを付けてプッシュ
git tag v1.2.3
git push --tags
# GitHub Actions が .alfredworkflow をビルドして GitHub Release を作成
```

## License

MIT — [LICENSE](LICENSE) を参照

---

*この文書には英語版（参照版）[README.md](README.md) があります。編集時は同一コミットで更新してください。*
