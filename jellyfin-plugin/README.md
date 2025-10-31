# Jellyfin TorrServer Plugin

Плагин для интеграции TorrServer с Jellyfin, позволяющий управлять торрентами прямо из интерфейса Jellyfin.

## 🎯 Возможности

### Управление торрентами
- ✅ Добавление торрентов через Jellyfin UI
- ✅ Просмотр списка активных торрентов с постерами
- ✅ Удаление торрентов
- ✅ Выбор категории (Movies/TV Shows) при добавлении
- ✅ Кастомная директория для .strm файлов
- ✅ Опция создания .strm файлов при добавлении

### Настройки кэша
- ✅ Выбор типа хранения: RAM (быстро) или Disk (постоянно)
- ✅ Настройка пути для кэша на диске
- ✅ Настройка размера кэша (32-2048 MB)
- ✅ Настройка предзагрузки (Reader Read Ahead, Preload Cache)
- ✅ Удаление кэша при удалении торрента
- ✅ Удаление кэша при остановке видео

### Интеграция
- ✅ Интеграция с TMDB для автоматического поиска постеров
- ✅ Автоматическое создание .strm файлов
- ✅ Полноценный UI прямо в Jellyfin Dashboard
- ✅ REST API для интеграции с другими приложениями

## 📦 Установка

### Вариант 1: Через репозиторий плагинов

1. Откройте Jellyfin Dashboard
2. Перейдите в **Plugins** → **Repositories**
3. Добавьте репозиторий:
   ```
   Name: TorrServer
   URL: https://raw.githubusercontent.com/Ernous/TorrServerJellyfin/master/jellyfin-plugin/manifest.json
   ```
4. Перейдите в **Catalog** и установите **TorrServer**
5. Перезапустите Jellyfin

### Вариант 2: Ручная установка

1. Скомпилируйте плагин:
   ```bash
   cd jellyfin-plugin/Jellyfin.Plugin.TorrServer
   dotnet build -c Release
   ```

2. Скопируйте DLL в папку плагинов Jellyfin:
   ```bash
   # Linux
   cp bin/Release/net8.0/Jellyfin.Plugin.TorrServer.dll /var/lib/jellyfin/plugins/TorrServer/

   # Windows
   copy bin\Release\net8.0\Jellyfin.Plugin.TorrServer.dll "C:\ProgramData\Jellyfin\Server\plugins\TorrServer\"

   # Docker
   docker cp bin/Release/net8.0/Jellyfin.Plugin.TorrServer.dll jellyfin:/config/plugins/TorrServer/
   ```

3. Перезапустите Jellyfin

## ⚙️ Настройка

1. Откройте **Dashboard** → **Plugins** → **TorrServer**
2. Настройте подключение к TorrServer:
   - **TorrServer URL**: `http://localhost:8090` (или ваш URL)
   - **Username/Password**: если настроена авторизация
3. Настройте .strm файлы:
   - **Auto-create .strm files**: включить для автоматического создания
   - **.strm Files Path**: `/media/strm` (путь где будут создаваться файлы)
   - **Default Category**: Movies или TV Shows
4. Настройте кэш:
   - **Cache Size**: размер кэша в MB (32-2048)
   - **Reader Read Ahead**: процент предзагрузки (50-100%)
   - **Preload Cache**: процент для предзагрузки перед воспроизведением
   - **Use Disk Cache**: использовать диск вместо RAM
   - **Remove Cache on Drop**: удалять кэш при удалении торрента
5. TMDB Integration:
   - **TMDB API Key**: ключ API от themoviedb.org

## 🚀 Использование

### Через API

Плагин предоставляет REST API для управления TorrServer:

#### Получить список торрентов
```bash
GET /TorrServer/torrents
```

#### Добавить торрент
```bash
POST /TorrServer/add
Content-Type: application/json

{
  "Link": "magnet:?xt=urn:btih:...",
  "Title": "Movie Title",
  "Poster": "https://image.tmdb.org/t/p/w500/poster.jpg",
  "Category": "Movies",
  "StrmDir": "custom-folder"
}
```

#### Удалить торрент
```bash
POST /TorrServer/remove
Content-Type: application/json

{
  "Hash": "torrent_hash_here"
}
```

#### Получить настройки
```bash
GET /TorrServer/settings
```

#### Обновить настройки
```bash
POST /TorrServer/settings
Content-Type: application/json

{
  "CacheSize": 128,
  "ReaderReadAHead": 95,
  "PreloadCache": 50,
  "UseDisk": true,
  "RemoveCacheOnDrop": false,
  "JlfnAddr": "/media/strm",
  "JlfnAutoCreate": true,
  "TorrServerHost": "http://localhost:8090",
  "TMDBApiKey": "your_api_key"
}
```

#### Поиск в TMDB
```bash
POST /TorrServer/search
Content-Type: application/json

{
  "Query": "Inception",
  "Language": "en",
  "Type": "multi"
}
```

### Через веб-интерфейс Jellyfin

1. Откройте **Dashboard** → **TorrServer** (в меню слева)
2. Нажмите **"Add Torrent"**
3. Вставьте magnet-ссылку или URL торрента
4. (Опционально) Введите название и нажмите **"Search TMDB"** для поиска постера
5. Выберите постер из результатов
6. Выберите категорию (Movies/TV Shows)
7. (Опционально) Укажите кастомную директорию для .strm файлов
8. Нажмите **"Add Torrent"**

Список торрентов обновляется автоматически каждые 5 секунд и показывает:
- Постер
- Название
- Размер
- Скорость загрузки/отдачи
- Количество пиров
- Прогресс загрузки

Для удаления торрента нажмите кнопку **"Remove"** на карточке торрента.

## 🔧 Разработка

### Требования

- .NET 8.0 SDK
- Jellyfin 10.9.0+

### Сборка

```bash
cd jellyfin-plugin/Jellyfin.Plugin.TorrServer
dotnet build
```

### Отладка

1. Установите плагин в Jellyfin
2. Запустите Jellyfin с отладкой:
   ```bash
   JELLYFIN_LOG_LEVEL=Debug jellyfin
   ```
3. Проверьте логи в `/var/log/jellyfin/`

## 📝 TODO

- [ ] Добавить кастомную страницу в Jellyfin UI для управления торрентами
- [ ] Поддержка автоматического обновления библиотеки после добавления торрента
- [ ] Уведомления о завершении загрузки
- [ ] Статистика использования
- [ ] Планировщик задач для автоматической очистки

## 🤝 Вклад

Pull requests приветствуются! Для крупных изменений сначала откройте issue.

## 📄 Лицензия

GPL-3.0

## 🔗 Ссылки

- [TorrServer](https://github.com/YouROK/TorrServer)
- [Jellyfin](https://jellyfin.org/)
- [Jellyfin Plugin Development](https://jellyfin.org/docs/general/server/plugins/)
