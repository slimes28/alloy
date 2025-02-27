---
canonical: https://grafana.com/docs/alloy/latest/set-up/migrate/from-prometheus/
aliases:
  - ../../tasks/migrate/from-prometheus/ # /docs/alloy/latest/tasks/migrate/from-prometheus/
description: Learn how to migrate from Prometheus to Grafana Alloy
menuTitle: Migrate from Prometheus
title: Migrate from Prometheus to Grafana Alloy
weight: 250
---

# Migrate from Prometheus to {{% param "FULL_PRODUCT_NAME" %}}

The built-in {{< param "PRODUCT_NAME" >}} convert command can migrate your [Prometheus][] configuration to an {{< param "PRODUCT_NAME" >}} configuration.

This topic describes how to:

* Convert a Prometheus configuration to an {{< param "PRODUCT_NAME" >}} configuration.
* Run a Prometheus configuration natively using {{< param "PRODUCT_NAME" >}}.

## Components used in this topic

* [`prometheus.scrape`][prometheus.scrape]
* [`prometheus.remote_write`][prometheus.remote_write]

## Before you begin

* You must have a Prometheus configuration.
* You must have a set of Prometheus applications ready to push telemetry data to {{< param "PRODUCT_NAME" >}}.
* You must be familiar with the concept of [Components][] in {{< param "PRODUCT_NAME" >}}.

## Convert a Prometheus configuration

To fully migrate your configuration from [Prometheus] to {{< param "PRODUCT_NAME" >}}, you must convert your Prometheus configuration into an {{< param "PRODUCT_NAME" >}} configuration.
This conversion allows you to take full advantage of the many additional features available in {{< param "PRODUCT_NAME" >}}.

In this task, you use the [convert][] CLI command to output an {{< param "PRODUCT_NAME" >}} configuration from a Prometheus configuration.

1. Open a terminal window and run the following command.

   ```shell
   alloy convert --source-format=prometheus --output=<OUTPUT_CONFIG_PATH> <INPUT_CONFIG_PATH>
   ```

   Replace the following:

   * _`<INPUT_CONFIG_PATH>`_: The full path to the Prometheus configuration.
   * _`<OUTPUT_CONFIG_PATH>`_: The full path to output the {{< param "PRODUCT_NAME" >}} configuration.

1. [Run][] {{< param "PRODUCT_NAME" >}} using the new {{< param "PRODUCT_NAME" >}} configuration from _`<OUTPUT_CONFIG_PATH>`_:

### Debugging

1. If the `convert` command can't convert a Prometheus configuration, diagnostic information is sent to `stderr`.\
   You can bypass any non-critical issues and output the {{< param "PRODUCT_NAME" >}} configuration using a best-effort conversion by including the `--bypass-errors` flag.

   {{< admonition type="caution" >}}
   If you bypass the errors, the behavior of the converted configuration may not match the original Prometheus configuration.
   Make sure you fully test the converted configuration before using it in a production environment.
   {{< /admonition >}}

   ```shell
   alloy convert --source-format=prometheus --bypass-errors --output=<OUTPUT_CONFIG_PATH> <INPUT_CONFIG_PATH>
   ```

   Replace the following:

   * _`<INPUT_CONFIG_PATH>`_: The full path to the Prometheus configuration.
   * _`<OUTPUT_CONFIG_PATH>`_: The full path to output the {{< param "PRODUCT_NAME" >}} configuration.

1. You can also output a diagnostic report by including the `--report` flag.

   ```shell
   alloy convert --source-format=prometheus --report=<OUTPUT_REPORT_PATH> --output=<OUTPUT_CONFIG_PATH> <INPUT_CONFIG_PATH>
   ```

   Replace the following:

   * _`<INPUT_CONFIG_PATH>`_: The full path to the Prometheus configuration.
   * _`<OUTPUT_CONFIG_PATH>`_: The full path to output the {{< param "PRODUCT_NAME" >}} configuration.
   * _`<OUTPUT_REPORT_PATH>`_: The output path for the report.

    Using the [example][] Prometheus configuration below, the diagnostic report provides the following information:

    ```plaintext
    (Info) Converted scrape_configs job_name "prometheus" into...
      A prometheus.scrape.prometheus component
    (Info) Converted 1 remote_write[s] "grafana-cloud" into...
      A prometheus.remote_write.default component
    ```

## Run a Prometheus configuration

If you're not ready to completely switch to an {{< param "PRODUCT_NAME" >}} configuration, you can run {{< param "PRODUCT_NAME" >}} using your Prometheus configuration.
The `--config.format=prometheus` flag tells {{< param "PRODUCT_NAME" >}} to convert your Prometheus configuration to an {{< param "PRODUCT_NAME" >}} configuration and load it directly without saving the new configuration.
This allows you to try {{< param "PRODUCT_NAME" >}} without modifying your Prometheus configuration infrastructure.

