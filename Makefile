# Compile bpf object
# compile static libbpf release
# compile userspace, linked against libbpf.a

OUTPUT_DIR ?= ./dist

LIBBPF_OBJDIR = $(abspath ./$(OUTPUT_DIR)/libbpf)
LIBBPF_OBJ := $(abspath $(LIBBPF_OBJDIR)/libbpf.a)

ARCH := $(shell uname -m | sed 's/x86_64/x86/')

CC = gcc
CLANG ?= clang
BPFTOOL ?= bpftool

CFLAGS = -g -O2 -Wall
LDFLAGS =
INCLUDES := -I$(OUTPUT_DIR)

LIBBPF_SRC := $(abspath ./libbpf/src)

CGO_CFLAGS_STATIC = "-I$(abspath $(OUTPUT_DIR))"
CGO_LDFLAGS_STATIC = "-lelf -lz $(LIBBPF_OBJ)"
CGO_EXTLDFLAGS_STATIC = '-w -extldflags "-static"'

GO_SRC := $(wildcard cmd/*.go)

default: $(OUTPUT_DIR)/weaver $(OUTPUT_DIR)/tester

$(OUTPUT_DIR)/libbpf:
	mkdir -p $@

$(LIBBPF_OBJ): $(LIBBPF_SRC) $(wildcard $(LIBBPF_SRC)/*.[ch]) | $(OUTPUT_DIR)/libbpf
	CC="$(CC)" CFLAGS="$(CFLAGS)" LD_FLAGS="$(LDFLAGS)" \
	   $(MAKE) -C $(LIBBPF_SRC) \
		BUILD_STATIC_ONLY=1 \
		OBJDIR=$(LIBBPF_OBJDIR) \
		DESTDIR=$(dir $(LIBBPF_OBJDIR)) \
		INCLUDEDIR= LIBDIR= UAPIDIR= install

$(OUTPUT_DIR)/weaver.bpf.o: weaver.bpf.c $(LIBBPF_OBJ) vmlinux.h | $(OUTPUT_DIR)
	$(CLANG) -g -O2 -target bpf -D__TARGET_ARCH_$(ARCH) $(INCLUDES) -c $(filter %.c,$^) -o $@

$(OUTPUT_DIR)/weaver: $(GO_SRC) $(OUTPUT_DIR)/weaver.bpf.o $(LIBBPF_OBJ)
	CC=$(CLANG) \
		CGO_CFLAGS=$(CGO_CFLAGS_STATIC) \
		CGO_LDFLAGS=$(CGO_LDFLAGS_STATIC) \
		go build \
		-tags netgo -ldflags $(CGO_EXTLDFLAGS_STATIC) \
		-o $@ ./cmd/...

$(OUTPUT_DIR)/tester:
	go build -o $@ ./tester/*.go

$(OUTPUT_DIR):
	mkdir -p $(OUTPUT_DIR)

clean:
	rm -rf $(OUTPUT_DIR)
