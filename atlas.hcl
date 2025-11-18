env "local" {
  url = getenv("POSTGRES_DSN")
  dev = "docker://postgres/18/dev"
  migration {
    dir = "file://migrations"
  }
  schema {
    src = "file://schema/schema.sql"
  }
}

env "test" {
  url = getenv("TEST_POSTGRES_DSN")
  dev = "docker://postgres/18/dev"
  migration {
    dir = "file://migrations"
  }
  schema {
    src = "file://schema/schema.sql"
  }
}