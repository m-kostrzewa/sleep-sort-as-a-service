release:
	docker login
	docker build -t kosiak/sleepsort-aas:latest .
	docker push kosiak/sleepsort-aas:latest
	docker logout
	kubectl create -f k8s.yml
	kubectl get service sleepsort-service -w
