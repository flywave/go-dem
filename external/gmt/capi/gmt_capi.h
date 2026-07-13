#ifndef GMT_CAPI_H
#define GMT_CAPI_H

#ifdef __cplusplus
extern "C" {
#endif

/* Initialize GMT session. Returns 0 on success. */
int gmt_begin(void);

/* Destroy GMT session. */
void gmt_end(void);

/*
 * Surface interpolation (splines in tension).
 * Reads XYZ from input_file, outputs grid to output_file.
 *   tension: 0 = minimum curvature, 1 = harmonic (default 0.25)
 *   xinc, yinc: grid spacing
 *   xmin, xmax, ymin, ymax: output region
 * Returns 0 on success.
 */
int gmt_surface(const char *input_file, const char *output_file,
                double tension, double xinc, double yinc,
                double xmin, double xmax, double ymin, double ymax);

/*
 * Grid filter.
 *   filter_type: e.g. "c100" (cosine 100km), "g500" (gaussian 500km)
 *   dist_flag: distance flag for grdfilter
 * Returns 0 on success.
 */
int gmt_grdfilter(const char *input_file, const char *output_file,
                  const char *filter_type, const char *dist_flag);

/*
 * Triangulate (Delaunay triangulation).
 * Reads XYZ from input_file, outputs grid to output_file.
 */
int gmt_triangulate(const char *input_file, const char *output_file,
                    double xinc, double yinc,
                    double xmin, double xmax, double ymin, double ymax);

/*
 * Blockmean (block average).
 */
int gmt_blockmean(const char *input_file, const char *output_file,
                  double xinc, double yinc,
                  double xmin, double xmax, double ymin, double ymax);

/*
 * Nearneighbor (nearest neighbor gridding).
 */
int gmt_nearneighbor(const char *input_file, const char *output_file,
                     double xinc, double yinc,
                     double xmin, double xmax, double ymin, double ymax,
                     double search_radius, int empty_value);

#ifdef __cplusplus
}
#endif

#endif /* GMT_CAPI_H */
