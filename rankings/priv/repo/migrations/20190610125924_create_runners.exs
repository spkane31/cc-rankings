defmodule Rankings.Repo.Migrations.CreateRunners do
  use Ecto.Migration

  def change do
    create table(:runners) do
      add :first_name, :string, null: false
      add :last_name, :string, null: false
      add :year, :string

      timestamps()
    end
  end
end
