defmodule Rankings.Repo.Migrations.CreateResults do
  use Ecto.Migration

  def change do
    create table(:results) do
      add :distance, :integer, null: false
      add :rating, :float, default: 0
      add :time, :string
    end
  end
end
