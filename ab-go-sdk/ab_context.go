package ab

import "fmt"

type ABContext struct {
	Factor       string
	ConditionCtx map[string]interface{}
	Userid       uint32
}

func NewABContext(factor string) *ABContext {
	return &ABContext{Factor: factor, ConditionCtx: map[string]interface{}{}}
}

func (this *ABContext) WithUserid(userid uint32) *ABContext {
	this.Userid = userid
	this.ConditionCtx["uid"] = userid
	return this
}

func (this *ABContext) WithConditions(conditions map[string]interface{}) *ABContext {
	for k, v := range conditions {
		this.WithCondition(k, v)
	}
	return this
}

func (this *ABContext) WithCondition(key string, value interface{}) *ABContext {
	if key == "uid" {
		fmt.Println("Cannot Set userid use WithCondition(s), Please use WithUserid")
		return this
	}
	this.ConditionCtx[key] = value
	return this
}
