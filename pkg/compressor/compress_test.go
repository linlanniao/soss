package compressor

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TestRaw = `apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    deployment.kubernetes.io/revision: "3"
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{},"name":"argo-server","namespace":"argo"},"spec":{"selector":{"matchLabels":{"app":"argo-server"}},"template":{"metadata":{"labels":{"app":"argo-server"}},"spec":{"containers":[{"args":["server"],"env":[],"image":"quay.io/argoproj/argocli:v3.5.4","name":"argo-server","ports":[{"containerPort":2746,"name":"web"}],"readinessProbe":{"httpGet":{"path":"/","port":2746,"scheme":"HTTPS"},"initialDelaySeconds":10,"periodSeconds":20},"securityContext":{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"readOnlyRootFilesystem":true,"runAsNonRoot":true},"volumeMounts":[{"mountPath":"/tmp","name":"tmp"}]}],"nodeSelector":{"kubernetes.io/os":"linux"},"securityContext":{"runAsNonRoot":true},"serviceAccountName":"argo-server","volumes":[{"emptyDir":{},"name":"tmp"}]}}}}
  creationTimestamp: "2024-02-27T04:00:47Z"
  generation: 3
  name: argo-server
  namespace: argo
  resourceVersion: "440410"
  uid: 5cd3999a-ddd1-40b1-ba18-4e7ae5dc07c2
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: argo-server
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: argo-server
    spec:
      containers:
      - args:
        - server
        - --auth-mode=server
        image: quay.io/argoproj/argocli:v3.5.4
        imagePullPolicy: IfNotPresent
        name: argo-server
        ports:
        - containerPort: 2746
          name: web
          protocol: TCP
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /
            port: 2746
            scheme: HTTPS
          initialDelaySeconds: 10
          periodSeconds: 20
          successThreshold: 1
          timeoutSeconds: 1
        resources: {}
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            drop:
            - ALL
          readOnlyRootFilesystem: true
          runAsNonRoot: true
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /tmp
          name: tmp
      dnsPolicy: ClusterFirst
      nodeSelector:
        kubernetes.io/os: linux
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext:
        runAsNonRoot: true
      serviceAccount: argo-server
      serviceAccountName: argo-server
      terminationGracePeriodSeconds: 30
      volumes:
      - emptyDir: {}
        name: tmp
status:
  availableReplicas: 1
  conditions:
  - lastTransitionTime: "2024-02-27T04:00:47Z"
    lastUpdateTime: "2024-02-27T08:17:16Z"
    message: ReplicaSet "argo-server-5fdfcf6fd6" has successfully progressed.
    reason: NewReplicaSetAvailable
    status: "True"
    type: Progressing
  - lastTransitionTime: "2024-03-22T08:42:00Z"
    lastUpdateTime: "2024-03-22T08:42:00Z"
    message: Deployment has minimum availability.
    reason: MinimumReplicasAvailable
    status: "True"
    type: Available
  observedGeneration: 3
  readyReplicas: 1
  replicas: 1
  updatedReplicas: 1
`

