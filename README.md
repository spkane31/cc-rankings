# CC-Rankings


[![Greenkeeper badge](https://badges.greenkeeper.io/spkane31/cc-rankings.svg)](https://greenkeeper.io/)
[![Build Status](https://travis-ci.org/spkane31/cc-rankings.svg?branch=master)](https://travis-ci.org/spkane31/cc-rankings)

This is a project which will rank collegiate cross country performances based on relative performances.

There are currently three parts
## Rankings
A web application built with Elixir and Phoenix
### TODO

## Scraper
A web scraping tool built with Go. It currently scrapes TFRRS.com for cross country results
### TODO

* Database Changes
  * [I] Add date, gender to results
  * [I] Create a graph database
  * [I] Add a ```counted``` column to results for results already taken care of in graph database
  * [I] 

* Web Scraper Changes
  * [I] create sub-folder for years, instead of grouping races together


* Database Insertion Changes
  * Check year and date to see if runner is now a year older instead of creating new one


* Data Analysis Changes

## Analysis
A tool built with Go to perform the heavy analysis of all data. 

### TODO
* [I] Data Analysis
* [W] Scraper
* [W] Web Client
* [I] CI/CD