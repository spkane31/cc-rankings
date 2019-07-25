defmodule RankingsWeb.CompareController do
  use RankingsWeb, :controller
  import Logger

  alias Rankings.{Runner, Repo}

  def index(conn, params) do
    query = Runner.list_runners(params)

    id1 = get_in(params, ["id1"])
    id2 = get_in(params, ["id2"])

    runner1 = find_runner(id1)
    runner2 = find_runner(id2)

    render(conn, "index.html", runners: query, id1: id1, id2: id2, runner1: runner1, runner2: runner2)
  end

  def first_runner(conn, params) do #%{"id1" => id1}) do
    IO.puts "def"

    query = Runner.list_runners(params)
    id1 = get_in(params, ["id1"])
    runner1 = find_runner(id1)

    render(conn, "second.html", runner1: runner1, runners: query)
  end

  def second_search(conn, params) do #%{"id1" => id1}) do
    IO.puts "second"

    query = Runner.list_runners(params)
    id1 = get_in(params, ["id1"])
    runner1 = find_runner(id1)
    IO.puts query

    render(conn, "second.html", runner1: runner1, runners: nil)
  end

  def show(conn, params) do
    IO.puts "SHOW"
    id1 = get_in(params, ["id1"])
    id2 = get_in(params, ["id2"])

    runner1 = Runner.get_runner(id1) |> Repo.preload(:team)
    runner2 = Runner.get_runner(id2) |> Repo.preload(:team)

    results1 = Runner.get_athlete_results(id1)
    results2 = Runner.get_athlete_results(id2)

    render(conn, "show.html", runner1: runner1, runner2: runner2, results1: results1, results2: results2)
  end

  def find_runner(id) do
    if id == nil do
      nil
    else
      Runner.get_runner(id) |> Repo.preload(:team)
    end
  end

end
