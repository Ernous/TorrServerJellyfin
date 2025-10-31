# Как создать релиз TorrServerJellyfin

## 🚀 Автоматический релиз через GitHub Actions

### Шаг 1: Создайте тег

```bash
# Убедитесь что все изменения закоммичены
git add .
git commit -m "Prepare release v1.0.0"

# Создайте тег с версией
git tag v1.0.0

# Запушьте тег на GitHub
git push origin v1.0.0
```

### Шаг 2: GitHub Actions автоматически:

1. ✅ Соберёт фронтенд (React)
2. ✅ Сгенерирует веб-ассеты
3. ✅ Скомпилирует TorrServer для всех платформ:
   - Linux AMD64
   - Linux ARM64
   - Windows AMD64
   - macOS AMD64
   - macOS ARM64
4. ✅ Создаст GitHub Release
5. ✅ Загрузит все бинарники в релиз

### Шаг 3: Проверьте релиз

Перейдите на: `https://github.com/Ernous/TorrServerJellyfin/releases`

Вы увидите новый релиз с бинарниками!

## 📦 Использование релизов в Dockerfile

После создания релиза, Dockerfile автоматически будет скачивать бинарники:

```dockerfile
# Скачает последний релиз
ARG TORRSERVER_VERSION=latest

# Или конкретную версию
ARG TORRSERVER_VERSION=v1.0.0
```

## 🔄 Обновление версии

### Для патч-релиза (1.0.0 → 1.0.1):
```bash
git tag v1.0.1
git push origin v1.0.1
```

### Для минор-релиза (1.0.0 → 1.1.0):
```bash
git tag v1.1.0
git push origin v1.1.0
```

### Для мажор-релиза (1.0.0 → 2.0.0):
```bash
git tag v2.0.0
git push origin v2.0.0
```

## 🛠️ Ручной релиз (если нужно)

Если GitHub Actions не работает, можно создать релиз вручную:

### 1. Соберите бинарники локально:

```bash
# Linux AMD64
cd server
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o TorrServer-linux-amd64 ./cmd

# Linux ARM64
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-w -s" -o TorrServer-linux-arm64 ./cmd

# Windows AMD64
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o TorrServer-windows-amd64.exe ./cmd

# macOS AMD64
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o TorrServer-darwin-amd64 ./cmd

# macOS ARM64 (M1/M2)
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-w -s" -o TorrServer-darwin-arm64 ./cmd
```

### 2. Создайте релиз на GitHub:

1. Перейдите: `https://github.com/Ernous/TorrServerJellyfin/releases/new`
2. Выберите тег или создайте новый
3. Заполните описание релиза
4. Загрузите все бинарники
5. Нажмите "Publish release"

## 📝 Changelog

Рекомендуется вести CHANGELOG.md:

```markdown
# Changelog

## [1.0.0] - 2025-10-31

### Added
- Jellyfin integration with automatic .strm file creation
- Custom directory selection for .strm files
- Automatic cleanup on torrent removal
- TMDB metadata integration

### Fixed
- Race condition in .strm file creation
- Path persistence in database

### Changed
- Updated to Go 1.23
- Improved error handling
```

## 🔍 Проверка релиза

После создания релиза, проверьте что бинарники работают:

```bash
# Скачайте бинарник
wget https://github.com/Ernous/TorrServerJellyfin/releases/latest/download/TorrServer-linux-amd64

# Сделайте исполняемым
chmod +x TorrServer-linux-amd64

# Запустите
./TorrServer-linux-amd64 --help
```

## 🐳 Тестирование Docker образа

После релиза, пересоберите Docker образ:

```bash
docker build -f Dockerfile.jellyfin -t torrserver-jellyfin:test .
docker run -p 8080:8080 torrserver-jellyfin:test
```

## 📊 Версионирование

Используйте [Semantic Versioning](https://semver.org/):

- **MAJOR** (1.0.0 → 2.0.0): Breaking changes
- **MINOR** (1.0.0 → 1.1.0): New features (backward compatible)
- **PATCH** (1.0.0 → 1.0.1): Bug fixes

## 🎯 Best Practices

1. ✅ Всегда тестируйте перед релизом
2. ✅ Пишите понятные release notes
3. ✅ Используйте семантическое версионирование
4. ✅ Создавайте changelog
5. ✅ Тегируйте коммиты перед релизом
6. ✅ Проверяйте что все бинарники работают

## 🔗 Полезные ссылки

- [GitHub Releases Documentation](https://docs.github.com/en/repositories/releasing-projects-on-github)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Semantic Versioning](https://semver.org/)
