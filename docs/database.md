# Database Design

## Races
* id
* name
* course
* distance
* gender
  * Note: There will be a different instance for the male and female races
* is_base
* average 
* std_dev
* correction_avg
* correction_graph
* one-to-many with Race Instance

## Race Instances
* id
* date
* race_id
* average
* std_dev
* many-to-one with Race
* valid - Need to add
  * a boolean of whether this instance counts towards the aggregate (weather,
  course change, timing error, etc.)

## Runners
* id
* first name - string
* last name - string
* team - string
* year - string
* team_id
* gender
* one-to-many w/ Results

## Teams
* name
* region (not sure how to do this one yet)
* conference (also not sure)
* runners: one-to-many relation

## Results
* id
* distance
* unit
* rating
* time
* scaled_time
* time_float
* many-to-one w/ Runners
* many-to-one w/ Race Instances

## Graph
TODO

## Logging into the Databse from the CL
```
sudo -u postgres psql
```