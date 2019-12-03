package main

import (
	"encoding/json"
	"fmt"
	"github.com/platinasystems/test"
	maas "github.com/platinasystems/tiles/pccserver/maas/models"
	"github.com/platinasystems/tiles/pccserver/models"
	"testing"
)

func reimageAllBrownNodes(t *testing.T) {
	t.Run("updateMAASInfo", updateMAASInfo)
	t.Run("reimageAllBrown", reimageAllBrown)
}

func updateMAASInfo(t *testing.T) {
	test.SkipIfDryRun(t)
	assert := test.Assert{t}
	var (
		body []byte
		resp HttpResp
		err  error
	)
	for _, i := range Env.Servers {
		var (
			node models.NodeWithKubernetes
			data []byte
		)
		addReq := nodeAddReq{
			Host:        i.HostIp,
			Id:          NodebyHostIP[i.HostIp],
			Bmc:         i.BMCIp,
			BmcUser:     i.BMCUser,
			BmcUsers:    i.BMCUsers,
			BmcPassword: i.BMCPass,
			AdminUser:   "admin",
			SSHKeys:     i.KeyId, //todo: get key for key manager
			Managed:     true,
		}
		endpoint := fmt.Sprintf("pccserver/node/update")
		if data, err = json.Marshal(addReq); err != nil {
			assert.Fatalf("invalid struct for node update request")
		}
		if resp, body, err = pccGateway("PUT", endpoint, data); err != nil {
			assert.Fatalf("%v\n%v\n", string(body), err)
			return
		}
		if resp.Status != 200 {
			assert.Fatalf("%v\n", string(body))
			fmt.Printf("Update node %v failed\n%v\n", i.HostIp, string(body))
			return
		}
		if err := json.Unmarshal(resp.Data, &node); err != nil {
			assert.Fatalf("%v\n%v\n", string(resp.Data), err)
			return
		}
	}

}

func reimageAllBrown(t *testing.T) {
	test.SkipIfDryRun(t)
	assert := test.Assert{t}
	var (
		body []byte
		resp HttpResp
		err  error
		//check      bool
		numServers uint64 = 0
	)
	nodesToCheck := make([]uint64, len(Env.Servers))

	for _, i := range Env.Servers {
		var (
			data      []byte
			nodesList []uint64 = []uint64{NodebyHostIP[i.HostIp]}
			keys      []string = i.KeyAlias
		)
		request := maas.MaasRequest{
			Nodes:     nodesList,
			Image:     "centos76",
			Locale:    "en-US",
			Timezone:  "PDT",
			AdminUser: "admin",
			SSHKeys:   keys,
		}
		fmt.Println(request)
		endpoint := fmt.Sprintf("maas/deployments")
		if data, err = json.Marshal(request); err != nil {
			assert.Fatalf("invalid struct for maas deployments request")
		}
		if resp, body, err = pccGateway("POST", endpoint, data); err != nil {
			assert.Fatalf("%v\n%v\n", string(body), err)
			return
		}
		if resp.Status != 200 {
			assert.Fatalf("%v\n", string(body))
			fmt.Printf("Deployment node %v failed\n%v\n", i.HostIp, string(body))
			return
		}
		//if err := json.Unmarshal(resp.Data, &node); err != nil {
		//	assert.Fatalf("%v\n%v\n", string(resp.Data), err)
		//	return
		//}
		nodesToCheck[numServers] = NodebyHostIP[i.HostIp]
		numServers++
	}

	//TODO
}