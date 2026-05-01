# AI_CONTEXT.md — alfred-password-generator

> このファイルは開発憲章（`docs/dev-charter/`）をこのプロジェクト向けにまとめたものです。
> AIツールはセッション開始時にこのファイルを読むことで、憲章全体を参照しなくても
> プロジェクトの方針を把握できます。

## Reference Order

AI はタスク開始時に以下の順で参照する:

1. `README.md`（概要・セットアップ）
2. `DEVELOPING.md`（ビルド・実装規約・命名規則）

必要に応じて以下を参照する（順不同）:
- `CONTRIBUTING.md`（PR・Issue ルール）
- `docs/architecture.md`（モジュール・コンポーネント構造）
- `docs/file-map.md`（ファイルレベルの依存関係 ※情報が足りない・古い場合は適宜探索し、追記・更新する）
- `docs/specification.md`（機能仕様・データフロー）
- `docs/ui-design.md`（UI 設計・コンポーネント仕様）

不明点は `docs/dev-charter/CHARTER_INDEX.md` → 該当ファイルの順で参照する。

---

## Project Overview

Alfred 5 Script Filter ワークフロー用の OSS テンプレート。
Python 3.9+、レイヤードアーキテクチャ、CI/CD 完備。
対象: 個人〜3人規模の開発チーム。ライセンス: MIT。

```
src/alfred/     ← Alfred SDK（response / router / cache / config / logger / safe_run）
src/app/        ← アプリケーション層（commands / services / clients）
workflow/       ← Alfred パッケージ（info.plist / scripts/entry.py / vendor/）
tests/          ← pytest テストスイート
scripts/        ← build.sh / dev.sh / release.sh / vendor.sh
```

詳細アーキテクチャ: `docs/architecture.md`

---

## Applied Charter Principles

### AIコンテキスト優先順位（AI_CONTEXT_HIERARCHY）

1. タスクコンテキスト（Issue / Pull Request）
2. **プロジェクトコンテキスト（このファイル・プロジェクトドキュメント）** ← ここ
3. 開発憲章（`docs/dev-charter/`）
4. グローバルコンテキスト

### 開発原則（PRINCIPLES）

#### 基本哲学
- **ローカルファースト** — Alfred ワークフローはオフラインで動作することを前提にする
- **インフラ最小化** — サーバーレス、外部依存なし（vendor/ に完結）
- **小さく始める** — 機能追加は必要性が確認されてから

#### コード設計
- **変更範囲は必要最小限** — Over-engineering しない
- **YAGNI** — 今必要ない機能は実装しない
- **DRY** — 2回の重複では抽象化しない。3回目で検討する
- **既存コードの再利用** — 新規実装前に類似機能がないか確認する
- **TODO/FIXME を残さない** — 実装するか、Issue として記録する（テンプレートの `# TODO:` コメントは「ユーザーが置き換える場所の目印」として例外的に許可）
- **既存パターンに従う** — 命名規則・アーキテクチャ・ディレクトリ構造を統一する

### AI協働ルール（AI_COLLABORATION_RULES）

#### AI行動原則
- **Scope 厳守** — 会話のタスク・ゴールを AI が勝手に変更しない
- **不明点は作業前に1回でまとめて質問する** — 重要な情報不足や曖昧さは質問する。軽微な不足は合理的な仮定で補い、仮定を明示する。推測で断定しない

#### コーディング前の確認必須項目
- ゴール（完了条件）
- 言語・FW・バージョン制約
- 新規 or 既存コード修正
- テストの要否
- 影響範囲

確認不要（既存コードに合わせて進める）: コードスタイル / ファイル配置 / 軽微な実装詳細

#### エラー対応
- **原因分析 → 修正方針説明 → 実装** の順で進める
- エラーログ・スタックトレースは全文確認してから対応
- 推測で修正しない（必要なら既存コードを確認する）
- デバッグ用の `print` 文は本番コードに残さない

#### 作業スタンス
- 大きな変更前に方針を説明してから着手する
- **不要な依存追加禁止** — 既存の依存で解決できないか先に検討する
- **ドキュメント同期** — 仕様・ルール・構成に変更が生じたとき、変更と同じ作業内で関連ドキュメントを更新する（対象: `docs/` 内ファイル、`AI_CONTEXT.md`、`README.md` 等）

#### dev-charter 変更ルール（dev-charter Modification Rules）

`docs/dev-charter/` 配下のファイルを**直接編集しない**。

- 変更が必要な場合は dev-charter リポジトリ本体に Issue を立て、`git subtree pull` でアップデートを取り込む
- `git subtree pull` によるアップデートのみ許可する
- このプロジェクト固有のルールは `AI_CONTEXT.md` または専用ファイルに記載する

#### 憲章の参照方法（Charter Lookup）

不明点が憲章に関係する場合は**全ファイルを検索せず**、以下の手順で参照する:

