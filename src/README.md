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
"true" || "false", depending on if account creation succeeded
```

#### Login to Account: `POST /login/createSession`
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
    status:         boolean,
    sessionId:      string
}
```

#### Get Session: `GET /session/{session_id}`
##### Return:
```
{
    expired:        boolean,
    ...
}
```

#### Post Logout: `POST /sessions/logout`
##### Body:
```
{
    status:         boolean,
    sessionId:      string
}
```

##### Response:
```
{
    status:         boolean,
    sessionId:      string
}
```


## Game

#### Get City for Account: `GET /city/{uuid}`
##### Return:
```
{
    population:     int,
    ...
}
```
