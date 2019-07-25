# Rankings

[![Greenkeeper badge](https://badges.greenkeeper.io/spkane31/cc-rankings.svg)](https://greenkeeper.io/)
[![Build Status](https://travis-ci.org/spkane31/cc-rankings.svg?branch=master)](https://travis-ci.org/spkane31/cc-rankings)

This is a project which will rank collegiate cross country performances based on relative performances.

To start your Phoenix server:

  * Install dependencies with `mix deps.get`
  * Create and migrate your database with `mix ecto.setup`
  * Install Node.js dependencies with `cd assets && npm install`
  * Start Phoenix endpoint with `mix phx.server`

Now you can visit [`localhost:4000`](http://localhost:4000) from your browser.

Ready to run in production? Please [check our deployment guides](https://hexdocs.pm/phoenix/deployment.html).

## Learn more

  * Official website: http://www.phoenixframework.org/
  * Guides: https://hexdocs.pm/phoenix/overview.html
  * Docs: https://hexdocs.pm/phoenix
  * Mailing list: http://groups.google.com/group/phoenix-talk
  * Source: https://github.com/phoenixframework/phoenix

## Useful Documents
  * https://devhints.io/phoenix-ecto


There are currently three parts
## Rankings
A web application built with Elixir and Phoenix


## Scraper
A web scraping tool built with Go. It currently scrapes TFRRS.com for cross country results

### TODO
* Database Changes
  * ...

* Web Scraper Changes
  * 

* Database Insertion Changes
  * Check year and date to see if runner is now a year older instead of creating new one
  * Reset the year, marks everyone as inactive and then increases everyone's year by one

* Data Analysis Changes

## Analysis
A tool built with Go to perform the heavy analysis of all data. I'd like to move this to Rust

### TODO
* [I] Data Analysis
* [W] Scraper
* [W] Web Client
* [I] CI/CD