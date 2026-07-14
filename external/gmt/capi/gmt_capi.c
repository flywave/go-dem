#include "gmt_capi.h"
#include <stdio.h>
#include <string.h>

void *GMT_Create_Session(const char *tag, unsigned int pad, unsigned int mode, int (*print_func)(FILE *, const char *));
int GMT_Destroy_Session(void *API);
int GMT_Call_Module(void *API, const char *module, int mode, void *args);
int GMT_Register_Module_Static(const char *name, void *func);

extern int GMT_surface(void *API, int mode, void *args);
extern int GMT_grdfilter(void *API, int mode, void *args);
extern int GMT_triangulate(void *API, int mode, void *args);
extern int GMT_blockmean(void *API, int mode, void *args);
extern int GMT_nearneighbor(void *API, int mode, void *args);
extern int GMT_grdmask(void *API, int mode, void *args);
extern int GMT_grdsample(void *API, int mode, void *args);

static void *g_api = NULL;

int gdemo_gmt_begin(void) {
    if (g_api) return 0;
    g_api = GMT_Create_Session("go-dem", 0, 0, NULL);
    if (g_api == NULL) return -1;
    /* Register static modules so GMT_Call_Module finds them without dlsym.
       Surface internally calls grdmask; grdfilter calls grdsample. */
    GMT_Register_Module_Static("GMT_surface", (void*)GMT_surface);
    GMT_Register_Module_Static("GMT_grdfilter", (void*)GMT_grdfilter);
    GMT_Register_Module_Static("GMT_triangulate", (void*)GMT_triangulate);
    GMT_Register_Module_Static("GMT_blockmean", (void*)GMT_blockmean);
    GMT_Register_Module_Static("GMT_nearneighbor", (void*)GMT_nearneighbor);
    GMT_Register_Module_Static("GMT_grdmask", (void*)GMT_grdmask);
    GMT_Register_Module_Static("GMT_grdsample", (void*)GMT_grdsample);
    return 0;
}

void gdemo_gmt_end(void) {
    if (g_api) { GMT_Destroy_Session(g_api); g_api = NULL; }
}

static int call_module(const char *name, const char *cmd) {
    if (gdemo_gmt_begin() != 0) return -1;
    return GMT_Call_Module(g_api, name, 0, (void*)cmd);
}

int gmt_surface(const char *in, const char *out,
                double tension, double xinc, double yinc,
                double xmin, double xmax, double ymin, double ymax) {
    char cmd[2048];
    snprintf(cmd, sizeof(cmd),
        "-G%s -I%.10g/%.10g -R%.10g/%.10g/%.10g/%.10g -T%.10g %s",
        out, xinc, yinc, xmin, xmax, ymin, ymax, tension, in);
    return call_module("surface", cmd);
}

int gmt_grdfilter(const char *in, const char *out,
                  const char *filter_type, const char *dist_flag) {
    char cmd[2048];
    snprintf(cmd, sizeof(cmd), "-G%s -F%s -D%s %s", out, filter_type, dist_flag, in);
    return call_module("grdfilter", cmd);
}

int gmt_triangulate(const char *in, const char *out,
                    double xinc, double yinc,
                    double xmin, double xmax, double ymin, double ymax) {
    char cmd[2048];
    snprintf(cmd, sizeof(cmd),
        "-G%s -I%.10g/%.10g -R%.10g/%.10g/%.10g/%.10g %s",
        out, xinc, yinc, xmin, xmax, ymin, ymax, in);
    return call_module("triangulate", cmd);
}

int gmt_blockmean(const char *in, const char *out,
                  double xinc, double yinc,
                  double xmin, double xmax, double ymin, double ymax) {
    char cmd[2048];
    snprintf(cmd, sizeof(cmd),
        "-G%s -I%.10g/%.10g -R%.10g/%.10g/%.10g/%.10g %s",
        out, xinc, yinc, xmin, xmax, ymin, ymax, in);
    return call_module("blockmean", cmd);
}

int gmt_nearneighbor(const char *in, const char *out,
                     double xinc, double yinc,
                     double xmin, double xmax, double ymin, double ymax,
                     double search_radius, int empty_value) {
    char cmd[2048];
    snprintf(cmd, sizeof(cmd),
        "-G%s -I%.10g/%.10g -R%.10g/%.10g/%.10g/%.10g -S%.10g -N%d %s",
        out, xinc, yinc, xmin, xmax, ymin, ymax,
        search_radius, empty_value, in);
    return call_module("nearneighbor", cmd);
}
