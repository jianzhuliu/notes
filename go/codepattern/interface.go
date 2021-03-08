/*
接口编程
*/

package codepattern

import "math"

//-------------接口---------------//
type Shape interface {
	Area() float64      //计算面积
	Perimeter() float64 //计算周长
}

//-------------长方形---------------//
type Rect struct {
	Width, Height float64
}

func (r *Rect) Area() float64 {
	return r.Width * r.Height
}

func (r *Rect) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

//-------------圆形---------------//
type Circle struct {
	Radius float64
}

func (c *Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c *Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}
