package main

import (
  "os"
  "net/http"
  "encoding/json"
  . "github.com/appuio/registry"
)


/*

{
    "kind": "List",
    "apiVersion": "v1",
    "metadata": {},
    "items": [
        {
            "kind": "ImageStream",
            "apiVersion": "v1",
            "metadata": {
                "name": "appuiojavaee7test",
                "namespace": "appuio-javaee7-example",
                "selfLink": "/oapi/v1/namespaces/appuio-javaee7-example/imagestreams/appuiojavaee7test",
                "uid": "f06af297-ec34-11e5-bcf0-001a4a026f57",
                "resourceVersion": "24186427",
                "creationTimestamp": "2016-03-17T11:39:42Z",
                "labels": {
                    "app": "appuiojavaee7test",
                    "application": "appuiojavaee7test",
                    "template": "jws-tomcat8-http-artifact"
                },
                "annotations": {
                    "description": "Keeps track of changes in the application image",
                    "openshift.io/generated-by": "OpenShiftNewApp",
                    "openshift.io/image.dockerRepositoryCheck": "2016-03-17T11:39:42Z"
                }
            },
            "spec": {},
            "status": {
                "dockerImageRepository": "172.30.15.22:5000/appuio-javaee7-example/appuiojavaee7test",
                "tags": [
                    {
                        "tag": "latest",
                        "items": [
                            {
                                "created": "2016-04-11T17:00:33Z",
                                "dockerImageReference": "172.30.15.22:5000/appuio-javaee7-example/appuiojavaee7test@sha256:9d5eeb11e6455540c95098504b8d32ce4dd86e7b9ee0662ac5146f3c4d162fb9",
                                "image": "sha256:9d5eeb11e6455540c95098504b8d32ce4dd86e7b9ee0662ac5146f3c4d162fb9"
                            },

*/


