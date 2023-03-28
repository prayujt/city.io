# REST API

There are protected and unprotected endpoints in the API. To make a request to a protected endpoint, you need to be authorized with a token. This is done by passing in the following header to the HTTP request: </br>
`Token: <JWT TOKEN>`
</br>
</br>

Endpoints that require authorization are marked with an `*`
</br>
GET endpoints that do not require this additional header, but will only give read access without authorization are marked with `**`

<details>

<summary>
Login
</summary>

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

#### Get Session: `GET /session` *
##### Response:

```
{
    status:         boolean
}
```

</br>

</details>

<details>

<summary>
Game
</summary>

<details>
<summary>
Cities
</summary>

#### Get Available Buildings: `GET /cities/buildings/available`

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

#### Get City for Account: `GET /cities` **

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

#### Get Territory for Account: `GET /cities/territory` *

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

#### Get Buildings for City: `GET /cities/buildings` **
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

#### Get Building for City: `GET /cities/buildings/{city_row}/{city_column}` **
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
    endTime:            string,
    upgradeCost:        double,
    upgradedProduction: double,
    upgradeTime:        int,
    upgradedHappiness:  int
}
```

</br>


#### Create Building in a City: `POST /cities/createBuilding` *

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

#### Upgrade Building in a City: `POST /cities/upgradeBuilding` *

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

</details>

<details>
<summary>
Armies
</summary>

#### Train Troops: `POST /armies/train` *

##### Body: 

```
{
    troopCount:         int
    cityName:           string (optional, if not given defaults to home city)
}
```

##### Response:
```
{
    status:             boolean
}
```

</br>

#### Move Troops: `POST /armies/move` *

##### Body: 

```
{
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

#### Get Marches: `GET /armies/marches` *

##### Response:
```
{
    fromCityName:       string,
    fromCityOwner:      string,
    toCityName:         string,
    toCityOwner:        string,
    armySize:           int,
    returning:          bool,
    attack:             bool,
    startTime:          string,
    endTime:            string
}
```

</br>

#### Get Single City Training: `GET /armies/training` *

##### Query Parameters (optional):

```
-   cityName:       string
```

##### Response:
```
{
    armySize:           int,
    startTime:          string,
    endTime:            string
}
```

</br>

#### Get Global Training: `GET /armies/training/global` *

##### Response:
```
[
    {
        cityName:           string,
        armySize:           int,
        startTime:          string,
        endTime:            string
    }
]
```

</br>


</details>

<details>
<summary>
Visit
</summary>

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
</details>

