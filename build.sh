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

output_name="${package_name}-${GOOS}-${GOARCH}"
zip_filename="${output_name}.zip"

if [ $GOOS = 'windows' ]; then
  output_name+='.exe'
fi

echo "Building $output_name..."
env GOOS=$GOOS GOARCH=$GOARCH go build --ldflags='-s -w' -o $output_name *.go
upx --ultra-brute $output_name &>/dev/null

if [ $? -ne 0 ]; then
  echo 'An error has occurred! Aborting the script execution...'
  exit 1
fi

zip $zip_filename $output_name
rm $output_name

mkdir build 2>/dev/null
mv $zip_filename build
