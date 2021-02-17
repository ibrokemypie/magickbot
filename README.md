# magickbot

A bot to allow users to run fun imagemagick commands on images from the fediverse.

## Installation

`GO111MODULE=on go get -u github.com/ibrokemypie/magickbot/cmd/magickbot`

Binary will be installed to your `GOPATH` (likely $HOME/go/bin).

Depends on go with modules support and imagemagick v7+.

## Usage

Run `magickbot`, first time running will prompt for oauth authentication for your bot user and instance.

Tag the bot either in a status containing media, a reply to a status containing media, or a reply to a status with no media to apply to the user's avatar. Include the command (eg. explode) in your status, optionally including the desired number of iterations (currently limited to 1 to 15 inclusive)

`command [1-15]`

### Commands

explode

implode

magik

### Sample config

`$HOME/.config/magickbot/config.yaml:`

```
instance:
    instance_url: https://mastodon.social;
    access_token: xxxxxxxxxxxxxxxxxxxxxxx
    visibility: public

last_mention_id: xxxxxx
```

### Todo

- More magick commands

- Command blacklist

- Command whitelist

- User whitelist

- User blacklist

- Local instance only mode

- Command list
