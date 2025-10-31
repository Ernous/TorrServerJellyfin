#!/bin/bash

echo "Starting Jellyfin + TorrServer..."

# Create necessary directories
mkdir -p /config/jellyfin /config/torrserver /cache/jellyfin /cache/torrserver /media/strm

# Configure TorrServer to use /media/strm for .strm files
cat > /config/torrserver/settings.json <<EOF
{
  "CacheSize": 67108864,
  "ReaderReadAHead": 95,
  "PreloadCache": 50,
  "UseDisk": true,
  "TorrentsSavePath": "/cache/torrserver",
  "RemoveCacheOnDrop": false,
  "JlfnAddr": "/media/strm",
  "JlfnAutoCreate": true,
  "TorrServerHost": "http://localhost:8090"
}
EOF

echo "Configuration created. Starting services with supervisor..."

# Start supervisor to manage both services
exec /usr/bin/supervisord -c /etc/supervisor/conf.d/supervisord.conf
