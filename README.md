# Gapbot

![trens](https://media1.tenor.com/images/9c4e40c1ee1511ef81092e9761f22930/tenor.gif)

A Discord bot for the Cary Academy Discord server

## Summary
This is intended to be a multipurpose Discord bot for the CA Discord server semi based off of [Nineball](https://github.com/morzack/nineball).
As of right now there aren't many planned commands.
If you want to add a suggestion, create an issue or PR, or contact @Valis#7360 or @Patchkat#9990 on Discord.
The bot is built using [discordgo](https://github.com/bwmarrin/discordgo).

## Commands

| Command | Description |
| --- | --- |
| help | get help on the bot |
| ping | ping the bot |
| avatar | display your avatar (may display others later) |

## Configuration
The default config is designed to run on a Raspberry Pi that's been set up for bots.
You will need to configure your system in order to test Gapbot.
The config file should be located at `$HOME/.config/gapbot/config.json`
The bot can be configured using the following fields:

| Field | Value |
| --- | --- |
| bot-prefix | the prefix for the bot to use |
| discord-key | the discord token for the bot |
| source-dir | the path to the folder that contains the source code and resources folder |

## Contributions
Contributions are welcome, but please make it easy for us.
Ideally all PRs would be well documented and tested.
Note that the bot key will not be public, so you need to [configure](#configuration) your environment to test it.


![Nagano Electric Railway](https://upload.wikimedia.org/wikipedia/commons/thumb/9/99/Nagaden_E1_at_Shinano-takehara.png/640px-Nagaden_E1_at_Shinano-takehara.png)
