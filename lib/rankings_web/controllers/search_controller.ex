defmodule RankingsWeb.SearchController do
  use RankingsWeb, :controller
  # alias Rankings.Repo
  # import Ecto.Query

  alias Rankings.{Team, Runner, Result, Race, RaceInstance}
  def index(conn, params) do

    first = get_in(params, ["first"])
    last = get_in(params, ["last"])
    team = get_in(params, ["team"])
    race = get_in(params, ["race"])
    if team == nil and race == nil do
      query = Runner.list_runners(params)
      render(conn, "index.html", query: query, runners: true, teams: false, races: false)
    else
      if team != nil do
        query = Team.list_teams(params)
        render(conn, "index.html", query: query, runners: false, teams: true, races: false)
      else
        query = Race.list_races(params)
        render(conn, "index.html", query: query, runners: false, teams: false, races: true)
      end
    end

    # runners = Runner.list_runners(params)
    # render(conn, "index.html", query: query)
  end

end
