# go-dem

**Go 语言实现的 DEM (Digital Elevation Model) 生成与处理引擎**

提供完整的 DEM 生产管线：数据获取 → 数据堆栈 → 格网化插值 → 后处理滤波 → 不确定性评估 → 可视化输出。

## 特性

### 核心技术

- **KDTree 空间索引** — 所有空间搜索算法 (IDW/CUBE/StepGrid/NaturalNeighbor) 使用 KDTree 替代暴力搜索，查询复杂度从 O(n) 优化到 O(log n)
- **NoData 安全处理** — `NoData` 使用 `*float64` 指针类型，配合 `CoalesceNoData()` 和 `IsNoDataValue()` 消除零值歧义

### 格网化引擎 (waffle/)

| 方法 | 说明 |
|------|------|
| IDW | 反距离加权插值，KDTree 加速，支持 PointsToGrid (C++) 和内存 IDW |
| Kriging | 高斯/指数/球状变差函数模型的普通克里金插值 |
| Linear | 基于 Delaunay 三角剖分的线性插值 (含网格空间索引) |
| Cubic | 基于 Delaunay 三角剖分的三次插值 |
| Nearest | 最近邻插值 |
| Natural Neighbor | Laplace (non-Sibsonian) 自然邻域插值 |
| Inpaint | Fast Marching Method (Telea 2004) 影像修复填充 |
| CUBE | TVU/THU 逐点不确定性建模 + 密度模式聚类 + 信息准则 (AIC-like) 假设选择 + 邻域假设传播 |
| Step Grid | 多分辨率 step-down: 局部点密度分位数动态层级选择 + KDTree 空间索引 + 不确定性加权跨层级融合 |

### 后处理引擎 (grits/)

| 滤镜 | 说明 |
|------|------|
| Gaussian | 可分离高斯模糊 |
| Median | 中值滤波去噪 |
| Clip | 矢量多边形裁剪 (Ray-Casting 算法) |
| Fill | 空洞填充 (边缘扩散 + IDW + GDAL FillNodata) |
| Blend | 重叠区域线性距离融合 |
| Morphology | 腐蚀/膨胀/开/闭运算 |
| Hydro | 水文 sink 填充 + 平坦区域梯度化 |
| Diff | DEM 差异检测 (支持重采样对齐) |
| Weights | 质量权重缓冲区 (Manhattan 距离变换) |
| Rivers | D8 流向 + 汇流累计 + 河网阈值提取 |
| Slope Filter | 坡度滤波 (Horn 二阶差分) |

### 垂向基准 (datum/)

- EGM84 / EGM96 / EGM2008 椭球高与正高互转
- VDatum 网格生成与 DEM 应用

### 点云过滤 (pointz/)

- 统计异常值过滤 (Statistical Outlier)
- 半径异常值过滤 (Radius Outlier)
- 体素/随机/径向降采样

### 可视化 (perspecto/)

- 山体阴影 (Hillshade)
- 坡度/坡向 (Slope/Aspect)
- 彩色晕渲 (Color Relief)
- 叠合渲染 (Shaded Relief)

### 不确定性估计 (uncertainty/)

- Split-Sample 交叉验证不确定性
- Proximity 距离基不确定性
- Combined 组合不确定性

### 数据获取 (fetch/)

- SRTM / GEBCO / Copernicus / NOAA Multibeam / USGS TNM / EMODNet / ArcticDEM

### 数据清单 (datalist/)

- 5 波段数据堆栈 (Elevation/Count/Weight/Uncertainty/SourceID)
- XYZ 点文件解析
- LAS 文件封装

### 输出 (output.go)

- GeoTIFF / COG
- 多波段 RGB
- 灵活的数据类型 / 压缩 / tiling 配置

## 架构

```
输入数据 (GeoTIFF/LAS/XYZ/...) 
    ↓
datalist/  数据堆栈 → 5 波段 (Z/Count/Weight/Uncertainty/SourceID)
    ↓
waffle/    格网化插值 → IDW/Kriging/TIN/NaturalNeighbor/CUBE...
    ↓
grits/     后处理滤波 → 平滑/裁剪/填充/融合/形态学...
    ↓
datum/     垂向基准变换 → EGM84/96/2008
    ↓
uncertainty/ 不确定性评估 → Split-Sample/Proximity
    ↓
perspecto/ 可视化 → Hillshade/Slope/ColorRelief
    ↓
output.go  输出 → GeoTIFF/COG
```

## 快速开始

```go
package main

import (
    "github.com/flywave/go-dem"
    "github.com/flywave/go-dem/waffle"
    "github.com/flywave/go-geo"
)

func main() {
    srs := geo.NewProj("EPSG:4326")
    region := dem.NewRegionFromBBox(-125, 40, -122, 43, srs, 0.0005, 0.0005)

    w, _ := waffle.New(dem.MethodIDW)

    result, _ := w.Run([]string{"input.tif"}, &waffle.Options{
        Region: region,
    })

    dem.CreateDEM(result.DEM, region, "output.tif", -9999)
}
```

## 依赖

| 库 | 用途 |
|----|------|
| [flywave-gdal](https://github.com/flywave/flywave-gdal) | GDAL 栅格 I/O |
| [flywave-pointcloud](https://github.com/flywave/flywave-pointcloud) | 点云处理 |
| [go-proj](https://github.com/flywave/go-proj) | 坐标参考系统变换 |
| [go-geoid](https://github.com/flywave/go-geoid) | EGM 垂向基准 |
| [go-geo](https://github.com/flywave/go-geo) | 地理空间类型 (Proj/Extent/Grid) |
| [go-kriging](https://github.com/flywave/go-kriging) | Kriging 插值 |
| [go-delaunay](https://github.com/flywave/go-delaunay) | Delaunay 三角剖分 |
| [go-geom](https://github.com/flywave/go-geom) | 几何类型 |
| [go-geos](https://github.com/flywave/go-geos) | GEOS 空间操作 |

## 测试

```bash
# 运行所有纯 Go 算法测试
go test ./dem ./grits/ ./perspecto/ ./uncertainty/ ./datalist/ ./waffle/

# 运行包含 CGo 的完整测试 (可能需要较长时间)
go test -timeout 120s ./...
```

## 许可

MIT
