# Database

## Race
* id
* race name
* race location / course
* distance
* gender
  * There will be a different instance for the male and female races
* average time
* count on finishers
* one-to-many with Race Instance

## Race Instance
* many-to-one with Race
* average time for the year
* number of participants
* valid
  * a boolean of whether this instance counts towards the aggregate (weather,
  course change, timing error, etc.)

## Runner
* id - primary key
* first name - string
* last name - string
* team - string
* year - string
### SQL Statement to Create the table
```
CREATE TABLE runner (
id SERIAL PRIMARY KEY,
first_name VARCHAR(255),
last_name VARCHAR(255),
team VARCHAR(255),
year VARCHAR(255)
);
```

## Logging into the Databse from the CL
```
sudo -u postgres psql
```