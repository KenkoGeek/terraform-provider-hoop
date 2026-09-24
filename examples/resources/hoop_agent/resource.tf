resource "hoop_agent" "example" {
  name = "terraform-agent"
  mode = "standard"
}

resource "hoop_connection" "example" {
  name     = "terraform-agent-bash"
  type     = "custom"
  agent_id = hoop_agent.example.id

  command = ["/bin/bash"]

  access_mode_runbooks = "enabled"
  access_mode_exec     = "enabled"
  access_mode_connect  = "disabled"
  access_schema        = "disabled"
}

output "hoop_agent_token" {
  value     = hoop_agent.example.token
  sensitive = true
}
