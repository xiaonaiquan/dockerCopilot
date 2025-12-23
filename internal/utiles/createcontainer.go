package utiles

import (
	"context"
	"dockerCopilot/internal/svc"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

type ComposeFile struct {
	Services map[string]struct {
		Image       string   `yaml:"image"`
		Ports       []string `yaml:"ports"`
		Environment []string `yaml:"environment"`
		NetworkMode string   `yaml:"network_mode"`
		Volumes     []string `yaml:"volumes"`
		User        string   `yaml:"user"`
		Privileged  bool     `yaml:"privileged"`
	} `yaml:"services"`
}

// parsePortMapping 解析端口映射字符串，返回 hostPort, containerPort, proto
func parsePortMapping(p string) (hostPort, containerPort, proto string) {
	proto = "tcp"
	portProto := p
	if strings.HasSuffix(p, "/tcp") {
		portProto = strings.TrimSuffix(p, "/tcp")
	} else if strings.HasSuffix(p, "/udp") {
		proto = "udp"
		portProto = strings.TrimSuffix(p, "/udp")
	}
	parts := strings.SplitN(portProto, ":", 2)
	if len(parts) == 2 {
		hostPort = parts[0]
		containerPort = parts[1]
	} else {
		containerPort = portProto
	}
	return
}

// CreateContainerAsync 启动异步任务，返回任务ID
func CreateContainerAsync(ctx *svc.ServiceContext, str string) (string, error) {
	var compose ComposeFile
	err := yaml.Unmarshal([]byte(str), &compose)
	if err != nil {
		return "", fmt.Errorf("yaml解析失败: %v", err)
	}
	if len(compose.Services) == 0 {
		return "", fmt.Errorf("未找到服务定义")
	}

	taskID := uuid.New().String()
	task := &Task{
		ID:      taskID,
		Status:  TaskPending,
		Log:     []string{"任务已提交"},
		Created: time.Now(),
		Done:    make(chan struct{}),
	}
	GlobalTaskManager.AddTask(task)

	go func() {
		GlobalTaskManager.UpdateTaskStatus(taskID, TaskRunning)
		cli := ctx.DockerClient
		dockerCtx := context.Background()
		for name, svc := range compose.Services {
			if svc.Image == "" {
				GlobalTaskManager.UpdateTaskLog(taskID, name+": 缺少镜像名，跳过")
				continue
			}
			GlobalTaskManager.UpdateTaskLog(taskID, name+": 开始拉取镜像...")
			reader, err := cli.ImagePull(dockerCtx, svc.Image, image.PullOptions{})
			if err != nil {
				GlobalTaskManager.UpdateTaskLog(taskID, name+": 拉取镜像失败 - "+err.Error())
				continue
			}
			if reader != nil {
				io.Copy(io.Discard, reader)
				_ = reader.Close()
			}
			GlobalTaskManager.UpdateTaskLog(taskID, name+": 镜像拉取完成")

			envs := svc.Environment
			exposedPorts := nat.PortSet{}
			portBindings := nat.PortMap{}
			for _, p := range svc.Ports {
				hostPort, containerPort, proto := parsePortMapping(p)
				np := nat.Port(containerPort + "/" + proto)
				exposedPorts[np] = struct{}{}
				if hostPort != "" {
					portBindings[np] = []nat.PortBinding{{HostPort: hostPort}}
				}
			}
			containerConfig := &container.Config{
				Image:        svc.Image,
				Env:          envs,
				ExposedPorts: exposedPorts,
				User:         svc.User,
			}
			hostConfig := &container.HostConfig{
				PortBindings: portBindings,
				Privileged:   svc.Privileged,
			}
			if svc.NetworkMode != "" {
				hostConfig.NetworkMode = container.NetworkMode(svc.NetworkMode)
			}
			if len(svc.Volumes) > 0 {
				hostConfig.Binds = svc.Volumes
			}
			GlobalTaskManager.UpdateTaskLog(taskID, name+": 创建容器...")
			created, err := cli.ContainerCreate(dockerCtx, containerConfig, hostConfig, nil, nil, name)
			if err != nil {
				GlobalTaskManager.UpdateTaskLog(taskID, name+": 创建失败 - "+err.Error())
				continue
			}
			err = cli.ContainerStart(dockerCtx, created.ID, container.StartOptions{})
			if err != nil {
				GlobalTaskManager.UpdateTaskLog(taskID, name+": 启动失败 - "+err.Error())
				continue
			}
			GlobalTaskManager.UpdateTaskLog(taskID, name+": 创建并启动成功 (ID: "+created.ID+")")
		}
		GlobalTaskManager.UpdateTaskStatus(taskID, TaskSuccess)
	}()

	return taskID, nil
}
