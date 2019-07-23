defmodule RankingsWeb.Router do
  use RankingsWeb, :router

  pipeline :browser do
    plug :accepts, ["html"]
    plug :fetch_session
    plug :fetch_flash
    plug :protect_from_forgery
    plug :put_secure_browser_headers
    plug RankingsWeb.Auth
  end

  pipeline :api do
    plug :accepts, ["json"]
  end

  scope "/", RankingsWeb do
    pipe_through :browser

    get "/", PageController, :index

    resources "/users", UserController, only: [:index, :show, :new, :create]
    resources "/sessions", SessionController, only: [:new, :create, :delete]
    resources "/runners", RunnerController, only:  [:index, :show]
    resources "/teams", TeamController, only: [:index, :show]
    resources "/races", RaceController,  only: [:index, :show]
    resources "/raceinstances", RaceInstanceController,  only: [:index, :show]
    resources "/results", ResultController,  only: [:index, :show]
    resources "/edges", EdgeController, only: [:index, :show]
    resources "/about", AboutController, only: [:index]
    # resources "/compare", CompareController, only: [:index, :show]
    resources "/virtual", VirtualController, only: [:index]


    #
    get "/compare", CompareController, :index
    get "/compare/:id1", CompareController, :first_runner
    get "/compare/:id1/search", CompareController, :second_search
    get "/compare/:id1/:id2", CompareController, :show


		get "/search", SearchController, :index
		post "/search", SearchController, :index

end

  # Other scopes may use custom stacks.
  # scope "/api", RankingsWeb do
  #   pipe_through :api
  # end
end
