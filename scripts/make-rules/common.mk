# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

#这段代码是一个Makefile脚本，它为构建和测试Go项目设置了各种变量和设置。
 #脚本以指定许可证的一些注释开始，并将SHELL设置为/bin/bash。
 #然后它包含一个公共生成文件，将common_SELF_DIR变量设置为包含当前生成文件的目录。
 #接下来，它将设置一些尚未定义的环境变量。ROOT_DIR被设置为项目根目录（比当前makefile高两级）。OUTPUT_DIR设置为ROOT_DIR的一个子目录，称为_OUTPUT，TOOLS_DIR和TMP_DIR是OUTPUT_IR的子目录。如果这些目录不存在，脚本将创建这些目录。
 #然后，脚本使用git标记设置项目的版本号，并检查git树是否脏（即是否有未提交的更改）。它还获取git提交哈希。
 #下一节将为各种变量设置一些默认值，例如COVERAGE和PLATFORMS。它还将PLATFORM变量设置为当前操作系统和体系结构，并根据是否指定了PLATFORM，将IMAGE_PLAT设置为适当的值。
 #还有各种命令行工具的设置，如FIND和XARGS，以及一些Makefile设置。
 #该脚本将githook脚本复制到.git/hooks目录中，并将CERTIFICATES变量设置为需要证书的组件列表。
 #最后，还有BLOCKER_TOOLS、CRITICAL_TOOLS和TRIVIAL_TOOLS的设置，这些设置指定了缺少工具的严重程度，并影响各种Makefile目标的成功。COMMA和SPACE变量也被定义用于列表构造。


SHELL := /bin/bash

# include the common make file
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/../.. && pwd -P))
endif
ifeq ($(origin OUTPUT_DIR),undefined)
OUTPUT_DIR := $(ROOT_DIR)/_output
$(shell mkdir -p $(OUTPUT_DIR))
endif
ifeq ($(origin TOOLS_DIR),undefined)
TOOLS_DIR := $(OUTPUT_DIR)/tools
$(shell mkdir -p $(TOOLS_DIR))
endif
ifeq ($(origin TMP_DIR),undefined)
TMP_DIR := $(OUTPUT_DIR)/tmp
$(shell mkdir -p $(TMP_DIR))
endif

# set the version number. you should not need to do this
# for the majority of scenarios.
ifeq ($(origin VERSION), undefined)
VERSION := $(shell git describe --tags --always --match='v*')
endif
# Check if the tree is dirty.  default to dirty
GIT_TREE_STATE:="dirty"
ifeq (, $(shell git status --porcelain 2>/dev/null))
	GIT_TREE_STATE="clean"
endif
GIT_COMMIT:=$(shell git rev-parse HEAD)

# Minimum test coverage
ifeq ($(origin COVERAGE),undefined)
COVERAGE := 60
endif

# The OS must be linux when building docker images
PLATFORMS ?= linux_amd64 linux_arm64
# The OS can be linux/windows/darwin when building binaries
# PLATFORMS ?= darwin_amd64 windows_amd64 linux_amd64 linux_arm64

# Set a specific PLATFORM
ifeq ($(origin PLATFORM), undefined)
	ifeq ($(origin GOOS), undefined)
		GOOS := $(shell go env GOOS)
	endif
	ifeq ($(origin GOARCH), undefined)
		GOARCH := $(shell go env GOARCH)
	endif
	PLATFORM := $(GOOS)_$(GOARCH)
	# Use linux as the default OS when building images
	IMAGE_PLAT := linux_$(GOARCH)
else
	GOOS := $(word 1, $(subst _, ,$(PLATFORM)))
	GOARCH := $(word 2, $(subst _, ,$(PLATFORM)))
	IMAGE_PLAT := $(PLATFORM)
endif

# Linux command settings
FIND := find . ! -path './third_party/*' ! -path './vendor/*'
XARGS := xargs --no-run-if-empty

# Makefile settings
ifndef V
MAKEFLAGS += --no-print-directory
endif

# Copy githook scripts when execute makefile
COPY_GITHOOK:=$(shell cp -f githooks/* .git/hooks/)

# Specify components which need certificate
ifeq ($(origin CERTIFICATES),undefined)
CERTIFICATES=iam-apiserver iam-authz-server admin
endif

# Specify tools severity, include: BLOCKER_TOOLS, CRITICAL_TOOLS, TRIVIAL_TOOLS.
# Missing BLOCKER_TOOLS can cause the CI flow execution failed, i.e. `make all` failed.
# Missing CRITICAL_TOOLS can lead to some necessary operations failed. i.e. `make release` failed.
# TRIVIAL_TOOLS are Optional tools, missing these tool have no affect.
BLOCKER_TOOLS ?= gsemver golines go-junit-report golangci-lint addlicense goimports codegen
CRITICAL_TOOLS ?= swagger mockgen gotests git-chglog github-release coscmd go-mod-outdated protoc-gen-go cfssl go-gitlint
TRIVIAL_TOOLS ?= depth go-callvis gothanks richgo rts kube-score coscli

COMMA := ,
SPACE :=
SPACE +=
