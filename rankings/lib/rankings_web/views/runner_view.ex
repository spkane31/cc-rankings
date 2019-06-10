defmodule RankingsWeb.RunnerView do
  use RankingsWeb, :view

  alias Rankings.Runner

  def get_first_name(%Runner{first_name: first}) do
    first
    |> String.split(" ")
    |> Enum.at(0)
  end
end
