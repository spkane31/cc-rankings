defmodule RankingsWeb.CompareController do
  use RankingsWeb, :controller

  alias Rankings.Runner

  def index(conn, params) do
    first = get_in(params, ["first"])
    last = get_in(params, ["last"])
    query = Runner.list_runners(params)
    render(conn, "index.html", runners: query)
  end

  def show(conn, %{"id1" => id1}) do
    render(conn, "index.html", runners: nil, id1: id1)
  end

  def show(conn, %{"id1" => id1, "id2" => id2}) do
    render(conn, "show.html", id1: id1, id2: id2)
  end
end
