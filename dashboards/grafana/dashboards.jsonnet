local grafana = import 'grafonnet/grafana.libsonnet';

{
  grafanaDashboards:: {
    empty_dashboard: grafana.dashboard.new('Empty Test Dashboard')
    .addPanel(
        grafana.graphPanel.new(
            title= 'TestDataGraph',
            datasource='TestDataDB',
        ),
        gridPos={
          x: 0,
          y: 0,
          w: 12,
          h: 9,
        }
    ),
  },
}