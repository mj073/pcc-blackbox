package pcc

import (
	"github.com/lib/pq"
	"github.com/platinasystems/tiles/pccserver/models"
	"encoding/json"
	"fmt"
	"time"
	"github.com/platinasystems/tiles/pccserver/storage/ceph"
	"strings"
)

type cephPoolTypes string
type cephPoolMeta struct {
	name string
	id uint64
}

const (
	CEPH_CLUSTER_NAME_PREFIX = "cephtest"
	CEPH_3_NODE_INSTALLATION_TIMEOUT = 1000
	CEPH_3_NODE_UNINSTALLATION_TIMEOUT = 300

	CEPH_POOL_CREATION_TIMEOUT = 300
	CEPH_POOL_DELETION_TIMEOUT = 300

	CEPH_FS_CREATION_TIMEOUT = 300
	CEPH_FS_DELETION_TIMEOUT = 300

	CEPH_POOL_FOR_K8S_TEST cephPoolTypes = "k8sTest"
	CEPH_POOL_FOR_RBD_TEST cephPoolTypes = "rbdTest"
	CEPH_POOL_FOR_CEPHFS_TEST cephPoolTypes = "cephFSTest"

	CEPH_POOL_K8S_1 = "k8sPool1"
	CEPH_POOL_K8S_2 = "k8sPool2"
	CEPH_POOL_RBD_1 = "rbdPool1"
	CEPH_POOL_RBD_2 = "rbdPool2"
	CEPH_POOL_METADATA = "metadataPool"
	CEPH_POOL_DEFAULT = "defaultPool"
	CEPH_POOL_DATA_1 = "dataPool1"
	CEPH_POOL_DATA_2 = "dataPool2"

	CEPH_FS_NAME = CEPH_CLUSTER_NAME_PREFIX + "_fs"

	// For Verification
	CEPH_CLUSTER_INSTALL_EVENT 	= "cephClusterInstall"
	CEPH_CLUSTER_UNINSTALL_EVENT 	= "cephClusterUninstall"
	CEPH_POOL_CREATE_EVENT 		= "cephPoolCreate"
	CEPH_POOL_DELETE_EVENT 		= "cephPoolDelete"
	CEPH_FS_CREATE_EVENT 		= "cephFSCreate"
	CEPH_FS_DELETE_EVENT 		= "cephFSDelete"
)

