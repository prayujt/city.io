CREATE TABLE Accounts (
    uuid varchar(50) PRIMARY KEY,
    username varchar(50) NOT NULL,
    password varchar(50) NOT NULL,
    UNIQUE(username)
);