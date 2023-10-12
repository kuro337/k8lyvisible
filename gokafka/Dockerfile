# Go Dockerfile

FROM golang:1.21-alpine3.18 as builder
WORKDIR /app
# Copy the go module and sum files to download dependencies
COPY src/go.mod src/go.sum /app/
RUN go mod download
# Copy the source files
COPY src/ /app/
# Build the application with optimizations
RUN go build -o /app/main .
# Use a smaller base image for the runtime
FROM golang:1.21-alpine3.18
# Copy the built binary from the builder stage
COPY --from=builder /app/main /app/main
# Expose the port the app runs on
EXPOSE 8080
# Specify the command to run on container start
CMD ["/app/main"]




# BUILD
# docker build -t gokafka:latest .
# RUN
# docker run --name gokafka -d -p 8080:8080 gokafka:latest
# Enter Container
# docker run --rm -it gokafka:latest /bin/sh

# ```
# - Pushing to Image Registry
# ```bash
# # PUSH ghcr
# docker login ghcr.io -u kuro337 -p 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
# docker tag  gokafka:latest ghcr.io/kuro337/gokafka:latest
# docker push ghcr.io/kuro337/gokafka:latest
# ```
# - Deployment file to launch 2 pods - and a service to make pod accessible from outside
# ```yaml
# # go-dep.yaml
# apiVersion: apps/v1
# kind: Deployment
# metadata:
# name: go-app-deployment
# spec:
# replicas: 2
# selector:
#   matchLabels:
#   app: go-app
# template:
#   metadata:
#   labels:
#     app: go-app
#   spec:
#   containers:
#     - name: go-app
#     image: ghcr.io/chinmay337/ gokafka:latest
#     ports:
#       - containerPort: 8080
#   imagePullSecrets:
#     - name: ghcr-secret
# ```
# - Service for Deployment
# ```yaml
# apiVersion: v1
# kind: Service
# metadata:
# name: go-app-service
# labels:
#   app: go-app
# spec:
# ports:
#   - port: 8080
#   targetPort: 8080
# selector:
#   app: go-app
# clusterIP: None # This makes it a headless service
# ```
# - Create Secret so the image can be pulled from ghcr
# ```bash
# kubectl create secret docker-registry ghcr-secret --docker-server=ghcr.io --docker-username=kuro337 --docker-password=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
# ```
# - Deploying to k8
# ```bash
# # Making sure there's a k3d Cluster active
# k3d cluster list
# k3d cluster create
# # Deploy Service and Deployment
# kubectl apply -f go-app-dep.yaml
# kubectl apply -f go-app-svc.yaml
# # Check status or issues
# kubectl get pods -l app=go-app
# # Describe pod to see errors
# kubectl describe pod go-app-deployment-9cd9f7886-9k82w
# # Verify deployment
# kubectl get deployments
# kubectl get pods
# kubectl get services
# # Testing Pods from localhost
# # Accessing first pod through localhost:8081
# kubectl port-forward go-app-deployment-57c8fcd655-tmwq4 8081:8080
# # Accessing second pod through localhost:8082
# kubectl port-forward go-app-deployment-57c8fcd655-9z9tv 8082:8080
# # Delete
# kubectl delete deployment go-app-deployment
# kubectl delete svc go-app-service
# ```