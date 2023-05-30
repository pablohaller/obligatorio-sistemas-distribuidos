CREATE USER replicator WITH REPLICATION ENCRYPTED PASSWORD 'my_replicator_password';
SELECT * FROM pg_create_physical_replication_slot('replication_slot_slave1');

CREATE TABLE IF NOT EXISTS Mediciones (
    datetime TIMESTAMP(0),
    sensor VARCHAR(255),
    sector VARCHAR(255),
    presion INT,
    PRIMARY KEY (datetime, sensor, sector)
);
