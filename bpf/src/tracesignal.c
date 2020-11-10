#include <argp.h>
#include <unistd.h>
#include "tracesignal.h"
#include "tracesignal.skel.h"

static struct env {
	pid_t pid;
    bool verbose;
} env = {};

static const struct argp_option opts[] = {
    { "pid", 'p', "PID", 0, "Process ID to trace"},
    { "verbose", 'v', NULL, 0, "Verbose debug output" },
    {},
};

static error_t parse_arg(int key, char *arg, struct argp_state *state)
{
    static int pos_args;
    long int pid;

    switch (key) {
        case 'p':
            errno = 0;
            pid = strtol(arg, NULL, 10);
            if (errno || pid <= 0) {
                fprintf(stderr, "INVALID PID: %s\n", arg);
            }
            env.pid = pid;
		    break;
        case 'v':
		    env.verbose = true;
		    break;
        case ARGP_KEY_ARG:
            if (pos_args++) {
                fprintf(stderr, "Unrecognized positional argument: %s\n", arg);
            }
            errno = 0;
            break;
        default:
            return 0;
    }
    return 0;
}

int libbpf_print_fn(enum libbpf_print_level level,
		    const char *format, va_list args)
{
	if (!env.verbose) {
		return 0;
    }
	return vfprintf(stderr, format, args);
}

static int handle_event(void *ctx, void *data, size_t len)
{
	if(len < sizeof(struct process_info)) {
		return -1;
	}

	struct process_info *s = (struct process_info*)data;
	printf("%d\t%d\t%s\n", s->pid, s->signal,s->comm);
	return 0;
}


void handle_lost_events(void *ctx, int cpu, __u64 lost_cnt)
{
	fprintf(stderr, "Lost %llu events on CPU #%d!\n", lost_cnt, cpu);
}

int main(int argc, char **argv) 
{
    static const struct argp argp = {
        .options = opts,
        .parser = parse_arg,
    };

struct bpf_program x = {};
    link = bpf_program__attach_uprobe();

    struct tracesignal_bpf *obj;
	int err;

    err = argp_parse(&argp, argc, argv, 0, NULL, NULL);
    if (err) {
        return err;
    }

    libbpf_set_print(libbpf_print_fn);
    
    obj = tracesignal_bpf__open();
	if (!obj) {
		fprintf(stderr, "failed to open and/or load BPF object\n");
		return 1;
	}

	obj->rodata->target_tgid = env.pid;

    err = tracesignal_bpf__load(obj);
    if (err) {
		fprintf(stderr, "failed to load BPF object: %d\n", err);
		goto cleanup;
    }

    err = tracesignal_bpf__attach(obj);
	if (err) {
		fprintf(stderr, "failed to attach BPF programs\n");
		goto cleanup;
	}

    struct ring_buffer *ringbuffer;
	int ringbuffer_fd;
    ringbuffer_fd = bpf_map__fd(obj->maps.ringbuf);

	ringbuffer = ring_buffer__new(ringbuffer_fd, handle_event, NULL, NULL);
   
    while (1) {
		// poll for new data with a timeout of -1 ms, waiting indefinitely
		ring_buffer__poll(ringbuffer, -1);
	}
cleanup:
	tracesignal_bpf__destroy(obj);
}