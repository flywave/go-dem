#include "gmt_capi.h"
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

typedef int (*gmt_module_func)(void *, int, void *);

void *GMT_Create_Session(const char *tag, unsigned int pad, unsigned int mode, int (*print_func)(FILE *, const char *));
int GMT_Destroy_Session(void *API);

/* Option list API (from libgmt) */
struct GMT_OPTION {
    struct GMT_OPTION *next;
    char option;
    char *arg;
};
struct GMT_OPTION *GMT_Make_Option(void *API, char option, const char *arg);
struct GMT_OPTION *GMT_Append_Option(void *API, struct GMT_OPTION *new_opt, struct GMT_OPTION *head);
void GMT_Destroy_Options(void *API, struct GMT_OPTION **head);

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

/* Build option list from space-separated command string, call module with mode=-1. */
static int call_module_opt(gmt_module_func func, const char *cmd) {
    if (gdemo_gmt_begin() != 0) return -1;
    char buf[2048], *save, *tok;
    struct GMT_OPTION *head = NULL;
    strncpy(buf, cmd, sizeof(buf));
    buf[sizeof(buf)-1] = '\0';
    save = buf;
    while ((tok = strtok_r(save, " ", &save))) {
        if (tok[0] == '-' && tok[1]) {
            char opt = tok[1];
            const char *arg = tok[2] ? tok + 2 : "";
            head = GMT_Append_Option(g_api, GMT_Make_Option(g_api, opt, arg), head);
        } else {
            head = GMT_Append_Option(g_api, GMT_Make_Option(g_api, '<', tok), head);
        }
    }
    int ret = (*func)(g_api, -1, head);
    GMT_Destroy_Options(g_api, &head);
    return ret;
}

int gmt_surface(const char *input_file, const char *output_file,
                double tension, double xinc, double yinc,
                double xmin, double xmax, double ymin, double ymax) {
    char cmd[2048];
    snprintf(cmd, sizeof(cmd),
        "-G%s -I%.10g/%.10g -R%.10g/%.10g/%.10g/%.10g -T%.10g %s",
        output_file, xinc, yinc, xmin, xmax, ymin, ymax, tension, input_file);
    return call_module_opt(GMT_surface, cmd);
}

int gmt_grdfilter(const char *input_file, const char *output_file,
                  const char *filter_type, const char *dist_flag) {
    char cmd[2048];
    snprintf(cmd, sizeof(cmd),
        "-G%s -F%s -D%s %s", output_file, filter_type, dist_flag, input_file);
    return call_module_opt(GMT_grdfilter, cmd);
}

int gmt_triangulate(const char *input_file, const char *output_file,
                    double xinc, double yinc,
                    double xmin, double xmax, double ymin, double ymax) {
    char cmd[2048];
    snprintf(cmd, sizeof(cmd),
        "-G%s -I%.10g/%.10g -R%.10g/%.10g/%.10g/%.10g %s",
        output_file, xinc, yinc, xmin, xmax, ymin, ymax, input_file);
    return call_module_opt(GMT_triangulate, cmd);
}

int gmt_blockmean(const char *input_file, const char *output_file,
                  double xinc, double yinc,
                  double xmin, double xmax, double ymin, double ymax) {
    char cmd[2048];
    snprintf(cmd, sizeof(cmd),
        "-G%s -I%.10g/%.10g -R%.10g/%.10g/%.10g/%.10g %s",
        output_file, xinc, yinc, xmin, xmax, ymin, ymax, input_file);
    return call_module_opt(GMT_blockmean, cmd);
}

int gmt_nearneighbor(const char *input_file, const char *output_file,
                     double xinc, double yinc,
                     double xmin, double xmax, double ymin, double ymax,
                     double search_radius, int empty_value) {
    char cmd[2048];
    snprintf(cmd, sizeof(cmd),
        "-G%s -I%.10g/%.10g -R%.10g/%.10g/%.10g/%.10g -S%.10g -N%d %s",
        output_file, xinc, yinc, xmin, xmax, ymin, ymax,
        search_radius, empty_value, input_file);
    return call_module_opt(GMT_nearneighbor, cmd);
}
