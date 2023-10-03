
# Changelog
## [1.4.0] - 2023-10-02
New:
 - Add DependOn to reduce alert overload, depended services only show as info if source is down

Change:
 - Improved time translation to include days
 - Update 'Change' field behavior, adding timestamp when empty

## [1.3.1] - 2023-09-18
New:
- Native ENV config
- TtlOkDelete, seconds before deleting a service with state "ok"

Change:
 - TTL behaviour, 4 config modes: Remove, PromoteOnce, PromoteOneStep, PromoteToOk (default)

## [1.3.0] - 2023-09-13
New:
 - k8s aware buddies and by extention, DNS aware buddies. DashGoat pods shoud be able to find each other, in the same namespace, provided there is a headless service called dashgoat-headless-svc.

## before 1.3.0
Will be add later
