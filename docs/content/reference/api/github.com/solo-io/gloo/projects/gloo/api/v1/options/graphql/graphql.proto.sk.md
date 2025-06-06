
---
title: "Graphql"
weight: 5
---

<!-- Code generated by solo-kit. DO NOT EDIT. -->


### Package: `graphql.options.gloo.solo.io` 
**Types:**


- [ServiceSpec](#servicespec)
- [Endpoint](#endpoint)
  



**Source File: [github.com/solo-io/gloo/projects/gloo/api/v1/options/graphql/graphql.proto](https://github.com/solo-io/gloo/blob/main/projects/gloo/api/v1/options/graphql/graphql.proto)**





---
### ServiceSpec

 
Deprecated: The GraphQL feature of Gloo Gateway will be removed in a future release.
Only supported in enterprise with the GraphQL addon. 
This is the service spec describing GraphQL upstreams. This will usually be filled
automatically via function discovery (if the upstream supports introspection).
If your upstream service is a GraphQL service, use this service spec (an empty
spec is fine).

```yaml
"endpoint": .graphql.options.gloo.solo.io.ServiceSpec.Endpoint

```

| Field | Type | Description |
| ----- | ---- | ----------- | 
| `endpoint` | [.graphql.options.gloo.solo.io.ServiceSpec.Endpoint](../graphql.proto.sk/#endpoint) | Endpoint provides the endpoint information, and how to call the GraphQL Server. This endpoint must be called via HTTP POST sending form data as mentioned in [the GraphQL Docs](https://graphql.org/learn/serving-over-http/#post-request). |




---
### Endpoint

 
Describes a GraphQL Endpoint information

```yaml
"url": string

```

| Field | Type | Description |
| ----- | ---- | ----------- | 
| `url` | `string` | The url for the graphql endpoint. Automation via Discovery only supports `http://<host>/graphql` ie: http://myurl.com/graphql. |





<!-- Start of HubSpot Embed Code -->
<script type="text/javascript" id="hs-script-loader" async defer src="//js.hs-scripts.com/5130874.js"></script>
<!-- End of HubSpot Embed Code -->
