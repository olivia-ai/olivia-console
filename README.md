<h1 align="center">
  <br>
  <img src="https://olivia-ai.org/img/icons/olivia-with-text.png" alt="Olivia's character" width="300">
  <img src="./olivia-cc.svg">
  <br>
</h1>

<h4 align="center">📟Console client for Olivia</h4>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/olivia-ai/olivia-console"><img src="https://goreportcard.com/badge/github.com/olivia-ai/olivia-console"></a>
  <br>
  <img src="https://github.com/olivia-ai/olivia/workflows/Format%20checker/badge.svg">
</p>

<p align="center">
  <a href="https://olivia-ai.org">Website</a> —
  <a href="https://discord.gg/wXDwTdy">Discord</a> —
  <a href="#getting-started">Getting started</a> —
  <a href="#contributors">Contributors</a> —
  <a href="#license">License</a>
</p>

# Getting started
Console client for [Olivia](https://github.com/olivia-ai/olivia)

## How to use it.
For the first time - you can simple run - ./main
It's will create new default config file, and new logfile.

### Example of config file:
```json
{
 "port": "8080",
 "host": "localhost",
 "debug_level": "error",
 "bot_name": "Olivia",
 "user_token": "52fdfc072182654f163f5f0f9a621d729566c74d10037c4d7bbb0407d1e2c64981855ad8681d0d86d1e91e00167939cb6694"
}
```

### Description:
* `bot_name` - name for your bot (default - Olivia)
* `debug_level` - verbosity (default - error)
* `host` - host where is running server part of olivia (default - localhost)
* `port` - the same about port
* `user_token` - your own token (default - generated on the first run)


# Contributors

<p align="center">
  <img alt="docker installation" height="85" src="https://i.imgur.com/6xr2zdp.png">
</p>

Thanks to @NerdDoc for the creation of this tool.

# Licence
Licensed under MIT