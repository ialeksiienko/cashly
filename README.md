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

Or you can use `.env` for critical vars:

```env
BOT_TOKEN=TELEGRAM_BOT_TOKEN
BOT_PASSWORD=123456

MONO_ENCRYPT_KEY=testtesttesttesttesttesttesttest

DB_PASS=pass
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

### Standalone Docker

```bash
docker build -t cashly-bot .

docker run -d \
  --name cashly \
  --env-file .env \
  -e CONFIG_PATH=config/config.yml \
  cashly-bot
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
│   └── bot/
│       └── bot.go              # Entry point
├── configs/                    # Config files
│   └── config.example.yml
├── internal/
│   ├── adapter/
│   │   ├── database/           # DB connection
│   │       ├── client.go
│   │       └── database.go
│   │   └── repository/         # DB repositories
│   │       ├── family/
│   │       ├── token/
│   │       └── user/
│   ├── app/                    # Bot setup
│   ├── config/                 # Load configs
│   ├── entity/                 # All entities
│   ├── handlers/               # All handlers
│   ├── middleware/             # All middlewares
│   ├── migration/              # DB migrations (goose)
│   ├── pkg/                    # App packages (internal)
│   ├── router/                 # Main router
│   ├── service/                # Business logic
│   │   ├── family/
│   │   │   └── mocks/
│   │   ├── token/
│   │   │   └── mocks/
│   │   └── user/
│   │        └── mocks/
│   ├── state/                  # State management
│   ├── usecase/                # Use cases
│   └── validate/               # Validate
├── pkg/                        # Custom packages
├── test/                       # All (services) tests
├── .env.example
├── .gitignore
├── .mockery.yml
├── docker-compose.yml
├── Dockerfile
├── family.example.json
├── go.mod
├── Makefile
└── README.md
```

## 🔒 Security

- ✅ Monobank tokens are encrypted before storage
- ✅ Password authentication with automatic timeout
- ✅ User whitelist via `family.json`
- ✅ Invite codes with limited validity (48 hours)
- ✅ Sensitive data excluded from logs

## ⚠️ **IMPORTANT**:

- Never commit `config/config.yml`, `family.json`, or `.env` to git!
- Always change the default `bot_password` before deploying!
- Use strong, unique passwords for production

## 🛠️ Development

### Makefile Commands

```bash
make run            # Run application
make goose-path     # Set migration dir
make goose-up       # Run migrations up
make goose-down     # Run migrations down
make docker-up      # Start with docker-compose
```

### Migration Structure

```bash
internal/migration/
├── 00001_users_table.sql
├── 00002_add_families_table.sql
└── ...
```

### Adding New Migration

```bash
goose -dir internal/migration create your_migration_name sql
```

## 🤝 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 👤 Author

**Illia Aleksiienko**

- GitHub: [@ialeksiienko](https://github.com/ialeksiienko)

## 🙏 Acknowledgments

- [Monobank](https://www.monobank.ua/) for the open API
- All contributors

---

⭐ If this project was helpful - give it a star on GitHub!
