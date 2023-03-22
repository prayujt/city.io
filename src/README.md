# REST API

## Login

#### Create Account: `POST /login/createAccount`
##### Body:
```
{
    username:       string,
    password:       string
}
```

##### Response:
```
{
    status:         boolean
}
```

</br>

#### Create Session: `POST /login/createSession`
##### Body:
```
{
    username:       string,
    password:       string
}
```

##### Response:
```
{
    sessionId:      string
}
```

</br>

#### Get Session: `GET /session/`
##### Response:

```
{
    status:         boolean
}
```

</br>


## Game

### Cities

#### Get Available Buildings: `GET /cities/buildings`

##### Response:
```
[
    {
        buildingType:           string,
        buildCost:              double,
        buildTime:              int,
        buildingProduction:     double,
        happinessChange:        int
    },
    ...
]
```

</br>

#### Get City for Account: `GET /cities`

##### Query Parameters (optional):

```
-   cityName:               string
```

##### Response:
```
{
    username:               string,
    balance:                double,
    population:             int,
    populationCapacity:     int,
    armySize:               int, (-1 if you are visiting a territory that you do not own)
    cityName:               string
}
```

</br>

#### Get Territory for Account: `GET /cities/territory`

##### Response:
```
[
    {
        cityName:               string,
        cityProduction:         double,
        armySize:               int,
    },
    ...
]
```

</br>

#### Get Buildings for City: `GET /cities/buildings`
##### Query Parameters (optional):

```
-   cityName:       string
```

##### Response:
```
{
    isOwner:        boolean,
    buildings:      [
                        {
                            buildingType:   string,
                            buildingLevel:  int,
                            cityRow:        int,
                            cityColumn:     int
                        },
                        ...
                    ]
}
```

</br>

#### Get Building for City: `GET /cities/buildings/{city_row}/{city_column}`
##### Query Parameters (optional):

```
-   cityName:           string
```

##### Response:
```

{
    buildingType:       string,
    buildingLevel:      int,
    buildingProduction: double,
    happinessChange:    int,
    startTime:          string,
    endTime:            string
}
```

</br>


#### Create Building in a City: `POST /cities/createBuilding`

##### Query Parameters (optional):

```
-   cityName:       string
```

##### Body:
```
{
    buildingType:   string,
    cityRow:        int,
    cityColumn:     int
}
```

##### Response:
```
{
    status:         boolean
}
```

</br>

#### Upgrade Building in a City: `POST /cities/upgradeBuilding`

##### Query Parameters (optional):

```
-   cityName:       string
```

##### Body:
```
{
    cityRow:        int,
    cityColumn:     int
}
```

##### Response:
```
{
    status:         boolean
}
```

</br>

### Armies

#### Train Troops: `POST /armies/train`

##### Body: 

```
{
    sessionId:          string,
    troopCount:         int
}
```

##### Response:
```
{
    status:             boolean
}
```

</br>

#### Move Troops: `POST /armies/move`

##### Body: 

```
{
    sessionId:          string,
    armySize:           int,
    fromCity:           string,
    toCity:             string
}
```

##### Response:
```
{
    status:             boolean
}
```

</br>


### Visit

#### Get City List: `GET /cities`
##### Response:
```
[
    {
        cityName:       string,
        cityOwner:      string
    },
    ...
]
```

</br>

#### Get Town List: `GET /towns`
##### Response:
```
[
    {
        cityName:       string,
        cityOwner:      string
    },
    ...
]

```

</br>

#### Get Leaderboard: `GET /leaderboard`
##### Response:
```
[
    {
        username:       string,
        balance:        double
    },
    ...
]
```


