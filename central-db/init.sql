CREATE ROLE IF NOT EXISTS repl_user;


CREATE TABLE IF NOT EXISTS Sectors (
    sector VARCHAR(15),
    coords text,
    PRIMARY KEY (sector)
);


CREATE TABLE IF NOT EXISTS Sensors (
    sensor VARCHAR(15),
    sector VARCHAR(15),
    min_pressure FLOAT,
    coord VARCHAR(50),
    FOREIGN KEY (sector) REFERENCES Sectors(sector),
    PRIMARY KEY (sector, sensor)
);

CREATE TABLE IF NOT EXISTS Measurements (
    datetime TIMESTAMP(0),
    sensor VARCHAR(15),
    sector VARCHAR(15),
    pressure FLOAT,
    FOREIGN KEY (sector, sensor) REFERENCES Sensors(sector, sensor),
    PRIMARY KEY (datetime, sensor, sector)
);



GRANT SELECT ON ALL TABLES IN SCHEMA public TO repl_user;
