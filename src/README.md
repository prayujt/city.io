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
    status:         boolean
}
```

</br>

#### Get Session: `GET /session/{session_id}`
##### Response:

```
{
    status:         boolean
}
```

</br>

#### Post Logout: `POST /sessions/logout`
##### Body:
```
{
    sessionId:      string
}
```

##### Response:
```
{
    status:         boolean
}
```

</br>

## Game

#### Get City for Account: `GET /cities/{session_id}`
##### Response:
```
{
    cityId:        string,
    population:    int,
    cityName:      string
}
```

</br>

#### Get Buildings for City: `GET /cities/{session_id}/buildings`
##### Query Parameters (optional):

```
-   username:       string
```

##### Response:
```
{
    isOwner:        boolean,
    buildings:      [{
                        buildingType:   string,
                        buildingLevel:  string,
                        buildingName:   string,
                        cityRow:        int,
                        cityColumn:     int
                    },
                    ...
                    ]
}
```

</br>

#### Get Building for City: `GET /cities/{session_id}/buildings/{city_row}/{city_column}`
##### Query Parameters (optional):

```
-   username:       string
```

##### Response:
```

		query = fmt.Sprintf("SELECT building_name, building_type, building_level, Buildings.city_id, building_production, happiness_change, start_time, end_time FROM Buildings NATURAL JOIN Building_Info LEFT JOIN Builds ON Buildings.city_id=Builds.city_id AND Buildings.city_row=Builds.city_row AND Buildings.city_column=Builds.city_column WHERE Buildings.city_id=(SELECT city_id FROM Accounts JOIN Cities ON player_id=city_owner WHERE username='%s') AND Buildings.city_row=%d AND Buildings.city_column=%d;", param[0], cityRow, cityColumn)
{
    buildingName:       string,
    buildingType:       string,
    buildingLevel:      int,
    buildingProduction: double,
    happinessChange:    int,
    startTime:          string,
    endTime:            string
}
```

</br>


#### Create Building in a City: `POST /cities/{session_id}/createBuilding`
##### Body:
```
{
    buildingName:   string,
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

#### Upgrade Building in a City: `POST /cities/{session_id}/upgradeBuilding`
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
