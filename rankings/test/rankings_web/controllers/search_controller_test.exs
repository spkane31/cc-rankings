defmodule RankingsWeb.SearchControllerTest do
  use RankingsWeb.ConnCase

  test "GET /search", %{conn: conn} do
    conn = get(conn, "/results")
    assert html_response(conn, 200) =~ "<!DOCTYPE html>"
  end

  test "POST /seach" do
    conn = post(build_conn(), "/search", [first: "Sean", last: "Kane"])
    assert html_response(conn, 200) =~ "<!DOCTYPE html>"
  end
end
