defmodule RankingsWeb.RaceInstanceController do
  use RankingsWeb, :controller

  alias Rankings.RaceInstance

  def index(conn, _params) do
    instances = RaceInstance.list_race_instances()
    render(conn, "index.html", instances: instances)
  end

  def show(conn, %{"id" => id}) do
    instance = RaceInstance.get_instance(id)
    render(conn, "show.html", instance: instance)
  end
end
