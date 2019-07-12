defmodule RankingsWeb.TeamControllerTest do
  use RankingsWeb.ConnCase

  test "GET /teams", %{conn: conn} do
    conn = get(conn, "/teams")
    assert html_response(conn, 200) =~ "<!DOCTYPE html>"
  end

  test "GET /teams/1", %{conn: conn} do
    conn = get(conn, "/teams")
    assert html_response(conn, 200) =~ "<!DOCTYPE html>"
  end
end
