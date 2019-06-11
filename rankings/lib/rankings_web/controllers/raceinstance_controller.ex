defmodule RankingsWeb.RaceInstanceController do
  use RankingsWeb, :controller

  alias Rankings.RaceInstance
  alias Rankings.Result
  alias Rankings.Repo

  def index(conn, _params) do
    instances = RaceInstance.list_race_instances()
    render(conn, "index.html", instances: instances)
  end

  def show(conn, %{"id" => id}) do
    instance = RaceInstance.get_instance(id) |> Repo.preload(:results)
    results = RaceInstance.get_results(id)
    render(conn, "show.html", instance: instance, results: results)
  end
end
