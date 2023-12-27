-- log into the database
-- run this commend to run this file
-- psql -h localhost -U postgres -d mynewdatabase -f init_customer.sql

-- create the table in the database
CREATE TABLE customer_inquiry(
    id SERIAL PRIMARY KEY,
    inquiry_id VARCHAR(36),
    email VARCHAR(255) NOT NULL, 
    time_inquired  TIMESTAMP WITH TIME ZONE NOT NULL,
    material   VARCHAR(255),
    loc VARCHAR(255),
    present BOOLEAN,
    price NUMERIC(10, 2),
    currency VARCHAR(3),
    data_sheet BYTEA
)

-- Verify the table creation through listing all tables
\dt

-- check the structure of the table
\d customer_inquiry
