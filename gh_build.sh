#!/bin/bash

set -e

rm -rf build

LATEST_TAG=$(git describe --tags)

gox -output 'build/{{.OS}}_{{.Arch}}/ttfm_bot' \
    --os "linux darwin windows" \
    -arch "386 amd64 arm" \
    -osarch '!darwin/386 !windows/arm !darwin/arm' \
    -ldflags="-X 'main.Version=${LATEST_TAG}'"

for rls in build/{linux,darwin}*; do \
    tar czf build/ttfm_bot-$(echo ${rls} | cut -f2 -d/).tgz -C ${rls} ttfm_bot; \
done

for rls in build/windows*; do \
    mv ${rls}/ttfm_bot.exe build/ttfm_bot-$(echo ${rls} | cut -f2 -d/).exe
done

for rls in build/{linux,darwin,windows}*; do \
    rm -rf ${rls}
done

ls -l build/

# GOOS=linux GOARCH=arm GOARM=7 go build -o ${BINARY_NAME}_rpi main.go