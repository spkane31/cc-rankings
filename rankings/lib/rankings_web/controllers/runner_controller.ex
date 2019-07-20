defmodule RankingsWeb.RunnerController do
  use RankingsWeb, :controller

  alias Rankings.Runner
  alias Rankings.Result
  alias Rankings.Repo
  alias Rankings.Edge

  def index(conn, _params) do
    runners = Runner.last_n_runners(25) |> Repo.preload(:team)
    render(conn, "index.html", runners: runners)
  end

  def show(conn, %{"id" => id}) do
    runner = Runner.get_runner(id) |> Repo.preload(:team)
    team = Runner.get_team_name(id)
    edges = Edge.get_runner_edges(id)
    results = Runner.get_athlete_results(id) |> Repo.preload([{:race_instance, :race}])
    render(conn, "show.html", runner: runner, team: team, results: results, edges: edges)
  end
end
