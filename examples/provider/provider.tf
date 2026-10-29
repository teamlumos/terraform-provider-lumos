terraform {
  required_providers {
    lumos = {
      source  = "teamlumos/lumos"
      version = "0.11.1"
    }
  }
}

provider "lumos" {
  server_url = "..." # Optional
}