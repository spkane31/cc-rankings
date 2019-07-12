defmodule Rankings.Repo.Migrations.FloatTimeDefault do
  use Ecto.Migration

  def change do
    alter table(:results) do
      modify :time_float, :float, default: 0
    end
  end
end
