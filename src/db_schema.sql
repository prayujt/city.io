CREATE TABLE Accounts (
    player_id VARCHAR(50) PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    password CHAR(64) NOT NULL,
    balance DOUBLE DEFAULT 2000000.0,
    CHECK (balance > 0),
    UNIQUE(username)
);

CREATE TABLE Cities (
    city_id VARCHAR(50) PRIMARY KEY,
    city_name VARCHAR(50) UNIQUE,
    city_owner VARCHAR(50),
    population INT DEFAULT 1000,
    population_capacity INT DEFAULT 1000,
    tax_rate DOUBLE DEFAULT 15.0,
    town BOOLEAN DEFAULT FALSE,
    FOREIGN KEY(city_owner) REFERENCES Accounts(player_id)
);

INSERT INTO Accounts VALUES ('neutral', 'Neutral', '', 2000000);

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
    building_type VARCHAR(50),
    building_level INT,
    city_id VARCHAR(50),
    city_row INT,
    city_column INT,
    FOREIGN KEY(city_id) REFERENCES Cities(city_id),
    FOREIGN KEY(building_type, building_level) REFERENCES Building_Info(building_type, building_level),
    PRIMARY KEY(city_id, city_row, city_column)
);

CREATE TABLE Builds (
    city_id VARCHAR(50),
    city_row INT,
    city_column INT,
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    FOREIGN KEY(city_id, city_row, city_column) REFERENCES Buildings(city_id, city_row, city_column),
    PRIMARY KEY(city_id, city_row, city_column)
);

CREATE TRIGGER Start_Build
AFTER INSERT ON Buildings
FOR EACH ROW
INSERT INTO Builds VALUES (
    NEW.city_id, NEW.city_row, NEW.city_column, NOW(), TIMESTAMPADD(SECOND, (SELECT build_time FROM Building_Info WHERE building_type=NEW.building_type AND building_level=NEW.building_level), NOW())
);

CREATE TRIGGER Start_Upgrade
AFTER UPDATE ON Buildings
FOR EACH ROW
INSERT INTO Builds VALUES (
    NEW.city_id, NEW.city_row, NEW.city_column, NOW(), TIMESTAMPADD(SECOND, (SELECT build_time FROM Building_Info WHERE building_type=NEW.building_type AND building_level=NEW.building_level), NOW())
);

CREATE EVENT Delete_Finished_Builds ON SCHEDULE EVERY 1 SECOND
STARTS '2023-01-01 00:00:00'
DO
DELETE FROM Builds WHERE end_time <= NOW();

CREATE EVENT Run_Production ON SCHEDULE EVERY 1 SECOND
STARTS '2023-01-01 00:00:00'
DO
UPDATE Accounts Set balance = balance + (SELECT SUM(building_production) FROM Buildings NATURAL JOIN Building_Info NATURAL JOIN Cities WHERE city_owner=player_id);

CREATE EVENT Run_Taxes ON SCHEDULE EVERY 1 SECOND
STARTS '2023-01-01 00:00:00'
DO
UPDATE Accounts Set balance = balance + (SELECT SUM(population * tax_rate / 86400) FROM Cities WHERE city_owner=player_id);

INSERT INTO Building_Info VALUES
('City Hall', 1, 0.0, 0, 100, 0.0, 0),
('Apartment', 1, 500.00, 2, 5000, 400000.00, 60),
('Apartment', 2, 750.00, 2, 7500, 800000.00, 300),
('Apartment', 3, 1000.00, 3, 10000, 1000000.00, 900),
('Apartment', 4, 2000.00, 3, 12500, 1500000.00, 1800),
('Apartment', 5, 3000.00, 3, 15000, 2000000.00, 3600),
('Apartment', 6, 5000.00, 4, 17500, 3000000.00, 9000),
('Apartment', 7, 7500.00, 4, 20000, 5000000.00, 18000),
('Apartment', 8, 1000.00, 4, 25000, 7500000.00, 36000),
('Apartment', 9, 1500.00, 4, 50000, 10000000.00, 72000),
('Apartment', 10, 2000.00, 5, 100000, 15000000.00, 144000),
('Hospital', 1, 1000.00, 5, 1000, 250000.00, 60),
('Hospital', 2, 2000.00, 5, 1000, 500000.00, 300),
('Hospital', 3, 3000.00, 5, 1000, 750000.00, 900),
('School', 1, 2000.00, 3, 500, 250000.00, 60),
('School', 2, 2000.00, 3, 500, 500000.00, 300),
('School', 3, 2000.00, 3, 500, 750000.00, 900),
('Supermarket', 1, 10000.00, 1, 250, 250000.00, 120),
('Supermarket', 2, 15000.00, 1, 250, 750000.00, 360),
('Supermarket', 3, 10000.00, 1, 250, 1500000.00, 900),
('Barracks', 1, 1000.00, 3, 500, 300000.00, 120),
('Test', 1, 0.0, 0, 0, 1.00, 1),
('Test', 2, 0.0, 0, 0, 1.00, 1);
