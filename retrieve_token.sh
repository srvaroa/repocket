#!/bin/bash
if [ $# -lt 1 ]; then
    echo "usage: $0 <consumer_key>"
    exit 1
fi

TOKEN=$(curl -X POST -d "consumer_key=$1&redirect_uri=localhost" https://getpocket.com/v3/oauth/request 2>/dev/null | cut -f2 -d=)

echo "* Go to this URL and authorize: https://getpocket.com/auth/authorize?request_token=$TOKEN&redirect_uri=localhost"
echo "* You may be redirected to an invalid page, no worries"

read -p "Press any key when done" x

RESP=$(curl -X POST -d "consumer_key=$1&code=$TOKEN" https://getpocket.com/v3/oauth/authorize 2>/dev/null)

echo "The access token is: " $(echo $RESP | cut -f1 -d'&' | cut -f2 -d=)
echo "The user is: " $(echo $RESP | cut -f2 -d'&' | cut -f2 -d=)
