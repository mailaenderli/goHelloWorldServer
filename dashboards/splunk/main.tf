provider "signalfx" {
  auth_token = "<+secrets.getValue("project.splunkToken")>"
  # If your organization uses a different realm
  # api_url = "https://api.us2.signalfx.com"
  # If your organization uses a custom URL
  # custom_app_url = "https://myorg.signalfx.com"
}

resource "signalfx_dashboard" "MiniProject0" {
  name            = "MiniProjectDashboard"
  dashboard_group = signalfx_dashboard_group.MiniProject.id

  time_range = "-30m"
}
