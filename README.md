# magickbot

A bot to allow users to run fun imagemagick commands on images from the fediverse.

## Installation

`GO111MODULE=on go get -u github.com/ibrokemypie/magickbot/cmd/magickbot`

Binary will be installed to your `GOPATH` (likely $HOME/go/bin).

Depends on go with modules support and imagemagick v7+.

## Usage

Run `magickbot`, first time running will prompt for oauth authentication for your bot user and instance.

Tag the bot either in a status containing media, a reply to a status containing media, a status containing mentions of users to apply to their avatars or a reply to a status with no media to apply to the user's avatar. Include the command (eg. explode) in your status, optionally include an argument. The only order that matters is argument must be after command.

`command [argument] [@user...]`

### Commands

help

explode [iterations]

implode [iterations]

magick [scale]

### Sample config

`$HOME/.config/magickbot/config.yaml:`

```
instance:
    instance_url: https://mastodon.social;
    access_token: xxxxxxxxxxxxxxxxxxxxxxx

last_mention_id: xxxxxx
max_pixels: 640000
max_iterations: 15
```

### Todo

- More magick commands

- Command blacklist

- User whitelist

- Local instance only mode