var (
	CEPH_INSTALLATION_SUCCESS_NOTIFICATION = "Ceph cluster has been deployed"
	CEPH_INSTALLATION_FAILED_NOTIFICATION_1 = "Ceph cluster [%+v] installation failed"
	CEPH_INSTALLATION_FAILED_NOTIFICATION_2 = "Unable to create ceph cluster "
	CEPH_INSTALLATION_FAILED_NOTIFICATION_3 = "Unable to store ceph cluster "
	CEPH_INSTALLATION_FAILED_NOTIFICATION_4 = "Unable to deploy ceph  cluster [%s] as there are no OSD nodes available"
	CEPH_INSTALLATION_INTERMEDIATE_NOTIFICATION_1 = "Ceph cluster installation begins"
	CEPH_INSTALLATION_INTERMEDIATE_NOTIFICATION_2 = "Successfully created network for ceph cluster"
	CEPH_INSTALLATION_INTERMEDIATE_NOTIFICATION_3 = "Create network failed for ceph cluster"
	CEPH_INSTALLATION_INTERMEDIATE_NOTIFICATION_4 = "Reachability check failed for ceph cluster"
	CEPH_INSTALLATION_INTERMEDIATE_NOTIFICATION_5 = "]. Cluster:[%+v]"
	CEPH_INSTALLATION_INTERMEDIATE_NOTIFICATION_6 = "Provisioning unused drives"
	CEPH_INSTALLATION_INTERMEDIATE_NOTIFICATION_7 = "Creating network for ceph cluster"
	CEPH_INSTALLATION_INTERMEDIATE_NOTIFICATION_8 = "Drive provisionig is finished"

	CEPH_UNINSTALLATION_SUCCESS_NOTIFICATION = "Successfully deleted network for ceph cluster"
	CEPH_UNINSTALLATION_FAILED_NOTIFICATION_1 = "Ceph cluster [%+v] uninstallation failed"
	CEPH_UNINSTALLATION_FAILED_NOTIFICATION_2 = "Unable to remove ceph cluster [%+v]"
	CEPH_UNINSTALLATION_INTERMEDIATE_NOTIFICATION_1 = "Ceph un-installation started"
	CEPH_UNINSTALLATION_INTERMEDIATE_NOTIFICATION_2 = "Deleting network for ceph cluster"
	CEPH_UNINSTALLATION_INTERMEDIATE_NOTIFICATION_3 = "ceph cluster has been removed from DB"
	CEPH_UNINSTALLATION_INTERMEDIATE_NOTIFICATION_4 = "Ceph cluster has been uninstalled but unable to remove it from database"

	CEPH_POOL_CREATION_SUCCESS_NOTIFICATION = fmt.Sprintf("Pool : [%%s] has been [%s] for cluster [%%s]", ceph.CEPH_POOL_DEPLOY_STATUS_COMPLETED)
	CEPH_POOL_CREATION_FAILED_NOTIFICATION = fmt.Sprintf("Pool : [%%s] has been [%s] for cluster [%%s]", ceph.CEPH_POOL_DEPLOY_STATUS_FAILED)
	CEPH_POOL_CREATION_INTERMEDIATE_NOTIFICATION_1 = "Creating Pool : [%s] for cluster [%s]"
	CEPH_POOL_DELETION_SUCCESS_NOTIFICATION = "Pool [%s] has been removed from DB"
	CEPH_POOL_DELETION_FAILED_NOTIFICATION = "Unable to remove pool [%s]"
	CEPH_POOL_DELETION_INTERMEDIATE_NOTIFICATION_1 = "Removing Pool : [%s] from cluster [%s]"

	CEPH_FS_CREATION_SUCCESS_NOTIFICATION = fmt.Sprintf("FS : [%%s] has been [%s] for cluster [%%s]", ceph.CEPH_FS_DEPLOY_STATUS_COMPLETED)
	CEPH_FS_CREATION_INTERMEDIATE_NOTIFICATION_1 = "Creating FS : [%s] for cluster [%s]"
	CEPH_FS_CREATION_FAILED_NOTIFICATION_1 = fmt.Sprintf("FS : [%%s] has been [%s] for cluster [%%s]", ceph.CEPH_FS_DEPLOY_STATUS_FAILED)
	CEPH_FS_CREATION_FAILED_NOTIFICATION_2 = "Unable to Create FS : [%s] for cluster [%s]"
	CEPH_FS_DELETION_SUCCESS_NOTIFICATION = "Ceph FS [%s] has been removed from DB"
	CEPH_FS_DELETION_FAILED_NOTIFICATION_1 = "Unable to remove ceph FS [%s]"
	CEPH_FS_DELETION_FAILED_NOTIFICATION_2 = "Unable to uninstall FS : [%s]"
	CEPH_FS_DELETION_INTERMEDIATE_NOTIFICATION_1 = "Removing FS : [%s]"

	CephPools = map[cephPoolTypes]map[string]uint64 {
		CEPH_POOL_FOR_K8S_TEST: map[string]uint64{CEPH_POOL_K8S_1: 0, CEPH_POOL_K8S_2: 0},
		CEPH_POOL_FOR_RBD_TEST: map[string]uint64{CEPH_POOL_RBD_1: 0, CEPH_POOL_RBD_2: 0},
		CEPH_POOL_FOR_CEPHFS_TEST: map[string]uint64{CEPH_POOL_METADATA: 0, CEPH_POOL_DEFAULT: 0, CEPH_POOL_DATA_1: 0, CEPH_POOL_DATA_2: 0},
	}
)
const (
	TestCreateCephCluster = "createCluster"
	TestDeleteCephCluster = "deleteCluster"
	TestCreateCephPools = "createPools"
	TestDeleteCephPools = "deletePools"
	TestCreateCephFS = "createFS"
	TestDeleteCephFS = "deleteFS"
)

