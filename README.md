![logo](logo.jpg)

Example of creating a DCA trading bot to Crypton Exchange. The bot gains a position, gradually averages it, and then creates a Take Profit order. Here is a simple trading cycle that the bot implements.

inspired by this [article](https://habr.com/ru/company/ruvds/blog/517234/).

## How it works

How DCA bots are useful:
1. It is a simplification of the way to invest: there is no need to study and use technical indicators, there is no need for complicated analysis.
2. No cost to run the bot. Runs on your computer.
3. No routine activities.
4. Allows you to reduce the average entry point to the market.

The bot can be used as the basis for creating your own trading bots.

## Configuring

Just fill `config.json` file.

More about the fields:

* ...
* ...
* ...

## Useful links

* [UDocs](https://udocs.gitbook.io/utopia-api/) - collection of all documentation about Utopia API.
* [Crypton Exchange](https://crp.is) - the exchange we work with.
* [CRP.IS API](https://crp.is/api-doc/) - API docs.

### Build from sources

just run

```bash
git clone https://github.com/Sagleft/utopist-trader
cd utopist-trader
go build
```

& run after `config.json` setup:

```bash
./bot
```

### How to build a bot for multiple platforms?

Create dir `build` and use the following script:

```bash
#!/usr/bin/env bash
platforms=("linux/386" "linux/amd64" "windows/386" "windows/amd64")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    output_name='build/'bot'_'$GOOS'-'$GOARCH

    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    echo "build for "$GOOS-$GOARCH".."
    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi

done
```

save it as `make.sh`, then run it:

```bash
bash make.sh
```
