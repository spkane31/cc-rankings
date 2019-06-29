defmodule RankingsWeb.RaceControllerTest do
  use RankingsWeb.ConnCase

  test "GET /race", %{conn: conn} do
    conn = get(conn, "/races")
    assert html_response(conn, 200) =~ "<!DOCTYPE html>"
  end

  test "GET /race/1", %{conn: conn} do
    conn = get(conn, "/races")
    assert html_response(conn, 200) =~ "<!DOCTYPE html>"
  end
end
