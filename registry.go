package registry

import "github.com/appuio/registry/Godeps/_workspace/src/github.com/pivotal-golang/bytefmt"
import "fmt"
import "sort"
// import "strings"

//import "strings"

type BySize []RegistryItem

/*type BySize struct {
  data *[]RegistryItem
}*/

func bySizeDesc(items []RegistryItem) sort.Interface {
	return sort.Reverse(BySize(items))
}

func (items BySize) Len() int {
	return len(items)
}

func (items BySize) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}

func (items BySize) Less(i, j int) bool {
	return items[i].Bytes() < items[j].Bytes()
}

type SortOrder func([]RegistryItem) sort.Interface

type RegistryItem interface {
	Name() string
	//  Bytes() uint64
	//  Size() string
	AddChild(newChild RegistryItem) RegistryItem
	RemoveChild(i int)
	RemoveEmpty()
	Bytes() uint64
	Layers() int
	Sort()
	Children() []RegistryItem
	CollectLayers(layers *map[string]*layer, path []string)
	//  Child(name string) RegistryItem
	String() string
}

type registryItem struct {
	orderBy  SortOrder
	name     string
	children []RegistryItem
	layers   map[string]*layer
}

/* func NewRegistryItem(name string, bytes uint64) *registryItem {
  return &registryItem{name: name, bytes: bytes}
} */

func (item *registryItem) Name() string {
	return item.name
}

func (item *registryItem) Bytes() uint64 {
	var result uint64
	layers := make(map[string]*layer)
	item.CollectLayers(&layers, []string{})
	for _, layer := range layers {
		result += layer.Bytes()
	}
	/*  for _, child := range item.children {
	    result += child.Bytes()
	  }*/
	return result
}

func (item *registryItem) Layers() int {
	layers := make(map[string]*layer)
	item.CollectLayers(&layers, []string{})
  return len(layers)
}

func (item *registryItem) CollectLayers(layers *map[string]*layer, path []string) {
	for _, child := range item.children {
		child.CollectLayers(layers, append(path, item.Name()))
	}
}

func (item *registryItem) RemoveEmpty() {
  i := 0
	for range item.children {    
    if item.children[i].Layers() == 0 {
//      fmt.Println(child.Name())
      item.RemoveChild(i)
    } else {
      item.children[i].RemoveEmpty()
      i += 1
    }
	}
}

func (item *registryItem) Sort() {
	if item.orderBy != nil {
		sort.Sort(item.orderBy(item.children))
	}
	for _, child := range item.children {
		child.Sort()
	}
}

/*
func (item *registryItem) Size() string {
  return item.size
}
*/
func (item *registryItem) Children() []RegistryItem {
	return item.children
}

func (item *registryItem) AddChild(newChild RegistryItem) RegistryItem {
	child := item.Child(newChild.Name())
	if child == nil {
		child = newChild
		item.children = append(item.children, child)
	}

	return child
}

func (item *registryItem) RemoveChild(i int) {
	item.children, item.children[len(item.children)-1] = append(item.children[:i], item.children[i+1:]...), nil
}

func (item *registryItem) String() string {
	return fmt.Sprintf("%s  size: %s  children: %d  layers: %d", item.name, bytefmt.ByteSize(item.Bytes()), len(item.children), item.Layers())
}

func (item *registryItem) Child(name string) RegistryItem {
	for _, child := range item.children {
		if child.Name() == name {
			return child
		}
	}
	return nil
}

/*func NewProject(bytes uint64) *project {
  return &project{Bytes: bytes, Size: bytefmt.ByteSize(bytes), Images: make(map[string]image)}
}*/

type Registry struct {
	registryItem
	//  Projects []*project
}

func (r *Registry) String() string {
	return fmt.Sprintf("projects: %d  size: %s  layers: %d", len(r.children), bytefmt.ByteSize(r.Bytes()), r.Layers())
}

func NewRegistry() *Registry {
	return &Registry{registryItem: registryItem{orderBy: bySizeDesc}}
}

type layer struct {
	registryItem
	bytes uint64
	tags map[string]struct{}
	cmd          []string
	containerCmd []string
}

func (r *rev) String() string {
	return fmt.Sprintf("%s  size: %s  layers: %d", r.created, bytefmt.ByteSize(r.Bytes()), len(r.children))
}

func (l *layer) String() string {
	var cmd []string
	if len(l.containerCmd) > 0 {
		cmd = l.containerCmd
	} else {
		cmd = l.cmd
	}
  
/*  tags := ""
  for path := range l.tags {
    tags += strings.Join(path[1:len(path) - 1], ",")
  } */
	return fmt.Sprintf("%s  size: %s  cmd: %s", l.name[7:19], bytefmt.ByteSize(l.bytes), cmd)
}

func (l *layer) Bytes() uint64 {
	return l.bytes
}

func (l *layer) Layers() int {
	return 1
}

func (l *layer) CollectLayers(layers *map[string]*layer, path []string) {
	(*layers)[l.name] = l
  l.tags[fmt.Sprintf("%s", path)] = struct{}{}
}

type image struct {
	registryItem
}

func (item *image) String() string {
  name := ""
  if projectName == "appuio-infra" || projectName == "dtschan" {
    name = item.name
  } else {
    name = string(item.name[0]) + "*******"
  }

	return fmt.Sprintf("%s  size: %s  tags: %d  layers: %d", name, bytefmt.ByteSize(item.Bytes()), len(item.children), item.Layers())
}

type tag struct {
	registryItem
}

func (item *tag) String() string {
	return fmt.Sprintf("%s  size: %s  revisions: %d  layers: %d", item.name, bytefmt.ByteSize(item.Bytes()), len(item.children), item.Layers())
}

type rev struct {
	registryItem
	created string
}

