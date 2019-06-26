defmodule :"Elixir.Rankings.Repo.Migrations.Add-statistics" do
  use Ecto.Migration

  def change do
    alter table(:race_instances) do
      add :average, :float, default: 0
      add :std_dev, :float, default: 0
    end

    alter table(:results) do
      add :scaled_time, :float, default: 0
    end

    alter table(:races) do
      add :is_base, :boolean, default: false
    end
  end
end
