# How to Work on This Project

## Dependencies
psql (PostgreSQL) - v10.7
nodejs - v8.10.0
elixir - v1.8.1
erlang - 21

## File Structure

assets
* Browser files like JavaScript and CSS
config
* Phoenix configuration goes into config
lib
    hello
    * supervision trees, long-running processes, application business logic
    hello_web
    * web-related code (controllers, views, and templates)
test
* tests (duh)

## Elixir Configuration

* lib
    * hello
    * hello_web
        * endpoint.ex
        * ...
    * hello.ex
    * hello_web.ex
* mix.exs
* mix.lock
* test