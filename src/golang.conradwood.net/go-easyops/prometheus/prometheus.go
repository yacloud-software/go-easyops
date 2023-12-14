/*
wrapper around prometheus SDK/Client.

the prometheus SDK API has changed occassionally in the past. This package, different to the
SDK also implements a timeout. That is, if a metric is not set for some time it will be removed from
the exported metrics. (The SDK exports the last set value forever).
*/
package prometheus