type CephConfiguration struct {
	ClusterName 	   string	   `json:"clusterName"`
	ClusterId	   uint64    	   `json:"clusterId"`
	HighAvailability   bool      	   `json:"highAvailability"`
	NumberOfNodes 	   int	  	   `json:"numberOfNodes"`
	PublicNetwork 	   string	   `json:"publicNetwork"`
	ClusterNetwork 	   string	   `json:"clusterNetwork"`
	IgwPolicy          string	   `json:"igwPolicy"`
	ControlCIDR        string          `json:"controlCIDR"`
	Tests 		   map[string]bool `json:"tests"`
	PccClient          *PccClient
}

func (config *CephConfiguration) GetCephClusterName() string {
	return config.ClusterName
}

func (config *CephConfiguration) SetCephClusterName(name string) {
	config.ClusterName = name
}

func (config *CephConfiguration) GetCephClusterId() uint64 {
	return config.ClusterId
}

func (config *CephConfiguration) SetCephClusterId(id uint64) {
	config.ClusterId = id
}

func (p *PccClient) ValidateCephConfig(config *CephConfiguration, identifier string) error{
	if config.ClusterName != "" {
		cluster, err := p.GetCephCluster(config.ClusterName)
		if err != nil {
			return err
		}
		config.ClusterId = cluster.Id
	} else {
		config.ClusterName = fmt.Sprintf("%s_%s", CEPH_CLUSTER_NAME_PREFIX, identifier)
	}
	return nil
}

type CreateCephClusterRequest struct {
	Name              string         `json:"name"`
	Nodes             []CephNodes    `json:"nodes"`
	Version           string         `json:"version"`
	Tags              pq.StringArray `json:"tags"`
	IgwPolicy          string	 `json:"igwPolicy"`
	ControlCIDR        string        `json:"controlCIDR"`

	models.CephClusterConfig `json:"config"`
}

type CephNodes struct {
	ID uint64
}

func (p *PccClient) CreateCephCluster(request CreateCephClusterRequest) (id uint64, err error) {
	var (
		body              []byte
		data              []byte
		resp              HttpResp
	)
	endpoint := fmt.Sprintf("pccserver/storage/ceph/cluster")
	if data, err = json.Marshal(request); err != nil {
		err = fmt.Errorf("Invalid struct for ceph creation..ERROR: %v", err)
	}else {
		if resp, body, err = p.pccGateway("POST", endpoint, data); err != nil {
			err = fmt.Errorf("%v\n%v\n", string(body), err)
		} else {
			if resp.Status != 200 {
				err = fmt.Errorf("%v\n", string(body))
			}
			time.Sleep(time.Second * 5)
			cluster, errGet := p.GetCephCluster(request.Name)
			if errGet == nil {
				if cluster != nil {
					id = cluster.Id
				}else {
					err = fmt.Errorf("Failed to get cluster")
				}
			} else {
				err = fmt.Errorf("Failed to get cluster")
			}
		}
	}
	return
}

func (p *PccClient) GetCephCluster(clusterName string) (cluster *models.CephCluster, err error){
	if clusters, err := p.GetAllCephClusters(); err != nil {
		return nil, err
	}else {
		if len(clusters) > 0 {
			for _, cluster = range clusters {
				if cluster.Name != clusterName {
					continue
				} else {
					return cluster, nil
				}
			}
		}else {
			err = fmt.Errorf("Ceph cluster %v not found", clusterName)
		}
	}
	return nil, err
}

