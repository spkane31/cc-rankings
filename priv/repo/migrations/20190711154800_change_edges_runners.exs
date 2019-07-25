defmodule Rankings.Repo.Migrations.ChangeEdgesRunners do
  use Ecto.Migration

  def change do
    alter table(:edges) do
      add :runner_id, references(:runners)
      add :difference, :float
      remove :count
      remove :total_time
      add :gender, :string
    end

    alter table(:results) do
      remove :added_to_graph
    end
  end
end
