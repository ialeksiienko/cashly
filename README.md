# 👨‍👩‍👧‍👦 Cashly Bot

A Telegram bot for managing family finances with Monobank (Ukrainian Bank) API integration. Create family groups, view card balances of all family members in one place, and control access through an invitation system.

## ✨ Features

- 👨‍👩‍👧‍👦 **Family Groups** - create and manage family unions
- 💳 **Balance Overview** - view card balances of all family members
- 🔐 **Invitation System** - secure invite codes with expiration time
- 👥 **Member Management** - add and remove participants
- 🏦 **Monobank Integration** - connect personal API tokens
- 🔒 **Authentication** - password-protected access with session timeout
- 👑 **Administration** - extended capabilities for family owners

## 🚀 Quick Start

### Requirements

- Go 1.21 or higher
- PostgreSQL 14+
- Telegram Bot Token (from [@BotFather](https://t.me/botfather))
- Monobank API token (for each user, get at [monobank.ua](https://api.monobank.ua/))

### Installation

1. **Clone the repository**

```bash
git clone https://github.com/ialeksiienko/cashly.git
cd cashly
```

2. **Install dependencies**

```bash
go mod download
```

3. **Configure the application**

You can use either **config.yml** or **environment variables**.

#### Using config.yml for local

```bash
cp config/config.example.yml config/config.yml
```

Edit `config/config.yml`:

```yaml
env: prod
bot:
  token: TELEGRAM_BOT_TOKEN
  long_poller: 10
  password: 123456
mono:
  encrypt_key: test
  api_url: https://api.monobank.ua/
db:
  user: admin
  pass: admin
  host: localhost
  port: 5432
  name: dbname
```

4. **Run migrations**

```bash
# install goose if not already installed
go install github.com/pressly/goose/v3/cmd/goose@latest

# run migrations
goose -dir internal/migration postgres "postgresql://admin:admin@localhost:5432/dbname?sslmode=disable" up
```

5. **Configure allowed users**

```bash
cp family.example.json family.json
```

Edit `family.json` and add Telegram IDs of allowed users:

```json
[
	{
		"firstname": "John",
		"id": 123456789
	}
]
```

> 💡 Find your Telegram ID via [@userinfobot](https://t.me/userinfobot)

6. **Run the bot**

```bash
go run cmd/main.go --config=config/config.yml
```

## 🐳 Docker Deploy

### Docker Compose (recommended)

1. **Create `.env` file**

```bash
cp .env.example .env
```

```env
# Telegram Configuration
BOT_TOKEN=token
BOT_PASSWORD=pass

# Database Configuration
DB_PASS=pass

# Monobank Configuration
MONO_ENCRYPT_KEY=key
```

2. **Start the application**

```bash
docker-compose up -d
```

3. **View logs**

```bash
docker-compose logs -f bot
```

## 📖 Usage

### First Run

1. Find your bot in Telegram
2. Send `/start`
3. Enter the password (from your config or ENV `BOT_PASSWORD`)
4. Choose an action:
   - **Create Family** - if you're the first one
   - **Join Family** - if you have an invite code
   - **Enter My Family** - if you're already in a family

### Basic Commands

- `/start` - main menu
- Password prompt appears automatically after session timeout

### Family Menu

- 💰 **View Balance** - total balance of all family cards
- 👥 **View Members** - list of members
- 🔑 **Add Monobank Token** - connect your cards
- 🗑️ **Remove Token** - disconnect cards
- 🚪 **Leave Family** - exit the group

### Admin Functions (family owner)

- 🎟️ **Create New Code** - generate invite code
- 🗑️ **Delete Member** - remove member from family
- ❌ **Delete Family** - complete group deletion

## 🏗️ Project Architecture

```
cashly/
├── cmd/
│   └── main.go                 # Entry point
├── config/
│   └── config.example.yml.     # Config yaml file
├── internal/
│   ├── adapter/
│   │   └── database/           # DB repositories
│   │       ├── familyrepo/
│   │       ├── tokenrepo/
│   │       └── userrepo/
│   ├── app/                    # Bot and database setup
│   ├── config/                 # Configuration
│   ├── delivery/
│   │   └── telegram/           # Telegram handlers
│   │       ├── handler/
│   │       └── router.go
│   ├── entity/                 # Domain models
│   ├── errorsx/                # Custom errors
│   ├── middleware/             # Middleware
│   ├── migration/              # DB migrations (goose)
│   ├── pkg/                    # Custom logger
│   ├── service/                # Business logic
│   │   ├── familyservice/
│   │   │   └── mocks/
│   │   ├── tokenservice/
│   │   └── userservice/
│   ├── session/                # State management
│   ├── usecase/                # Use cases
│   └── validate/               # Validate
├── .env.example
├── .mockery.yml
├── docker-compose.yml
├── Dockerfile
├── family.example.json
├── go.mod
└── Makefile
```
