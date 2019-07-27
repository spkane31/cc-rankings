use Mix.Config

# Configure your database
config :rankings, Rankings.Repo,
  username: "postgres",
  password: "postgres",
  database: "rankings_test",
  hostname: "localhost",
  pool: Ecto.Adapters.SQL.Sandbox

# We don't run a server during test. If one is required,
# you can enable the server option below.
config :rankings, RankingsWeb.Endpoint,
  http: [port: 5432],
  server: false

# Print only warnings and errors during test
config :logger, level: :warn
