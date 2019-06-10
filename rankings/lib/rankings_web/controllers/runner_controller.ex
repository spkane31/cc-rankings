defmodule RankingsWeb.RunnerController do
  use RankingsWeb, :controller

  alias Rankings.Runner

  def index(conn, _params) do
    runners = Runner.list_runners()
    render(conn, "index.html", runners: runners)
  end

  def show(conn, %{"id" => id}) do
    runner = Runner.get_runner(id)
    render(conn, "show.html", runner: runner)
  end
end