In this task, you use the [run][] CLI command to run {{< param "PRODUCT_NAME" >}} using a Prometheus configuration.

[Run][run] {{< param "PRODUCT_NAME" >}} and include the command line flag `--config.format=prometheus`.
Your configuration file must be a valid Prometheus configuration file rather than an {{< param "PRODUCT_NAME" >}} configuration file.

### Debug

1. You can follow the convert CLI command [debugging][] instructions to generate a diagnostic report.

1. Refer to [Debug {{< param "PRODUCT_NAME" >}}][DebuggingUI] for more information about a running {{< param "PRODUCT_NAME" >}}.

1. If your Prometheus configuration can't be converted and loaded directly into {{< param "PRODUCT_NAME" >}}, diagnostic information is sent to `stderr`.
   You can bypass any non-critical issues and start {{< param "PRODUCT_NAME" >}} by including the `--config.bypass-conversion-errors` flag in addition to `--config.format=prometheus`.

   {{< admonition type="caution" >}}
   If you bypass the errors, the behavior of the converted configuration may not match the original Prometheus configuration.
   Don't use this flag in a production environment.
   {{< /admonition >}}

## Example

This example demonstrates converting a Prometheus configuration file to an {{< param "PRODUCT_NAME" >}} configuration file.

The following Prometheus configuration file provides the input for the conversion.

```yaml
global:
  scrape_timeout:    45s

scrape_configs:
  - job_name: "prometheus"
    static_configs:
      - targets: ["localhost:12345"]

remote_write:
  - name: "grafana-cloud"
    url: "https://prometheus-us-central1.grafana.net/api/prom/push"
    basic_auth:
      username: <USERNAME>
      password: <PASSWORD>
```

The convert command takes the YAML file as input and outputs a [{{< param "PRODUCT_NAME" >}} configuration][configuration]] file.

```shell
alloy convert --source-format=prometheus --output=<OUTPUT_CONFIG_PATH> <INPUT_CONFIG_PATH>
```

Replace the following:

* _`<INPUT_CONFIG_PATH>`_: The full path to the Prometheus configuration.
* _`<OUTPUT_CONFIG_PATH>`_: The full path to output the {{< param "PRODUCT_NAME" >}} configuration.

The new {{< param "PRODUCT_NAME" >}} configuration file looks like this:

```alloy
prometheus.scrape "prometheus" {
  targets = [{
    __address__ = "localhost:12345",
  }]
  forward_to     = [prometheus.remote_write.default.receiver]
  job_name       = "prometheus"
  scrape_timeout = "45s"
}

prometheus.remote_write "default" {
  endpoint {
    name = "grafana-cloud"
    url  = "https://prometheus-us-central1.grafana.net/api/prom/push"

    basic_auth {
      username = "USERNAME"
      password = "PASSWORD"
    }

    queue_config {
      capacity             = 2500
      max_shards           = 200
      max_samples_per_send = 500
    }

    metadata_config {
      max_samples_per_send = 500
    }
  }
}
```

## Limitations

Configuration conversion is done on a best-effort basis. {{< param "PRODUCT_NAME" >}} issues warnings or errors where the conversion can't be performed.

After the configuration is converted, review the {{< param "PRODUCT_NAME" >}} configuration file created and verify that it's correct before starting to use it in a production environment.

The following list is specific to the convert command and not {{< param "PRODUCT_NAME" >}}:

* The following configurations aren't available for conversion to {{< param "PRODUCT_NAME" >}}: `rule_files`, `alerting`, `remote_read`, `storage`, and `tracing`.
  Any additional unsupported features are returned as errors during conversion.
* Check if you are using any extra command line arguments with Prometheus that aren't present in your configuration file. For example, `--web.listen-address`.
* Meta-monitoring metrics exposed by {{< param "PRODUCT_NAME" >}} usually match Prometheus meta-monitoring metrics but uses a different name.
  Make sure that you use the new metric names, for example, in your alerts and dashboards queries.
* The logs produced by {{< param "PRODUCT_NAME" >}} differ from those produced by Prometheus.
* {{< param "PRODUCT_NAME" >}} exposes the {{< param "PRODUCT_NAME" >}} [UI][].

[Prometheus]: https://prometheus.io/docs/prometheus/latest/configuration/configuration/
[debugging]: #debugging
[example]: #example
[prometheus.scrape]: ../../../reference/components/prometheus/prometheus.scrape/
[prometheus.remote_write]: ../../../reference/components/prometheus/prometheus.remote_write/
[Components]: ../../../get-started/components/
[convert]: ../../../reference/cli/convert/
[run]: ../../../reference/cli/run/
[DebuggingUI]: ../../../troubleshoot/debug/
[configuration]: ../../../get-started/configuration-syntax/
[UI]: ../../../troubleshoot/debug/#alloy-ui
