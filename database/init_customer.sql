-- log into the database
-- run this commend to run this file
-- psql -h localhost -U postgres -d mynewdatabase -f init_customer.sql

-- create the table in the database
CREATE TABLE Customer_Inquiry(
    ID SERIAL PRIMARY KEY,
    Inquiry_ID VARCHAR(36)
    Email VARCHAR(255) NOT NULL, 
    Time_Inquired  TIMESTAMP NOT NULL,
    Material   VARCHAR(255),
    Loc VARCHAR(255),
    Present BOOLEAN,
    Price NUMERIC(10, 2),
    Currency VARCHAR(3),
    Data_Sheet BYTEA
)

-- Verify the table creation through listing all tables
\dt

-- check the structure of the table
\d Customer_Inquiry
