FROM alpine:3.8

ENV KUSTOMIZE_VER 1.0.11
ENV KUBECTL_VER 1.6.4

RUN set -x && \
    apk update && apk upgrade && \
    apk add --no-cache curl bash git ca-certificates openssh-client && \

    chmod +x /usr/bin/git && \
    /usr/bin/git version && \

    curl -L https://github.com/kubernetes-sigs/kustomize/releases/download/v${KUSTOMIZE_VER}/kustomize_${KUSTOMIZE_VER}_linux_amd64 -o /usr/bin/kustomize && \
    chmod +x /usr/bin/kustomize && \

    curl -L https://storage.googleapis.com/kubernetes-release/release/v${KUBECTL_VER}/bin/linux/amd64/kubectl -o /usr/bin/kubectl && \
    chmod +x /usr/bin/kubectl
