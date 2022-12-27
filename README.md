# Asynchronously checking weather 

## Used
- Mongo
- Redis
- Postgres 
- JWT

## Todo
- gRPC
- Dockerfile 
- Swagger
- Migration in docker-compose

## Start work
Fill config.yaml with your OpenWeatherAPI token.

## Register 

Send POST request: http://ADDR:PORT/register 

```
{
    "user_login":"hello",
    "user_password":"world"
}
```

## Login 

Send POST request: http://ADDR:PORT/login

```
{
    "user_login":"hello",
    "user_password":"world"
}
```

## Get weather forecast 

Send POST request http://ADDR:PORT/api/forecast

```
{
    "city":"tokyo"
}
```

```
{
    "weather": [
        {
            "description": "переменная облачность"
        }
    ],
    "main": {
        "temp": 26.01,
        "feels_lie": 0,
        "temp_min": 23.38,
        "temp_max": 26.86
    },
    "sys": {
        "country": "JP",
        "sunrise": 1660939409,
        "sunset": 1660987550
    },
    "name": "Япония"
}
```
Also supported multiple requests 

e.g

```
{
    "city":"москва, калининград"
}
```

```
{
    "weather": [
        {
            "description": "пасмурно"
        }
    ],
    "main": {
        "temp": 24.88,
        "feels_lie": 0,
        "temp_min": 24.12,
        "temp_max": 25.4
    },
    "sys": {
        "country": "RU",
        "sunrise": 1660874928,
        "sunset": 1660928277
    },
    "name": "Замоскворечье"
}
{
    "weather": [
        {
            "description": "облачно с прояснениями"
        }
    ],
    "main": {
        "temp": 29.95,
        "feels_lie": 0,
        "temp_min": 29.95,
        "temp_max": 29.95
    },
    "sys": {
        "country": "RU",
        "sunrise": 1660879235,
        "sunset": 1660932181
    },
    "name": "Калининград"
}

```