/*func NewTag() *tag {
  return &tag{Tags: make(map[string]tag)}
}*/

type project struct {
	registryItem
	//  Images []*image
}

//var projects map[string]int
var projectName string

func (item *project) String() string {
  projectName = item.name
  name := ""
  if item.name == "dtschan" || item.name == "appuio-infra" {
    name = item.name
  } else {
    name = string(item.name[0]) + "*******"
  }
    
	return fmt.Sprintf("%s  size: %s  images: %d  layers: %d", name, bytefmt.ByteSize(item.Bytes()), len(item.children), item.Layers())

/*    if projects == nil { projects = make(map[string]int) }
    index := projects[item.name]
    if index == 0 {
      index = len(projects) + 1
      projects[item.name] = index
    }
  	return fmt.Sprintf("%s  size: %s  images: %d  layers: %d", "project", index, bytefmt.ByteSize(item.Bytes()), len(item.children), item.Layers())*/
}

/*func (parentItem *RegistryItem) AddChild(child *RegistryItem) {
  parent := Parent(parentItem)
  *parent.Children() = append(*parent.Children(), child)
  parentItem.Bytes += child.Bytes
  parentItem.Size = bytefmt.ByteSize(parentItem.Bytes)
}*/

/*  AddBytes(bytes uint64) {

  } */

func (reg *Registry) AddManifest(projectName string, imageName string, tagName string, revName string, revCreated string, manifest *Manifest) {
	/*  if (layers == nil) {
	    layers = make(map[string]*layer)
	  }*/

	//  registry.AddBytes(bytes)

	/*  var t *tag
	    var l *layer
	    if l = layers[layerName]; l != nil && bytes > 0 {
	      t = &tag{registryItem: registryItem{name: "_shared_", orderBy: bySizeDesc}}
	    } else {
	      t = &tag{registryItem: registryItem{name: tagName}}
	      l = &layer{registryItem: registryItem{name: layerName}, bytes: bytes, blobSum: blobSum, cmd: strings.Replace(cmd, "/bin/sh\",\"-c\",\"#(nop) ", "", 1)}
	      layers[layerName] = l
	    }
	*/

	/*  proj := reg.AddChild(&project{registryItem: registryItem{name: projectName, orderBy: bySizeDesc}})

	    img := proj.AddChild(&image{registryItem: registryItem{name: imageName, orderBy: bySizeDesc}})

	    img.AddChild(t)

	    t.AddChild(l)*/

	proj := reg.AddChild(&project{registryItem: registryItem{name: projectName, orderBy: bySizeDesc}})

	img := proj.AddChild(&image{registryItem: registryItem{name: imageName, orderBy: bySizeDesc}})

	tag := img.AddChild(&tag{registryItem: registryItem{name: tagName}})

	rev := tag.AddChild(&rev{registryItem: registryItem{name: revName}, created: revCreated})

	// strings.Replace(cmd, "/bin/sh\",\"-c\",\"#(nop) ", "", 1)
	for _, l := range manifest.History {
		var cmd, containerCmd []string
  	if projectName == "dtschan" || projectName == "appuio-infra" {
      if l.Config != nil {
	  		cmd = l.Config.Cmd
  		}
	  	if l.Container_Config != nil {
		  	containerCmd = l.Container_Config.Cmd
  		}
    }
		rev.AddChild(&layer{registryItem: registryItem{name: l.BlobSum}, bytes: l.Size, cmd: cmd, containerCmd: containerCmd, tags: make(map[string]struct{})})
	}

	/*

	  img := &image {
	    registryItem: *NewRegistryItem(imageName, bytes),
	  }
	  proj.AddChild(img)

	  t := &tag {
	    registryItem: *NewRegistryItem(tagName, bytes),
	  }
	  img.AddChild(t)

	  l := &layer {
	    registryItem: *NewRegistryItem(layerName, bytes),
	  }
	  t.AddChild(l) */

	//  proj.AddBytes(bytes)

	/*  img, ok := proj.Images[imageName]
	    if !ok {
	      img = &image {
	        Tags: make(map[string]*tag),
	      }
	      proj.Images[imageName] = img
	    }
	    img.Bytes += bytes
	    img.Size = bytefmt.ByteSize(img.Bytes) */

	/*  t, ok := img.Tags[tagName]
	  if !ok {
	    t = &tag {
	//      Layers: make(map[string]*layer),
	    }
	   img.Tags[tagName] = t
	//   img.Tags = append(img.Tags, t)
	  }
	  t.Bytes += bytes
	  t.Size = bytefmt.ByteSize(t.Bytes)

	//    t.Layers[layerName] = l
	   l := &layer {
	      Name: layerName,
	      BlobSum: blobSum,
	      Bytes: bytes,
	    }

	    t.Layers = append(t.Layers, l)
	//  } else {
	    // add to _shared tag
	//  }
	  l.Bytes = bytes
	  l.Size = bytefmt.ByteSize(l.Bytes) */
}

func (reg *Registry) Deduplicate() {
	for _, proj := range reg.children {
		proj.(*project).deduplicate()
	}
}

func (proj *project) deduplicate() {
	layers := make(map[string]int)

	for _, img := range proj.Children() {
		for _, t := range img.Children() {
			for _, l := range t.Children() {
				if l.Bytes() > 0 {
					layers[l.Name()]++
				}
			}
		}
	}

	shared := proj.AddChild(&image{registryItem: registryItem{name: "_shared_", orderBy: bySizeDesc}})
	for _, img := range proj.Children() {
		for _, t := range img.Children() {
			for i := len(t.Children()) - 1; i >= 0; i-- {
				l := t.Children()[i]
				if layers[l.Name()] > 1 {
					shared.AddChild(l)
					t.RemoveChild(i)
				}
			}
		}
	}
}
