# magickbot

A bot to allow users to run fun imagemagick commands on images from the fediverse.

## Installation

`go get -u github.com/ibrokemypie/magickbot/cmd/magickbot`

Binary will be installed to your `GOPATH` (likely $HOME/go/bin).

## Usage

Run `magickbot`, first time running will prompt for oauth authentication for your bot user and instance.

### Commands

None

### Sample config

`$HOME/.config/magickbot/config.yaml:`

```
instance:
    instance_url: https://mastodon.social;
    access_token: xxxxxxxxxxxxxxxxxxxxxxx
    visibility: public
```

### Todo

- Media download

- Media upload

- Magick commands

- Command blacklist

- Command whitelist

- User whitelist

- User blacklist

- Local instance only mode
