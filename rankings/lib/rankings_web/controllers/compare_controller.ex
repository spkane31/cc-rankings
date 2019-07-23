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

  def first_runner(conn, %{"id1" => id1}) do
    IO.puts "def"

    # query = Runner.list_runners(params)
    # id1 = get_in(params, ["id1"])
    runner1 = find_runner(id1)

    render(conn, "second.html", runner1: runner1, runners: nil)
  end

  def second_search(conn, %{"id1" => id1}, params) do
    IO.puts "second"

    query = Runner.list_runners(params)
    # id1 = get_in(params, ["id1"])
    runner1 = find_runner(id1)

    render(conn, "second.html", runner1: runner1, runners: query)
  end

  def show(conn, %{"id1" => id1}, params) do
    IO.puts "ghi"
    # id1 = get_in(params, ["id1"])
    id2 = get_in(params, ["id2"])
    render(conn, "index.html", id1: id1, id1: id2)
  end

  def find_runner(id) do
    if id == nil do
      nil
    else
      Runner.get_runner(id) |> Repo.preload(:team)
    end
  end

end
