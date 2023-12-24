data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "./scripts/loader.go",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/latest"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}