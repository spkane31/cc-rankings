defmodule RankingsWeb.Router do
  use RankingsWeb, :router

  pipeline :browser do
    plug :accepts, ["html"]
    plug :fetch_session
    plug :fetch_flash
    plug :protect_from_forgery
    plug :put_secure_browser_headers
  end

  pipeline :api do
    plug :accepts, ["json"]
  end

  scope "/", RankingsWeb do
    pipe_through :browser

    get "/users", UserController, :index
    get "/users/:id", UserController, :show
    get "/runners", RunnerController, :index
    get "/runners/:id", RunnerController, :show
    get "/", PageController, :index
  end

  # Other scopes may use custom stacks.
  # scope "/api", RankingsWeb do
  #   pipe_through :api
  # end
end
