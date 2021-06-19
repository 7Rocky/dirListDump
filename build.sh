#!/usr/bin/env bash

package_name=$1
arch_name=${2/-/\/}

if [[ -z "$package_name" ]] || [[ -z "$arch_name" ]]; then
  echo "usage: $0 <package-name> <arch-name>"
  exit 1
fi

arch_split=(${arch_name//\// })

GOOS=${arch_split[0]}
GOARCH=${arch_split[1]}

zip_filename="${package_name}-${GOOS}-${GOARCH}.zip"

if [ $GOOS = 'windows' ]; then
  package_name+='.exe'
fi

env GOOS=$GOOS GOARCH=$GOARCH go build --ldflags='-s -w' -o $package_name *.go
upx --ultra-brute $package_name &>/dev/null

if [ $? -ne 0 ]; then
  echo 'An error has occurred! Aborting the script execution...'
  exit 1
fi

zip $zip_filename $package_name
rm $package_name
