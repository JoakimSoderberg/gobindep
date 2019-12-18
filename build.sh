#!/bin/bash

package_name=gobindep
platforms=("windows/amd64" "linux/amd64" "linux/arm64" "darwin/amd64")

mkdir -p build/

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=${package_name}'-'${GOOS}'-'${GOARCH}
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi  

    GOOS=${GOOS} GOARCH=${GOARCH} go build -v -o build/${output_name} ${package}
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done
