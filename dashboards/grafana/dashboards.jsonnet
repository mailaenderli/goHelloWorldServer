local grafana = import 'grafonnet/grafana.libsonnet';

{
  grafanaDashboards:: {
    empty_dashboard: grafana.dashboard.new('Empty Test Dashboard')
    .addPanel(
        grafana.singlestat.new(
            'TestData',
            format='s',
            datasource='TestDataDB',
            span=2,
            valueName='current',
        )
    ),
  },
}