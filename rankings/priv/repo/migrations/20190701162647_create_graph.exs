defmodule Rankings.Repo.Migrations.CreateGraph do
  use Ecto.Migration

  def change do
    alter table(:race_instances) do
      add :valid, :boolean
      remove :date
      add :date, :date
    end

    create table(:edges) do
      add :from_race_id, references(:races)
      add :to_race_id, references(:races)
      add :total_time, :float
      add :count, :int

      timestamps()
    end

    alter table(:results) do
      add :added_to_graph, :boolean, default: false
      add :date, :date
      add :gender, :string
      add :place, :int
    end

    alter table(:runners) do
      add :elo_rating, :float, default: 1400
      add :speed_rating, :float, default: 0.0
    end

    create table(:elo_ratings) do
      add :runner_id, references(:runners)
      add :rating, :float
      add :most_recent, :boolean
      add :date, :date
    end
  end
end
