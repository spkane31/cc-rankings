defmodule Rankings.Repo.Migrations.CreateTeams do
  use Ecto.Migration

  def change do
    create table(:teams) do
      add :name, :string, null: false
      add :region, :string
      add :conference, :string

      timestamps()
    end
  end
end
