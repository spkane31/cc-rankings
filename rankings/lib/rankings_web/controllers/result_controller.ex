defmodule RankingsWeb.ResultController do
  use RankingsWeb, :controller

  alias Rankings.Result
  alias Rankings.Repo

  def index(conn, _params) do
    results = Result.last_n_results(25) |> Repo.preload([{:runner, :team}, {:race_instance, :race}])
    render(conn, "index.html", results: results)
  end

  def show(conn, %{"id" => id}) do
    result = Result.get_result(id) |> Repo.preload([{:runner, :team}, {:race_instance, :race}])
    render(conn, "show.html", result: result)
  end
end
