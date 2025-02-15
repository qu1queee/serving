//go:build e2e
// +build e2e

/*
Copyright 2019 The Knative Authors

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

package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	pkgtest "knative.dev/pkg/test"
	"knative.dev/pkg/test/spoof"
	"knative.dev/serving/test"
	v1test "knative.dev/serving/test/v1"

	rtesting "knative.dev/serving/pkg/testing/v1"
)

func TestCustomResourcesLimits(t *testing.T) {
	t.Parallel()
	clients := test.Setup(t)

	t.Log("Creating a new Route and Configuration")
	withResources := rtesting.WithResourceRequirements(corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse("350Mi"),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse("350Mi"),
		},
	})

	names := test.ResourceNames{
		Service: test.ObjectNameForTest(t),
		Image:   test.Autoscale,
	}

	test.EnsureTearDown(t, clients, &names)

	objects, err := v1test.CreateServiceReady(t, clients, &names, withResources)
	if err != nil {
		t.Fatalf("Failed to create initial Service %v: %v", names.Service, err)
	}
	endpoint := objects.Route.Status.URL.URL()

	_, err = pkgtest.CheckEndpointState(
		context.Background(),
		clients.KubeClient,
		t.Logf,
		endpoint,
		spoof.MatchesAllOf(spoof.IsStatusOK),
		"ResourceTestServesText",
		test.ServingFlags.ResolvableDomain,
		test.AddRootCAtoTransport(context.Background(), t.Logf, clients, test.ServingFlags.HTTPS),
		spoof.WithHeader(test.ServingFlags.RequestHeader()))
	if err != nil {
		t.Fatalf("Error probing %s: %v", endpoint, err)
	}

	sendPostRequest := func(resolvableDomain bool, url *url.URL) (*spoof.Response, error) {
		t.Log("Request", url)
		client, err := pkgtest.NewSpoofingClient(context.Background(), clients.KubeClient, t.Logf, url.Hostname(), resolvableDomain, test.AddRootCAtoTransport(context.Background(), t.Logf, clients, test.ServingFlags.HTTPS))
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodPost, url.String(), nil)
		if err != nil {
			return nil, err
		}
		spoof.WithHeader(test.ServingFlags.RequestHeader())(req)
		return client.Do(req)
	}

	pokeCowForMB := func(mb int) error {
		u, _ := url.Parse(endpoint.String())
		q := u.Query()
		q.Set("bloat", strconv.Itoa(mb))
		u.RawQuery = q.Encode()
		response, err := sendPostRequest(test.ServingFlags.ResolvableDomain, u)
		if err != nil {
			return err
		}
		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("StatusCode = %d, want %d", response.StatusCode, http.StatusOK)
		}
		return nil
	}

	t.Log("Querying the application to see if the memory limits are enforced.")
	if err := pokeCowForMB(100); err != nil {
		t.Fatalf("Didn't get a response from bloating cow with %d MBs of Memory: %v", 100, err)
	}

	if err := pokeCowForMB(200); err != nil {
		t.Fatalf("Didn't get a response from bloating cow with %d MBs of Memory: %v", 200, err)
	}

	// ExceedingMemoryLimitSize defaults to 500.
	// Allows override the memory usage to get a non-200 resposne because the serverless platform
	// MAY automatically adjust the resource limits.
	// See https://github.com/knative/specs/blob/main/specs/serving/runtime-contract.md#memory-and-cpu-limits
	exceedingMemory := test.ServingFlags.ExceedingMemoryLimitSize
	if err := pokeCowForMB(exceedingMemory); err == nil {
		t.Fatalf("We shouldn't have got a response from bloating cow with %d MBs of Memory: %v", exceedingMemory, err)
	}
}
