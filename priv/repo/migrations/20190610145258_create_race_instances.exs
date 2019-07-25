defmodule Rankings.Repo.Migrations.CreateRaceInstances do
  use Ecto.Migration

  def change do
    create table(:race_instances) do
      add :date, :string, null: false
      add :race_id, references(:races)

      timestamps()
    end
  end
end
