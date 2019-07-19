defmodule RankingsWeb.EdgeController do
  use RankingsWeb, :controller

  alias Rankings.{Edge, Runner, RaceInstance, Race, Repo}

  def index(conn, _params) do
    edges = Edge.slowest_n_edges(100)
    render(conn, "index.html", edges: edges)
  end

  def show(conn, %{"id" => id}) do
    edge = Edge.get_edge(id)
    sim_edges = Edge.get_edges(edge.from_race_id, edge.to_race_id)
    count = Edge.count_edges(edge.from_race_id, edge.to_race_id)
    avg = Edge.get_avg(edge.from_race_id, edge.to_race_id)
    render(conn, "show.html", edge: edge, sim_edges: sim_edges, count: count, avg: avg)
  end
end
