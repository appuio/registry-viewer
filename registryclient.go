package main

import (
  "fmt"
  "io/ioutil"
  "net/http"
)

type RegistryClient struct {
  Registry string
  Username string
  Password string  
  httpClient *http.Client
}

func (client *RegistryClient) DeleteBlob(blob string) bool {
  status, _ := client.request("DELETE", "/admin/blobs/sha256:%s", blob) 

  if status != 200 && status != 204 && status != 404 {
    return false
  }

  return true
}

func (client *RegistryClient) DeleteLayer(namespace, name, image string) bool {
  status, _ := client.request("DELETE", "/v2/%s/%s/blobs/sha256:%s", namespace, name, image) 

    // fmt.Println(status)
  if status != 200 && status != 204 && status != 404 {
    return false
  }

  return true
}


func (client *RegistryClient) request(method string, urlFormat string, a ...interface{}) (status int, response string) {
  url := fmt.Sprintf(urlFormat, a...)

  req, err := http.NewRequest(method, client.Registry + url, nil)
  req.SetBasicAuth(client.Username, client.Password)

  // fmt.Println(client.Registry + url)

  if client.httpClient == nil {
    client.httpClient = &http.Client{}
  }

  resp, err := client.httpClient.Do(req)
  if err != nil {
    return 500, err.Error()
  }
  defer resp.Body.Close()

  body, _ := ioutil.ReadAll(resp.Body)
  return resp.StatusCode, string(body)
}
