

# SyedBot
Discord Bot that does various things

**Main Features**

- Discord command to Tweet
- Anilist anime info
- Time until
- Current time in city
- Stock symbol data 

# Installation & Setup
## 1. 
```
go get https://github.com/Monko2k/SyedBot
make
./syed
```

## 2. Configuration
Create a `config.json` at `config/config.json`
```
{
    "DiscordToken": "",
    "Twitter": {
        "Token": "",
        "TokenSecret": "",
        "Key": "",
        "KeySecret": ""
    },
    "Geocode": "",
    "TimeZoneDB": ""
}
```
API keys can be obtained from:
API | Link
------------ | -------------
DiscordToken | https://discordapp.com/developers/applications
Twitter | https://developer.twitter.com/en/docs/twitter-api
Geocode | https://geocode.xyz/api
TimeZoneDB | https://timezonedb.com/api

Note that leaving empty string is fine as the associated commands will cease to function.


# Commands
## Anime
**anime** `title` 
> Display info for a show

**anirand** `type:value1-value2` `type:value`
> Returns a random anime, optional parameters are available to limit the selection of anime to random with.
The `type` can be either upper/lower case of the character.

Type | Example value(s)
------------ | -------------
y (Year) | `2021` or `1940-2022`
g (Genres) | `Action` `` delimit by any combination of `,` or ` `
s (Score) | `22` or `60-80`

**anistaff** `name` 
> Returns the detail of the person that participated in anime.

**anichar** `name` 
> **anistaff** but for characters in anime.

## Crypto
**crypto** `name` 
> Returns the current value of the crypto in USD (from CoinGecko API).
- `name` is the full name of the crypto, for example instead of symbol `BCH` the user need to type `bitcoin cash`.
- `name` is case insensitive.

## Stock
**stock** `ticker symbol` 
> Returns the Previous Close and Market Price of the requested ticker/stock symbol.

## Time
**time**
> Returns the current UTC time.

**time in** `City Name`
> Returns the current time on `City Name` Time Zone.
- `City Name` or `District Name` alone are also sometimes viable.
- Format: `address`, `town`, `postcode country`

**time until** `time`
> Returns the time until the input `time` where `time` is in UTC 
- Different length of the `time` string will translate as follows:
- `h` `hh` `hmm` `hhmm` where `h` is hour `m` is minute

## Twitter

**tweet** `message` `(optional) Media link`
> Greater than 3 (exclusive) of ✅ and ✅ + 2 > 🖕 is required to tweet the requested message

## Misc

**choose** `option1`, `option2`, `option3` ...
> For decidophobia needs

**github**
> Returns https://github.com/Monko2k/SyedBot

**wholesome**
> Supercalifragilisticexpialidocious wholesome command

**whitecatify**
> holy shit guys

**piss**
> shid

**salam**
> salam