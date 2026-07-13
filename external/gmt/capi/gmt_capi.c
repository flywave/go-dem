#include "gmt_capi.h"
#include <stdio.h>
#include <string.h>

/* GMT C API declarations (from gmt.h) */
typedef struct GMTAPI_CTRL GMTAPI_CTRL;

GMTAPI_CTRL *GMT_Create_Session(const char *tag, unsigned int pad, unsigned int mode, int (*print_func)(FILE *, const char *));
int GMT_Destroy_Session(GMTAPI_CTRL *API);
int GMT_Call_Module(GMTAPI_CTRL *API, const char *module, int mode, const char *args);

static GMTAPI_CTRL *g_api = NULL;

int gmt_begin(void) {
    if (g_api) return 0; /* already initialized */
    g_api = GMT_Create_Session("go-dem", 0, 0, NULL);
    return (g_api == NULL) ? -1 : 0;
}

void gmt_end(void) {
    if (g_api) {
        GMT_Destroy_Session(g_api);
        g_api = NULL;
    }
}

int gmt_surface(const char *input_file, const char *output_file,
                double tension, double xinc, double yinc,
                double xmin, double xmax, double ymin, double ymax) {
    if (gmt_begin() != 0) return -1;
    char cmd[1024];
    snprintf(cmd, sizeof(cmd),
        "%s -G%s -I%.16g/%.16g -R%.16g/%.16g/%.16g/%.16g -T%.16g -V",
        input_file, output_file, xinc, yinc, xmin, xmax, ymin, ymax, tension);
    return GMT_Call_Module(g_api, "surface", 0, cmd);
}

int gmt_grdfilter(const char *input_file, const char *output_file,
                  const char *filter_type, const char *dist_flag) {
    if (gmt_begin() != 0) return -1;
    char cmd[1024];
    snprintf(cmd, sizeof(cmd),
        "%s -G%s -F%s -D%s -V",
        input_file, output_file, filter_type, dist_flag);
    return GMT_Call_Module(g_api, "grdfilter", 0, cmd);
}

int gmt_triangulate(const char *input_file, const char *output_file,
                    double xinc, double yinc,
                    double xmin, double xmax, double ymin, double ymax) {
    if (gmt_begin() != 0) return -1;
    char cmd[1024];
    snprintf(cmd, sizeof(cmd),
        "%s -G%s -I%.16g/%.16g -R%.16g/%.16g/%.16g/%.16g -V",
        input_file, output_file, xinc, yinc, xmin, xmax, ymin, ymax);
    return GMT_Call_Module(g_api, "triangulate", 0, cmd);
}

int gmt_blockmean(const char *input_file, const char *output_file,
                  double xinc, double yinc,
                  double xmin, double xmax, double ymin, double ymax) {
    if (gmt_begin() != 0) return -1;
    char cmd[1024];
    snprintf(cmd, sizeof(cmd),
        "%s -G%s -I%.16g/%.16g -R%.16g/%.16g/%.16g/%.16g -V",
        input_file, output_file, xinc, yinc, xmin, xmax, ymin, ymax);
    return GMT_Call_Module(g_api, "blockmean", 0, cmd);
}

int gmt_nearneighbor(const char *input_file, const char *output_file,
                     double xinc, double yinc,
                     double xmin, double xmax, double ymin, double ymax,
                     double search_radius, int empty_value) {
    if (gmt_begin() != 0) return -1;
    char cmd[1024];
    snprintf(cmd, sizeof(cmd),
        "%s -G%s -I%.16g/%.16g -R%.16g/%.16g/%.16g/%.16g -S%.16g -N%d -V",
        input_file, output_file, xinc, yinc, xmin, xmax, ymin, ymax,
        search_radius, empty_value);
    return GMT_Call_Module(g_api, "nearneighbor", 0, cmd);
}