func (p *PccClient) GetAllCephClusters() (clusterList []*models.CephCluster, err error){
	var (
		body              []byte
		data              []byte
		resp              HttpResp
	)
	endpoint := fmt.Sprintf("pccserver/storage/ceph/cluster")
	if resp, body, err = p.pccGateway("GET", endpoint, data); err != nil {
		return nil, fmt.Errorf("%v\n%v\n", string(body), err)
	}
	if resp.Status != 200 {
		fmt.Printf("Ceph status check failed:\n%v\n", string(body))
		return nil, fmt.Errorf("%v\n", string(body))
	}
	err = json.Unmarshal(resp.Data, &clusterList)
	if err != nil {
		return nil, fmt.Errorf("JSON unmarshal failed for status check..ERROR: %v", err)
	}
	return clusterList, nil
}

func (p *PccClient) DeleteCephCluster(id uint64) (err error){
	var (
		body              []byte
		data              []byte
		resp              HttpResp
	)
	/*cluster, err := p.GetCephCluster(name)
	if err != nil {
		return err
	}*/
	endpoint := fmt.Sprintf("pccserver/storage/ceph/cluster/%v", id)
	if resp, body, err = p.pccGateway("DELETE", endpoint, data); err != nil {
		return fmt.Errorf("%v\n%v\n", string(body), err)
	}
	if resp.Status != 200 {
		return fmt.Errorf("%v\n", string(body))
	}
	return nil
}

type CreateCephPoolRequest struct {
	Name          string         `json:"name,omitempty"`
	PoolType      string         `json:"pool_type,omitempty"`
	Size          uint64         `json:"size,omitempty"`
	Quota         uint64         `json:"quota,omitempty"`
	QuotaUnit     string         `json:"quota_unit,omitempty"`
	Tags          pq.StringArray `json:"tags,omitempty"`
	CephClusterId uint64         `json:"ceph_cluster_id,omitempty"`
	FailureDomain int            `json:"failure_domain,omitempty"`
}

func (p *PccClient) GetCephPool(name string, cephClusterId uint64) (cephPool *models.CephPool, err error) {
	pools, err := p.GetAllCephPools(cephClusterId)
	if err == nil {
		if len(pools) > 0 {
			for _, pool := range pools {
				if pool.Name != name {
					continue
				}
				return pool, nil
			}
		}else {
			err = fmt.Errorf("No ceph pool found")
		}
	}
	return nil, err
}

func (p *PccClient) GetAllCephPools(cephClusterId uint64) (cephPools []*models.CephPool, err error) {
	if cephClusterId != 0{
		var (
			body              []byte
			data              []byte
			resp              HttpResp
		)
		endpoint := fmt.Sprintf("pccserver/storage/ceph/cluster/%v/pools", cephClusterId)
		if resp, body, err = p.pccGateway("GET", endpoint, data); err != nil {
			err = fmt.Errorf("%v\n%v\n", string(body), err)
		}else if resp.Status != 200 {
			fmt.Printf("Ceph pools get list failed:\n%v\n", string(body))
			err = fmt.Errorf("%v\n", string(body))
		}else {
			err = json.Unmarshal(resp.Data, &cephPools)
			if err != nil {
				err = fmt.Errorf("JSON unmarshal failed for status check..ERROR: %v", err)
			}
		}
	}else {
		err = fmt.Errorf("ERROR: No Ceph Cluster found")
	}
	return
}

func (p *PccClient) GetAvailableCephPools() (cephPools []*models.CephPool, err error) {
	return
}

func (p *PccClient) CreateCephPool(request CreateCephPoolRequest) (id uint64, err error) {
	var (
		body              []byte
		data              []byte
		resp              HttpResp
	)
	endpoint := fmt.Sprintf("pccserver/storage/ceph/pool")
	if data, err = json.Marshal(request); err != nil {
		err = fmt.Errorf("Invalid struct for ceph pool creation..ERROR: %v", err)
	}else {
		if resp, body, err = p.pccGateway("POST", endpoint, data); err != nil {
			err = fmt.Errorf("%v\n%v\n", string(body), err)
		} else {
			if resp.Status != 200 {
				fmt.Printf("Ceph Pool creation failed:\n%v\n", string(body))
				err = fmt.Errorf("%v\n", string(body))
			}
			time.Sleep(time.Second * 5)
			cephPool, errGet := p.GetCephPool(request.Name, request.CephClusterId)
			if errGet == nil {
				if cephPool != nil {
					id = cephPool.Id
				}else {
					err = fmt.Errorf("Failed to get ceph pool")
				}
			}
		}
	}
	return
}

