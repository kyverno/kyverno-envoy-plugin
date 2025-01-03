# kyverno-envoy-plugin

A plugin to enforce kyverno policies with Envoy. This plugin allows applying Kyverno policies to APIs managed by Envoy.

## Overview 

[Envoy](https://www.envoyproxy.io/docs/envoy/latest/intro/what_is_envoy) is a L7 proxy and communication bus designed for large modern service oriented architectures . Envoy (v1.7.0+) supports an External Authorization filter which calls an authorization service to check if the incoming request is authorized or not. [External Authorization filter](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/security/ext_authz_filter.html) feature will help us to make a decision based on Kyverno policies . 

**WARNING: ⚠️ Kyverno-envoy-plugin is in development stage.**


## 📙 Documentation

Kyverno Envoy plugin installation and reference documents are available [here](https://kyverno.github.io/kyverno-envoy-plugin)

👉 **[Quick Start](https://kyverno.github.io/kyverno-envoy-plugin/latest/quick-start/)**

👉 **[Installation](https://kyverno.github.io/kyverno-envoy-plugin/latest/quick-start/authz-server/)**

## RoadMap

For detailed information on our planned features and upcoming updates, please [view our Roadmap](./ROADMAP.md).

## 🙋‍♂️ Getting Help

We are here to help!

👉 For feature requests and bugs, file an [issue](https://github.com/kyverno/kyverno-envoy-plugin/issues).

👉 For discussions or questions, join the [Kyverno Slack channel](https://slack.k8s.io/#kyverno).

👉 To get notified on updates ⭐️ [star this repository](https://github.com/kyverno/kyverno-envoy-plugin/stargazers).

## ➕ Contributing

Thanks for your interest in contributing to Kyverno! Here are some steps to help get you started:

✔ Look through the [good first issues](https://github.com/kyverno/kyverno-envoy-plugin/labels/good%20first%20issue) list. Add a comment with `/assign` to request the assignment of the issue.

✔ Check out the Kyverno [Community page](https://kyverno.io/community/) for other ways to get involved.

## License

Copyright 2023, the Kyverno project. All rights reserved. kyverno-envoy-plugin is licensed under the [Apache License 2.0](LICENSE).
