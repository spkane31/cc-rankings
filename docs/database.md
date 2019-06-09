# Database Design

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

## Team
* Name
* Region (not sure how to do this one yet)
* Conference (also not sure)
* runners: one-to-one relation


### SQL Statement to Create the table
```
CREATE TABLE runners (
id SERIAL PRIMARY KEY,
first_name VARCHAR(255) NOT NULL,
last_name VARCHAR(255) NOT NULL,
year VARCHAR(255) 
);

CREATE TABLE teams (
id SERIAL PRIMARY KEY,
name VARCHAR(255) NOT NULL,
region VARCHAR(255),
conference VARCHAR(255),
CONSTRAINT fk_runner_id FOREIGN KEY (id) REFERENCES runners (id)

CREATE TABLE races (
id SERIAL PRIMARY KEY,
name VARCHAR(255) NOT NULL,
course VARCHAR(255) NOT NULL,
distance INT NOT NULL,
gender VARCHAR(255) NOT NULL,
correction REAL);

CREATE TABLE race_instances (
id SERIAL PRIMARY KEY,
date VARCHAR(255) NOT NULL,
race_id INT NOT NULL,
FOREIGN KEY (race_id) REFERENCES races(id) ON DELETE CASCADE
);

CREATE TABLE results (
id SERIAL PRIMARY KEY,
distance int not null,
rating real not null,
time VARCHAR(255),
race_instance_id INT NOT NULL,
FOREIGN KEY (race_instance_id) REFERENCES race_instances(id) ON DELETE CASCADE,
runner_id INT NOT NULL,
FOREIGN KEY (runner_id) REFERENCES runners(id) ON DELETE CASCADE
);

```

## Logging into the Databse from the CL
```
sudo -u postgres psql
```