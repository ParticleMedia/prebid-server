package prometheusmetrics

import (
	"strconv"

	"github.com/prebid/prebid-server/openrtb_ext"
	"github.com/prebid/prebid-server/pbsmetrics"
)

func actionsAsString() []string {
	values := pbsmetrics.RequestActions()
	valuesAsString := make([]string, len(values))
	for i, v := range values {
		valuesAsString[i] = string(v)
	}
	return valuesAsString
}

func adaptersAsString() []string {
	values := openrtb_ext.BidderList()
	valuesAsString := make([]string, len(values))
	for i, v := range values {
		valuesAsString[i] = string(v)
	}
	return valuesAsString
}

func adapterErrorsAsString() []string {
	values := pbsmetrics.AdapterErrors()
	valuesAsString := make([]string, len(values))
	for i, v := range values {
		valuesAsString[i] = string(v)
	}
	return valuesAsString
}

func boolValuesAsString() []string {
	return []string{
		strconv.FormatBool(true),
		strconv.FormatBool(false),
	}
}

func cookieTypesAsString() []string {
	values := pbsmetrics.CookieTypes()
	valuesAsString := make([]string, len(values))
	for i, v := range values {
		valuesAsString[i] = string(v)
	}
	return valuesAsString
}

func cacheResultsAsString() []string {
	values := pbsmetrics.CacheResults()
	valuesAsString := make([]string, len(values))
	for i, v := range values {
		valuesAsString[i] = string(v)
	}
	return valuesAsString
}

func ifaTypesAsString() []string {
	values := pbsmetrics.IfaTypes()
	valuesAsString := make([]string, len(values))
	for i, v := range values {
		valuesAsString[i] = string(v)
	}
	return valuesAsString
}

func geoTypesAsString() []string {
	values := pbsmetrics.GeoTypes()
	valuesAsString := make([]string, len(values))
	for i, v := range values {
		valuesAsString[i] = string(v)
	}
	return valuesAsString
}

func IPTypesAsString() []string {
	values := pbsmetrics.IPTypes()
	valuesAsString := make([]string, len(values))
	for i, v := range values {
		valuesAsString[i] = string(v)
	}
	return valuesAsString
}

func requestStatusesAsString() []string {
	values := pbsmetrics.RequestStatuses()
	valuesAsString := make([]string, len(values))
	for i, v := range values {
		valuesAsString[i] = string(v)
	}
	return valuesAsString
}

func requestTypesAsString() []string {
	values := pbsmetrics.RequestTypes()
	valuesAsString := make([]string, len(values))
	for i, v := range values {
		valuesAsString[i] = string(v)
	}
	return valuesAsString
}

func storedDataTypesAsString() []string {
	values := pbsmetrics.StoredDataTypes()
	valuesAsString := make([]string, len(values))
	for i, v := range values {
		valuesAsString[i] = string(v)
	}
	return valuesAsString
}

func storedDataFetchTypesAsString() []string {
	values := pbsmetrics.StoredDataFetchTypes()
	valuesAsString := make([]string, len(values))
	for i, v := range values {
		valuesAsString[i] = string(v)
	}
	return valuesAsString
}

func storedDataErrorsAsString() []string {
	values := pbsmetrics.StoredDataErrors()
	valuesAsString := make([]string, len(values))
	for i, v := range values {
		valuesAsString[i] = string(v)
	}
	return valuesAsString
}

func tcfVersionsAsString() []string {
	values := pbsmetrics.TCFVersions()
	valuesAsString := make([]string, len(values))
	for i, v := range values {
		valuesAsString[i] = string(v)
	}
	return valuesAsString
}