func (p *PccClient) DeleteCephPool(id uint64) (err error) {
	var (
		body              []byte
		data              []byte
		resp              HttpResp
	)
	endpoint := fmt.Sprintf("pccserver/storage/ceph/pool/%v", id)
	if resp, body, err = p.pccGateway("DELETE", endpoint, data); err != nil {
		return fmt.Errorf("%v\n%v\n", string(body), err)
	}
	if resp.Status != 200 {
		return fmt.Errorf("%v\n", string(body))
	}
	return nil
}

type CreateCephFSRequest struct {
	Name          string         `json:"name,omitempty"`
	Tags          pq.StringArray `json:"tags"`
	CephClusterId uint64         `json:"ceph_cluster_id,omitempty"`
	MetadataPool  *models.CephPool   `json:"metadata_pool,omitempty"`
	DefaultPool   *models.CephPool         `json:"default_pool,omitempty"`
	DataPools     []*models.CephPool       `json:"data_pools,omitempty"`
	MaxMDS        int            `json:"max_mds,omitempty"`
}

func (p *PccClient) GetCephFS(name string, cephClusterId uint64) (cephFS *models.CephFS, err error) {
	cephFSList, err := p.GetAllCephFS(cephClusterId)
	if err == nil {
		if len(cephFSList) > 0 {
			for _, fs := range cephFSList {
				if fs.Name != name {
					continue
				}
				return fs, nil
			}
		}else {
			err = fmt.Errorf("No ceph pool found")
		}
	}
	return nil, err
}

func (p *PccClient) GetAllCephFS(cephClusterId uint64) (cephFSList []*models.CephFS, err error) {
	if cephClusterId != 0{
		var (
			body              []byte
			data              []byte
			resp              HttpResp
		)
		endpoint := fmt.Sprintf("pccserver/storage/ceph/cluster/%v/fs", cephClusterId)
		if resp, body, err = p.pccGateway("GET", endpoint, data); err != nil {
			err = fmt.Errorf("%v\n%v\n", string(body), err)
		}else if resp.Status != 200 {
			fmt.Printf("Ceph FS get list failed:\n%v\n", string(body))
			err = fmt.Errorf("%v\n", string(body))
		}else {
			err = json.Unmarshal(resp.Data, &cephFSList)
			if err != nil {
				err = fmt.Errorf("JSON unmarshal failed for status check..ERROR: %v", err)
			}
		}
	}else {
		err = fmt.Errorf("ERROR: No Ceph Cluster found")
	}
	return
}

func (p *PccClient) CreateCephFS(request CreateCephFSRequest) (id uint64, err error) {
	var (
		body              []byte
		data              []byte
		resp              HttpResp
	)
	endpoint := fmt.Sprintf("pccserver/storage/ceph/fs")
	if data, err = json.Marshal(request); err != nil {
		err = fmt.Errorf("Invalid struct for ceph fs creation..ERROR: %v", err)
	}else {
		if resp, body, err = p.pccGateway("POST", endpoint, data); err != nil {
			err = fmt.Errorf("%v\n%v\n", string(body), err)
		} else {
			if resp.Status != 200 {
				fmt.Printf("Ceph FS creation failed:\n%v\n", string(body))
				err = fmt.Errorf("%v\n", string(body))
			}
			time.Sleep(time.Second * 5)
			cephFS, errGet := p.GetCephFS(request.Name, request.CephClusterId)
			if errGet == nil {
				if cephFS != nil {
					id = cephFS.Id
				}else {
					err = fmt.Errorf("Failed to get ceph fs")
				}
			}
		}
	}
	return
}

