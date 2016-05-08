#!/bin/bash

URL='https://api.telegram.org/bot'$TELEGRAM_API_TOKEN
MSG_URL=$URL'/sendMessage'

VERSION=$1
RELEASE_URL=https://github.com/leominov/self-monitoring/releases/tag/$VERSION

echo "Sending release notice"
RELEASE_URL=$(curl -s http://tinyurl.com/api-create.php?url=$RELEASE_URL)
MSG_TEXT="🚀 Gomon v$VERSION released $RELEASE_URL"
res=$(curl -s "$MSG_URL" -F "chat_id=$TELEGRAM_CHAT_ID" -F "text=$MSG_TEXT" -F "parse_mode=html")
