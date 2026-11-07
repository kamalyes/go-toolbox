#!/usr/bin/env bash
#
# @Author: kamalyes 501893067@qq.com
# @Description: 把父模块与全部嵌套子模块按同一个版本号打 tag
#
# Go 的子模块 tag 带目录前缀，所以 `v0.20.0` 与 `pkg/protox/v0.20.0` 是两个不同的
# tag 名，可以落在同一个 commit 上 —— 版本号一致因此是做得到的，代价只是 tag 数量
#
# [EN] Go namespaces nested-module tags by directory, so `v0.20.0` and
# `pkg/protox/v0.20.0` are distinct tag names and can sit on the same commit
# Keeping one version number across every module is therefore possible; the only
# cost is the number of tags
#
# 顺序上有一个硬约束：有些子模块 `require` 父模块，它们的 go.mod 必须在打 tag 之前
# 就已经写上新版本号，否则下游取到的子模块会指着一个旧的父版本
# 所以流程是「先改 require → 再提交 → 最后在这一个 commit 上打齐所有 tag」
#
# [EN] One hard ordering constraint: submodules that require the parent must already
# carry the new version in their go.mod before tagging, or consumers get a submodule
# pointing at a stale parent, hence rewrite -> commit -> tag the single commit
#
# Copyright (c) 2026 by kamalyes, All Rights Reserved.

set -euo pipefail

VERSION="${1:-}"
PARENT_MODULE="github.com/kamalyes/go-toolbox"

usage() {
	cat >&2 <<-EOF
		用法: scripts/release.sh vX.Y.Z [--push]

		  vX.Y.Z   本次发布的版本号，父模块与所有子模块共用这一个号
		  --push   打完 tag 后直接推送；缺省只在本地打好并打印推送命令

		[EN] usage: scripts/release.sh vX.Y.Z [--push]
	EOF
	exit 2
}

[[ -n "$VERSION" ]] || usage
[[ "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]] || {
	echo "✗ 版本号格式应为 vX.Y.Z，收到: $VERSION" >&2
	exit 1
}

cd "$(dirname "$0")/.."
ROOT="$(pwd)"

# ---- 前置检查 | Preflight ----

[[ -z "$(git status --porcelain)" ]] || {
	echo "✗ 工作区不干净，先提交或清理再发布" >&2
	git status --short >&2
	exit 1
}

git rev-parse "$VERSION" >/dev/null 2>&1 && {
	echo "✗ tag $VERSION 已存在" >&2
	exit 1
}

# 主版本号 >= 2 时 module path 要带 /v2 后缀，tag 形态也随之改变，本脚本不覆盖
# [EN] Major >= 2 requires a /vN suffix in module paths and different tag shapes
[[ "${VERSION%%.*}" == "v0" || "${VERSION%%.*}" == "v1" ]] || {
	echo "✗ $VERSION 是 v2+，module path 需要加 /vN 后缀，本脚本不处理" >&2
	exit 1
}

# ---- 模块清单以 go.work 为唯一来源，新增子模块不必改本脚本 ----
# [EN] go.work is the single source of truth for the module list
# 用 while read 而不是 mapfile：mapfile 是 bash 4+ 的内建，而 macOS 自带的是 bash 3.2
# [EN] while-read instead of mapfile: mapfile needs bash 4+, macOS ships bash 3.2
DIRS=()
while IFS= read -r line; do
	DIRS+=("$line")
done < <(go work edit -json | python3 -c '
import json,sys
for m in json.load(sys.stdin)["Use"]:
    p = m["DiskPath"]
    print(p[2:] if p.startswith("./") else ("" if p == "." else p))
')

echo "本次发布 ${VERSION}，共 $(( ${#DIRS[@]} )) 个模块："
for d in "${DIRS[@]}"; do
	echo "  • ${d:-（父模块）}"
done
echo

# ---- 1. 把依赖父模块的子模块的 require 改到新版本 ----
# 子模块各自带 `replace ../` 只对本地构建生效，下游会忽略，所以 require 的版本号
# 必须是真实可下载的 tag 本步只改 go.mod，不动 go.sum ——
# replace 指向本地目录时父模块不进 go.sum，没有哈希要更新
#
# [EN] Only go.mod changes: with a local replace the parent never enters go.sum
BUMPED=()
for d in "${DIRS[@]}"; do
	[[ -n "$d" ]] || continue
	if grep -q "^[[:space:]]*${PARENT_MODULE} v" "$d/go.mod"; then
		(cd "$d" && go mod edit -require="${PARENT_MODULE}@${VERSION}")
		BUMPED+=("$d")
	fi
done

if (( ${#BUMPED[@]} )); then
	echo "已把以下子模块的父模块 require 改到 ${VERSION}："
	# bash 3.2 下 set -u 展开空数组会报 unbound，所以统一用 ${a[@]+...} 形式
	# [EN] Expanding an empty array under set -u errors on bash 3.2, hence this form
	printf '  • %s\n' ${BUMPED[@]+"${BUMPED[@]}"}
else
	echo "没有子模块 require 父模块，无需改写"
fi
echo

# ---- 2. 改完先验证，再提交 ----
# 此刻 require 指向的 tag 还不存在，本地靠 replace 才编得过；
# 所以这里只能 build/test，不能跑 go mod tidy（tidy 会去网络找那个还没有的 tag）
#
# [EN] The required tag does not exist yet and only the local replace makes this
# build, so do not run `go mod tidy` here: tidy would go to the network for it
echo "验证 ${#DIRS[@]} 个模块……"
for d in "${DIRS[@]}"; do
	label="${d:-（父模块）}"
	printf '  %-26s' "$label"
	if (cd "$ROOT/${d:-.}" && GOWORK=off go build ./... >/dev/null 2>&1); then
		echo "build ok"
	else
		echo "build 失败 ✗"
		(cd "$ROOT/${d:-.}" && GOWORK=off go build ./... 2>&1 | head -5)
		echo "✗ 已中止，go.mod 的改动留在工作区供检查" >&2
		exit 1
	fi
done
echo

if (( ${#BUMPED[@]} )); then
	for d in ${BUMPED[@]+"${BUMPED[@]}"}; do git add "$d/go.mod"; done
	git commit -m "🔖 release: ${VERSION}，父模块与 $(( ${#DIRS[@]} - 1 )) 个子模块共用同一版本号"
	echo "已提交 require 改动"
fi

# ---- 3. 在同一个 commit 上打齐所有 tag ----
# 父模块 `vX.Y.Z`，子模块 `<相对目录>/vX.Y.Z`，互不冲突
# [EN] Parent gets vX.Y.Z, each submodule gets <dir>/vX.Y.Z on the same commit
SHA="$(git rev-parse HEAD)"
TAGS=()
for d in "${DIRS[@]}"; do
	if [[ -z "$d" ]]; then
		TAGS+=("$VERSION")
	else
		TAGS+=("$d/$VERSION")
	fi
done

for t in "${TAGS[@]}"; do
	git tag -a "$t" -m "$t" "$SHA"
	echo "  ✓ $t"
done
echo
echo "已在 $SHA 上打好 ${#TAGS[@]} 个 tag"

# ---- 4. 推送：对外动作，默认不自动做 ----
# [EN] Pushing is outward-facing, so it is opt-in
if [[ "${2:-}" == "--push" ]]; then
	git push origin HEAD
	git push origin "${TAGS[@]}"
	echo "已推送 commit 与全部 tag"
else
	echo
	echo "本地已就绪，推送命令："
	echo "  git push origin HEAD"
	echo "  git push origin ${TAGS[*]}"
fi