func (p *PccClient) DeleteCephFS(id uint64) (err error) {
	var (
		body              []byte
		data              []byte
		resp              HttpResp
	)
	endpoint := fmt.Sprintf("pccserver/storage/ceph/fs/%v", id)
	if resp, body, err = p.pccGateway("DELETE", endpoint, data); err != nil {
		return fmt.Errorf("%v\n%v\n", string(body), err)
	}
	if resp.Status != 200 {
		return fmt.Errorf("%v\n", string(body))
	}
	return nil
}

func (config *CephConfiguration) VerifyCeph(startTime time.Time, action string, name string) (s Status, err error){
	s = config.PccClient.Verify(startTime, config.getCephVerifier(action, name))

	failed := !( strings.Contains(s.Msg, CEPH_INSTALLATION_SUCCESS_NOTIFICATION) ||
		strings.Contains(s.Msg, CEPH_UNINSTALLATION_SUCCESS_NOTIFICATION) ||
		strings.Contains(s.Msg, fmt.Sprintf(CEPH_POOL_CREATION_SUCCESS_NOTIFICATION, name, config.ClusterName)) ||
		strings.Contains(s.Msg, fmt.Sprintf(CEPH_POOL_DELETION_SUCCESS_NOTIFICATION, name)) ||
		strings.Contains(s.Msg, fmt.Sprintf(CEPH_FS_CREATION_SUCCESS_NOTIFICATION, name, config.ClusterName)) ||
		strings.Contains(s.Msg, fmt.Sprintf(CEPH_FS_DELETION_SUCCESS_NOTIFICATION, name)) )
	if s.IsError || failed {
		err = fmt.Errorf("%v", s.Msg)
	}
	return
}

type CephVerifier struct {
	timeout time.Duration
	events EventsToCheck
	breakLoop chan bool
}

func (v *CephVerifier) GetTimeout() time.Duration {
	return v.timeout
}

func (v *CephVerifier) GetEventsToCheck() EventsToCheck {
	return v.events
}

func (v *CephVerifier) GetBreakLoopChan() chan bool {
	return v.breakLoop
}

