#!/usr/bin/env bash
PROJECT_DIR=$1
COMMAND=$2

if [[ ${COMMAND} != "build" && ${COMMAND} != "push" ]]; then
  echo "$COMMAND is invalid command. (Required build or push)." 1>&2
  exit 1
fi

CURRENT_BRANCH=`git rev-parse --abbrev-ref @`
IMAGE_TAG=${CURRENT_BRANCH/\//_}

# 変更があったdockerイメージを取得
if [ ${CURRENT_BRANCH} = "master" ]; then
  # 現在がmasterであれば、直前のコミットと比較
  TARGET="HEAD^ HEAD"
else
  # masterブランチ以外であれば、origin/masterの最新と比較
  TARGET="origin/master"
fi
git diff ${TARGET} --name-only  | awk '/^pkg/' | awk '{sub("pkg/", "", $0); print $0}' | awk '{print substr($0, 0, index($0, "/") -1)}' | awk  '!a[$0]++' > check.tmp

for pkgname in `cat check.tmp`; do
  if [[ ${COMMAND} == "build" ]]; then
    # test
    ${PROJECT_DIR}/bin/bazel query //... | grep "//pkg/$pkgname" | xargs ${PROJECT_DIR}/bin/bazel test --define IMAGE_TAG=${IMAGE_TAG} --local_resources=4096,2.0,1.0
    # build
    ${PROJECT_DIR}/bin/bazel query //... | awk "/^\/\/pkg\/$pkgname:$pkgname$/" | xargs ${PROJECT_DIR}/bin/bazel build --define IMAGE_TAG=${IMAGE_TAG} --local_resources=4096,2.0,1.0
  elif [[ ${COMMAND} == "push" ]]; then
    ${PROJECT_DIR}/bin/bazel query //... | awk "/^\/\/pkg\/$pkgname:container_push$/"  | xargs ${PROJECT_DIR}/bin/bazel run --define IMAGE_TAG=${IMAGE_TAG}
  fi
done

#rm check.tmp
