package ab

import (
	"strings"

	"github.com/antonmedv/expr"
	"github.com/spaolacci/murmur3"
)

type BucketService struct {
	service *ZoneLayerService
}

func NewBucketService(layers []*Layer) *BucketService {
	service := NewZoneLayerService()
	layersIndex := map[string]*ZoneLayerService{}
	for _, layer := range layers {
		zone := layer.Zone
		zoneLayer := layersIndex[zone]
		if zoneLayer == nil {
			zoneLayer = NewZoneLayerService()
			layersIndex[zone] = zoneLayer
			var parentZoneLayer *ZoneLayerService
			zones := strings.Split(zone, ".")
			for i := len(zones) - 1; i > 0 && parentZoneLayer == nil; i-- {
				parentZoneLayer = layersIndex[strings.Join(zones[:i], ".")]
			}
			if parentZoneLayer != nil {
				if parentZoneLayer.Left == nil {
					parentZoneLayer.Left = NewZoneLayerService()
				}
				parentZoneLayer.Left.Selfs = append(parentZoneLayer.Left.Selfs, zoneLayer)
			} else {
				service.Selfs = append(service.Selfs, zoneLayer)
			}
		}
		zoneLayer.Selfs = append(zoneLayer.Selfs, layer)
	}
	return &BucketService{service: service}
}

func (this *BucketService) AB(ctx *ABContext, result *ABResult) {
	this.service.ForEachLayer(func(layer *Layer) bool {
		if result.ContainsLayer(layer.Name) {
			return true
		}
		platform, _ := ctx.ConditionCtx["platform"].(string)
		if layer.Platform != "" && layer.Platform != platform {
			return false
		}
		n := murmur3.Sum32([]byte(layer.shufflePrefix + "@" + ctx.Factor))
		version := layer.GetBucket(n)
		if version == nil {
			return false
		}
		if version.Exp.Conditions != "" {
			if ctx.ConditionCtx == nil {
				return false
			}
			out, err := expr.Eval(version.Exp.Conditions, ctx.ConditionCtx)
			conditionSuccess, conditionBool := out.(bool)
			if err != nil || !conditionSuccess || !conditionBool {
				return false
			}
		}
		cohort := version.Exp.Cohort
		if cohort != nil {
			if ctx.Userid == 0 || !cohort.Contains(ctx.Userid) {
				return false
			}
		}
		result.MergeVersion(version)
		return true
	})
}