func handler(w http.ResponseWriter, r *http.Request) {
//  referencedManifests := make(map[string]struct{})
//  layers := make(map[string][]ImageStreamMetadata)
  registry := NewRegistry()

  registryUrl := os.Getenv("REGISTRY_URL")
  username := os.Getenv("OPENSHIFT_USER")
  password := os.Getenv("OPENSHIFT_PASSWORD")

  if username != "" && password != "" {
    Sh("oc login --username='%s' --password='%s'", username, password).CheckErrors()
  }

  token := Sh("oc whoami -t").Stdout()


//  client := RegistryClient{Registry: *registryPtr, Username: *usernamePtr, Password: token}

  
  var imageStreams ImageStreamList
  json.Unmarshal(Sh("oc get is -o json --all-namespaces").StdoutBytes(), &imageStreams)
//  imageStreams := imageStreamList["items"]
//  fmt.Println(imageStreams.Kind)
//  fmt.Println(len(imageStreams.Items))

  imageStreams.LoadManifests(registryUrl, username, token)

  for _, is := range imageStreams.Items {
    for _, tag := range is.Status.Tags {
      for _, rev := range tag.Items {
         registry.AddManifest(is.Metadata.Namespace, is.Metadata.Name, tag.Tag, rev.Image, rev.Created, rev.Manifest)
      }
    }  
  }

//  b, _ := json.Marshal(imageStreams)
//  fmt.Println(string(b))

//  fmt.Println(imageStreams)
//  imageStreams.Items[0].Metadata.Namespace, imageStreams.Items[0].Metadata.Name, imageStreams.Items[0].Status.Tags[0].Items[0].Image

//  allLayers := sh("oc exec -n default `oc get pod -n default -l deploymentconfig=docker-registry -o jsonpath='{..metadata.name}'` -- find /registry -path \"/registry/docker/registry/v2/repositories/*/*/_layers/sha256/*\" -type f -printf '%h\n' | sed -e 's,.*/\\([^/]\\+\\),\\1,'").StdoutLines()
//  allLayerStrings := sh("oc exec -n default `oc get pod -n default -l deploymentconfig=docker-registry -o jsonpath='{..metadata.name}'` -- find /registry -path \"/registry/docker/registry/v2/repositories/*/*/_layers/sha256/*\" -type f -printf '%h\n'").StdoutLines()
//  allLayers := make([]string, len(allLayerStrings))

/*  re := regexp.MustCompile("/registry/docker/registry/v2/repositories/([^/]+)/([^/]+)/_layers/sha256/([0-9a-z]+)")
  for i, layer := range allLayerStrings {
//    fmt.Println(layer)
    matches := re.FindStringSubmatch(layer)
    allLayers[i] = matches[3]
    layers[matches[3]] = append(layers[matches[3]], ImageStreamMetadata{ Namespace: matches[1], Name: matches[2] })
  }
  allLayerStrings = []string{}

   sort.Strings(allLayers)

  for _, imageStream := range imageStreams.Items {
//    fmt.Printf("%s/%s\n", imageStream.Metadata.Namespace, imageStream.Metadata.Name)
    for _, tag := range imageStream.Status.Tags {
      for _, rev := range tag.Items {
        for _, layer := range rev.Manifest.FsLayers {
//          referencedManifests[layer.BlobSum] = struct{}{}
          delete(layers,strings.Replace(layer.BlobSum, "sha256:", "", 1))
        }
      }
    }
  }*/



/*       for _, layer := range allLayers {
         if _, ok := layers[layer]; ok {
           for _, metadata := range layers[layer] {
             fmt.Printf("/exports/vdb8/docker/registry/v2/repositories/%s/%s/_layers/sha256/%s\n", metadata.Namespace, metadata.Name, layer)
  //           client.DeleteLayer(metadata.Namespace, metadata.Name, layer)
           }
//           client.DeleteBlob(layer)
         }


//         delete(layers,strings.Replace(layer.BlobSum, "sha256:", "", 1))
       }*/


  
/*
  repos := sh("oc get is -o json --all-namespaces |jq -r '.items[].status.dockerImageRepository'|sed -e 's,^.*:5000/,,'").StdoutLines()
//  fmt.Println(err)

   for _, repo := range repos {
     proc := sh("docker-ls tags --basic-auth --registry %s --user %s --password %s --json %s | jq -r '.Tags[]'", registryUrl, username, token, repo)
     if proc.Err() != nil {
   //    fmt.Println(err.Error() + "\n" + stderr)
       continue
     }    
     for _, tag := range proc.StdoutLines() {
       manifestJson := sh("docker-ls tag --basic-auth --registry %s --user %s --password %s --json --raw-manifest %s:%s", registryUrl, username, token, repo, tag).Stdout()
       manifest := Manifest{}
       json.Unmarshal([]byte(manifestJson), &manifest)
       for _, entry := range manifest.History {         
         json.Unmarshal([]byte(entry.V1Compatibility), &entry)
         entry.V1Compatibility = ""
       }
       for _, layer := range manifest.FsLayers {
        // fmt.Println(strings.Replace(layer.BlobSum, "sha256:", "", 1))
         delete(layers,strings.Replace(layer.BlobSum, "sha256:", "", 1))
       }
*/
//      registry.addTag(project, image, tag, id, Size, id, cmd)

/*       for _, line := range sh("docker-ls tag --basic-auth --registry http://172.30.15.22:5000 --user %s --password %s --json --raw-manifest %s:%s 2>/dev/null | jq -r .history[].v1Compatibility | jq -r -s '.[] | \"\\(.id)\t\\(.Size)\t\\(.container_config.Cmd)\t\\(.config.Cmd)\"'", *usernamePtr, token, repo, tag).StdoutLines() {
         if line == "" {
           continue
         }
         splitLine := strings.Split(line, "\t")
         id := splitLine[0]
         size := splitLine[1]
         cmd1 := splitLine[2]
         cmd2 := splitLine[3]
         cmd := ""
         if cmd1 != "null" {
             cmd = cmd1
           } else {
             cmd = cmd2
           }
//         layers[id].Repos[repo] = true
         project := strings.Split(repo, "/")[0]
         image := strings.Split(repo, "/")[1]
         Size, _ := strconv.ParseUint(size, 10, 0)
         registry.addTag(project, image, tag, id, Size, id, cmd)
//         projects[project] = 
       } */
//     }    
// resp, err := http.Get("http://example.com/")

/*  for id, layer := range layers {
    if len(layer.Repos) > 1 {
      fmt.Printf("%s %d %s %s\n", id, layer.Size, layer.Cmd, layer.Repos)
    }
  }*/
//  }

//  registry.Deduplicate()
  registry.Sort()
  RegistryTmpl(w, registry)
}

func main() {
  http.HandleFunc("/", handler)
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/srv"))))
  http.ListenAndServe(":8080", nil)
}
