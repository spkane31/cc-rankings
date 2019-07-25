defmodule Rankings.Repo.Migrations.AddTimeFloat do
  use Ecto.Migration

  def change do
    alter table(:results) do
      add :time_float, :float
    end

  end
end
