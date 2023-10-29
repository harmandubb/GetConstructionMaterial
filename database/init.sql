-- init_db.sql

-- to initialize the database based on this script 
-- psql -h hostname -U username -d databasename -f init_db.sql
-- hostname == localhost 
-- username == postgres
-- databasename is the name that you have given to your desired database

-- Creating a table to store product information
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100),
    data_sheet BYTEA,
    pictures BYTEA[],
    price NUMERIC(10, 2)
);

-- Optional: Inserting some sample data
INSERT INTO products (name, category, price)
VALUES
('Laptop', 'Electronics', 999.99),
('Coffee Maker', 'Appliances', 59.99);

--Testing commands for the database 
-- 1. Connect to the database
-- psql -h hostname -U username -d databasename

-- 2. List Tables 
-- \dt

-- 3. Describe Table Structure
-- \d tablename 

-- 4. Check for Sample Data 
-- SELECT * FROM tablename;