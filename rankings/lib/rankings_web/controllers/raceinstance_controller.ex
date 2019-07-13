defmodule RankingsWeb.RaceInstanceController do
  use RankingsWeb, :controller

  alias Rankings.RaceInstance
  alias Rankings.Repo

  def index(conn, _params) do
    instances = RaceInstance.get_n_instances(25) |> Repo.preload(:race)
    render(conn, "index.html", instances: instances)
  end

  # alias Rankings.Result
  def show(conn, %{"id" => id}) do
    instance = RaceInstance.get_instance(id) |> Repo.preload(:race)
    results = RaceInstance.get_results(id) |> Repo.preload([{:runner, :team}, {:race_instance, :race}])
    render(conn, "show.html", results: results, instance: instance)
  end
end