func (config *CephConfiguration) getCephVerifier(action string, name string) (v *CephVerifier){
	switch action {
	case CEPH_CLUSTER_INSTALL_EVENT:
		v = &CephVerifier{
			timeout: CEPH_3_NODE_INSTALLATION_TIMEOUT,
			events: EventsToCheck{
				CEPH_INSTALLATION_SUCCESS_NOTIFICATION : true,
				fmt.Sprintf(CEPH_INSTALLATION_FAILED_NOTIFICATION_1, config.GetCephClusterName()): true,
				CEPH_INSTALLATION_FAILED_NOTIFICATION_2: true,
				CEPH_INSTALLATION_FAILED_NOTIFICATION_3: true,
				fmt.Sprintf(CEPH_INSTALLATION_FAILED_NOTIFICATION_4, config.GetCephClusterName()): true,
				CEPH_INSTALLATION_INTERMEDIATE_NOTIFICATION_1: false,
				CEPH_INSTALLATION_INTERMEDIATE_NOTIFICATION_2: false,
				CEPH_INSTALLATION_INTERMEDIATE_NOTIFICATION_3: false,
				CEPH_INSTALLATION_INTERMEDIATE_NOTIFICATION_4: false,
				fmt.Sprintf(CEPH_INSTALLATION_INTERMEDIATE_NOTIFICATION_5, config.GetCephClusterName()): false,
				CEPH_INSTALLATION_INTERMEDIATE_NOTIFICATION_6: false,
				CEPH_INSTALLATION_INTERMEDIATE_NOTIFICATION_7: false,
				CEPH_INSTALLATION_INTERMEDIATE_NOTIFICATION_8: false,
			},
			breakLoop: make(chan bool),
		}
	case CEPH_CLUSTER_UNINSTALL_EVENT:
		v = &CephVerifier{
			timeout: CEPH_3_NODE_UNINSTALLATION_TIMEOUT,
			events: EventsToCheck{
				CEPH_UNINSTALLATION_SUCCESS_NOTIFICATION: true,
				fmt.Sprintf(CEPH_UNINSTALLATION_FAILED_NOTIFICATION_1, config.GetCephClusterName()): true,
				fmt.Sprintf(CEPH_UNINSTALLATION_FAILED_NOTIFICATION_2, config.GetCephClusterName()): true,
				CEPH_UNINSTALLATION_INTERMEDIATE_NOTIFICATION_1: false,
				CEPH_UNINSTALLATION_INTERMEDIATE_NOTIFICATION_2: false,
				CEPH_UNINSTALLATION_INTERMEDIATE_NOTIFICATION_3: false,
				CEPH_UNINSTALLATION_INTERMEDIATE_NOTIFICATION_4: false,
			},
			breakLoop: make(chan bool),
		}
	case CEPH_POOL_CREATE_EVENT:
		v = &CephVerifier{
			timeout: CEPH_POOL_CREATION_TIMEOUT,
			events: EventsToCheck{
				fmt.Sprintf(CEPH_POOL_CREATION_SUCCESS_NOTIFICATION, name, config.GetCephClusterName()): true,
				fmt.Sprintf(CEPH_POOL_CREATION_FAILED_NOTIFICATION, name, config.GetCephClusterName()): true,
				fmt.Sprintf(CEPH_POOL_CREATION_INTERMEDIATE_NOTIFICATION_1, name, config.GetCephClusterName()): false,
			},
			breakLoop: make(chan bool),
		}
	case CEPH_POOL_DELETE_EVENT:
		v = &CephVerifier{
			timeout: CEPH_POOL_DELETION_TIMEOUT,
			events: EventsToCheck{
				fmt.Sprintf(CEPH_POOL_DELETION_SUCCESS_NOTIFICATION, name): true,
				fmt.Sprintf(CEPH_POOL_DELETION_FAILED_NOTIFICATION, name): true,
				fmt.Sprintf(CEPH_POOL_DELETION_INTERMEDIATE_NOTIFICATION_1, name, config.GetCephClusterName()): false,
			},
			breakLoop: make(chan bool),
		}
	case CEPH_FS_CREATE_EVENT:
		v = &CephVerifier{
			timeout: CEPH_FS_CREATION_TIMEOUT,
			events: EventsToCheck{
				fmt.Sprintf(CEPH_FS_CREATION_SUCCESS_NOTIFICATION, name, config.GetCephClusterName()): true,
				fmt.Sprintf(CEPH_FS_CREATION_FAILED_NOTIFICATION_1, name, config.GetCephClusterName()): true,
				fmt.Sprintf(CEPH_FS_CREATION_FAILED_NOTIFICATION_2, name, config.GetCephClusterName()): true,
				fmt.Sprintf(CEPH_FS_CREATION_INTERMEDIATE_NOTIFICATION_1, name, config.GetCephClusterName()): false,
			},
			breakLoop: make(chan bool),
		}
	case CEPH_FS_DELETE_EVENT:
		v = &CephVerifier{
			timeout: CEPH_FS_DELETION_TIMEOUT,
			events: EventsToCheck{
				fmt.Sprintf(CEPH_FS_DELETION_SUCCESS_NOTIFICATION, name): true,
				fmt.Sprintf(CEPH_FS_DELETION_FAILED_NOTIFICATION_1, name): true,
				fmt.Sprintf(CEPH_FS_DELETION_FAILED_NOTIFICATION_2, name): true,
				fmt.Sprintf(CEPH_FS_DELETION_INTERMEDIATE_NOTIFICATION_1, name): false,
			},
			breakLoop: make(chan bool),
		}
	}
	return
}