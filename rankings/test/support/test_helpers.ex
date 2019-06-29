defmodule Rankings.TestHelpers do
  alias Rankings.{
    Runner,
    Result,
    Team,
    RaceInstance,
    Race,
  }

  def runner_fixture(attrs \\ %{}) do
    name = "runner#{System.unique_integer([:positive])}"

    {:ok, runner} =
      attrs
      |> Enum.into(%{
        name: "Some Runner"
      })
      |> Repo.insert()
    runner
  end
end
