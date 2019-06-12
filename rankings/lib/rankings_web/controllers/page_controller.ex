defmodule RankingsWeb.PageController do
  use RankingsWeb, :controller
  alias Rankings.Repo

  alias Rankings.{Team, Runner, Result, Race, RaceInstance}
  def index(conn, _params) do
    teams = length(Repo.all(Team))
    runners = length(Repo.all(Runner))
    results = length(Repo.all(Result))
    races = length(Repo.all(Race))
    instances = length(Repo.all(RaceInstance))
    render(conn, "index.html", teams: teams, runners: runners, results: results, races: races, instances: instances)
  end
end
