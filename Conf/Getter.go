package conf

import (
	"errors"
	"strconv"
)

func (c *Conf) Bool(key string) (bool, error) {
	if v, ok := c.m[key]; ok {
		switch v.(type) {
		case string:
			return strconv.ParseBool(v.(string))
		case int:
			if v.(int) > 0 {
				return true, nil
			} else {
				return false, nil
			}
		default:
			return false, errors.New("conf: bool syntaxError.")
		}
	} else {
		return false, errors.New("conf key is not exist")
	}
}

func (c *Conf) DefaultBool(key string, defaultval bool) bool { //给出解析错误时的默认值
	if v, err := c.Bool(key); err != nil {
		return defaultval
	} else {
		return v
	}
}

func (c *Conf) Int(key string) (int, error) {
	if v, ok := c.m[key]; ok {
		switch v.(type) {
		case string:
			return strconv.Atoi(v.(string))
		case int:
			return v.(int), nil
		case int64:
			return int(v.(int64)), nil
		case int32:
			return int(v.(int32)), nil
		default:
			return 0, errors.New("conf: Int syntaxError.")
		}
	} else {
		return 0, errors.New("conf key is not exist")
	}
}

func (c *Conf) DefaultInt(key string, defaultval int) int {
	if v, err := c.Int(key); err != nil {
		return defaultval
	} else {
		return v
	}
}

func (c *Conf) Int64(key string) (int64, error) {
	v, err := c.Int(key)
	return int64(v), err
}

func (c *Conf) DefaultInt64(key string, defaultval int64) int64 {
	if v, err := c.Int64(key); err != nil {
		return defaultval
	} else {
		return v
	}
}

func (c *Conf) Float64(key string) (float64, error) {
	if v, ok := c.m[key]; ok {
		switch v.(type) {
		case string:
			return strconv.ParseFloat(v.(string), 64)
		case float64:
			return v.(float64), nil
		case float32:
			return float64(v.(float32)), nil
		case int:
			return float64(v.(int)), nil
		case int32:
			return float64(v.(int32)), nil
		case int64:
			return float64(v.(int64)), nil
		default:
			return 0.0, errors.New("conf: Float64 syntaxError.")

		}
	} else {
		return 0.0, errors.New("conf key is not exist")
	}
}

func (c *Conf) DefaultFloat64(key string, defaultval float64) float64 {
	if v, err := c.Float64(key); err != nil {
		return defaultval
	} else {
		return v
	}
}

func (c *Conf) String(key string) (string, error) {
	if v, ok := c.m[key]; ok {
		switch v.(type) {
		case string:
			return v.(string), nil
		case int:
			return strconv.FormatInt(int64(v.(int)), 10), nil
		case int64:
			return strconv.FormatInt(v.(int64), 10), nil
		case int32:
			return strconv.FormatInt(int64(v.(int32)), 10), nil
		case float64:
			return strconv.FormatFloat(v.(float64), 'f', -1, 64), nil
		case float32:
			return strconv.FormatFloat(float64(v.(float32)), 'f', -1, 32), nil
		default:
			return "", errors.New("conf: String syntaxError.")
		}
	} else {
		return "", errors.New("conf key is not exist")
	}
}

func (c *Conf) DefaultString(key string, defaultval string) string {
	if v, err := c.String(key); err != nil {
		return defaultval
	} else {
		return v
	}
}

func (c *Conf) Strings(key string) ([]string, error) {
	if v, ok := c.m[key]; ok {
		switch v.(type) {
		case []interface{}:
			ret := []string{}
			for _, vi := range v.([]interface{}) {
				switch vi.(type) {
				case string:
					ret = append(ret, vi.(string))
				case int:
					ret = append(ret, strconv.FormatInt(int64(vi.(int)), 10))
				case int64:
					ret = append(ret, strconv.FormatInt(vi.(int64), 10))
				case int32:
					ret = append(ret, strconv.FormatInt(int64(vi.(int32)), 10))
				case float64:
					ret = append(ret, strconv.FormatFloat(vi.(float64), 'f', -1, 64))
				case float32:
					ret = append(ret, strconv.FormatFloat(float64(vi.(float32)), 'f', -1, 32))
				}
			}
			return ret, nil
		case []string:
			return v.([]string), nil
		case []int:
			ret := []string{}
			for _, vi := range v.([]int) {
				ret = append(ret, strconv.FormatInt(int64(vi), 10))
			}
			return ret, nil
		case []int32:
			ret := []string{}
			for _, vi := range v.([]int32) {
				ret = append(ret, strconv.FormatInt(int64(vi), 10))
			}
			return ret, nil
		case []int64:
			ret := []string{}
			for _, vi := range v.([]int64) {
				ret = append(ret, strconv.FormatInt(vi, 10))
			}
			return ret, nil
		case []float32:
			ret := []string{}
			for _, vi := range v.([]float32) {
				ret = append(ret, strconv.FormatFloat(float64(vi), 'f', -1, 64))
			}
			return ret, nil
		case []float64:
			ret := []string{}
			for _, vi := range v.([]float64) {
				ret = append(ret, strconv.FormatFloat(vi, 'f', -1, 64))
			}
			return ret, nil
		default:
			return nil, errors.New("conf: Strings syntaxError.")
		}
	} else {
		return nil, errors.New("conf key is not exist")
	}
}

