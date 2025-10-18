#!/bin/bash

BIN_NAME="ug"
SRC_PATH="./$BIN_NAME"
DEST_PATH="/usr/local/bin/$BIN_NAME"

echo "🔨 Go 빌드 중: $BIN_NAME"
go build -o "$SRC_PATH"
if [ $? -ne 0 ]; then
    echo "❌ Go 빌드 실패"
    exit 1
fi

echo "📂 /usr/local/bin으로 복사: $DEST_PATH"
sudo cp "$SRC_PATH" "$DEST_PATH"

echo "⚡ 실행 권한 설정: $DEST_PATH"
sudo chmod +x "$DEST_PATH"

echo "✅ 완료: $BIN_NAME 설치 완료! 터미널에서 '$BIN_NAME'로 실행 가능합니다."
echo "✅ 'ug --help' 실"
