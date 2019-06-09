defmodule Rankings.Accounts do
  @moduledoc """
  The Accounts context
  """

  alias Rankings.Accounts.User

  def list_users do
    [
      %User{id: "1", name: "Sean", username: "spkane31"},
      %User{id: "2", name: "Jack", username: "jackrandall"},
      %User{id: "3", name: "Kyle", username: "kyleklingler"}
    ]
  end

  def get_user(id) do
    Enum.find(list_users(), fn map -> map.id == id end)
  end

  def get_user_by(params) do
    Enum.find(list_users(), fn map ->
      Enum.all?(params, fn {key, val} -> Map.get(map, key) == val end)
    end)
  end
end
