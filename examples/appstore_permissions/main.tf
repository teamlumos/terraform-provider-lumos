terraform {
  required_providers {
    lumos = {
      source = "lumos.com/tf/lumos"
    }
  }
}

provider "lumos" {
  api_token = "lsk_uPOw7AtKDHE5QhQMOSxPrw7gAH_J1HeNgidlgPx4kp7Q"
  base_url  = "http://localhost:8000"
}

// TODO Replace UUIDs with your local UUIDs For Testing
locals {
  terraform_test_app_id = "<FILL IN ID>" # Find UUID like so: https://www.loom.com/share/770d0a4900e7430f8ec8e552819eafdc?sid=61920d66-1a98-4b64-b8b1-844d422d6930
  group_1_id            = "<FILL IN ID>" # Find UUID like so: https://www.loom.com/share/b313034f7ad84d83b0195fe1b0f370f8?sid=ff1220e9-874c-4021-b915-d9be509ee1c2
  group_2_id            = "<FILL IN ID>"
  group_3_id            = "<FILL IN ID>"
  user_1_id             = "<FILL IN ID>" # Find UUID like so: https://www.loom.com/share/1413b71a5f8f4692a01f986a95b3e6fc?sid=5def1c2a-6cae-4c9b-b2d2-b944131ab5b8
  user_2_id             = "<FILL IN ID>"
  user_3_id             = "<FILL IN ID>"
}


resource "lumos_pre_approval_rule" "pre_approval_one" {
  app_id        = "b4bb7f9c-e287-4cda-878d-8638a85206f9"
  justification = "CHANGEtesting!!!"

  preapproved_groups = [
    { id = "93b9b293-592b-f24e0-5a77-2b49dd9942f6" }
  ]
  preapproved_permissions = [
    {
      id = "77974394-b239-1d96-801f-c77d83e7d92f"
    }
  ]
}
