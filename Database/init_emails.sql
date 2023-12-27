-- log into the database
-- run this commend to run this file
-- psql -h localhost -U postgres -d mynewdatabase -f init_emails.sql

-- create the table in the database
CREATE TABLE emails(
    id SERIAL PRIMARY KEY,
    inquiry_id VARCHAR(36),
    client_email VARCHAR(255) NOT NULL, 
    time_sent  TIMESTAMP WITH TIME ZONE NOT NULL,
    material   VARCHAR(255),
    supplier_map_id VARCHAR(255),
    supplier_name VARCHAR(255),
    supplier_lat REAL,
    supplier_lng REAL,
    supplier_email VARCHAR(255),
    sent_out BOOLEAN,
    present BOOLEAN,
    price NUMERIC(10, 2),
    currency VARCHAR(3),
    data_sheet BYTEA
)

-- Verify the table creation through listing all tables
\dt

-- check the structure of the table
\d emails