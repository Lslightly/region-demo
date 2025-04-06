package regiondemo_test

import "github.com/Lslightly/region-demo/region"

type MyStruct struct{}

func ExampleBasic() {
	var keep *int
	region.Do(func() {
		w := new(int)
		x := new(MyStruct)
		y := make([]int, 10)
		z := make(map[string]string)
		*w = use3(x, y, z)
		keep = w // w is unbound from the region.
	})
	_ = *keep
}

func use3(*MyStruct, []int, map[string]string) int {
	return 0
}

func ExampleNestedRegion() {
	region.Do(func() {
		z := new(MyStruct)
		var y *MyStruct
		region.Do(func() {
			x := new(MyStruct)
			use2(x, z)
			y = x
		})
		use(y)
	})
}

func use2(*MyStruct, *MyStruct) {}

func use(*MyStruct) {}
