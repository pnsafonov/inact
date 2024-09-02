#include <stdio.h>
#include <utmpx.h>
#include <time.h>
#include <string.h>
#include <stdlib.h>

// FreeBSD
//
// struct utmpx {
//        short           ut_type;        /* Type of entry. */
//        struct timeval  ut_tv;          /* Time entry was made. */
//        char            ut_id[8];       /* Record identifier. */
//        pid_t           ut_pid;         /* Process ID. */
//        char            ut_user[32];    /* User login name. */
//        char            ut_line[16];    /* Device name. */
//#if __BSD_VISIBLE
//        char            ut_host[128];   /* Remote hostname. */
//#else
//        char            __ut_host[128];
//#endif
//        char            __ut_spare[64];
//};

#define DEF_DAYS 7

typedef struct context
{
	int days;
	int mins;
} context_t;

void do_job(context_t *ctx)
{
	struct utmpx *utp = NULL;
	int count = 0;
	time_t now;
	time_t before;

	now = time(NULL);

	before = now - (60 * 60 * 24 * DEF_DAYS);
	if (ctx->mins > 0)
	{
		before = now - (60 * ctx->mins > 0);
	}
	if (ctx->mins == 0 && ctx->days > 0)
	{
		before = now - (60 * 60 * 24 * ctx->days);
	}

	while ((utp = getutxent()))
	{
		if (utp->ut_type != USER_PROCESS)
		{
			continue;
		}
		if (utp->ut_tv.tv_sec < before)
		{
			continue;
		}

//		time_t time0 = utp->ut_tv.tv_sec;
//		char* time0_str = NULL;
//
//		time0_str = ctime(&time0);
//		time0_str = strtok(time0_str,  "\n");
//		printf("i = %d, ut_type = %d, ut_tv.tv_sec = %d, tv_sec_str = %s, ut_id = %s, ut_pid = %d, ut_user = %s, ut_line = %s\n",
//			   count,
//			   utp->ut_type,
//			   utp->ut_tv.tv_sec,
//			   time0_str,
//			   utp->ut_id,
//			   utp->ut_pid,
//			   utp->ut_user,
//			   utp->ut_line);
		count++;
	}

	if (count != 0)
	{
		printf("where is %d recent logins, no shutdown, exit\n", count);
		exit(0);
	}

#define PATH_SIZE_0 1024
	char *path0;
	char *path1 = "/sbin";
	char path2[PATH_SIZE_0] = {0};
	int len1 = PATH_SIZE_0;
	int len2 = PATH_SIZE_0;

//  set env to find shutdown command
	path0 = getenv("PATH");
	if (path0 != NULL)
	{
		path1 = path0;
	}
	strncpy(path2, path1, PATH_SIZE_0);
	len1 = strnlen(path2, PATH_SIZE_0);
	if (len1 < PATH_SIZE_0)
	{
		len1++; // null termination of c-string
	}
	len2 = PATH_SIZE_0 - len1;

	strncat(path2, ":/sbin:/bin:/usr/sbin:/usr/bin:/usr/local/sbin:/usr/local/bin:/root/bin", len2);
	setenv("PATH", path2, 1);

#if defined(linux)
#define SHUTDOWN_CMD "shutdown -h now"
#endif

#if defined(__DragonFly__) || defined(__FreeBSD__) || defined(__OpenBSD__) || defined(__NetBSD__)
#define SHUTDOWN_CMD "shutdown -p now"
#endif

	printf("do shutdown, command: %s\n", SHUTDOWN_CMD);
	system(SHUTDOWN_CMD);
}

void print_help()
{
	printf("inact, version 1.0.0\n"
		   "Check last logins time, and shutdown if no recent logins.\n"
		   "\n"
		   "    -d <days>    days count to check, default 7\n"
		   "    -m <mins>    mins count to check\n"
		   "\n"
		   "    -h           print help and exit\n");
}

void print_ctx(context_t *ctx)
{
	time_t now;
	char *now_str;

	now = time(NULL);
	now_str = ctime(&now);

	printf("inact started at %s\n", now_str);
	printf("days = %d\n", ctx->days);
	printf("mins = %d\n", ctx->mins);
}

void parse_args(context_t *ctx, int argc, const char **argv)
{
	int l0 = argc;
	for (int i = 1; i < l0; i++)
	{
		const char *arg = argv[i];
		int i0;
		char *arg0;
		if (strcmp(arg, "-m") == 0)
		{
			i0 = i + 1;
			if (i0 < l0)
			{
				ctx->mins = atoi(argv[i0]);
			}
			continue;
		}
		if (strcmp(arg, "-d") == 0)
		{i0 = i + 1;
			if (i0 < l0)
			{
				ctx->days = atoi(argv[i0]);
			}
			continue;
		}
		if (strcmp(arg, "-h") == 0)
		{
			print_help();
			exit(0);
		}
	}

	print_ctx(ctx);
}

int main(int argc, const char **argv, const char **envp)
{
	context_t ctx;

	ctx.days = DEF_DAYS;
	ctx.mins = 0;

	parse_args(&ctx, argc, argv);
	do_job(&ctx);

	return 0;
}
