defmodule RankingsWeb.RunnerController do
  use RankingsWeb, :controller

  alias Rankings.{Runner, Repo, Edge}

  plug :authenticate when action in [:edit, :delete, :create]

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

  def edit(conn, %{"id" => id}) do
    runner = Runner.get_runner(id) |> Repo.preload(:team)
    render(conn, "edit.html", runner: runner)
  end

  defp authenticate(conn, _opts) do
    if conn.assigns.current_user do
      conn
    else
      conn
      |> put_flash(:error, "You must be logged in to access that page")
      |> redirect(to: Routes.page_path(conn, :index))
      |> halt()
    end
  end
end
