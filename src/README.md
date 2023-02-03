# REST API

## Login

#### Login to Account: `POST /login/verifyAccount`
##### Body:
```
{
    username: string, 
    password: string
}
```

##### Response:
```
"true" || "false", depending on if account credentials are correct
```


#### Create Account: `POST /login/createAccount`
##### Body:
```
{
    username: string, 
    password: string
}
```

##### Response:
```
"true" || "false", depending on if account creation succeeded
```


## Game

#### Get City for Account: `GET /city/{username}`
##### Return:
```
{
    ...
}
```
