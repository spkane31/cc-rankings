defmodule RankingsWeb.RaceInstanceControllerTest do
  use RankingsWeb.ConnCase

  test "GET /raceinstances", %{conn: conn} do
    conn = get(conn, "/raceinstances")
    assert html_response(conn, 200) =~ "<!DOCTYPE html>"
  end

  test "GET /raceinstances/1", %{conn: conn} do
    conn = get(conn, "/raceinstances")
    assert html_response(conn, 200) =~ "<!DOCTYPE html>"
  end
end
