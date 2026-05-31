.PHONY:
	rebuild logs shell down k8s-start k8s-stop k8s-build k8s-load k8s-deploy k8s-delete k8s-status

rebuild:
	docker compose down -v
	docker compose up -d --build

rebuild-cache:
	docker compose down
	docker compose up -d --build

down:
	docker compose down -v

logs:
	docker compose logs -f

k8s-start:
	minikube start

k8s-stop:
	minikube stop

k8s-build:
	eval $$(minikube docker-env) && \
	docker build -t inventory-service:local -f inventory-service/Dockerfile . && \
	docker build -t payment-service:local -f payment-service/Dockerfile . && \
	docker build -t order-service:local -f order-service/Dockerfile .

k8s-deploy:
	kubectl apply -k kubernetes/overlays/local

k8s-delete:
	kubectl delete -k kubernetes/overlays/local

k8s-status:
	kubectl get all -n observability

k8s-urls:
	@echo "Grafana:           http://localhost:3000"
	@echo "Jaeger:            http://localhost:16686"
	@echo "inventory-service: http://localhost:8001"
	@echo "order-service:     http://localhost:8002"
	@echo "payment-service:   http://localhost:8003"

k8s-forward:
	kubectl port-forward -n observability svc/grafana 3000:3000 & \
	kubectl port-forward -n observability svc/jaeger 16686:16686 & \
	kubectl port-forward -n observability svc/inventory-service 8001:8080 & \
	kubectl port-forward -n observability svc/order-service 8002:8080 & \
	kubectl port-forward -n observability svc/payment-service 8003:8080

# kubectl port-forward -n observability svc/order-service 8002:8080