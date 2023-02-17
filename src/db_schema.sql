CREATE TABLE Accounts (
    player_id VARCHAR(50) PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    password CHAR(64) NOT NULL,
    balance DOUBLE,
    UNIQUE(username)
);

CREATE TABLE Cities (
    city_id VARCHAR(50) PRIMARY KEY,
    city_name VARCHAR(50),
    city_owner VARCHAR(50),
    population INT DEFAULT 1000,
    population_capacity INT DEFAULT 1000,
    tax_rate DOUBLE DEFAULT 15.0,
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

CREATE TABLE Sessions (
    session_id VARCHAR(50) PRIMARY KEY,
    player_id VARCHAR(50),
    expires_on TIMESTAMP,
    FOREIGN KEY(player_id) REFERENCES Accounts(player_id)
);

CREATE TABLE Building_Info (
    building_type VARCHAR(50),
    building_level INT,
    building_production DOUBLE,
    happiness_change INT,
    population_capacity_change INT,
    build_cost DOUBLE,
    build_time INT,
    PRIMARY KEY(building_type, building_level)
);

CREATE TABLE Buildings (
    building_name VARCHAR(50),
    building_type VARCHAR(50),
    building_level INT,
    city_id VARCHAR(50),
    city_row INT,
    city_column INT,
    FOREIGN KEY(city_id) REFERENCES Cities(city_id),
    FOREIGN KEY(building_type, building_level) REFERENCES Building_Info(building_type, building_level),
    PRIMARY KEY(city_id, city_row, city_column)
);

INSERT INTO Building_Info VALUES
--building type, level,production, hapiness, capacity, cost, time
('City Hall', 1, 0.0, 0, 100, 0.0, 0),
('Apartment', 1, 500.00, 2, 5000, 10000.00, 1),
('Supermarket', 1, 4000.00, 2, 400, 10000.00, 1),
('Hospital', 1, 1000.00, 5, 1000, 10000.00 , 1),
('School', 1, 2000.00, 4, 500, 10000.00, 1);

