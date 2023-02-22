![logo](logo.jpg)

Example of creating a DCA trading bot to Crypton Exchange. The bot gains a position, gradually averages it, and then creates a Take Profit order. Here is a simple trading cycle that the bot implements.

Inspired by this [article](https://habr.com/ru/company/ruvds/blog/517234/).

The finished compiled version can be found on the :arrow_down: [https://github.com/Sagleft/utopist-trader/releases](releases) page.

![screenshot](screenshot.jpg)

# How it works

How DCA bots are useful:
1. It is a simplification of the way to invest: there is no need to study and use technical indicators, there is no need for complicated analysis.
2. No cost to run the bot. Runs on your computer.
3. No routine activities.
4. Allows you to reduce the average entry point to the market.

The bot can be used as the basis for creating your own trading bots.

## Bot strategy

Trading is divided into equal time intervals. Let's say we chose the strategy - to buy. Then we can choose a trading pair, for example `crp_usdt`. We must have a balance in USDT. The bot will buy some coins in each of the intervals, and then place a Take Profit order to sell them. If the Take Profit order is closed, then the trading cycle is complete and a new cycle begins.

If the price changes, the bot calculates whether to buy more, less or do nothing. The more the price rises, the less the bot buys, the more the price falls, the more the bot will buy.

This strategy is suitable both for buying a coin at the lowest price in time, and for the gradual sale. Or you can just use the bot to make money on volatility.

## Configuring

Just fill `config.json` file.

More about the fields:

* `strategy` - `buy` or `sell`.
* `profitPercent` - what % the bot will place Take Profit orders.
* `tradePair` - trading pair.
* `deposit` - the maximum amount that will be used by the bot in the trading cycle.
* `exchange` - data to connect to the exchange.
* `intervalTimeoutSeconds` - timeout between intervals when the bot will act.
* `intervalDepositMaxPercent` - what is the maximum % of the deposit a bot is allowed to use in one interval.
* `noWait` - whether to execute the bot action immediately at startup.
* `debug` - enabling/disabling detailed logs.

## Useful links

* [Forum thread](https://talk.u.is/viewtopic.php?pid=5267) - an opportunity to ask questions and discuss this bot.
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

just run `bash build.sh`

## TODO

1. recording of statistics on trading cycles;
2. offset support for sending new orders after completing a lap;
3. perhaps write an AI to self-learn from collected statistics.
