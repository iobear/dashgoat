
# Changelog
## [v1.4.3-Dev] - 2023-12-31
New:
 - Prometheus /metrics endpoint.
 - Metrics history via /metricshistory/\<host>\<service>/\<hours>
 - HTML service timeline [jpg](doc/history_dev.jpg)

Change:
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

## before v1.3.0
Will be add later.