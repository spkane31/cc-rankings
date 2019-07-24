defmodule RankingsWeb.ResultControllerTest do
  use RankingsWeb.ConnCase

  test "GET /results", %{conn: conn} do
    conn = get(conn, "/results")
    assert html_response(conn, 200) =~ "<!DOCTYPE html>"
  end

  test "GET /results/1", %{conn: conn} do
    conn = get(conn, "/results/1")
    assert html_response(conn, 200) =~ "<!DOCTYPE html>"
  end
end
