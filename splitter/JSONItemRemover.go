package splitter

import "bytes"

//JSONItemRemover interface helper untuk clean up json. ini untuk remove sub data . ini untuk menghapus sub data. tujuan nya json sudah di pecah dalam file lain
type JSONItemRemover interface {
	//IsRemovedPath
	IsRemovedPath(nodePath string) bool
	//Reset reset catalog item remove
	Reset()
	//AddRangeToRemove tambahkan range untuk di remove dari string / array
	AddRangeToRemove(startIndex int, endIndex int)
	//MakeCleanedUpString membuat string dengan item array yang di masukan di kosongkan
	MakeCleanedUpString(dataToClean string) []byte
	//MakeCleanedUpByte membuat bye yang sudah di cleanup
	MakeCleanedUpByte(dataToClean []byte) []byte
	//RegisterRemovedPath tambah path untuk di remove
	RegisterRemovedPath(thePath string)
}

//containerJSONItemRemover struct container untuk JSONItemRemover
type containerJSONItemRemover struct {
	//removedDataIndex dummy map. container index of item to remove
	removedDataIndex map[int]bool
	//removedPath path to remove from json string
	removedPath map[string]bool
	//removedPathCount count removed item
	removedPathCount int
}

//NewJSONItemRemover instantiate data remover
func NewJSONItemRemover(removedPath []string) JSONItemRemover {
	removedDataIndex1 := make(map[int]bool)
	removedPath1 := make(map[string]bool)
	var rmCount int
	if removedPath != nil {
		rmCount = len(removedPath)
		for _, val := range removedPath {
			removedPath1[val] = true

		}
	}
	return containerJSONItemRemover{removedDataIndex: removedDataIndex1, removedPath: removedPath1, removedPathCount: rmCount}
}

//RegisterRemovedPath tambah path untuk di remove
func (p containerJSONItemRemover) RegisterRemovedPath(thePath string) {
	if len(thePath) > 0 {
		p.removedPath[thePath] = true
		p.removedPathCount = p.removedPathCount + 1
	}

}

//MakeCleanedUpString membuat string dengan item array yang di masukan di kosongkan
func (p containerJSONItemRemover) MakeCleanedUpString(dataToClean string) []byte {
	if p.removedPathCount == 0 {
		return []byte(dataToClean)
	}
	var containerX bytes.Buffer
	for i := 0; i < len(dataToClean); i++ {
		if !p.removedDataIndex[i] {
			containerX.WriteByte(dataToClean[i])
		}
	}
	return containerX.Bytes()
}

//MakeCleanedUpByte membuat bye yang sudah di cleanup
func (p containerJSONItemRemover) MakeCleanedUpByte(dataToClean []byte) []byte {
	defer func() {
		if err := recover(); err != nil {
			println(err)
		}
	}()
	if p.removedPathCount == 0 {
		return dataToClean
	}
	var containerX bytes.Buffer
	for i := 0; i < len(dataToClean); i++ {
		if !p.removedDataIndex[i] {
			containerX.WriteByte(dataToClean[i])
		}
	}
	return containerX.Bytes()
}

//IsRemovedPath check is path include in removed item
func (p containerJSONItemRemover) IsRemovedPath(nodePath string) bool {
	return p.removedPath[nodePath]
}

//AddRangeToRemove add range to removed index. dari start sampai dengan end akan di remove dari string
func (p containerJSONItemRemover) AddRangeToRemove(startIndex int, endIndex int) {
	for i := startIndex; i <= endIndex; i++ {
		p.removedDataIndex[i] = true
	}
}

//Reset empty the map
func (p containerJSONItemRemover) Reset() {
	p.removedDataIndex = make(map[int]bool)
}
