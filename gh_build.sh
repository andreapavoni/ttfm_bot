#!/bin/bash

set -e

rm -rf build

LATEST_TAG=$(git describe --tags)

gox -output 'build/{{.OS}}_{{.Arch}}/ttfm_bot' \
    --os "linux darwin windows" \
    -arch "386 amd64" \
    -osarch '!darwin/386' \
    -ldflags="-X 'main.Version=${LATEST_TAG}'"

# build for raspberry pi 3+
GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-X 'main.Version=${LATEST_TAG}'" -o build/linux_arm7/ttfm_bot main.go
tar czf build/ttfm_bot-linux_arm7.tgz -C build/linux_arm7/ ttfm_bot; \

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