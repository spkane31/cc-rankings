defmodule RankingsWeb.RaceInstanceControllerTest do
  use RankingsWeb.ConnCase

  test "GET /raceinstancess", %{conn: conn} do
    conn = get(conn, "/raceinstances")
    assert html_response(conn, 200) =~ "<!DOCTYPE html>"
  end

  test "GET /raceinstancess/1", %{conn: conn} do
    conn = get(conn, "/raceinstances")
    assert html_response(conn, 200) =~ "<!DOCTYPE html>"
  end
end
