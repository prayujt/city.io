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


#### Create Account: `POST /login/createAccount`
##### Body:
```
{
    uuid:     string,
    username: string, 
    password: string
}
```


## Game

#### Get City for Account: `GET /city/{uuid}`
##### Return:
```
{
    ...
}
```
