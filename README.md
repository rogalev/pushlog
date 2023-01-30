Pushlog is a simple server for receiving and sending logs. The http server is used as the recipient of the messages. The data is stored in RAM. A telegram bot is used to send logs. However, it is enough to implement the "receiver", "storage" and "publisher" interfaces and you can receive messages, for example, from a database or queue server, store them in the file system and send them as SMS messages.

Installation
---
```
> git clone https://github.com/rogalev/pushlog.git 
> cd pushlog
> make build
```
The compiled file will be located in the pushlog/.bin directory

**Before executing the "make build" command, please make sure that golang is installed on your working machine**


Usage
---
#### Run application

```
> ./pushlog --config="..."
```

| Flag   | Default value | Description         |
|:-------|:-------------:|:--------------------|
| config |  config.json  | Path to config file |


#### Generate config file


```
> ./pushlog genconfig --config="..."
```

| Flag      | Default value | Description                                                                                           |
|:----------|:-------------:|:------------------------------------------------------------------------------------------------------|
| config    |  config.json  | Path to output config file. At the moment, only json notation is supported for the configuration file |


#### Watch telegram updates

```
> ./pushlog tgupdates --token="..." --update=0
```

| Flag   | Default value | Description            |
|:-------|:-------------:|:-----------------------|
| token  |               | Telegram bot token     |
| update |       0       | Telegram update offset |


#### Send log to pushlog server
```
> curl -X POST localhost:8000/push -d "key=MESSAGE_KEY&expiration=EXIRATION_TIME&body=LOG_MESSAGE" 
```

| Value      | Default value | Description                                                                                                     |
|:-----------|:-------------:|:----------------------------------------------------------------------------------------------------------------|
| key        |               | The key for determining the uniqueness of the message log for message caching                                   |
| expiration |       0       | The time (in seconds) of finding the message key in the cache to avoid sending the same messages multiple times |
| body       |               | Message text                                                                                                    |
