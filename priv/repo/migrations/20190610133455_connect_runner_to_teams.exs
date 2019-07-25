defmodule Rankings.Repo.Migrations.ConnectRunnerToTeams do
  use Ecto.Migration

  def change do
    alter table(:runners) do
      add :team_id, references(:teams)
    end
  end
end
