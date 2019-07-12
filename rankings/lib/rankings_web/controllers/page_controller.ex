defmodule RankingsWeb.PageController do
  use RankingsWeb, :controller
  alias Rankings.Repo
  import Ecto.Query

  alias Rankings.{Team, Runner, Result, Race, RaceInstance}
  def index(conn, _params) do
    runners = Repo.one(from r in Runner, select: count(r.id))
    teams = Repo.one(from t in Team, select: count(t.id))
    results = Repo.one(from r in Result, select: count(r.id))
    races = Repo.one(from r in Race, select: count(r.id))
    instances = Repo.one(from i in RaceInstance, select: count(i.id))

    # teams = length(Repo.all(Team))
    # runners = length(Repo.all(Runner))
    # results = length(Repo.all(Result))
    # races = length(Repo.all(Race))
    # instances = length(Repo.all(RaceInstance))
    render(conn, "index.html", teams: teams, runners: runners, results: results, races: races, instances: instances)
  end
end
