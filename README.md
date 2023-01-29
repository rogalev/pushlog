### Installation

```
go get github.com/rogalev/pushlog
```

### Usage

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

