INCLUDE_DIRECTORIES("${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/iris"
)

set(gdal_iris_frmt_SOURCE_FILES
	"${CMAKE_CURRENT_SOURCE_DIR}/gdal/frmts/iris/irisdataset.cpp"
)

ADD_LIBRARY(gdal_iris_frmt STATIC
            ${gdal_iris_frmt_SOURCE_FILES}
          )

SOURCE_GROUP("src" FILES ${gdal_iris_frmt_SOURCE_FILES})
