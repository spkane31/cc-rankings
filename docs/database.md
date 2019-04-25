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