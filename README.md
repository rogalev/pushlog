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

