local grafana = import 'grafonnet/grafana.libsonnet';

{
  grafanaDashboards:: {
    empty_dashboard: grafana.dashboard.new('Empty Test Dashboard')
    .addPanel(
        grafana.singlestat.new(
            'uptime',
            format='s',
            datasource='Prometheus',
            span=2,
            valueName='current',
        )
        .addTarget(
            grafana.prometheus.target(
                'time() - process_start_time_seconds{env="$env", job="$job", instance="$instance"}',
            )
        ), gridPos={
            x: 0,
            y: 0,
            w: 24,
            h: 3,
        }
    ),
  },
}