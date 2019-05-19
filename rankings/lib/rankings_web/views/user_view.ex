defmodule RankingsWeb.UserView do
  use RankingsWeb, :view

  alias Rankings.Accounts

  def first_name(%Accounts.User{name: name}) do
    name
    |> String.split(" ")
    |> Enum.at(0)
  end
end
