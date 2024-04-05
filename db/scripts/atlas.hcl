data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./loader.go",
  ]
}

variable "url" {
  type    = string
  default = "docker://postgres/15/dev"
}

env "dev" {
  src = data.external_schema.gorm.url
  url = var.url
  migration {
    dir = "file://../migrations?format=goose"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}