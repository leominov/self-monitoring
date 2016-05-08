#!/bin/bash

URL='https://api.telegram.org/bot'$TELEGRAM_API_TOKEN
CHAT_ID=$TELEGRAM_CHAT_ID
MSG_URL=$URL'/sendMessage'

VERSION=$1
RELEASE_URL=https://github.com/leominov/self-monitoring/releases/tag/$VERSION
MSG_TEXT="🚀 Gomon v$VERSION released
$RELEASE_URL"

echo "Sending release notice"
res=$(curl -s "$MSG_URL" -F "chat_id=$CHAT_ID" -F "text=$MSG_TEXT" -F "parse_mode=html")
