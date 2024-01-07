# Following are the steps to install postgress
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo service postgresql start

# initialize the database
createdb -U postgres mynewdatabase
psql -U postgres -d mynewdatabase -f ./Database/init_customer.sql
psql -U postgres -d mynewdatabase -f ./Database/init_emails.sql