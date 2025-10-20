env "local" {
  url = env("POSTGRES_DSN")
  dev = "docker://postgres/17/dev"
  migration {
    dir = "file://migrations"
  }
}