# nodelocaldns-injector

## nodelocaldns 解决什么问题

- nodelocaldns 解决的具体问题，请详见官方文档[Using NodeLocal DNSCache in Kubernetes clusters](https://kubernetes.io/docs/tasks/administer-cluster/nodelocaldns/)，需要注意：部署细节与kube-proxy的mode相关。编排业务时，dnsConfig需要与 [dnsPolicy](https://kubernetes.io/zh/docs/concepts/services-networking/dns-pod-service/#pod-%E7%9A%84-dns-%E8%AE%BE%E5%AE%9A) 配置使用
- 我们的验证是在tke上进行的，但这个文档[在 TKE 集群中使用 NodeLocal DNS Cache](https://cloud.tencent.com/document/product/457/40613) 不能工作，对比官方文档，缺少了一些东西。

## 本应用要解决的问题

- 当 nodelocaldns 部署完成后，我们不想在kubelet上进行变更，因为这会对节点部署流程有侵入。
- 在应用编排中通过dnsConfig和dnsPolicy更好控制些，而且这个机制本身也需要进行灰度验证。
- 有多个集群，只在某个特定的集群上部署了nodelocaldns，应用的编排只有一套，会应用于多个集群。于是编排和具体的集群的细节绑定了，虽然可以通过helm的模板系统引入变量、逻辑来处理，但仍引入了复杂度。

## 本应用如何解决上面的问题

- 通过 mutating admission webhook 在pod创建时，修改dnsConfig, dnsPolicy 的值。
- 如何标记应用？通过 MutatingWebhookConfiguration 支持的规则。如namespace打label, 应用编排中(Deployment等)添加特定的webhook定义的annotations。如istio 注入sidecar就是这样做的。

### 技术层面的更多细节请详见

- [准入控制](https://kubernetes.io/zh/docs/reference/access-authn-authz/extensible-admission-controllers/)
- [k8s实现的webhook的样例](https://github.com/kubernetes/kubernetes/blob/v1.13.0/test/images/webhook/main.go)
- [istio-sidecar-injector 的实现](https://github.com/istio/istio/tree/master/pkg/kube/inject)
