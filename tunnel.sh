#!/bin/bash
# Start cloudflared tunnel and save URL
# Usage: ./tunnel.sh [start|stop|url]

TUNNEL_LOG="/tmp/cloudflared_lingvo.log"
PIDFILE="/tmp/cloudflared_lingvo.pid"

case "${1:-start}" in
  start)
    nohup /tmp/cloudflared tunnel --url http://localhost:8080 --no-autoupdate > "$TUNNEL_LOG" 2>&1 &
    echo $! > "$PIDFILE"
    echo "Tunnel started (PID $!)"
    sleep 6
    URL=$(grep -o 'https://[^ ]*\.trycloudflare\.com' "$TUNNEL_LOG" | head -1)
    echo "URL: $URL"
    echo "$URL" > /tmp/cloudflared_lingvo_url.txt
    ;;
  stop)
    if [ -f "$PIDFILE" ]; then
      kill $(cat "$PIDFILE") 2>/dev/null
      rm -f "$PIDFILE"
    fi
    pkill -f "cloudflared.*:8080" 2>/dev/null || true
    echo "Tunnel stopped"
    ;;
  url)
    if [ -f /tmp/cloudflared_lingvo_url.txt ]; then
      cat /tmp/cloudflared_lingvo_url.txt
    else
      echo "No saved URL. Start tunnel first."
    fi
    ;;
  *)
    echo "Usage: $0 [start|stop|url]"
    ;;
esac
