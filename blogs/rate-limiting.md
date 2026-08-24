---
title: Rate limiting outbound traffic
date: 2026-05-13
summary: Explicit control over how often and how many requests leave your system.
---

# Rate limiting outbound traffic

Whether you are shipping a side project or B2B SaaS, calling third-party
APIs is routine. It is also easy to get wrong: retries stack, concurrency
spikes, and small mistakes in scheduling can create noisy traffic or
flaky behavior.

If you care about scale, many users, or simply doing the right thing
upstream, you want explicit control over how often and how many requests
leave your system. Most public APIs publish expectations — requests per
minute, payload limits, concurrent connections — often tied to your
client id.

## The naive fix

Exponential backoff, or a counter for “requests per minute,” looks
enough until it is not. Real systems add concurrency, overlapping async
work, uneven latency, and retries — and you still need observability
when something backs up or fails.

You can build a time-based queue with concurrency caps yourself. Or you
can use a well-traveled option like [bottleneck](https://github.com/SGrondin/bottleneck):
configure limits that match the provider’s docs, pass a function that
returns a promise, and let the limiter schedule execution. You get
predictable throughput instead of ad hoc timers.

## Typical knobs

Names vary slightly by version; check the docs for yours.

- **maxConcurrent** — cap parallel in-flight calls
- **minTime** — minimum spacing between starts (smooths bursts)
- **highWater** — how deep the queue may grow before you apply a
  strategy (for example, rejecting new work instead of unbounded
  memory growth)

## Shape

Wrap `schedule` in try / catch so you can map provider errors, transform
responses, and return clean results to your own clients.

```js
import Bottleneck from "bottleneck";

const limiter = new Bottleneck({
  maxConcurrent: 5,
  minTime: 200,
  highWater: 100,
  strategy: Bottleneck.strategy.OVERFLOW,
});

export async function callProvider(payload) {
  try {
    const raw = await limiter.schedule(() => fetchFromProvider(payload));
    return transform(raw);
  } catch (err) {
    throw mapProviderError(err);
  }
}
```

For distributed servers, bottleneck also supports shared scheduler state
via Redis. The docs mention that all the features are supported in
Clustering mode.

## Errors and hooks

Use the limiter’s lifecycle hooks / events — around failures and retries,
depending on version — so logging, metrics, and user-facing messages
stay consistent. That keeps “what happened” visible without sprinkling
one-off handlers everywhere.
