defmodule RankingsWeb.RaceInstanceController do
  use RankingsWeb, :controller

  alias Rankings.RaceInstance
  # alias Rankings.Repo

  def index(conn, _params) do
    instances = RaceInstance.list_race_instances()
    render(conn, "index.html", instances: instances)
  end

  alias Rankings.Result
  def show(conn, %{"id" => id}) do
    results = RaceInstance.get_results(id)
    render(conn, "show.html", results: results)
  end
end
