package lab

import (
    "os"
    "os/signal"
    "syscall"
    "fmt"
    "errors"
    "net"
    "strings"
    "strconv"
    consulapi "github.com/hashicorp/consul/api"
    redisapi "gopkg.in/redis.v4"
    rabbitapi "github.com/streadway/amqp"
)

var FullName string
var ShortName string
var ConsulIP string
var redisIP string
var rabbitIP string

func registerEscape() {
    c := make(chan os.Signal, 2)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)

    go func() {
        <-c
        cleanup()
        os.Exit(1)
    }()
}

func cleanup() {
    client := getConsulClient(ConsulIP)
    if client != nil {
        if err := client.Agent().ServiceDeregister(FullName); err != nil {fmt.Println(err)}
    }
    fmt.Println("Deregistered "+ShortName+" from Consul server")
}

func RegisterMe(consulIP string, urlsToRegister []string, myPort string) (bool,error) {
    registerEscape()
    tagsToRegister := getAllURL(urlsToRegister)
    myIP := getLocalIP()
    if myIP == "" {os.Exit(401)}
    myID := getHostNumber(myIP)
    myShortName := "localhost-"+myID
    myFullName := myShortName+"-"+myIP+":"+myPort
    FullName = myFullName
    ShortName = myShortName
    ConsulIP = consulIP
    myPortInt,_ := strconv.Atoi(myPort)
    myURL := "http://"+myIP+":"+myPort
    redisIP = GetServiceAddress("redis")

    service := &consulapi.AgentServiceRegistration{
        ID:    myFullName,
        Name:    myShortName,
        Port:    myPortInt,
        Address: myIP,
        Tags:    tagsToRegister,
        Check:    &consulapi.AgentServiceCheck{
            HTTP:        myURL+"/health",
            Interval:    "5s",
            Timeout:    "6s",
        },
    }

    client := getConsulClient(ConsulIP)
    if client == nil {
        return false, errors.New("Unable to attach to Consul server.")
    }
    if err := client.Agent().ServiceRegister(service); err != nil {
        fmt.Println(err)
        return false, err
    }
    fmt.Printf("Registered service %q in consul with %d tags\n", myShortName, len(tagsToRegister))

    return true, nil
}

func AddQueue(message string) {
    sendRabbitMQ(message+ShortName)
}

func sendRabbitMQ(message string) {
    if rabbitIP == "" {
        rabbitIP = getServiceAddress("rabbitmq")
    }

    conn, err := rabbitapi.Dial("amqp://guest:guest@"+rabbitIP+":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare("hello", false, false, false, false, nil,)
	failOnError(err, "Failed to declare a queue")

	err = ch.Publish("", q.Name, false, false, rabbitapi.Publishing{ContentType: "text/plain",Body: []byte(message),})
	fmt.Printf(" [x] Sent %s\n", message)
	failOnError(err, "Failed to publish a message")
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
	}
}

func GetDBStartCount() string {
    return getRedisCount("STARTCOUNTER")
}

func GetPageCount() string {
    return getRedisCount("PAGECOUNTER")
}

func getRedisCount(key string) string {
    client := redisapi.NewClient(&redisapi.Options{
        Addr:     redisIP+":6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
    if err := client.Incr(key).Err(); err != nil {return ""}
    n, err := client.Get(key).Int64()
    if err != nil {return ""}
    return strconv.FormatInt(n, 10)
}

func getConsulClient(consulIP string) *consulapi.Client {
    config := &consulapi.Config{Address: ConsulIP+":8500", Scheme: "http", Token: ""}
    consul, err := consulapi.NewClient(config)
    if err == nil {
        return consul
    } else {
        fmt.Println(err)
        return nil
    }
}

func GetServiceAddress(key string) string {
    return getServiceAddress(key)
}

func GetServerCount(key string) string {
    client := getConsulClient(ConsulIP)
    services, err := client.Agent().Services()
    count := 0
    if err == nil {
        for _, service := range services {
            if strings.Contains(service.Service, key) {
                count += 1
                //fmt.Println("ID: "+service.ID)
                //fmt.Println("Service: "+service.Service)
                //fmt.Println("Tags: "+strconv.Itoa(len(service.Tags)))
                //fmt.Println("Port: "+strconv.Itoa(service.Port))
                //fmt.Println("Address: "+service.Address)
                //fmt.Println("Override: "+strconv.FormatBool(service.EnableTagOverride))
            }
        }
    }
    return strconv.Itoa(count)
}

func getServiceAddress(key string) string {
    client := getConsulClient(ConsulIP)
    services, err := client.Agent().Services()
    if err == nil {
        for _, service := range services {
            if service.ID == key {
                //fmt.Println("ID: "+service.ID)
                //fmt.Println("Service: "+service.Service)
                //fmt.Println("Tags: "+strconv.Itoa(len(service.Tags)))
                //fmt.Println("Port: "+strconv.Itoa(service.Port))
                //fmt.Println("Address: "+service.Address)
                //fmt.Println("Override: "+strconv.FormatBool(service.EnableTagOverride))
                return service.Address
            }
        }
    }
    return ""
}

func getHostNumber(ip string) string {
    numb := strings.Split(ip,".")
    if len(numb) == 4 {
        return numb[3]
    }
    return "0"
}

func getLocalIP() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return ""
    }
    for _, address := range addrs {
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return ""
}

func getAllURL(urls []string) []string {
    tags := []string{}
    for _,url := range urls {
        tags = append(tags, "urlprefix-"+url)
    }
    return tags
}
