CREATE TABLE Accounts (
    player_id VARCHAR(50) PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    password CHAR(64) NOT NULL,
    balance DOUBLE DEFAULT 2000000.0 NOT NULL,
    CHECK (balance > 0),
    UNIQUE(username)
);

CREATE TABLE Cities (
    city_id VARCHAR(50) PRIMARY KEY,
    city_name VARCHAR(50) UNIQUE,
    city_owner VARCHAR(50),
    population INT DEFAULT 1000000,
    population_capacity INT DEFAULT 1000,
    tax_rate DOUBLE DEFAULT 15.0,
    army_size INT DEFAULT 0,
    town BOOLEAN DEFAULT FALSE,
    CHECK (army_size >= 0),
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
SET NEW.city_name=CONCAT(IF(NEW.town=0, 'City ', 'Town '), NEW.city_id);

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

CREATE TABLE Marches (
    march_id VARCHAR(50) PRIMARY KEY,
    from_city VARCHAR(50) NOT NULL,
    to_city VARCHAR(50) NOT NULL,
    army_size VARCHAR(50),
    attack BOOLEAN,
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    FOREIGN KEY(from_city) REFERENCES Cities(city_id),
    FOREIGN KEY(to_city) REFERENCES Cities(city_id)
);

CREATE TABLE Training (
    city_id VARCHAR(50) PRIMARY KEY,
    army_size INT,
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    FOREIGN KEY(city_id) REFERENCES Cities(city_id)
);

CREATE TABLE Battles (
    battle_id VARCHAR(50) PRIMARY KEY,
    from_city VARCHAR(50),
    from_city_owner VARCHAR(50),
    to_city VARCHAR(50),
    to_city_owner VARCHAR(50),
    battle_time TIMESTAMP,
    attacker_army_size INT,
    defender_army_size INT,
    attack_victory BOOLEAN,
    amount_looted DOUBLE,
    FOREIGN KEY(from_city) REFERENCES Cities(city_id),
    FOREIGN KEY(to_city) REFERENCES Cities(city_id),
    FOREIGN KEY(from_city_owner) REFERENCES Accounts(username),
    FOREIGN KEY(to_city_owner) REFERENCES Accounts(username)
);

CREATE VIEW Building_Ownership
AS
SELECT *
FROM Cities
     NATURAL JOIN Buildings
     NATURAL JOIN Building_Info
     JOIN Accounts ON city_owner=player_id;

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

CREATE TRIGGER Buy_Troops
BEFORE INSERT ON Training
FOR EACH ROW
UPDATE Accounts SET balance=balance-(NEW.army_size * 1000) WHERE player_id=(SELECT city_owner FROM Cities WHERE city_id=NEW.city_id);

CREATE TRIGGER Train_Troops
AFTER DELETE ON Training
FOR EACH ROW
UPDATE Cities SET army_size=army_size+OLD.army_size WHERE city_id=OLD.city_id;

CREATE EVENT Delete_Finished_Builds ON SCHEDULE EVERY 1 SECOND
STARTS '2023-01-01 00:00:00'
DO
DELETE FROM Builds WHERE end_time <= NOW();

CREATE EVENT Run_Production ON SCHEDULE EVERY 1 SECOND
STARTS '2023-01-01 00:00:00'
DO
UPDATE Accounts SET balance = balance + (SELECT IF(SUM(building_production) IS NULL, 0, SUM(building_production/3600)) FROM Buildings NATURAL JOIN Building_Info NATURAL JOIN Cities LEFT JOIN Builds ON Buildings.city_id=Builds.city_id AND Buildings.city_row=Builds.city_row AND Buildings.city_column=Builds.city_column WHERE start_time IS NULL AND city_owner=player_id) WHERE username != 'Neutral';

CREATE EVENT Run_Taxes ON SCHEDULE EVERY 1 SECOND
STARTS '2023-01-01 00:00:00'
DO
UPDATE Accounts Set balance = balance + (SELECT SUM(population * tax_rate / 86400) FROM Cities WHERE city_owner=player_id) WHERE username != 'Neutral';

CREATE EVENT Finish_Troop_Training ON SCHEDULE EVERY 1 SECOND
STARTS '2023-01-01 00:00:00'
DO
DELETE FROM Training WHERE end_time <= NOW();

CREATE EVENT Generate_Town ON SCHEDULE EVERY 1 DAY
STARTS '2023-01-01 00:00:00'
DO
INSERT INTO Cities (city_id, city_owner, population, population_capacity, town) VALUES
(uuid(), 'neutral', (SELECT RAND()*(1000000-50000)+50000), 1000000, 1);

INSERT INTO Building_Info VALUES
('City Hall', 1, 0.0, 0, 100, 0.0, 0),
('Apartment', 1, 1000.00, 2, 5000, 400000.00, 90),
('Apartment', 2, 2000.00, 2, 7500, 800000.00, 270),
('Apartment', 3, 4000.00, 3, 10000, 1000000.00, 900),
('Apartment', 4, 10000.00, 3, 12500, 1500000.00, 1800),
('Apartment', 5, 20000.00, 3, 15000, 2000000.00, 3600),
('Apartment', 6, 50000.00, 4, 17500, 3000000.00, 9000),
('Apartment', 7, 75000.00, 4, 20000, 5000000.00, 18000),
('Apartment', 8, 125000.00, 4, 25000, 7500000.00, 36000),
('Apartment', 9, 150000.00, 4, 50000, 10000000.00, 72000),
('Apartment', 10, 175000.00, 5, 100000, 15000000.00, 144000),
('Hospital', 1, 2000.00, 5, 1000, 250000.00, 60),
('Hospital', 2, 6000.00, 5, 1000, 500000.00, 300),
('Hospital', 3, 13000.00, 5, 1000, 750000.00, 900),
('Hospital', 4, 20000.00, 5, 1000, 1000000.00, 1800),
('Hospital', 5, 40000.00, 5, 1000, 5000000.00, 9000),
('Hospital', 6, 60000.00, 5, 1000, 7500000.00, 18000),
('Hospital', 7, 100000.00, 5, 1000, 10000000.00, 30000),
('Hospital', 8, 200000.00, 5, 1000, 30000000.00, 70000),
('Hospital', 9, 400000.00, 5, 1000, 50000000.00, 100000),
('Hospital', 10, 600000.00, 5, 1000, 750000000.00, 300000),
('School', 1, 2000.00, 3, 500, 250000.00, 60),
('School', 2, 2000.00, 3, 500, 500000.00, 300),
('School', 3, 2000.00, 3, 500, 7500000.00, 900),
('School', 4, 20000.00, 3, 500, 10000000.00, 2700),
('School', 5, 30000.00, 3, 500, 40000000.00, 5400),
('School', 6, 600000.00, 3, 500, 75000000.00, 20000),
('School', 7, 800000.00, 3, 500, 100000000.00, 60000),
('School', 8, 1000000.00, 3, 500, 300000000.00, 100000),
('School', 9, 2000000.00, 3, 500, 600000000.00, 300000),
('School', 10, 9999999.99, 3, 500, 1000000000.00, 600000),
('Supermarket', 1, 5000.00, 1, 250, 250000.00, 120),
('Supermarket', 2, 7500.00, 1, 250, 750000.00, 360),
('Supermarket', 3, 10000.00, 1, 250, 1500000.00, 900),
('Supermarket', 4, 20000.00, 1, 250, 2500000.00, 1800),
('Supermarket', 5, 30000.00, 1, 250, 5000000.00, 3600),
('Supermarket', 6, 40000.00, 1, 250, 5000000.00, 4800),
('Supermarket', 7, 50000.00, 1, 250, 5000000.00, 6000),
('Supermarket', 8, 60000.00, 1, 250, 5000000.00, 7200),
('Supermarket', 9, 70000.00, 1, 250, 5000000.00, 8400),
('Supermarket', 10, 80000.00, 1, 250, 5000000.00, 9600),
('Barracks', 1, 1000.00, 3, 500, 300000.00, 120),
('Barracks', 2, 2000.00, 3, 500, 600000.00, 360),
('Barracks', 3, 3000.00, 4, 500, 900000.00, 720),
('Barracks', 4, 3000.00, 4, 500, 1000000.00, 7200),
('Barracks', 5, 3000.00, 4, 500, 2000000.00, 14400),
('Barracks', 6, 3000.00, 4, 500, 4000000.00, 60000),
('Barracks', 7, 3000.00, 4, 500, 6000000.00, 100000),
('Barracks', 8, 3000.00, 4, 500, 9000000.00, 300000),
('Barracks', 9, 3000.00, 4, 500, 10000000.00, 600000),
('Barracks', 10, 3000.00, 4, 500, 50000000.00, 600000),
('Test', 1, 0.0, 0, 0, 1.00, 1),
('Test', 2, 0.0, 0, 0, 1.00, 1);

INSERT INTO Cities (city_id, city_owner, population, population_capacity, town) VALUES
(uuid(), 'neutral', (SELECT RAND()*(1000000-50000)+50000), 1000000, 1),
(uuid(), 'neutral', (SELECT RAND()*(1000000-50000)+50000), 1000000, 1),
(uuid(), 'neutral', (SELECT RAND()*(1000000-50000)+50000), 1000000, 1),
(uuid(), 'neutral', (SELECT RAND()*(1000000-50000)+50000), 1000000, 1),
(uuid(), 'neutral', (SELECT RAND()*(1000000-50000)+50000), 1000000, 1),
(uuid(), 'neutral', (SELECT RAND()*(1000000-50000)+50000), 1000000, 1),
(uuid(), 'neutral', (SELECT RAND()*(1000000-50000)+50000), 1000000, 1),
(uuid(), 'neutral', (SELECT RAND()*(1000000-50000)+50000), 1000000, 1),
(uuid(), 'neutral', (SELECT RAND()*(1000000-50000)+50000), 1000000, 1),
(uuid(), 'neutral', (SELECT RAND()*(1000000-50000)+50000), 1000000, 1),
(uuid(), 'neutral', (SELECT RAND()*(1000000-50000)+50000), 1000000, 1),
(uuid(), 'neutral', (SELECT RAND()*(1000000-50000)+50000), 1000000, 1),
(uuid(), 'neutral', (SELECT RAND()*(1000000-50000)+50000), 1000000, 1);

DELIMITER &&
CREATE PROCEDURE reset_tests ()
BEGIN
    DELETE FROM city_io.Battles WHERE from_city IN (SELECT city_id FROM city_io.Cities JOIN city_io.Accounts ON city_owner=player_id WHERE username='User200') OR to_city IN (SELECT city_id FROM city_io.Cities JOIN city_io.Accounts ON city_owner=player_id WHERE username='User200');
    DELETE FROM city_io.Training WHERE city_id IN (SELECT city_id FROM city_io.Cities JOIN city_io.Accounts ON city_owner=player_id WHERE username IN ('User200', 'train'));
    DELETE FROM city_io.Builds WHERE city_id IN (SELECT city_id FROM city_io.Cities JOIN city_io.Accounts ON city_owner=player_id WHERE username IN ('User200', 'train'));
    DELETE FROM city_io.Buildings WHERE city_id IN (SELECT city_id FROM city_io.Cities JOIN city_io.Accounts ON city_owner=player_id WHERE username IN ('User200', 'train'));
    DELETE FROM city_io.Cities WHERE city_id=(SELECT city_id FROM (SELECT * FROM city_io.Cities) AS TempCities JOIN city_io.Accounts ON city_owner=player_id WHERE username='User200');
    DELETE FROM city_io.Accounts WHERE username='User200';
END &&
DELIMITER ;
