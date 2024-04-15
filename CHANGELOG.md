# Changelog
## [v1.7.dev] - 2024-04-15
Add:
 - Config element logformat <txt/json>
 - Heartbeat via HTTP GET (uri) /heartbeat/:urikey/:host/:nextupdatesec/:tags
 - Alertmanager Webhook
 - ack (acknowledge) to ServiceState
 - Listen for Azure functions port

Fix:
 - Missing version number on dashboard
 - webpath prefix

## [v1.6.0] - 2024-03-18
Add:
 - Slog logging
 - PagerDuty push
 - Report state change

Fix:
 - Add missing buddynsconfig option to config file

Change:
 - Upgrade Go + dependencies
 - config element nsconfig is now called buddynsconfig

Deprecated:
 - weblog config element

## [v1.5.4] - 2024-03-06
Fix:
 - Logic error in relations to DisableMetrics
 - Wrong state of timeoutprobe

Add:
 - Basic Prometheus docs

## [v1.5.3] - 2024-03-04
Change:
 - Change versioning to follow git tags

## [v1.5.2] - 2024-02-26
Fix:
 - JS dashboard update error
 - JS lowercase undefined error

Change:
 - Made a common module, moved ServiceState to common for other tools

## [v1.5.1] - 2024-02-04
New:
 - Embed webfiles to single binary

## [v1.5.0] - 2024-02-02
New:
 - Prometheus /metrics endpoint.
 - Metrics history with Prometheus backend, via /metricshistory/\<host>\<service>/\<hours>
 - Metrics timeline for every service entry.

Change:
 - New table layout with \<div\>
 - Optimize DOM update.

## [v1.4.2] - 2023-12-14
Fix:
 - Check for nextupdatesec, always beeing 19sec.

Change:
 - ttlHousekeeping() to more readable code.

## [v1.4.1] - 2023-11-29
New:
 - Add status favicon.

## [v1.4.0] - 2023-10-02
New:
 - Add DependOn to reduce alert overload, depended services only show as info if source is down.

Change:
 - Improved time translation to include days.
 - Update 'Change' field behavior, adding timestamp when empty.

## [v1.3.1] - 2023-09-18
New:
- Native ENV config, instead of translation via Dockerfile.
- TtlOkDelete, seconds before deleting a service with state "ok".

Change:
 - TTL behaviour, 4 config modes: Remove, PromoteOnce, PromoteOneStep, PromoteToOk (default).

## [v1.3.0] - 2023-09-13
New:
 - k8s aware buddies and by extention, DNS aware buddies. DashGoat pods shoud be able to find each other, in the same namespace, provided there is a headless service called dashgoat-headless-svc.

## [v1.2.9] - 2023-03-01
Change:
 - Code now under MPL 2.0
 - Default port is now :2000

## [v1.2.8] - 2023-02-27
Change:
 - Config code cleanup

## [v1.2.6] - 2022-07-02
Change:
 - TTL feature, add option to only change state to ok, instead of beeing deleted

## [v1.2.0] - 2022-03-18
New:
 - Buddy cluster feature, share state with your buddy instance over HTTP.

## before v1.x.0
Will be added later.
