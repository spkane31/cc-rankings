defmodule RankingsWeb.CompareControllerTest do
  use RankingsWeb.ConnCase

  test "GET /compare", %{conn: conn} do
    conn = get(conn, "/compare")
    assert html_response(conn, 200) =~ "<!DOCTYPE html>"
  end

  test "POST /compare" do
    conn = post(build_conn(), "/compare", [first: "Sean", last: "Kane"])
    assert html_response(conn, 200) =~ "<!DOCTYPE html>"
  end

  test "GET /compare/:id1", %{conn: conn} do
    conn = get(conn, "/compare/1")
    assert html_response(conn, 200) =~ "<!DOCTYPE html>"
  end

  test "POST /compare/:id1" do
    conn = post(build_conn(), "/compare/1", [first: "Sean", last: "Kane"])
    assert html_response(conn, 200) =~ "<!DOCTYPE html>"
  end

  test "GET /compare/:id1/:id2", %{conn: conn} do
    conn = get(conn, "compare/1/2")
    assert html_response(conn, 200) =~ "<!DOCTYPE html>"
  end
end
