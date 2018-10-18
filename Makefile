build:
	docker build --no-cache -t cephfs-quota:0.1 .
	docker tag quota-setter:0.1  peterbupt/quotasetter:0.1
deploy-only:
	kubectl create ns cephquota || true
	kubectl delete -f deployment || true
	kubectl apply -f deployment
clean:
	kubectl delete ns cephquota || true
deploy: build deploy-only

