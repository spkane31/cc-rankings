# Troubleshooting

## PostgreSQL connection
Update and install PostgreSQL
```
sudo apt-get update
sudo apt-get install postgresql-10.4
```
Force the postgres user to have the password 'postgres'
```
sudo -u postgres psql -c "ALTER USER postgres PASSWORD 'postgres';"
```
Start the PostgreSQL server
```
sudo service postgresql start
```
Stop the PostgreSQL server
```
sudo service postgresql stop
```