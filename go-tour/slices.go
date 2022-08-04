package main
func Pic(dx, dy int) [][]uint8 {
	s := make([][]uint8,dy)
	for i := range s{
		s[i] = make([]uint8, dx)
		for j := range s[i] {
			switch {
				case i%(dx/8) >= (dx/16) && j % (dx/8) >=(dx/16) :
					s[i][j] = uint8(1)
				case i%(dx/8) >= (dx/16) || j % (dx/8) >=(dx/16) :
					s[i][j] = uint8(255)
				default :
					s[i][j] = uint8(1)
			}
		}
	}
	return s
}
func test4() {
	pic.Show(Pic)
}
