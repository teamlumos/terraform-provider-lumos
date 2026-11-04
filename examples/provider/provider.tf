terraform {
  required_providers {
    lumos = {
      source  = "teamlumos/lumos"
      version = "~> 1.0"
    }
  }
}

provider "lumos" {
  server_url = "..." # Optional
}