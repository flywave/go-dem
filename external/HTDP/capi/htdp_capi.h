#ifndef HTDP_CAPI_H
#define HTDP_CAPI_H

#ifdef __cplusplus
extern "C" {
#endif

/* HTDP version info */
const char* htdp_version(void);

/* Get the velocity grid path. Returns 0 on success, -1 on error. */
int htdp_set_grid_path(const char* path);

/*
 * Transform a single point between reference frames across time.
 *
 * Parameters:
 *   lat, lon   - latitude/longitude in degrees
 *   h          - ellipsoidal height in meters
 *   src_id     - HTDP source reference frame ID (1-24)
 *   src_epoch  - source epoch in decimal years
 *   dst_id     - HTDP destination reference frame ID (1-24)
 *   dst_epoch  - destination epoch in decimal years
 *   out_lat    - output latitude in degrees
 *   out_lon    - output longitude in degrees
 *   out_h      - output ellipsoidal height in meters
 *
 * Returns 0 on success, non-zero on error.
 */
int htdp_transform(
    double lat, double lon, double h,
    int src_id, double src_epoch,
    int dst_id, double dst_epoch,
    double* out_lat, double* out_lon, double* out_h
);

/*
 * Get velocity at a point.
 *
 * Parameters:
 *   lat, lon   - latitude/longitude in degrees
 *   h          - ellipsoidal height in meters
 *   vn         - north velocity (m/yr)
 *   ve         - east velocity (m/yr)
 *   vu         - up velocity (m/yr)
 *
 * Returns 0 on success, non-zero on error.
 */
int htdp_velocity(
    double lat, double lon, double h,
    double* vn, double* ve, double* vu
);

/*
 * Convert between geodetic and geocentric coordinates.
 */
void htdp_geodetic_to_xyz(double lat, double lon, double h,
                          double* x, double* y, double* z);

void htdp_xyz_to_geodetic(double x, double y, double z,
                          double* lat, double* lon, double* h);

#ifdef __cplusplus
}
#endif

#endif /* HTDP_CAPI_H */
