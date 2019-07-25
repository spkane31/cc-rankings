defmodule RankingsWeb.VirtualController do
  use RankingsWeb, :controller

  def index(conn, _params) do
    render(conn, "index.html")
  end
end
