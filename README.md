# todoapp-discordbot

[github-task-controller](https://github.com/sikigasa/github-task-controller) と連携する Discord Bot です。  
スラッシュコマンドを使って TODO・タスク・プロジェクトの管理が Discord 上から行えます。

## 機能

### TODO 管理 (`/todo`)

| コマンド | 説明 |
|----------|------|
| `/todo create <title> [description]` | TODO を作成 |
| `/todo list` | TODO 一覧を表示 |
| `/todo get <id>` | TODO の詳細を取得 |
| `/todo complete <id>` | TODO を完了にする |
| `/todo update <id> [title] [description]` | TODO を更新 |
| `/todo delete <id>` | TODO を削除 |

### タスク管理 (`/task`)

| コマンド | 説明 |
|----------|------|
| `/task create <project_id> <title> [description] [status] [priority]` | タスクを作成 |
| `/task list <project_id>` | プロジェクトのタスク一覧を表示 |
| `/task get <id>` | タスクの詳細を取得 |
| `/task update <id> [title] [description] [status] [priority]` | タスクを更新 |
| `/task delete <id>` | タスクを削除 |

### プロジェクト管理 (`/project`)

| コマンド | 説明 |
|----------|------|
| `/project create <title> [description]` | プロジェクトを作成 |
| `/project list` | プロジェクト一覧を表示 |
| `/project get <id>` | プロジェクトの詳細を取得 |
| `/project delete <id>` | プロジェクトを削除 |

## セットアップ

### 前提条件

- Go 1.21 以上
- [github-task-controller](https://github.com/sikigasa/github-task-controller) が起動していること
- Discord Bot Token ([Discord Developer Portal](https://discord.com/developers/applications) で作成)

### 1. Discord Bot の作成

1. [Discord Developer Portal](https://discord.com/developers/applications) にアクセス
2. 「New Application」でアプリを作成
3. 「Bot」セクションでトークンをコピー
4. 「OAuth2 > URL Generator」で以下のスコープを選択:
   - `bot`
   - `applications.commands`
5. Bot Permissions で以下を選択:
   - Send Messages
   - Use Slash Commands
6. 生成された URL からサーバーに招待

### 2. 環境変数の設定

```bash
cp .env.example .env
```

`.env` ファイルを編集:

```env
DISCORD_TOKEN=your_discord_bot_token
API_BASE_URL=http://localhost:8080
AUTH_COOKIE=your_auth_session_cookie
DEFAULT_USER_ID=your_user_id
```

| 変数 | 説明 |
|------|------|
| `DISCORD_TOKEN` | Discord Bot のトークン |
| `API_BASE_URL` | github-task-controller の URL (デフォルト: `http://localhost:8080`) |
| `AUTH_COOKIE` | github-task-controller にログイン後の `auth-session` Cookie 値 |
| `DEFAULT_USER_ID` | プロジェクト操作に使用するユーザー ID |

### 3. ビルド & 実行

```bash
# ビルド
go build -o bin/bot ./cmd/bot

# 実行
./bin/bot
```

または直接実行:

```bash
go run ./cmd/bot
```

## ステータス & 優先度

### タスクステータス

| 値 | 表示 |
|----|------|
| 0 | 📋 To Do |
| 1 | 🔄 In Progress |
| 2 | ✅ Done |

### タスク優先度

| 値 | 表示 |
|----|------|
| 0 | 🟢 Low |
| 1 | 🟡 Medium |
| 2 | 🔴 High |

## ライセンス

MIT
