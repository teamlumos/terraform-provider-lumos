terraform {
  required_providers {
    lumos = {
      source  = "teamlumos/lumos"
      version = "0.10.5"
    }
  }
}

provider "lumos" {
  server_url = "..." # Optional
}