1. `docs/dev-charter/CHARTER_INDEX.md` を読み、該当トピックのファイルを特定する
2. 特定したファイル（原則 1〜2 件）のみを読む
3. 参照後にユーザーへ提案・確認を行う

推測で断定せず、憲章を参照してからユーザーに提案・質問する。

#### GitHub Operations

Issue を作成する場合は、必ずリポジトリオーナーを `assignee` に設定する。

```bash
gh issue create --title "..." --body "..." --assignee @me
```

### 言語ポリシー（LANGUAGE_POLICY）

OSS プロジェクトのため、**公開面は英語を主言語**とする。
**日本語版が編集の起点（正本）であり、英語版はその翻訳として同期する。**

| 対象 | 言語 |
|---|---|
| `README.md` | 英語（参照版） |
| `README-jp.md` | 日本語（正本） |
| コミットメッセージ | 英語 |
| Issue / PR のタイトルと本文 | 英語 |
| 公開 API / public 関数 docstring | 英語 |
| examples/ のコメント | 英語 |
| エラーメッセージ・ログ | 英語 |
| private 関数・実装詳細のコメント | 日本語 OK |
| 変数名・識別子 | 英語 |

両言語ファイルが存在する場合: **日本語を正本として編集し、英語はそれに合わせて同一コミットで更新する**。

### プロジェクトライフサイクル（PROJECT_LIFECYCLE）

- 規模: 個人〜3人。アジャイルで迅速な意思決定
- **コミット粒度** — 機能単位・動作確認 OK 後にコミット
- **コミットメッセージ** — Conventional Commits 形式（feat / fix / refactor / docs / chore）
- **WIP 禁止** — 動作しないコードはコミットしない

### セキュリティ（SECURITY_POLICY）

#### 二層構造
1. **個人 git フック**（`~/.config/git/hooks/pre-commit`）— 開発者個人のマシン全体に適用
2. **per-repo pre-commit フック**（`.pre-commit-config.yaml`）— チーム強制・CI でも動作

#### 自動ブロック項目
- `anonymous` のままコミット（個人 git フック側で対応。per-repo フックでは検知しない）
- `.env` ファイルのコミット（`.env.example` は許可）
- SSH 秘密鍵・クラウドトークン（gitleaks で検知）
- ローカル絶対パスのハードコード（環境依存コードの防止。`.md`・`docs/` は allowlist で除外）
- 500 KB を超えるファイル

#### 手動遵守事項
- API キー・パスワードをコードに書かない（Alfred の暗号化キーチェーンを使う）
- 誤ってコミットしたシークレットは、履歴から削除した上で即座にローテーションする
- AI に秘密情報を含むファイルやコードを渡さない
- AI が生成したコードは必ずレビューしてからコミットする
- AI との会話ログをリポジトリにコミットしない

#### コードレビュー
- `main` に到達するコミットは可能な限り他の開発者がレビューする（個人開発の場合は PR を経由してセルフレビューする）
- 認証・認可・暗号化・データアクセスに関わる変更はセキュリティレビューを必須とする

詳細: `CONTRIBUTING.md`（Security セクション）、`docs/dev-charter/SECURITY_POLICY.md`

### UIガイドライン（UI_GUIDELINES）

Alfred Script Filter のレスポンス（JSON items）に適用するルール:

- **Unicode 絵文字禁止** — Alfred の結果アイテムの `title` / `subtitle` に Unicode 絵文字を使わない
  - 代替: ASCII 記号（`>`、`*`、`[x]` など）または何も使わない
- アイコンは `workflow/icon.png` で制御する（PNG ファイル）
- 外観モード（ライト/ダーク）は Alfred が制御するため、ワークフロー側での対応は不要

### マネタイズ（MONETIZATION_POLICY）

OSS プロジェクトのため、以下の方式を採用:

- **Buy Me a Coffee**: https://www.buymeacoffee.com/y.marui
- **GitHub Sponsors**: リポジトリの Sponsors 機能（`.github/FUNDING.yml` 設定済み）

README.md の末尾に Buy Me a Coffee バッジを掲載する。
マネタイズを本格検討する場合は `MONETIZATION.md` を作成し、このファイルに概要を追記する。

### ローカライゼーション（LOCALIZATION_POLICY）

Alfred ワークフローは現時点では UI テキストのローカライゼーション機能を持たない。
将来的に対応する場合の優先言語順:

1. ユーザー設定
2. システム言語設定
3. 英語（デフォルト）

対応候補言語: 日本語 / 英語 / 中国語 / ヒンディー語 / スペイン語 / フランス語 / ポルトガル語

---

## Project-Specific Rules

### アーキテクチャ制約

- `workflow/scripts/entry.py` は Alfred が実行する**唯一のファイル**。ビジネスロジックを書かない
- `src/alfred/` は Alfred SDK ヘルパーのみ — アプリケーションロジックは不可
- Commands → Services → Clients の順に呼ぶ。レイヤーをスキップしない
- すべての `output()` 呼び出しは `alfred.response.output()` を経由する
- `main()` は必ず `safe_run()` でラップする（未捕捉例外 = Alfred が空白表示になる）

