terraform {
  required_providers {
    signalfx = {
      source  = "splunk-terraform/signalfx"
      version = "~> 6.11.0"
    }
  }
}

variable "splunkToken" {
  type = string
}

provider "signalfx" {
  auth_token = var.splunkToken
  # If your organization uses a different realm
  api_url = "https://api.eu0.signalfx.com"
  # If your organization uses a custom URL
  # custom_app_url = "https://myorg.signalfx.com"
}

resource "signalfx_dashboard_group" "MiniProject" {
  name        = "Raphael's MiniProject"
  description = "Cool dashboard group"
}

resource "signalfx_dashboard" "MiniProject0" {
  name            = "MiniProjectDashboard"
  dashboard_group = signalfx_dashboard_group.MiniProject.id

  time_range = "-30m"

  chart {
    chart_id = signalfx_time_chart.tracesCount.id
    width    = 12
    height   = 1
  }

}

resource "signalfx_time_chart" "tracesCount" {
  name = "Count of traces"

  program_text = <<-EOF
        data('traces.count').publish(label='tracesCount')
        EOF
}
