# Changelog
## [1.7.14] - 2025-05-15
Add:
 - Search parameter in your GET request

## [1.7.13] - 2025-04-17
Change:
 - Update dependencies

## [1.7.12] - 2025-03-13
Change:
 - Go update, plus dependencies

## [1.7.11] - 2024-12-19
Change:
 - Go Echo update

## [v1.7.10] - 2024-12-13
Change:
 - Go update dependencies

## [v1.7.9] - 2024-07-26
Change:
 - Improved error handling in buddy config
 - Updated dependencies

Fix:
 - Resolved issue with incorrect buddy timestamp

## [v1.7.8] - 2024-06-09
Fix:
 - Buddy config unintended behavior, not ignoring own instance name from config file
 - Buddy config unintended behavior, not using dashGoat 'updatekey' when none given for buddy.
 - Some spelling

Change:
 - Upgrade Go to v1.22.4

## [v1.7.7] - 2024-06-03
Change:
 - Upgrade Go dependencies
 - JSON return from API

## [v1.7.6] - 2024-05-29
Change:
 - buddyDownStatusMsg to buddyDownStatus in config

Add:
 - probeTimeoutStatus in config

Fix:
 - Change time, after ttl timeout

## [v1.7.5] - 2024-05-26
Fix:
 - Pagerduty timeout

## [v1.7.4] - 2024-05-16
Fix:
 - Heartbeat change behavior
 - Tests to check for change behavior
 - Prometheus config error

Add:
 - Pagerduty retry

## [v1.7.3] - 2024-05-01
Change:
 - Low Pagerduty Timeout
 - Upgrade Echo dependencies

## [v1.7.2] - 2024-04-23
Fix:
 - Buddy state inconsistencies
 - Buddy state race problems, making the app reset dataset
 - Change timestamp not updated correctly

## [v1.7.1] - 2024-04-22
Change:
 - Upgrade Go dependencies

## [v1.7.0] - 2024-04-18
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
