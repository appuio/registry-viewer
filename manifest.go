package registry

import (
  "fmt"
//  "os"
  "encoding/json"
  "net/http"
  "io/ioutil"
)

type FsLayer struct {
  BlobSum string
}

type DockerConfig struct {
  Cmd []string  
}

type HistoryEntry struct {
  Id string
  BlobSum string
  Size uint64
  Config *DockerConfig
  Container_Config *DockerConfig
  V1Compatibility string
}

type Manifest struct {
  FsLayers []*FsLayer
  History []*HistoryEntry
}

var client http.Client

/*init() {
 := &http.Client{}
}*/

func (manifest *Manifest) load(registry, username, password, namespace, name, image string) {
  url := fmt.Sprintf("%s/v2/%s/%s/manifests/%s", registry, namespace, name, image)
  req, err := http.NewRequest("GET", url, nil)
  req.SetBasicAuth(username, password)

 // client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    panic(err)
  }
  defer resp.Body.Close()

  body, _ := ioutil.ReadAll(resp.Body)
  if resp.StatusCode != 200 && resp.StatusCode != 404 {
    panic(resp.Status + " " + string(body))
  }

  // fmt.Println(string(body))
  json.Unmarshal(body, manifest)
  for i, entry := range manifest.History {
    entry.BlobSum = manifest.FsLayers[i].BlobSum
    json.Unmarshal([]byte(entry.V1Compatibility), entry)
    entry.V1Compatibility = ""
//    result, _ := json.Marshal(manifest)
//    fmt.Fprintf(os.Stderr, "%s\n", string(result))
  }

/*  for _, layer := range manifest.FsLayers {
    // fmt.Println(strings.Replace(layer.BlobSum, "sha256:", "", 1))
    delete(layers,strings.Replace(layer.BlobSum, "sha256:", "", 1))
  }*/
}
/*
func deleteLayer(registry, username, password, namespace, name, image string) bool {
  url := fmt.Sprintf("%s/v2/%s/%s/blobs/sha256:%s", registry, namespace, name, image)
  req, err := http.NewRequest("DELETE", url, nil)
  req.SetBasicAuth(username, password)

 // client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    panic(err)
  }
  io.Copy(ioutil.Discard, resp.Body)
  resp.Body.Close()

  // body, _ := ioutil.ReadAll(resp.Body)
  if resp.StatusCode != 200 && resp.StatusCode != 404 {
    return false
   //  panic(resp.Status + " " + string(body))
  }

  return true
}
*/