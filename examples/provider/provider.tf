terraform {
  required_providers {
    lumos = {
      source  = "teamlumos/lumos"
      version = "0.10.4"
    }
  }
}

provider "lumos" {
  server_url = "..." # Optional
}