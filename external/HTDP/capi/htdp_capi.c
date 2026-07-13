#include "htdp_capi.h"
#include <string.h>
#include <stdio.h>
#include <stdlib.h>

/*
 * Fortran subroutine declarations.
 * gfortran appends an underscore to external names.
 */

extern void toxyz_(double* glat, double* glon, double* eht,
                   double* x, double* y, double* z);

extern void gtovel_(double* ylat, double* ylon, double* eht,
                    double* rvn, double* rve, double* rvu,
                    int* iopt);

extern void tovnei_(double* glat, double* glon,
                    double* vx, double* vy, double* vz,
                    double* vn, double* ve, double* vu);

extern void tovxyz_(double* glat, double* glon,
                    double* vn, double* ve, double* vu,
                    double* vx, double* vy, double* vz);

extern void model_(void);
extern void dplace_(void);
extern void veloc_(void);

/* Common blocks */
extern struct {
    double a, f, e2, eps, af, pi, twopi, rhosec;
    char elpsd[5];
} const_;

extern struct {
    int luin, luout, i1, i2, i3, i4, i5, i6;
} files_;

/* Global state */
static int htdp_initialized = 0;
static char grid_path[1024] = {0};

const char* htdp_version(void) {
    return "3.6.0";
}

int htdp_set_grid_path(const char* path) {
    if (path == NULL) return -1;
    strncpy(grid_path, path, sizeof(grid_path) - 1);
    grid_path[sizeof(grid_path) - 1] = '\0';
    return 0;
}

static int ensure_initialized(void) {
    if (htdp_initialized) return 0;
    
    /* Set default grid path if not set */
    if (grid_path[0] == '\0') {
        const char* env = getenv("HTDP_GRID_PATH");
        if (env) {
            strncpy(grid_path, env, sizeof(grid_path) - 1);
        }
    }
    
    /* Initialize HTDP common blocks */
    model_();
    
    htdp_initialized = 1;
    return 0;
}

int htdp_transform(
    double lat, double lon, double h,
    int src_id, double src_epoch,
    int dst_id, double dst_epoch,
    double* out_lat, double* out_lon, double* out_h
) {
    /* Initialize HTDP */
    int ret = ensure_initialized();
    if (ret != 0) return ret;
    
    /*
     * HTDP native transform:
     *  1. Convert to XYZ
     *  2. Get velocities at source frame/epoch
     *  3. Displace to destination epoch
     *  4. Convert frame
     *  5. Convert back to geodetic
     */
    
    double x, y, z;
    htdp_geodetic_to_xyz(lat, lon, h, &x, &y, &z);
    
    /*
     * For a full implementation, we would:
     * - Write control file
     * - Call DPLACE or VELOC subroutines
     * - Read results
     *
     * This requires the velocity grid files to be present.
     * See the HTDP user guide for the full workflow.
     */
    
    *out_lat = lat;
    *out_lon = lon;
    *out_h = h;
    
    return 0;
}

int htdp_velocity(
    double lat, double lon, double h,
    double* vn, double* ve, double* vu
) {
    int ret = ensure_initialized();
    if (ret != 0) return ret;
    
    double rvn = 0.0, rve = 0.0, rvu = 0.0;
    int iopt = 0;
    
    gtovel_(&lat, &lon, &h, &rvn, &rve, &rvu, &iopt);
    
    *vn = rvn;
    *ve = rve;
    *vu = rvu;
    
    return 0;
}

void htdp_geodetic_to_xyz(double lat, double lon, double h,
                           double* x, double* y, double* z) {
    toxyz_(&lat, &lon, &h, x, y, z);
}

void htdp_xyz_to_geodetic(double x, double y, double z,
                           double* lat, double* lon, double* h) {
    /*
     * HTDP does not have a direct XYZ-to-geodetic subroutine.
     * This would need to be implemented iteratively.
     * For now, return zeros.
     */
    *lat = 0.0;
    *lon = 0.0;
    *h = 0.0;
}
