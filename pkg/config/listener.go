/*
 * Copyright 2022 The flomesh.io Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"context"
	"fmt"
	"github.com/flomesh-io/ErieCanal/pkg/commons"
	"github.com/flomesh-io/ErieCanal/pkg/kube"
	"github.com/flomesh-io/ErieCanal/pkg/repo"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

type meshCfgChangeListenerForIngress struct {
	k8sApi      *kube.K8sAPI
	configStore *Store
}

var _ MeshConfigChangeListener = &meshCfgChangeListenerForIngress{}

func (l meshCfgChangeListenerForIngress) OnConfigCreate(cfg *MeshConfig) {
	l.onUpdate(nil, cfg)
}

func (l meshCfgChangeListenerForIngress) OnConfigUpdate(oldCfg, cfg *MeshConfig) {
	l.onUpdate(oldCfg, cfg)
}

func (l meshCfgChangeListenerForIngress) OnConfigDelete(cfg *MeshConfig) {
	l.onUpdate(cfg, nil)
}

func (l meshCfgChangeListenerForIngress) onUpdate(oldCfg, cfg *MeshConfig) {
	if oldCfg == nil {
		oldCfg = l.configStore.MeshConfig.GetConfig()
	}

	if cfg == nil { // cfg is deleted
		cfg = &MeshConfig{}
	}

	klog.V(5).Infof("Operator Config is updated, new values: %#v", l.configStore.MeshConfig)
	//klog.V(5).Infof("Old RepoBaseURL = %q", oldCfg.RepoBaseURL())
	//klog.V(5).Infof("New RepoBaseURL = %q", cfg.RepoBaseURL())
	klog.V(5).Infof("Old IngressCodebasePath = %q", oldCfg.IngressCodebasePath())
	klog.V(5).Infof("New IngressCodebasePath = %q", cfg.IngressCodebasePath())

	// if ingress codebase path is changed, we need to edit ingress-controller deployment
	if oldCfg.IngressCodebasePath() != cfg.IngressCodebasePath() {
		l.updateIngressController(cfg)
	}
}

func (l meshCfgChangeListenerForIngress) updateIngressController(mc *MeshConfig) {
	// patch the deployment spec template triggers the action of rollout restart like with kubectl
	patch := fmt.Sprintf(
		`{"spec": {"template":{"metadata": {"annotations": {"kubectl.kubernetes.io/restartedAt": "%s"}}}}}`,
		time.Now().Format(commons.ProxyProfileLastUpdatedTimeFormat),
	)
	klog.V(5).Infof("patch = %s", patch)

	selector := labels.SelectorFromSet(
		map[string]string{
			"app.kubernetes.io/component": "controller",
			"app.kubernetes.io/instance":  "erie-canal-ingress-pipy",
		},
	)
	ingressList, err := l.k8sApi.Client.AppsV1().
		Deployments(v1.NamespaceAll).
		List(context.TODO(), metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		klog.Errorf("Error listing all ingress-pipy instances: %s", err)
		return
	}

	for _, ing := range ingressList.Items {
		_, err := l.k8sApi.Client.AppsV1().
			Deployments(ing.Namespace).
			Patch(context.TODO(), ing.Name, types.StrategicMergePatchType, []byte(patch), metav1.PatchOptions{})
		if err != nil {
			klog.Errorf("Patch deployment %s/%s error, %s", ing.Namespace, ing.Name, err)
		}
	}
}

type meshCfgChangeListenerForBasicConfig struct {
	client      client.Client
	k8sApi      *kube.K8sAPI
	configStore *Store
}

func (l meshCfgChangeListenerForBasicConfig) OnConfigCreate(cfg *MeshConfig) {
	// TODO: implement it if needed
}

func (l meshCfgChangeListenerForBasicConfig) OnConfigUpdate(oldCfg, cfg *MeshConfig) {
	klog.V(5).Infof("Updating basic config ...")

	if cfg.Ingress.Enabled &&
		(oldCfg.Ingress.HTTP.Enabled != cfg.Ingress.HTTP.Enabled ||
			oldCfg.Ingress.HTTP.Listen != cfg.Ingress.HTTP.Listen) {
		if err := UpdateIngressHTTPConfig(commons.DefaultIngressBasePath, repo.NewRepoClient(cfg.RepoRootURL()), cfg); err != nil {
			klog.Errorf("Failed to update HTTP config: %s", err)
		}
	}

	if oldCfg.Ingress.TLS.Enabled != cfg.Ingress.TLS.Enabled ||
		oldCfg.Ingress.TLS.Listen != cfg.Ingress.TLS.Listen ||
		oldCfg.Ingress.TLS.MTLS != cfg.Ingress.TLS.MTLS {
		if err := UpdateIngressTLSConfig(commons.DefaultIngressBasePath, repo.NewRepoClient(cfg.RepoRootURL()), cfg); err != nil {
			klog.Errorf("Failed to update TLS config: %s", err)
		}
	}
}

func (l meshCfgChangeListenerForBasicConfig) OnConfigDelete(cfg *MeshConfig) {
	// TODO: implement it if needed
}

var _ MeshConfigChangeListener = &meshCfgChangeListenerForBasicConfig{}
