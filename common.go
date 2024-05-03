package test

/*
#cgo CFLAGS: -g -Wall
#cgo LDFLAGS: -L../ -l:./src/jalibs/libjacommon_a-jalockutil.o -l:./src/libs/zbxcommon/str.o -l:./src/jalibs/libjacommon_a-jafile.o -l:./src/libs/zbxlog/log.o -l:./tools/json-c-0.9/json_object.o -l:./tools/json-c-0.9/printbuf.o -l:./tools/json-c-0.9/arraylist.o -l:./tools/json-c-0.9/linkhash.o -l:./src/libs/zbxsys/threads.o -l:./src/libs/zbxcommon/misc.o -l:./src/libs/zbxsys/mutexs.o -l:./src/jalibs/libjacommon_a-jajobfile.o

char	*CONFIG_LOG_FILE	 = "/var/log/jobarranger/jobarg_server.log";
int CONFIG_DB_CON_COUNT = 10;
char	*CONFIG_TMPDIR = "/var/lib/jobarranger/tmp/";
int	CONFIG_LOG_FILE_SIZE	= 1;
int JA_FILE_PATH_LEN = 260;
char	*CONFIG_FILE		 = "/home/zuki/Documents/dev/work/jobarranger-6.0.5.1/jaconf/jobarg_server.conf";

const char	title_message[] = "Job Arranger Server";
const char	*progname = NULL;
const char	usage_message[] = "[-hV] [-c <file>]";
const char	*help_message[] = {
	"Options:",
	"  -c --config <file>   Absolute path to the configuration file",
	"  -f --foreground     Run Job Arranger server in foreground",
	"",
	"Other options:",
	"  -h --help            Give this help",
	"  -V --version         Display version number",
	NULL
};
*/
import "C"
