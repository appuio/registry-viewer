package registry

type ImageStreamList struct {
  Items []*ImageStream
}

type ImageStream struct {
  Metadata *ImageStreamMetadata
  Status *ImageStreamStatus
}

type ImageStreamMetadata struct {
  Name string
  Namespace string
}

type ImageStreamStatus struct {
  Tags []*ImageStreamTag
}

type ImageStreamTag struct {
  Tag string
  Items []*ImageStreamTagRevision
}

type ImageStreamTagRevision struct {
  Image string
  Created string
  Manifest *Manifest
}

func (isl *ImageStreamList) LoadManifests(registry, username, password string) {
  for _, is := range isl.Items {
    is.LoadManifests(registry, username, password)
//    fmt.Println(is)
//    break
  }
}

func (is *ImageStream) LoadManifests(registry, username, password string) {
  for _, tag := range is.Status.Tags {
    for _, rev := range tag.Items {
      rev.Manifest = new(Manifest)
      rev.Manifest.load(registry, username, password, is.Metadata.Namespace, is.Metadata.Name, rev.Image)
    }
  }
}
