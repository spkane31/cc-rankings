defmodule RankingsWeb.RunnerController do
  use RankingsWeb, :controller

  alias Rankings.Runner
  alias Rankings.Result
  alias Rankings.Repo

  def index(conn, _params) do
    runners = Runner.list_runners()
    render(conn, "index.html", runners: runners)
  end

  def show(conn, %{"id" => id}) do
    runner = Runner.get_runner(id) |> Repo.preload(:team)
    # runner = Repo.preload(runner, [:team, :results])
    team = Runner.get_team_name(id)
    results = Runner.get_athlete_results(id)#Repo.get(Runner, id) |> Repo.preload([{:results, :runner}]) #Runner.get_results(id) |> Repo.preload({:race_instance, :runner})
    render(conn, "show.html", runner: runner, team: team, results: results)
  end
end
