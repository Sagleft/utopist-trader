#!/usr/bin/env bash
platforms=("linux/386" "linux/amd64" "windows/386" "windows/amd64")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    output_name='bot_'$GOOS'-'$GOARCH
    arch_name=$output_name'.zip'

    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    echo "build for "$GOOS-$GOARCH".."
    env GOOS=$GOOS GOARCH=$GOARCH go build -o 'build/'$output_name
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi

    cd build
    zip -1 $arch_name $output_name 'config.json'
    cd ..

done
