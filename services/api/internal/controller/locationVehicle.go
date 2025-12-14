package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	model "tj/pkg/model"
)

type VehicleHandler struct {
	DB *gorm.DB
}

func NewVehicleHandler(dbConn *gorm.DB) *VehicleHandler {
	return &VehicleHandler{DB: dbConn}
}

func (h *VehicleHandler) GetLastLocation(c *gin.Context) {
	vehicleId := c.Param("vehicle_id")
	if vehicleId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle id not found"})
		return
	}

	var loc model.VehicleLocation
	err := h.DB.
		Where("vehicle_id = ?", vehicleId).
		Order("timestamp DESC").
		Limit(1).
		Take(&loc).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "location not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	resp := model.VehicleLocation{
		Id: loc.Id,
		MQTTLocationStruct: model.MQTTLocationStruct{
			VehicleId: loc.VehicleId,
			Latitude:  loc.Latitude,
			Longitude: loc.Longitude,
			Timestamp: loc.Timestamp,
		},
		CreatedAt: loc.CreatedAt,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *VehicleHandler) GetHistory(c *gin.Context) {
	var (
        start *int64
        end   *int64
    )

	vehicleID := c.Param("vehicle_id")
	startStr := c.Query("start")
    if startStr != "" {
        v, err := strconv.ParseInt(startStr, 10, 64)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start"})
            return
        }

        start = &v
    }

	endStr := c.Query("end")
    if endStr != "" {
        v, err := strconv.ParseInt(endStr, 10, 64)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end"})
            return
        }

        end = &v
    }

	limitStr := c.DefaultQuery("limit", "10")
    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit <= 0 {
        limit = 10
    }

	offsetStr := c.DefaultQuery("offset", "0")
    offset, err := strconv.Atoi(offsetStr)
    if err != nil || offset < 0 {
        offset = 0
    }

	q := h.DB.Where("vehicle_id = ?", vehicleID)

    if start != nil {
        q = q.Where("timestamp >= ?", *start)
    }
    if end != nil {
        q = q.Where("timestamp <= ?", *end)
    }

	var rows []model.VehicleLocation
	if err := q.Order("timestamp ASC").Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
        return
    }

	resp := make([]model.VehicleLocation, 0, len(rows))
	for _, r := range rows {
		resp = append(resp, model.VehicleLocation{
			Id: r.Id,
			MQTTLocationStruct: model.MQTTLocationStruct{
				VehicleId: r.VehicleId,
				Latitude:  r.Latitude,
				Longitude: r.Longitude,
				Timestamp: r.Timestamp,
			},
			CreatedAt: r.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, resp)
}