const TestRaw2 = `apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    deployment.kubernetes.io/revision: "1"
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{},"name":"workflow-controller","namespace":"argo"},"spec":{"selector":{"matchLabels":{"app":"workflow-controller"}},"template":{"metadata":{"labels":{"app":"workflow-controller"}},"spec":{"containers":[{"args":[],"command":["workflow-controller"],"env":[{"name":"LEADER_ELECTION_IDENTITY","valueFrom":{"fieldRef":{"apiVersion":"v1","fieldPath":"metadata.name"}}}],"image":"quay.io/argoproj/workflow-controller:v3.5.4","livenessProbe":{"failureThreshold":3,"httpGet":{"path":"/healthz","port":6060},"initialDelaySeconds":90,"periodSeconds":60,"timeoutSeconds":30},"name":"workflow-controller","ports":[{"containerPort":9090,"name":"metrics"},{"containerPort":6060}],"securityContext":{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"readOnlyRootFilesystem":true,"runAsNonRoot":true}}],"nodeSelector":{"kubernetes.io/os":"linux"},"priorityClassName":"workflow-controller","securityContext":{"runAsNonRoot":true},"serviceAccountName":"argo"}}}}
  creationTimestamp: "2024-02-27T04:00:47Z"
  generation: 1
  name: workflow-controller
  namespace: argo
  resourceVersion: "78424"
  uid: 1c7f9d58-fe2c-42fd-a09f-659fc203a8bc
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: workflow-controller
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: workflow-controller
    spec:
      containers:
      - command:
        - workflow-controller
        env:
        - name: LEADER_ELECTION_IDENTITY
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.name
        image: quay.io/argoproj/workflow-controller:v3.5.4
        imagePullPolicy: IfNotPresent
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /healthz
            port: 6060
            scheme: HTTP
          initialDelaySeconds: 90
          periodSeconds: 60
          successThreshold: 1
          timeoutSeconds: 30
        name: workflow-controller
        ports:
        - containerPort: 9090
          name: metrics
          protocol: TCP
        - containerPort: 6060
          protocol: TCP
        resources: {}
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            drop:
            - ALL
          readOnlyRootFilesystem: true
          runAsNonRoot: true
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      nodeSelector:
        kubernetes.io/os: linux
      priorityClassName: workflow-controller
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext:
        runAsNonRoot: true
      serviceAccount: argo
      serviceAccountName: argo
      terminationGracePeriodSeconds: 30
status:
  availableReplicas: 1
  conditions:
  - lastTransitionTime: "2024-02-27T04:01:13Z"
    lastUpdateTime: "2024-02-27T04:01:13Z"
    message: Deployment has minimum availability.
    reason: MinimumReplicasAvailable
    status: "True"
    type: Available
  - lastTransitionTime: "2024-02-27T04:00:47Z"
    lastUpdateTime: "2024-02-27T04:01:13Z"
    message: ReplicaSet "workflow-controller-5df59665c9" has successfully progressed.
    reason: NewReplicaSetAvailable
    status: "True"
    type: Progressing
  observedGeneration: 1
  readyReplicas: 1
  replicas: 1
  updatedReplicas: 1
`

func TestCompressDecompressBytes(t *testing.T) {
	t.Parallel()

	var (
		raw = []byte(TestRaw)
	)

	b1, err := CompressS2Bytes(raw)
	if err != nil {
		t.Fatalf("failed to compress bytes: %v", err)
	}

	raw2, err := DecompressS2Bytes(b1)
	if err != nil {
		t.Fatalf("failed to decompress bytes: %v", err)
	}

	if !bytes.Equal(raw, raw2) {
		t.Fatalf("bytes are not equal")
	}

}

func BenchmarkCompressS2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c, _ := CompressS2Str(TestRaw2)
		_, _ = DecompressS2Str(c)
	}
}

func string2File(content string, saveTo string) (string, error) {
	p := filepath.Dir(saveTo)
	err := os.MkdirAll(p, 0755)
	if err != nil {
		return "", fmt.Errorf("error creating dir: %v", err)
	}

	file, err := os.Create(saveTo)
	if err != nil {
		return "", fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Write the string content to the file
	_, err = file.WriteString(content)
	if err != nil {
		return "", fmt.Errorf("error writing to file: %v", err)
	}
	return filepath.Abs(saveTo)
}

func TestCreateTgzArchive(t *testing.T) {
	raws := []struct {
		str string
		out string
	}{
		{"aaa", "aaa.txt"},
		{"bbb", "bbb/bbb.txt"},
		{"ccc", "ccc/cccc/ccc.txt"},
	}
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "charts")
	_ = os.Mkdir(subDir, 0755)
	t.Logf(subDir)
	for _, r := range raws {
		d := filepath.Join(subDir, r.out)
		p, err := string2File(r.str, d)
		assert.NoError(t, err)
		t.Log(p)
	}

	err := CreateTgzArchive(subDir, filepath.Join(tmpDir, "test.tgz"))
	assert.NoError(t, err)
}
