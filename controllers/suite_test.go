/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
/*
 * SPDX-FileCopyrightText: 2024 Samir Zeort <samir.zeort@sap.com>
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package controllers_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	operatorkymaprojectiov1alpha1 "github.com/kyma-project/cfapi/api/v1alpha1"
	controllers "github.com/kyma-project/cfapi/controllers"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	k8sClient  client.Client                //nolint:gochecknoglobals
	k8sManager manager.Manager              //nolint:gochecknoglobals
	testEnv    *envtest.Environment         //nolint:gochecknoglobals
	ctx        context.Context              //nolint:gochecknoglobals
	cancel     context.CancelFunc           //nolint:gochecknoglobals
	reconciler *controllers.CFAPIReconciler //nolint:gochecknoglobals
)

const (
	testChartPath               = "./test/busybox"
	rateLimiterBurstDefault     = 200
	rateLimiterFrequencyDefault = 30
	failureBaseDelayDefault     = 1 * time.Second
	failureMaxDelayDefault      = 1000 * time.Second
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	ctx, cancel = context.WithCancel(context.Background())
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	rateLimiter := controllers.RateLimiter{
		Burst:           rateLimiterBurstDefault,
		Frequency:       rateLimiterFrequencyDefault,
		BaseDelay:       failureBaseDelayDefault,
		FailureMaxDelay: failureMaxDelayDefault,
	}

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = controllers.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = controllers.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	k8sManager, err = ctrl.NewManager(
		cfg, ctrl.Options{
			Scheme: scheme.Scheme,
		})
	Expect(err).ToNot(HaveOccurred())

	reconciler = &controllers.CFAPIReconciler{
		Client:             k8sManager.GetClient(),
		Scheme:             scheme.Scheme,
		EventRecorder:      k8sManager.GetEventRecorderFor("tests"),
		FinalState:         operatorkymaprojectiov1alpha1.StateReady,
		FinalDeletionState: operatorkymaprojectiov1alpha1.StateDeleting,
	}

	err = reconciler.SetupWithManager(k8sManager, rateLimiter)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred(), "failed to run manager")
	}()
})

var _ = AfterSuite(func() {
	By("canceling the context for the manager to shutdown")
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