### テスト規約

- `src/app/`（commands / services）をテスト対象とする — 純粋 Python
- 外部 I/O は行わない（`passgen_service` は stdlib のみ使用）
- `conftest.py` が Alfred 環境変数を tmp ディレクトリに自動設定する
- Alfred SDK ヘルパーのテストは `tests/test_alfred.py`

詳細な開発フロー・命名規則・コードレビュー手順は `DEVELOPING.md` を参照する。

### Python 開発環境（PYTHON_TOOLCHAIN）

| 役割 | ツール |
|---|---|
| Python バージョン管理 | pyenv |
| パッケージ管理・仮想環境・スクリプト実行 | uv |
| Linter / Formatter | ruff |
| 型チェック | mypy（strict モード） |
| テスト | pytest |

新しい Python プロジェクトを立ち上げる場合、または依存関係を変更する場合はこのツールチェーンに従う。

### Alfred 実行環境（RUNTIME）

Alfred からスクリプトを実行する際は `use_uv` 変数で実行方法を切り替える。

- `use_uv` は Config Builder の checkbox（デフォルト: ON）として定義する
- **ON かつ `uv` が PATH に存在する場合**: `uv run python` で実行
- **OFF または `uv` が存在しない場合**: `python3` で実行

`workflow/scripts/entry.py` を呼ぶスクリプト行のパターン:

```bash
[ "${use_uv:-1}" = "1" ] && command -v uv >/dev/null 2>&1 && exec uv run python scripts/entry.py "$1" || exec python3 scripts/entry.py "$1"
```

### 設定管理（CONFIG_BUILDER）

- **ユーザーが設定する値はすべて Config Builder に入れる** — `workflow/info.plist` の `userconfigurationconfig` 配列に追加する
- Alfred の `variables` キー（environment variable）は使わない。Config Builder で代替できる場合は必ず Config Builder を使う
- Config Builder の値は Alfred がスクリプト実行時に環境変数として自動で渡すため、スクリプト側では `os.environ` で読める
- 新しい設定項目を追加するときは以下の型から選ぶ: `textfield` / `checkbox` / `select` / `file` / `password`

### コードスタイル

- コメントは **「なぜそうするか」のみ** 書く。コードから自明な処理には書かない
- ruff（linter）+ ruff format（formatter）、行長 100
- すべての public 関数に型ヒント必須
- 各モジュール先頭に `from __future__ import annotations`
- mypy strict モード（`pyproject.toml` 参照）

命名規則・コミットメッセージ形式・PR チェックリストは `CONTRIBUTING.md` を参照する。

### パフォーマンス

- Script Filter のレスポンスタイム目標: **100ms 未満**
- ネットワーク呼び出しには `alfred.cache.Cache` を使用する
- キャッシュ TTL デフォルト: 300s（5分）

### 依存管理

- ランタイム依存 → `requirements.txt` → `workflow/vendor/` にベンダリング（`make vendor`）
- 開発依存 → `pyproject.toml [dependency-groups]`（uv 推奨形式）
- ランタイム依存は最小限に保つ（パッケージ追加 = ワークフローサイズ増加）

---

## AI Tool Assignments

| ツール | 担当 |
|---|---|
| Claude Code | プロジェクト立ち上げ・大規模変更・アーキテクチャ設計・リファクタリング |
| GitHub Copilot | バグ修正・細かな実装補助・単体テスト作成 |
| Gemini CLI | プライバシーポリシー作成・更新 / ストア説明文 / 審査用ドキュメント / プロジェクト全体のドキュメント管理 |

### AI並用時のルール
- Claude Code 作業中は Copilot 提案を**参考程度**に（盲目的に受け入れない）
- Copilot の提案がプロジェクト規約に反する場合は無視し、Claude Code でレビュー後採用

---

## Prohibited Actions

- シークレット・認証情報・`.env` ファイルのコミット
- pre-commit フックのスキップ（`--no-verify` 禁止）
- `workflow/scripts/entry.py` へのビジネスロジックの追加
- レイヤーをスキップした呼び出し（例: Command が Client を直接呼ぶ）
- テストでの実際の HTTP 通信
- デバッグ用 `print` 文の本番コードへの残置
- Alfred 結果アイテムへの Unicode 絵文字の使用
- ハードコードされた絶対パス（`$HOME` を使う）
- Config Builder で代替できる設定を Alfred の `variables` キー（environment variable）に直接書くこと
- AI に秘密情報を含むファイルやコードを渡すこと
- AI との会話ログのリポジトリへのコミット

---

## 開発コマンド・リリース手順

開発フロー・コマンド一覧・リリース手順は `CONTRIBUTING.md` を参照する。

---

*このファイルは `docs/dev-charter/` の内容をプロジェクト向けにまとめたものです。
憲章が更新された場合（`git subtree pull` 後）は、このファイルも更新してください。*
