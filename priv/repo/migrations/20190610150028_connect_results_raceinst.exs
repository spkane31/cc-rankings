defmodule Rankings.Repo.Migrations.ConnectResultsRaceinst do
  use Ecto.Migration

  def change do
    alter table(:results) do
      add :race_instance_id, references(:race_instances)
    end
  end
end
