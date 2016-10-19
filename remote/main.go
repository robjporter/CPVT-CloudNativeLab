package main

import (
    "fmt"

    "github.com/fsouza/go-dockerclient"
)

func main() {
    endpoint := "unix:///var/run/docker.sock"
    client, _ := docker.NewClient(endpoint)
    imgs, _ := client.ListImages(docker.ListImagesOptions{All: false})
    for _, img := range imgs {
        fmt.Println("ID: ", img.ID)
        fmt.Println("RepoTags: ", img.RepoTags)
        fmt.Println("Created: ", img.Created)
        fmt.Println("Size: ", img.Size)
        fmt.Println("VirtualSize: ", img.VirtualSize)
        fmt.Println("ParentId: ", img.ParentID)
        fmt.Println("==============================================================")
    }

    containers, _ := client.ListContainers(docker.ListContainersOptions{All:true})
    for _, container := range containers {
        fmt.Println("Container ID:", container.ID)
        fmt.Println("Container Image:", container.Image)
        fmt.Println("Container Command:", container.Command)
        fmt.Println("Container Created:", container.Created)
        fmt.Println("Container State:", container.State)
        fmt.Println("Container Status:", container.Status)
        fmt.Println("Container Ports:", container.Ports)
        fmt.Println("Container Size Raw:", container.SizeRw)
        fmt.Println("Container Size FS:", container.SizeRootFs)
        fmt.Println("Container Names:", container.Names)
        fmt.Println("Container Labels:", container.Labels)
        fmt.Println("Container Network:", container.Networks)
        fmt.Println("Container Mounts:", container.Mounts)
        fmt.Println("==============================================================")
    }
    fmt.Println(len(containers))

    err := client.Stats(docker.StatsOptions{})
    fmt.Println(err)
}
