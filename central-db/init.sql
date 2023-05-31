GRANT SELECT ON ALL TABLES IN SCHEMA public TO repl_user;

CREATE TABLE IF NOT EXISTS Mediciones (
    datetime TIMESTAMP(0),
    sensor VARCHAR(255),
    sector VARCHAR(255),
    presion INT,
    PRIMARY KEY (datetime, sensor, sector)
);