func (c *Conf) DefaultStrings(key string, defaultval []string) []string {
	if v, err := c.Strings(key); err != nil {
		return defaultval
	} else {
		return v
	}
}

func (c *Conf) Ints(key string) ([]int, error) {
	if v, ok := c.m[key]; ok {
		switch v.(type) {
		case []interface{}:
			ret := []int{}
			for _, vi := range v.([]interface{}) {
				switch vi.(type) {
				case string:
					vv, _ := strconv.Atoi(vi.(string))
					ret = append(ret, vv)
				case int:
					ret = append(ret, vi.(int))
				case int64:
					ret = append(ret, int(vi.(int64)))
				case int32:
					ret = append(ret, int(vi.(int32)))
				}
			}
			return ret, nil
		case []string:
			ret := []int{}
			for _, vi := range v.([]string) {
				vv, _ := strconv.Atoi(vi)
				ret = append(ret, vv)
			}
			return ret, nil
		case []int:
			return v.([]int), nil
		case []int32:
			ret := []int{}
			for _, vi := range v.([]int32) {
				ret = append(ret, int(vi))
			}
			return ret, nil
		case []int64:
			ret := []int{}
			for _, vi := range v.([]int64) {
				ret = append(ret, int(vi))
			}
			return ret, nil
		case []float32:
			ret := []int{}
			for _, vi := range v.([]float32) {
				ret = append(ret, int(vi))
			}
			return ret, nil
		case []float64:
			ret := []int{}
			for _, vi := range v.([]float64) {
				ret = append(ret, int(vi))
			}
			return ret, nil
		default:
			return nil, errors.New("conf: Ints syntaxError.")
		}
	} else {
		return nil, errors.New("conf key is not exist")
	}
}

func (c *Conf) DefaultInts(key string, defaultval []int) []int {
	if v, err := c.Ints(key); err != nil {
		return defaultval
	} else {
		return v
	}
}

func (c *Conf) Float64s(key string) ([]float64, error) {
	if v, ok := c.m[key]; ok {
		switch v.(type) {
		case []interface{}:
			ret := []float64{}
			for _, vi := range v.([]interface{}) {
				switch vi.(type) {
				case string:
					vv, _ := strconv.ParseFloat(vi.(string), 64)
					ret = append(ret, vv)
				case int:
					ret = append(ret, float64(vi.(int)))
				case int64:
					ret = append(ret, float64(vi.(int64)))
				case int32:
					ret = append(ret, float64(vi.(int32)))
				case float64:
					ret = append(ret, vi.(float64))
				case float32:
					ret = append(ret, float64(vi.(float32)))
				}
			}
			return ret, nil
		case []string:
			ret := []float64{}
			for _, vi := range v.([]string) {
				vv, _ := strconv.ParseFloat(vi, 64)
				ret = append(ret, vv)
			}
			return ret, nil
		case []float64:
			return v.([]float64), nil
		case []float32:
			ret := []float64{}
			for _, vi := range v.([]float32) {
				ret = append(ret, float64(vi))
			}
			return ret, nil
		case []int:
			ret := []float64{}
			for _, vi := range v.([]int) {
				ret = append(ret, float64(vi))
			}
			return ret, nil
		case []int64:
			ret := []float64{}
			for _, vi := range v.([]int64) {
				ret = append(ret, float64(vi))
			}
			return ret, nil
		case []int32:
			ret := []float64{}
			for _, vi := range v.([]int32) {
				ret = append(ret, float64(vi))
			}
			return ret, nil
		default:
			return nil, errors.New("conf: Float64s syntaxError.")
		}
	} else {
		return nil, errors.New("conf key is not exist")
	}
}

func (c *Conf) DefaultFloat64s(key string, defaultval []float64) []float64 {
	if v, err := c.Float64s(key); err != nil {
		return defaultval
	} else {
		return v
	}
}
