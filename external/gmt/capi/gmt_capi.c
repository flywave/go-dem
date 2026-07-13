#include "gmt_capi.h"
#include <stdio.h>
#include <string.h>

/* Direct module function declarations (bypass GMT_Call_Module shared lib lookup) */
typedef int (*gmt_module_func)(void *, int, void *);

/* GMT API functions */
void *GMT_Create_Session(const char *tag, unsigned int pad, unsigned int mode, int (*print_func)(FILE *, const char *));
int GMT_Destroy_Session(void *API);

/* Core module functions compiled into libgmt.a */
extern int GMT_surface(void *API, int mode, void *args);
extern int GMT_grdfilter(void *API, int mode, void *args);
extern int GMT_triangulate(void *API, int mode, void *args);
extern int GMT_blockmean(void *API, int mode, void *args);
extern int GMT_nearneighbor(void *API, int mode, void *args);

static void *g_api = NULL;

int gdemo_gmt_begin(void) {
    if (g_api) return 0;
    g_api = GMT_Create_Session("go-dem", 0, 0, NULL);
    return (g_api == NULL) ? -1 : 0;
}

void gdemo_gmt_end(void) {
    if (g_api) {
        GMT_Destroy_Session(g_api);
        g_api = NULL;
    }
}

/* Helper: tokenize command string into argv array and call module */
static int call_module(gmt_module_func func, const char *cmd) {
    if (gdemo_gmt_begin() != 0) return -1;
    /* Parse command into argv array */
    char buf[2048];
    strncpy(buf, cmd, sizeof(buf)-1);
    char *argv[128];
    int argc = 0;
    char *token = strtok(buf, " ");
    while (token && argc < 126) {
        argv[argc++] = token;
        token = strtok(NULL, " ");
    }
    argv[argc] = NULL;
    return (*func)(g_api, argc, argv);
}

int gmt_surface(const char *input_file, const char *output_file,
                double tension, double xinc, double yinc,
                double xmin, double xmax, double ymin, double ymax) {
    char cmd[1024];
    snprintf(cmd, sizeof(cmd),
        "%s -G%s -I%.16g/%.16g -R%.16g/%.16g/%.16g/%.16g -T%.16g",
        input_file, output_file, xinc, yinc, xmin, xmax, ymin, ymax, tension);
    return call_module(GMT_surface, cmd);
}

int gmt_grdfilter(const char *input_file, const char *output_file,
                  const char *filter_type, const char *dist_flag) {
    char cmd[1024];
    snprintf(cmd, sizeof(cmd),
        "%s -G%s -F%s -D%s",
        input_file, output_file, filter_type, dist_flag);
    return call_module(GMT_grdfilter, cmd);
}

int gmt_triangulate(const char *input_file, const char *output_file,
                    double xinc, double yinc,
                    double xmin, double xmax, double ymin, double ymax) {
    char cmd[1024];
    snprintf(cmd, sizeof(cmd),
        "%s -G%s -I%.16g/%.16g -R%.16g/%.16g/%.16g/%.16g",
        input_file, output_file, xinc, yinc, xmin, xmax, ymin, ymax);
    return call_module(GMT_triangulate, cmd);
}

int gmt_blockmean(const char *input_file, const char *output_file,
                  double xinc, double yinc,
                  double xmin, double xmax, double ymin, double ymax) {
    char cmd[1024];
    snprintf(cmd, sizeof(cmd),
        "%s -G%s -I%.16g/%.16g -R%.16g/%.16g/%.16g/%.16g",
        input_file, output_file, xinc, yinc, xmin, xmax, ymin, ymax);
    return call_module(GMT_blockmean, cmd);
}

int gmt_nearneighbor(const char *input_file, const char *output_file,
                     double xinc, double yinc,
                     double xmin, double xmax, double ymin, double ymax,
                     double search_radius, int empty_value) {
    char cmd[1024];
    snprintf(cmd, sizeof(cmd),
        "%s -G%s -I%.16g/%.16g -R%.16g/%.16g/%.16g/%.16g -S%.16g -N%d",
        input_file, output_file, xinc, yinc, xmin, xmax, ymin, ymax,
        search_radius, empty_value);
    return call_module(GMT_nearneighbor, cmd);
}
