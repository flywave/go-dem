
INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/../libsqlite3" )


INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/sqlite"
)

set(gdal_ogr_sqlite_frmt_SOURCE_FILES
          "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/sqlite/ogrsqliteutility.cpp"
          "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/sqlite/ogrsqlitedriver.cpp"
          "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/sqlite/ogrsqlitedatasource.cpp"
          "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/sqlite/rasterlite2.cpp"
          "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/sqlite/ogrsqlitelayer.cpp"
          "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/sqlite/ogrsqlitevfs.cpp"
          "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/sqlite/ogrsqliteselectlayer.cpp"
          "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/sqlite/ogrsqliteviewlayer.cpp"
          "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/sqlite/ogrsqlitetablelayer.cpp"
          "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/sqlite/ogrsqliteapiroutines.c"
          "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/sqlite/ogrsqliteexecutesql.cpp"
          "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/sqlite/ogrsqlitevirtualogr.cpp"
          "${CMAKE_CURRENT_SOURCE_DIR}/gdal/ogr/ogrsf_frmts/sqlite/ogrsqlitesinglefeaturelayer.cpp"
)

ADD_LIBRARY(gdal_ogr_sqlite_frmt STATIC
            ${gdal_ogr_sqlite_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_ogr_sqlite_frmt_SOURCE_FILES})
