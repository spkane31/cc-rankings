defmodule Rankings.Repo.Migrations.AddAvgStdDevRaces do
  use Ecto.Migration

  def change do
    alter table(:races) do
      add :average, :float, default: 0
      add :std_dev, :float, default: 0
      remove :correction
      add :correction_avg, :float, default: 0
      add :correction_graph, :float, default: 0
    end
  end
end
