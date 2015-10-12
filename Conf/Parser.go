package Conf

import (
    "strconv"
	"errors"
)

func (c *Conf) Bool(key string) (bool, error) {
    if vs,ok1 := c.m[key].(string) ; ok1 {
	    return strconv.ParseBool(vs)
	}else vi,ok2 := c.m[key].(int);ok2{
	    if vi == 1 {
		   return true , nil
		}else if vi==0{
		   return false , nil
		}
	}
	return false , errors.New("conf: bool syntaxError.")
}

func (c *Conf) DefaultBool(key string, defaultval bool) bool { //给出解析错误时的默认值
	if v, err := c.Bool(key); err != nil {
		return defaultval
	} else {
		return v
	}
}

func (c *Conf) Int(key string) (int, error) {
    if vs,ok1 := c.m[key].(string) ; ok1 {
	   return strconv.Atoi(vs)
	}else vi,ok2 := c.m[key].(int);ok2{
	   return vi , nil
	}
	return 0,  errors.New("conf: Int syntaxError.")
}

func (c *Conf) DefaultInt(key string, defaultval int) int {
	if v, err := c.Int(key); err != nil {
		return defaultval
	} else {
		return v
	}
}

func (c *Conf) Int64(key string) (int64, error) {
    v,err := c.Int(key string) 
    return int64(v),err
}

func (c *Conf) DefaultInt64(key string, defaultval int64) int64 {
	if v, err := c.Int64(key); err != nil {
		return defaultval
	} else {
		return v
	}
}

func (c *Conf) Float64(key string) (float64, error) {
    if vs,ok1 := c.m[key].(string) ; ok1 {
	    return strconv.ParseFloat(vs, 64)
    }else if vf,ok2 := c.m[key].(float64);ok2{
	    return vf , nil
	}
	return 0.0,  errors.New("conf: Float64 syntaxError.")
}

func (c *Conf) DefaultFloat64(key string, defaultval float64) float64 {
	if v, err := c.Float64(key); err != nil {
		return defaultval
	} else {
		return v
	}
}

func (c *Conf) String(key string) string {
	if vs,ok := c.m[key].(string) ; ok {
	   return vs
	}
	return ""
}

func (c *Conf) DefaultString(key string, defaultval string) string {
	if v := c.String(key); v == "" {
		return defaultval
	} else {
		return v
	}
}

/*
func (c *Conf) Strings(key string) []string {
	return strings.Split(c[key], ";")
}

func (c *Conf) DefaultStrings(key string, defaultval []string) []string {
	if v := c.Strings(key); len(v) == 0 {
		return defaultval
	} else {
		return v
	}
}
*/
