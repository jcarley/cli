# simpletable

I really just wanted a very simple way to print out a sorted struct as a table, without any fancy formatting.

See [https://godoc.org/github.com/forana/simpletable](https://godoc.org/github.com/forana/simpletable) for docs.

## Example

```go
package main

import (
	"github.com/forana/simpletable"
)

func main() {
	type mountain struct {
		Name   string
		Height int
		Range  string
	}

	mountains := []mountain{
		mountain{Name: "Mount Everest", Height: 8848, Range: "Himalayas"},
		mountain{Name: "K2", Height: 8611, Range: "Karakoram"},
		mountain{Name: "Kangchenjunga", Height: 8586, Range: "Himalayas"},
		mountain{Name: "Lhotse", Height: 8516, Range: "Himalayas"},
		mountain{Name: "Makalu", Height: 8462, Range: "Himalayas"},
		mountain{Name: "Cho Oyu", Height: 8201, Range: "Himalayas"},
		mountain{Name: "Dhaulagiri", Height: 8167, Range: "Himalayas"},
		mountain{Name: "Manaslu", Height: 8156, Range: "Himalayas"},
		mountain{Name: "Nanga Parbat", Height: 8125, Range: "Himalayas"},
		mountain{Name: "Annapurna", Height: 8091, Range: "Himalayas"},
		mountain{Name: "Gasherbrum I", Height: 8068, Range: "Karakoram"},
		mountain{Name: "Broad Peak", Height: 8047, Range: "Karakoram"},
		mountain{Name: "Gasherbrum II", Height: 8035, Range: "Karakoram"},
		mountain{Name: "Shishapangma", Height: 8012, Range: "Himalayas"},
	}

	table, err := simpletable.New(simpletable.HeadersForType(mountain{}), mountains)
	if err != nil {
		panic(err)
	}

	table.Sort(0) // Name
	table.Print()
}
```

Prints:

```
NAME          HEIGHT RANGE
Annapurna     8091   Himalayas
Broad Peak    8047   Karakoram
Cho Oyu       8201   Himalayas
Dhaulagiri    8167   Himalayas
Gasherbrum I  8068   Karakoram
Gasherbrum II 8035   Karakoram
K2            8611   Karakoram
Kangchenjunga 8586   Himalayas
Lhotse        8516   Himalayas
Makalu        8462   Himalayas
Manaslu       8156   Himalayas
Mount Everest 8848   Himalayas
Nanga Parbat  8125   Himalayas
Shishapangma  8012   Himalayas
```
