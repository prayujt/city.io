CREATE TABLE Accounts (
    player_id varchar(50) PRIMARY KEY,
    username varchar(50) NOT NULL,
    password varchar(50) NOT NULL,
    UNIQUE(username)
);

CREATE TABLE Cities (
    city_id varchar(50) PRIMARY KEY,
    city_name varchar(50),
    city_owner varchar(50),
    population int DEFAULT 0,
    FOREIGN KEY(city_owner) REFERENCES Accounts(player_id)
);

CREATE TRIGGER Create_City
AFTER INSERT ON Accounts
FOR EACH ROW
INSERT INTO Cities (city_id, city_owner) SELECT uuid(), NEW.player_id;

CREATE TRIGGER City_Name
BEFORE INSERT ON Cities
FOR EACH ROW
SET NEW.city_name=CONCAT('City ', NEW.city_id);

CREATE TABLE Buildings (
    building_name varchar(50) PRIMARY KEY,
    building_value int,
    city_id varchar(50),
    FOREIGN KEY(city_id) REFERENCES Cities(city_id)
);

CREATE TABLE Sessions (
    session_id varchar(50) PRIMARY KEY,
    player_id varchar(50),
    expires_on timestamp,
    FOREIGN KEY(player_id) REFERENCES Accounts(player_id)
);
