defmodule Rankings.Repo.Migrations.ConnectRunnerResult do
  use Ecto.Migration

  def change do
    alter table(:results) do
      add :runner_id, references(:runners)
    end
  end
end
