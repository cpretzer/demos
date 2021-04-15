package main

import (
	"log"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

const (
	l5d_injected = "l5d-injected"
	l5d_dc       = "l5d-dc"
)

var catalogServices []*api.CatalogService

func main() {
	log.Printf("hello %s", "world")

	// Get a new client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}

	c := client.Catalog()
	queryOptions, err := createQueryOptions(
		"", "l5d-dc", false, true,
		"60s", "120s", 10, "waitHash",
		"80s", "", "", nil, 0, false, false, "")

	if err != nil {
		panic(err)
	}

	s, q, err := c.Services(queryOptions)
	if err != nil {
		log.Printf("Services error %s", err)
	}

	log.Printf("query meta %+v", q)

	extractLinkerdServices(s, c)

	p := createWatchPlan(c)

	err = p.RunWithConfig("http://localhost:8500", api.DefaultConfig())
	if err != nil {
		log.Fatalf("Error creating the plan %s", err)
	}

}

func extractLinkerdServices(sm map[string][]string, c *api.Catalog) {
	log.Printf("Catalog Services before extraction \n [%+v]", catalogServices)
	for name, meta := range sm {
		// log.Printf("service: %s -> %s", name, meta)
		for _, tag := range meta {
			// log.Printf("tag at %d: %s", i, tag)
			if tag == l5d_injected {
				log.Printf("%s is injected with linkerd", name)
				cs, q, err := c.Service(name, "", nil)
				if err != nil {
					log.Printf("Error getting service %s: [%s]", name, err)
				}

				log.Printf("Service QueryMeta %+v", q)

				log.Printf("Service %s: \n", name)

				for _, csvc := range cs {
					log.Printf("\nService: %+v", csvc)
					log.Printf("\nTags: %+v", csvc.ServiceTags)
					for i, st := range csvc.ServiceTags {
						log.Printf("Tag index [%d]: [%s]", i, st)

						if st == l5d_injected {
							log.Printf("This is a Linkerd Service")
							catalogServices = append(catalogServices, csvc)
						}
					}
				}
			}

		}
	}
	log.Printf("Catalog Services after extraction \n [%+v]", catalogServices)

}

func createQueryOptions(namespace string, datacenter string, allowStale bool,
	useCache bool, maxAge string, staleIfError string,
	waitIndex uint64, waitHash string, waitTime string, token string,
	near string, nodeMeta map[string]string, relayFactor uint8, localOnly bool,
	connect bool, filter string) (*api.QueryOptions, error) {

	maxAgeDuration, _ := time.ParseDuration(maxAge)
	staleIfErrorDuration, _ := time.ParseDuration(staleIfError)
	waitTimeDuration, _ := time.ParseDuration(waitTime)

	return &api.QueryOptions{
		Namespace:    namespace,
		Datacenter:   datacenter,
		AllowStale:   allowStale,
		UseCache:     useCache,
		MaxAge:       maxAgeDuration,
		StaleIfError: staleIfErrorDuration,
		WaitIndex:    waitIndex,
		WaitHash:     waitHash,
		WaitTime:     waitTimeDuration,
		Token:        token,
		Near:         near,
		NodeMeta:     nodeMeta,
		RelayFactor:  relayFactor,
		LocalOnly:    localOnly,
		Connect:      connect,
		Filter:       filter,
	}, nil
}

func createWatchPlan(c *api.Catalog) *watch.Plan {

	params := make(map[string]interface{})

	params["datacenter"] = l5d_dc
	params["type"] = "services"

	// params["watcher"] = "services"

	p, err := watch.Parse(params)
	if err != nil {
		log.Printf("Error creating plan %s", err)
		return nil
	}

	log.Printf("Created plan ")

	p.Handler = func(val uint64, data interface{}) {
		log.Printf("HandlerFunc called with value [%d]\n and data \n[%+v]", val, data)

		for k, v := range data.(map[string][]string) {
			log.Printf("Data at [%s] is [%+v]", k, v)
		}

		extractLinkerdServices(data.(map[string][]string), c)
	}

	// p.Watcher = func(p *watch.Plan) (watch.BlockingParamVal, interface{}, error) {
	// 	log.Printf("WatcherFunc called for plan\n %+v", p)
	// 	return watch.WaitHashVal("hash"), nil, nil
	// }

	return p
}
