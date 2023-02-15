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
`"true" || "false"`, depending on if account creation succeeded

</br>

#### Get Session: `GET /session/{session_id}`
##### Response:
`"true" || "false"`, depending on if session is valid and has not expired

</br>

#### Post Logout: `POST /sessions/logout`
##### Body:
```
{
    sessionId:      string
}
```

##### Response:
`"true" || "false"`, depending on if session logout succeeded

</br>

## Game

#### Get City for Account: `GET /city/{uuid}`
##### Return:
```
{
    population:     int,
    ...
}
```

