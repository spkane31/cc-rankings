defmodule RankingsWeb.RunnerControllerTest do
  use RankingsWeb.ConnCase

  test "GET /runners", %{conn: conn} do
    conn = get(conn, "/runners")
    assert html_response(conn, 200) =~ "<!DOCTYPE html>"
  end

  test "GET /runners/1", %{conn: conn} do
    conn = get(conn, "/runners")
    assert html_response(conn, 200) =~ "<!DOCTYPE html>"
  end
end
