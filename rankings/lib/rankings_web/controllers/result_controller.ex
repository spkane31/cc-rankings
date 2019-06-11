defmodule RankingsWeb.ResultController do
  use RankingsWeb, :controller

  alias Rankings.Result

  def index(conn, _params) do
    results = Result.list_results()
    render(conn, "index.html", results: results)
  end

  def show(conn, %{"id" => id}) do
    result = Result.get_result(id)
    render(conn, "show.html", result: result)
  end
end
