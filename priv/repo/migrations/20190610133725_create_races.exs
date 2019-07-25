defmodule Rankings.Repo.Migrations.CreateRaces do
  use Ecto.Migration

  def change do
    create table(:races) do
      add :name, :string, null: false
      add :course, :string, null: false
      add :distance, :integer
      add :gender, :string
      add :correction, :float

      timestamps()
    end
  end
end
