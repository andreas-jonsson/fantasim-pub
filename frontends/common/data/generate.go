// +build ignore

/*
Copyright (C) 2017-2018 Andreas T Jonsson

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/shurcooL/vfsgen"
)

func main() {
	const outputFile = "../common/data/data.go"
	err := vfsgen.Generate(http.Dir("../common/data/src"), vfsgen.Options{
		Filename:     outputFile,
		PackageName:  "data",
		VariableName: "FS",
	})
	if err != nil {
		log.Fatalln(err)
	}

	data, err := ioutil.ReadFile(outputFile)
	if err != nil {
		log.Fatalln(err)
	}

	dataStr := strings.Replace(string(data), "vfsgen۰", "vfsgen_", -1)
	dataStr = strings.Replace(dataStr, "io.SeekStart", "os.SEEK_SET", -1)
	dataStr = strings.Replace(dataStr, "io.SeekCurrent", "os.SEEK_CUR", -1)

	if err := ioutil.WriteFile(outputFile, []byte(strings.Replace(dataStr, "io.SeekEnd", "os.SEEK_END", -1)), 0644); err != nil {
		log.Fatalln(err)
	}
}
